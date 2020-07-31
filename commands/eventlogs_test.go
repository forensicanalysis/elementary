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

package commands

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/forensicanalysis/elementary/daggy"
	"github.com/forensicanalysis/forensicstore"
)

func TestEventlogsPlugin_Run(t *testing.T) {
	log.Println("Start setup")
	storeDir, err := setup("example2.forensicstore")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Setup done")
	defer cleanup(storeDir)

	example2 := filepath.Join(storeDir, "example2.forensicstore")

	type args struct {
		url  string
		args []string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{"eventlogs Test", args{example2, nil}, 806, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := eventlogs()

			command.Flags().Set("format", "none")
			command.Flags().Set("add-to-store", "true")
			command.SetArgs(append(tt.args.args, tt.args.url))
			err = command.Execute()

			if (err != nil) != tt.wantErr {
				t.Fatalf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			store, teardown, err := forensicstore.Open(tt.args.url)
			if err != nil {
				t.Fatalf("forensicstore.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer teardown()

			elements, err := store.Select(daggy.Filter{{"type": "eventlog"}})
			if err != nil {
				t.Fatalf("store.Select() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(elements) != tt.wantCount {
				t.Fatalf("len(elements) = %v, wantCount %v", len(elements), tt.wantCount)
			}
		})
	}
}
