package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
)

type CSVOutput struct {
	dest      io.Writer
	csvWriter *csv.Writer
	headers   []string
}

func NewCSVOutput(dest io.Writer, headers []string) *CSVOutput {
	o := &CSVOutput{dest: dest, csvWriter: csv.NewWriter(dest), headers: headers}
	err := o.csvWriter.Write(o.headers)
	if err != nil {
		log.Println(err)
	}
	return o
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
