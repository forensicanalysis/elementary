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
	"encoding/json"
	"log"
	"os"

	"crawshaw.io/sqlite"
	"github.com/spf13/cobra"

	"github.com/forensicanalysis/forensicstore"
)

func bulkSearch() *cobra.Command {
	var file string
	bulkSearchCommand := &cobra.Command{
		Use:   "bulk-search <forensicstore>",
		Short: "Bulk search indicators",
		Args:  RequireStore,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("run bulk-search %s", args)

			store, teardown, err := forensicstore.Open(args[0])
			if err != nil {
				return err
			}
			defer teardown()

			file, err := os.Open(file) // #nosec
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			output := newOutputWriterStore(cmd, store, &outputConfig{Header: []string{"ioc", "count"}})

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
				output.writeLine(element) // nolint: errcheck
			}
			if err := scanner.Err(); err != nil {
				return err
			}
			output.WriteFooter()
			return nil
		},
	}
	addOutputFlags(bulkSearchCommand)
	bulkSearchCommand.Flags().StringVar(&file, "file", "", "file with IOCs")
	_ = bulkSearchCommand.MarkFlagRequired("file")
	return bulkSearchCommand
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
