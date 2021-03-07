package output

import (
	"io"
)

var _ Writer = &JSONOutput{}

type JSONOutput struct {
	moreElements bool
	dest         io.Writer
}

func NewJsonOutput(dest io.Writer) *JSONOutput {
	o := &JSONOutput{dest: dest}
	o.dest.Write([]byte("[")) // nolint: errcheck
	return o
}

func (o *JSONOutput) WriteHeader([]string) {}

func (j *JSONOutput) WriteLine(element []byte) {
	if j.moreElements {
		j.dest.Write([]byte(",")) // nolint: errcheck
	} else {
		j.moreElements = true
	}
	j.dest.Write(element) // nolint: errcheck
}

func (j *JSONOutput) WriteFooter() {
	j.dest.Write([]byte("]")) // nolint: errcheck
}
