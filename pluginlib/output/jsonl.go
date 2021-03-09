package output

import (
	"fmt"
	"io"
	"log"
)

type JSONLOutput struct {
	dest io.Writer
}

func NewJsonlOutput(dest io.Writer) *JSONLOutput {
	return &JSONLOutput{dest: dest}
}

func (o *JSONLOutput) WriteLine(element []byte) {
	_, err := fmt.Fprintln(o.dest, string(element))
	if err != nil {
		log.Println(err)
	}
}
