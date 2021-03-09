package output

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

type TableOutput struct {
	dest        io.Writer
	tableWriter *tablewriter.Table
	headers     []string
}

func NewTableOutput(dest io.Writer, headers []string) *TableOutput {
	o := &TableOutput{dest: dest, tableWriter: tablewriter.NewWriter(dest), headers: headers}
	o.tableWriter.SetHeader(o.headers)
	return o
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
