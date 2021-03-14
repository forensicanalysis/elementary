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
	"bufio"
	"encoding/json"
	"log"
	"os"

	"crawshaw.io/sqlite"

	"github.com/forensicanalysis/elementary/pluginlib"
)

var _ pluginlib.Plugin = &BulkSearch{}

type BulkSearch struct {
	parameter pluginlib.ParameterList
}

func (b *BulkSearch) Name() string {
	return "bulk-search"
}

func (b *BulkSearch) Short() string {
	return "Bulk search indicators"
}

func (b *BulkSearch) Parameter() pluginlib.ParameterList {
	if b.parameter == nil {
		b.parameter = pluginlib.ParameterList{
			{Name: "forensicstore", Type: pluginlib.Path, Description: "forensicstore", Required: true, Argument: true},
			{Name: "file", Type: pluginlib.Path, Description: "file with IOCs", Required: true},
		}
	}
	return b.parameter
}

func (b *BulkSearch) Output() *pluginlib.Config {
	return &pluginlib.Config{Header: []string{"type", "ioc", "count"}}
}

func (b *BulkSearch) Run(p pluginlib.Plugin, out pluginlib.LineWriter) error {
	store, teardown, err := getForensicStore(p)
	if err != nil {
		return err
	}
	defer teardown()

	iocListPath := p.Parameter().StringValue("file")

	file, err := os.Open(iocListPath) // #nosec
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ioc := scanner.Text()
		if ioc == "" {
			continue
		}
		log.Println("search", ioc)

		element, err := getSearchCount(store.Connection(), ioc)
		if err != nil {
			return err
		}
		out.WriteLine(element) // nolint: errcheck
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func getSearchCount(conn *sqlite.Conn, term string) ([]byte, error) {
	stmt, err := conn.Prepare("SELECT count(json) as count FROM elements WHERE json LIKE $search")
	if err != nil {
		return nil, err
	}
	stmt.SetText("$search", "%"+term+"%")
	if _, err := stmt.Step(); err != nil {
		return nil, err
	}

	b, err := json.Marshal(struct {
		Type  string `json:"type"`
		IOC   string `json:"ioc"`
		Count int64  `json:"count"`
	}{
		Type:  "bulksearch",
		IOC:   term,
		Count: stmt.GetInt64("count"),
	})
	if err != nil {
		return nil, err
	}

	return b, stmt.Finalize()
}
