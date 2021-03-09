package pluginlib

import (
	"bytes"
	"log"

	"github.com/tidwall/gjson"
)

type LineWriter interface {
	WriteLine([]byte)
}

type MultiLineWriter struct {
	LineWriter []LineWriter
}

func (m *MultiLineWriter) WriteLine(b []byte) {
	for _, lw := range m.LineWriter {
		lw.WriteLine(b)
	}
}

type LineWriterBuffer struct {
	// firstLine bool
	buffer *bytes.Buffer
	Writer LineWriter
}

func (o *LineWriterBuffer) Write(b []byte) (n int, err error) {
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

func (o *LineWriterBuffer) writeLine(element []byte) { // nolint: gocyclo
	element = bytes.TrimSpace(element)

	/*
		if o.firstLine {
			o.firstLine = false
			config := &Config{}
			err := json.Unmarshal(element, config)
			if err != nil || len(config.Header) == 0 {
				log.Printf("could not unmarshal config: %s, '%s'", err, element)
				return
			}

			// o.Writer.WriteHeader(config.Header)
			return
		}
	*/

	// print to output
	if !gjson.ValidBytes(element) {
		log.Println(element)
		return
	}

	o.Writer.WriteLine(element)
}

func (o *LineWriterBuffer) WriteFooter() {
	if o.buffer.Len() > 0 {
		o.writeLine(o.buffer.Bytes())
	}

	// o.Writer.WriteFooter()
}
