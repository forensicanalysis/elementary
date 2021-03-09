package output

import (
	"io"
)

type JSONOutput struct {
	moreElements bool
	dest         io.Writer
}

func NewJSONOutput(dest io.Writer) *JSONOutput {
	o := &JSONOutput{dest: dest}
	o.dest.Write([]byte("[")) // nolint: errcheck
	return o
}

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
