package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
)

var _ Writer = &CSVOutput{}

type CSVOutput struct {
	dest      io.Writer
	csvWriter *csv.Writer
	headers   []string
}

func NewCSVOutput(dest io.Writer) *CSVOutput {
	o := &CSVOutput{dest: dest, csvWriter: csv.NewWriter(dest)}
	return o
}

func (o *CSVOutput) WriteHeader(headers []string) {
	o.headers = headers
	err := o.csvWriter.Write(o.headers)
	if err != nil {
		log.Println(err)
	}
}

func (o *CSVOutput) WriteLine(element []byte) {
	if o.headers == nil {
		fmt.Fprintln(o.dest, string(element)) // nolint: errcheck
		return
	}
	err := o.csvWriter.Write(getColumns(o.headers, element))
	if err != nil {
		fmt.Fprintln(o.dest, string(element)) // nolint: errcheck
		log.Println(err)
	}
}

func (o *CSVOutput) WriteFooter() {
	o.csvWriter.Flush()
}
