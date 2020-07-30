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

package daggy

import (
	"log"
	"testing"
)

/*
func TestParse(t *testing.T) {
	workflow := &Workflow{
		Tasks: []Task{
			{Command: "hotfixes"},
			{Command: "networking"},
			{Command: "prefetch"},
			{Command: "run-keys"},
			{Command: "services"},
			{Command: "shimcache"},
			{Command: "software"},
		},
	}

	type args struct {
		workflowFile string
	}
	tests := []struct {
		name    string
		args    args
		want    *Workflow
		wantErr bool
	}{
		{"Parse example-workflow.yml", args{"../test/data/test.yml"}, workflow, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.workflowFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("Parse() got = %#v, want %#v", got, tt.want)
				return
			}

			got.graph = nil

			if !reflect.DeepEqual(got.Tasks, tt.want.Tasks) {
				t.Errorf("Parse() got = %#v, want %#v", got.Tasks, tt.want.Tasks)
			}
		})
	}
}
*/

func Test_setupLogging(t *testing.T) {
	setupLogging()
	log.Print("test")
}
