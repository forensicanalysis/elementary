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

package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tidwall/gjson"
)

type Writer interface {
	WriteHeader(header []string)
	WriteLine(element []byte)
	WriteFooter()
}

type Config struct {
	Header []string `json:"header,omitempty"`
}

type MainWriter struct {
	firstLine bool

	buffer *bytes.Buffer

	dest    io.Writer
	writers []Writer
}

func New(path, format string, config *Config) *MainWriter {
	output := &MainWriter{
		buffer: &bytes.Buffer{},
	}

	var dest io.Writer
	if path != "" {
		f, err := os.Create(path)
		if err != nil {
			log.Println(err)
		}
		dest = f
	} else {
		dest = os.Stdout
	}

	switch format {
	case "table":
		output.writers = append(output.writers, NewTableOutput(dest))
	case "csv":
		output.writers = append(output.writers, NewCSVOutput(dest))
	case "jsonl":
		output.writers = append(output.writers, NewJsonlOutput(dest))
	case "json":
		output.writers = append(output.writers, NewJsonOutput(dest))
	case "none":
	default:
		log.Printf("unknown writer %s\n", format)
	}

	if config == nil {
		output.firstLine = true
	} else {
		output.writeHeaderConfig(config)
	}

	return output
}

func (o *MainWriter) Add(writer Writer) {
	o.writers = append(o.writers, writer)
}

func (o *MainWriter) writeHeaderLine(line []byte) {
	config := &Config{}
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

func (o *MainWriter) Write(b []byte) (n int, err error) {
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

func (o *MainWriter) writeHeaderConfig(outConfig *Config) {
	o.firstLine = false

	for _, writer := range o.writers {
		writer.WriteHeader(outConfig.Header)
	}
}

func (o *MainWriter) WriteLine(element []byte) { // nolint: gocyclo
	element = bytes.TrimSpace(element)

	if o.firstLine {
		o.writeHeaderLine(element)
		return
	}

	// print to output
	if !gjson.ValidBytes(element) {
		_, err := fmt.Fprintln(o.dest, string(element))
		if err != nil {
			log.Println(err)
		}
	}

	for _, writer := range o.writers {
		writer.WriteLine(element)
	}
}

func (o *MainWriter) WriteFooter() {
	if o.buffer.Len() > 0 {
		o.WriteLine(o.buffer.Bytes())
	}

	for _, writer := range o.writers {
		writer.WriteFooter()
	}

	out := o.dest
	if closer, ok := out.(io.Closer); ok && out != os.Stdout {
		closer.Close()
	}
}

func getColumns(headers []string, element []byte) []string {
	var columns []string
	for _, header := range headers {
		value := gjson.GetBytes(element, header)
		if value.Exists() {
			columns = append(columns, value.String())
		} else {
			columns = append(columns, "")
		}
	}
	return columns
}
