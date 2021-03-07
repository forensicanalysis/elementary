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

	"github.com/forensicanalysis/elementary/plugin"
	"github.com/forensicanalysis/elementary/plugin/output"
	"github.com/forensicanalysis/forensicstore"
)

func bulkSearch() plugin.Plugin {
	return &command{
		name:  "bulk-search",
		short: "Bulk search indicators",
		parameter: plugin.ParameterList{
			{Name: "file", Type: plugin.Path, Description: "file with IOCs", Required: true},
			ForensicStore, AddToStore, output.File, output.Format,
		},
		run: func(cmd plugin.Plugin) error {
			log.Printf("run bulk-search")

			path := cmd.Parameter().StringValue("forensicstore")
			iocListPath := cmd.Parameter().StringValue("file")

			store, teardown, err := forensicstore.Open(path)
			if err != nil {
				return err
			}
			defer teardown()

			file, err := os.Open(iocListPath) // #nosec
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			out := setupOut(cmd, store, []string{"ioc", "count"})
			defer out.WriteFooter()

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
		},
	}
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
