// Copyright (c) 2019 Siemens AG
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
	"io/ioutil"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

// A Task is a single element in a workflow yml file.
type Task struct {
	ID        uuid.UUID              `yaml:"id"`
	Command   string                 `yaml:"command"`
	Arguments map[string]interface{} `yaml:"arguments"`
}

// Workflow can be used to parse workflow yml files.
type Workflow struct {
	Tasks []Task `yaml:"tasks"`
}

// Parse reads a workflow file.
func Parse(workflowFile string) (*Workflow, error) {
	// parse the yaml definition
	data, err := ioutil.ReadFile(workflowFile) // #nosec
	if err != nil {
		return nil, err
	}
	workflow := Workflow{}
	err = yaml.Unmarshal(data, &workflow)
	return &workflow, err
}
