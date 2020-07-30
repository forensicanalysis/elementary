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
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/tidwall/gjson"
)

func TestExportTimesketch(t *testing.T) {
	log.Println("Start setup")
	storeDir, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Setup done")
	defer cleanup(storeDir)

	example1 := filepath.Join(storeDir, "example1.forensicstore")

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
		{"export timesketch", args{example1, []string{}}, 1746, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := exportTimesketch()

			command.Flags().Set("format", "jsonl")
			command.Flags().Set("output", filepath.Join(storeDir, "out.jsonl"))
			command.SetArgs(append(tt.args.args, tt.args.url))
			err = command.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			file, err := os.Open(filepath.Join(storeDir, "out.jsonl"))
			if err != nil {
				t.Fatal(err)
			}
			fileScanner := bufio.NewScanner(file)
			lineCount := 0
			for fileScanner.Scan() {
				if lineCount < 10 {
					fmt.Println(fileScanner.Text())
				}
				lineCount++
			}
			err = file.Close()
			if err != nil {
				t.Fatal(err)
			}

			if lineCount != tt.wantCount {
				t.Errorf("Run() error, wrong number of resuls = %d, want %d", lineCount, tt.wantCount)
			}
		})
	}
}

func Test_jsonToText(t *testing.T) {
	type args struct {
		element gjson.Result
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"dict", args{gjson.Get(`{"a": "b"}`, "@this")}, "a: b"},
		{"list", args{gjson.Get(`["a", "b"]`, "@this")}, "a, b"},
		{"complex 1", args{gjson.Get(`{"a": ["b", "c"]}`, "@this")}, "a: b, c"},
		{"complex 2", args{gjson.Get(`{"a": ["b", "c"], "x": [1, 2]}`, "@this")}, "a: b, c; x: 1, 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonToText(&tt.args.element); got != tt.want {
				t.Errorf("jsonToText() = %v, want %v", got, tt.want)
			}
		})
	}
}
