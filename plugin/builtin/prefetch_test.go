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

package builtin

import (
	"log"
	"path/filepath"
	"testing"
)

func TestPrefetchPlugin_Run(t *testing.T) {
	log.Println("Start setup")
	storeDir, err := setup("example1.forensicstore")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Setup done")
	defer cleanup(storeDir)

	example1 := filepath.Join(storeDir, "example1.forensicstore")

	type args struct {
		url    string
		filter []string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		// {"Prefetch", args{"example1.forensicstore", nil, nil}, 261, false},
		{"Prefetch with Filter", args{example1, []string{"origin.path=%artifactcollector%"}}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlw := &testLineWriter{}
			command := &Prefetch{}

			command.Parameter().Set("filter", tt.args.filter)
			command.Parameter().Set("forensicstore", tt.args.url)
			err = command.Run(command, tlw)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tlw.lines) != tt.wantCount {
				t.Errorf("len(elements) = %v, wantCount %v", len(tlw.lines), tt.wantCount)
			}
		})
	}
}
