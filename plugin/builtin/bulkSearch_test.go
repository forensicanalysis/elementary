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
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/forensicanalysis/elementary/plugin"

	"github.com/tidwall/gjson"

	"github.com/forensicanalysis/forensicstore"
)

func TestBulkSearch(t *testing.T) {
	log.Println("Start setup")
	storeDir, err := setup("example1.forensicstore")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Setup done")
	defer cleanup(storeDir)

	example1 := filepath.Join(storeDir, "example1.forensicstore")

	iocFile := filepath.Join(storeDir, "ioc.txt")
	ioutil.WriteFile(iocFile, []byte("exe"), 0755)

	type args struct {
		url  string
		file string
	}
	tests := []struct {
		name        string
		args        args
		wantResults int
		wantCount   int64
		wantErr     bool
	}{
		{"ioc search", args{example1, iocFile}, 1, 529, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := bulkSearch()

			command.Parameter().Set("format", "none")
			command.Parameter().Set("add-to-store", true)
			command.Parameter().Set("file", tt.args.file)
			command.Parameter().Set("forensicstore", tt.args.url)

			err = command.Run(command)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			store, teardown, err := forensicstore.Open(tt.args.url)
			if err != nil {
				t.Fatal(err)
			}
			defer teardown()
			elements, err := store.Select(plugin.Filter{{"type": "bulksearch"}})
			if err != nil {
				t.Fatal(err)
			}

			if len(elements) != tt.wantResults {
				t.Errorf("Run() error, wrong number of resuls = %d, want %d", len(elements), tt.wantResults)
			}

			count := gjson.GetBytes(elements[0], "count").Int()
			if count != tt.wantCount {
				t.Errorf("Run() error, wrong count of resuls = %d, want %d", count, tt.wantCount)
			}
		})
	}
}
