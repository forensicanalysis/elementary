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
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/forensicanalysis/elementary/daggy"
	"github.com/forensicanalysis/forensicstore"
)

func exportTimesketch() *cobra.Command {
	var filtersets []string
	outputCommand := &cobra.Command{
		Use:   "export-timesketch <forensicstore>",
		Short: "Export in timesketch jsonl format",
		Args:  RequireStore,
		RunE: func(cmd *cobra.Command, args []string) error {
			return exportStore(args[0], extractFilter(filtersets), cmd)
		},
		Annotations: map[string]string{"plugin_property_flags": "ex"},
	}
	addOutputFlags(outputCommand)
	outputCommand.Flags().StringArrayVar(&filtersets, "filter", nil, "filter processed events")
	return outputCommand
}

func exportStore(url string, filter daggy.Filter, cmd *cobra.Command) error {
	store, teardown, err := forensicstore.Open(url)
	if err != nil {
		return err
	}
	defer teardown()

	elements, err := store.Select(filter)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return nil
	}

	output := newOutputWriterStore(cmd, store, &outputConfig{Header: []string{"message", "datetime", "timestamp_desc"}})

	for _, element := range elements {
		element := element
		gjson.GetBytes(element, "@this").ForEach(func(key, value gjson.Result) bool {
			field := key.String()
			if field == "atime" || field == "ctime" || field == "mtime" || strings.HasSuffix(field, "_time") {
				t, err := time.Parse(time.RFC3339Nano, value.String())
				if err != nil {
					return true
				}

				jsonResult := gjson.GetBytes(element, "@this")

				b, err := json.Marshal(struct {
					Type          string `json:"type"`
					Message       string `json:"message"`
					Datetime      string `json:"datetime"`
					TimestampDesc string `json:"timestamp_desc"`
				}{
					Type:          "timesketch",
					Message:       jsonToText(&jsonResult),
					Datetime:      t.UTC().Format(time.RFC3339Nano),
					TimestampDesc: field,
				})
				if err != nil {
					log.Println(err)
					return true
				}
				output.writeLine(b) // nolint: errcheck
			}
			return true
		})
	}

	output.WriteFooter()
	return nil
}

func jsonToText(element *gjson.Result) string {
	switch {
	case element.IsObject():
		var parts []string
		element.ForEach(func(key, value gjson.Result) bool {
			parts = append(parts, fmt.Sprintf("%s: %s", key.String(), jsonToText(&value)))
			return true
		})
		return strings.Join(parts, "; ")
	case element.IsArray():
		var parts []string
		element.ForEach(func(_, value gjson.Result) bool {
			parts = append(parts, jsonToText(&value))
			return true
		})
		return strings.Join(parts, ", ")
	default:
		return element.String()
	}
}
