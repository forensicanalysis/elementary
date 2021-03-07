package output

import (
	"fmt"
	"io"
	"log"
)

var _ Writer = &JSONLOutput{}

type JSONLOutput struct {
	dest io.Writer
}

func NewJsonlOutput(dest io.Writer) *JSONLOutput {
	o := &JSONLOutput{dest: dest}
	o.dest.Write([]byte("[")) // nolint: errcheck
	return o
}

func (o *JSONLOutput) WriteHeader([]string) {}

func (o *JSONLOutput) WriteLine(element []byte) {
	_, err := fmt.Fprintln(o.dest, string(element))
	if err != nil {
		log.Println(err)
	}
}

func (o *JSONLOutput) WriteFooter() {
	o.dest.Write([]byte("]")) // nolint: errcheck
}
