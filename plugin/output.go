// Copyright (c) 2020 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

package plugin

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"

	"github.com/forensicanalysis/forensicstore"
)

type format int

const (
	tableFormat format = iota
	csvFormat
	jsonlFormat
	noneFormat
	jsonFormat
)

func fromString(s string) format {
	for i, f := range []string{"table", "csv", "jsonl", "none", "json"} {
		if s == f {
			return format(i)
		}
	}
	return tableFormat
}

type OutputConfig struct {
	Header []string `json:"header,omitempty"`
}

type OutputWriter struct {
	format       format
	store        *forensicstore.ForensicStore
	config       *OutputConfig
	firstLine    bool
	moreElements bool
	cmd          Plugin
	dest         io.Writer

	buffer *bytes.Buffer

	tableWriter *tablewriter.Table
	csvWriter   *csv.Writer
}

func newOutputWriter(store *forensicstore.ForensicStore, cmd Plugin) *OutputWriter {
	o, format, addToStore := parseOutputFlags(cmd)
	outStore := store
	if !addToStore {
		outStore = nil
	}

	output := &OutputWriter{
		format: format,
		store:  outStore,
		cmd:    cmd,
		buffer: &bytes.Buffer{},
	}

	if o != "" {
		f, err := os.Create(o)
		if err != nil {
			log.Println(err)
		}
		output.dest = f
	} else {
		output.dest = os.Stdout
	}

	switch format {
	case csvFormat:
		output.csvWriter = csv.NewWriter(output.dest)
	case tableFormat:
		output.tableWriter = tablewriter.NewWriter(output.dest)
	case jsonFormat:
		output.dest.Write([]byte("[")) // nolint: errcheck
	}

	return output
}

func NewOutputWriterStore(cmd Plugin, store *forensicstore.ForensicStore, config *OutputConfig) *OutputWriter {
	o := newOutputWriter(store, cmd)
	o.writeHeaderConfig(config)
	return o
}

func NewOutputWriterURL(cmd Plugin, url string) (*OutputWriter, func() error) {
	var store *forensicstore.ForensicStore
	teardown := func() error { return nil }
	if cmd.Parameter().BoolValue("add-to-store") {
		var err error
		store, teardown, err = forensicstore.Open(url)
		if err != nil {
			store = nil
		}
	}
	o := newOutputWriter(store, cmd)
	o.firstLine = true
	return o, teardown
}

func (o *OutputWriter) writeHeaderLine(line []byte) {
	config := &OutputConfig{}
	err := json.Unmarshal(line, config)
	if err != nil || len(config.Header) == 0 {
		log.Printf("could not unmarshal config: %s, '%s'", err, line)
		_, err = fmt.Fprintln(o.dest, string(line))
		if err != nil {
			log.Println(err)
		}
		return
	}

	o.writeHeaderConfig(config)
}

func (o *OutputWriter) writeHeaderConfig(outConfig *OutputConfig) {
	o.config = outConfig
	o.firstLine = false

	switch o.format {
	case tableFormat:
		o.tableWriter.SetHeader(o.config.Header)
	case csvFormat:
		err := o.csvWriter.Write(o.config.Header)
		if err != nil {
			log.Println(err)
		}
	case jsonlFormat, noneFormat, jsonFormat:
	default:
		log.Println("unknown output format:", o.format)
	}
}

func (o *OutputWriter) Write(b []byte) (n int, err error) {
	n = len(b)
	for {
		i := bytes.IndexByte(b, '\n')
		if i < 0 { // no more newlines
			o.buffer.Write(b)
			break
		}

		o.WriteLine(append(o.buffer.Bytes(), b[:i]...))
		b = b[i+1:]
		o.buffer.Reset()
	}
	return n, nil
}

func (o *OutputWriter) WriteLine(element []byte) { // nolint: gocyclo
	element = bytes.TrimSpace(element)

	if o.firstLine {
		o.writeHeaderLine(element)
		return
	}

	// print to output
	switch {
	case !gjson.ValidBytes(element) ||
		o.format == jsonlFormat ||
		(o.format == tableFormat && o.config == nil) ||
		(o.format == csvFormat && o.config == nil):
		_, err := fmt.Fprintln(o.dest, string(element))
		if err != nil {
			log.Println(err)
		}
	case o.format == tableFormat:
		o.tableWriter.Append(o.getColumns(element))
	case o.format == csvFormat:
		err := o.csvWriter.Write(o.getColumns(element))
		if err != nil {
			fmt.Fprintln(o.dest, string(element)) // nolint: errcheck
			log.Println(err)
		}
	case o.format == jsonFormat:
		if o.moreElements {
			o.dest.Write([]byte(",")) // nolint: errcheck
		} else {
			o.moreElements = true
		}
		o.dest.Write(element) // nolint: errcheck
	}

	// add to forensicstore
	if o.store != nil {
		_, err := o.store.Insert(element)
		if err != nil {
			log.Println(err, string(element))
		}
	}
}

func (o *OutputWriter) getColumns(element forensicstore.JSONElement) []string {
	var columns []string
	for _, header := range o.config.Header {
		value := gjson.GetBytes(element, header)
		if value.Exists() {
			columns = append(columns, value.String())
		} else {
			columns = append(columns, "")
		}
	}
	return columns
}

func (o *OutputWriter) WriteFooter() {
	if o.buffer.Len() > 0 {
		o.WriteLine(o.buffer.Bytes())
	}

	switch o.format {
	case csvFormat:
		o.csvWriter.Flush()
	case tableFormat:
		if o.tableWriter.NumLines() > 0 {
			o.tableWriter.Render()
		}
	case jsonFormat:
		o.dest.Write([]byte("]")) // nolint: errcheck
	}

	out := o.dest
	if closer, ok := out.(io.Closer); ok && out != os.Stdout {
		closer.Close()
	}
}

func parseOutputFlags(cmd Plugin) (string, format, bool) {
	output := cmd.Parameter().StringValue("output")

	formatString := cmd.Parameter().StringValue("format")
	format := fromString(formatString)

	addToStore := false
	if !HasAnnotation(cmd, Exporter) {
		addToStore = cmd.Parameter().BoolValue("add-to-store")
	}

	return output, format, addToStore
}

func OutputParameter(cmd Plugin) []*Parameter {
	parameter := []*Parameter{
		{
			Name:        "output",
			Description: "choose an output file",
			Type:        Path,
			Value:       "",
			Required:    false,
		},
		{Name: "format",
			Description: "choose output format [csv, jsonl, table, json, none]",
			Type:        String,
			Value:       "table",
			Required:    false,
		},
	}

	if !HasAnnotation(cmd, Exporter) {
		parameter = append(parameter,
			&Parameter{
				Name:        "add-to-store",
				Description: "choose an output file",
				Type:        Bool,
				Value:       false,
				Required:    false,
			},
		)
	}

	return parameter
}
