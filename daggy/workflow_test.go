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
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/forensicanalysis/forensicstore"
	"github.com/otiai10/copy"
)

func setup() (storeDir string, err error) {
	tempDir, err := ioutil.TempDir("", "forensicstoreprocesstest")
	if err != nil {
		return "", err
	}
	storeDir = filepath.Join(tempDir, "test")
	err = os.MkdirAll(storeDir, 0755)
	if err != nil {
		return "", err
	}

	err = copy.Copy(filepath.Join("..", "test"), storeDir)
	if err != nil {
		return "", err
	}

	return storeDir, nil
}

func cleanup(folders ...string) (err error) {
	for _, folder := range folders {
		err := os.RemoveAll(folder)
		if err != nil {
			return err
		}
	}
	return nil
}

var _ Command = &testCommand{}

type testCommand struct {
	name string
	run  func(command Command) error
}

func (t *testCommand) Name() string {
	return t.name
}

func (t *testCommand) Short() string {
	return t.name
}

func (t *testCommand) Parameter() ParameterList {
	return nil
}

func (t *testCommand) Run(command Command) error {
	return t.run(command)
}

func (t *testCommand) Annotations() []Annotation {
	return nil
}

func Test_processTask(t *testing.T) {
	log.Println("Start setup")
	storeDir, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Setup done")
	defer cleanup(storeDir)

	type args struct {
		task Task
	}
	tests := []struct {
		name      string
		storeName string
		args      args
		wantType  string
		wantCount int
		wantErr   bool
	}{
		{"dummy plugin", "example1.forensicstore", args{Task{Command: "example"}}, "example", 0, false},
		{"command not existing", "example1.forensicstore", args{Task{Command: "foo"}}, "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workflow := &Workflow{Tasks: []Task{tt.args.task}}

			plugins := []Command{&testCommand{
				name: "example",
				run: func(cmd Command) error {
					return nil
				},
			}}

			engine := New(plugins)

			if err := engine.Run(workflow, filepath.Join(storeDir, tt.storeName)); (err != nil) != tt.wantErr {
				t.Errorf("runTask() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				store, teardown, err := forensicstore.New(filepath.Join(storeDir, tt.storeName))
				if err != nil {
					t.Fatal(err)
				}
				defer teardown()

				log.Println("Start select")
				if tt.wantCount > 0 {
					elements, err := store.Select(Filter{{"type": tt.wantType}})
					if err != nil {
						t.Fatal(err)
					}
					if tt.wantCount != len(elements) {
						t.Errorf("runTask() error, wrong number of resuls = %d, want %d (%v)", len(elements), tt.wantCount, len(elements))
					}
				}
			}
		})
	}
}

func Test_toCmdline(t *testing.T) {
	var i interface{}
	i = []map[string]string{
		{"foo": "bar", "bar": "baz"},
		{"a": "b"},
	}

	type args struct {
		name string
		i    interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"filter", args{"filter", i}, []string{"--filter", "bar=baz,foo=bar", "--filter", "a=b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toCmdline(tt.args.name, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toCmdline() = %v, want %v", got, tt.want)
			}
		})
	}
}
