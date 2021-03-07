package output

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

var _ Writer = &TableOutput{}

type TableOutput struct {
	dest        io.Writer
	tableWriter *tablewriter.Table
	headers     []string
}

func NewTableOutput(dest io.Writer) *TableOutput {
	return &TableOutput{dest: dest, tableWriter: tablewriter.NewWriter(dest)}
}

func (o *TableOutput) WriteHeader(headers []string) {
	o.headers = headers
	o.tableWriter.SetHeader(o.headers)
}

func (o *TableOutput) WriteLine(element []byte) {
	if o.headers == nil {
		fmt.Fprintln(o.dest, string(element)) // nolint: errcheck
		return
	}
	o.tableWriter.Append(getColumns(o.headers, element))
}

func (o *TableOutput) WriteFooter() {
	if o.tableWriter.NumLines() > 0 {
		o.tableWriter.Render()
	}
}
