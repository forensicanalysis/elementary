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

package commands

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
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

type outputConfig struct {
	Header []string `json:"header,omitempty"`
}

type outputWriter struct {
	format       format
	store        *forensicstore.ForensicStore
	config       *outputConfig
	firstLine    bool
	moreElements bool
	cmd          *cobra.Command

	buffer *bytes.Buffer

	tableWriter *tablewriter.Table
	csvWriter   *csv.Writer
}

func newOutputWriter(store *forensicstore.ForensicStore, cmd *cobra.Command) *outputWriter {
	format, addToStore := parseOutputFlags(cmd)
	outStore := store
	if !addToStore {
		outStore = nil
	}

	output := &outputWriter{
		format: format,
		store:  outStore,
		cmd:    cmd,
		buffer: &bytes.Buffer{},
	}

	switch format {
	case csvFormat:
		output.csvWriter = csv.NewWriter(cmd.OutOrStdout())
	case tableFormat:
		output.tableWriter = tablewriter.NewWriter(cmd.OutOrStdout())
	case jsonFormat:
		cmd.OutOrStdout().Write([]byte("[")) // nolint: errcheck
	}

	return output
}

func newOutputWriterStore(cmd *cobra.Command, store *forensicstore.ForensicStore, config *outputConfig) *outputWriter {
	o := newOutputWriter(store, cmd)
	o.writeHeaderConfig(config)
	return o
}

func newOutputWriterURL(cmd *cobra.Command, url string) (*outputWriter, func() error) {
	var store *forensicstore.ForensicStore
	teardown := func() error { return nil }
	addToStore, err := cmd.Flags().GetBool("add-to-store")
	if err == nil && addToStore {
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

func (o *outputWriter) writeHeaderLine(line []byte) {
	config := &outputConfig{}
	err := json.Unmarshal(line, config)
	if err != nil || len(config.Header) == 0 {
		log.Printf("could not unmarshal config: %s, '%s'", err, line)
		_, err = fmt.Fprintln(o.cmd.OutOrStdout(), string(line))
		if err != nil {
			log.Println(err)
		}
		return
	}

	o.writeHeaderConfig(config)
}

func (o *outputWriter) writeHeaderConfig(outConfig *outputConfig) {
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

func (o *outputWriter) Write(b []byte) (n int, err error) {
	n = len(b)
	for {
		i := bytes.IndexByte(b, '\n')
		if i < 0 { // no more newlines
			o.buffer.Write(b)
			break
		}

		o.writeLine(append(o.buffer.Bytes(), b[:i]...))
		b = b[i+1:]
		o.buffer.Reset()
	}
	return n, nil
}

func (o *outputWriter) writeLine(element []byte) { // nolint: gocyclo
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
		_, err := fmt.Fprintln(o.cmd.OutOrStdout(), string(element))
		if err != nil {
			log.Println(err)
		}
	case o.format == tableFormat:
		o.tableWriter.Append(o.getColumns(element))
	case o.format == csvFormat:
		err := o.csvWriter.Write(o.getColumns(element))
		if err != nil {
			fmt.Fprintln(o.cmd.OutOrStdout(), string(element)) // nolint: errcheck
			log.Println(err)
		}
	case o.format == jsonFormat:
		if o.moreElements {
			o.cmd.OutOrStdout().Write([]byte(",")) // nolint: errcheck
		} else {
			o.moreElements = true
		}
		o.cmd.OutOrStdout().Write(element) // nolint: errcheck
	}

	// add to forensicstore
	if o.store != nil {
		_, err := o.store.Insert(element)
		if err != nil {
			log.Println(err, string(element))
		}
	}
}

func (o *outputWriter) getColumns(element forensicstore.JSONElement) []string {
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

func (o *outputWriter) WriteFooter() {
	if o.buffer.Len() > 0 {
		o.writeLine(o.buffer.Bytes())
	}

	switch o.format {
	case csvFormat:
		o.csvWriter.Flush()
	case tableFormat:
		if o.tableWriter.NumLines() > 0 {
			o.tableWriter.Render()
		}
	case jsonFormat:
		o.cmd.OutOrStdout().Write([]byte("]")) // nolint: errcheck
	}

	out := o.cmd.OutOrStdout()
	if closer, ok := out.(io.Closer); ok && out != os.Stdout {
		closer.Close()
	}
}

func addOutputFlags(cmd *cobra.Command) {
	cmd.Flags().String("output", "", "choose an output file")
	cmd.Flags().String("format", "table", "choose output format [csv, jsonl, table, json, none]")

	if cmd.Annotations != nil {
		if properties, ok := cmd.Annotations["plugin_property_flags"]; ok {
			if strings.Contains(properties, "ex") {
				return
			}
		}
	}
	cmd.Flags().Bool("add-to-store", false, "additionally save output to store")
}

func parseOutputFlags(cmd *cobra.Command) (format, bool) {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		log.Println(err)
	}
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			log.Println(err)
		}
		cmd.SetOut(f)
	}

	formatString, err := cmd.Flags().GetString("format")
	if err != nil {
		log.Println(err)
	}
	format := fromString(formatString)

	addToStore, err := cmd.Flags().GetBool("add-to-store")
	if err != nil {
		log.Println(err)
	}

	return format, addToStore
}
