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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/forensicstore"

	"github.com/tidwall/gjson"
)

var _ pluginlib.Plugin = &ExportTimesketch{}

type ExportTimesketch struct {
	parameter pluginlib.ParameterList
}

func (e *ExportTimesketch) Name() string {
	return "export-timesketch"
}

func (e *ExportTimesketch) Short() string {
	return "Export in timesketch jsonl format"
}

func (e *ExportTimesketch) Parameter() pluginlib.ParameterList {
	if e.parameter == nil {
		e.parameter = pluginlib.ParameterList{
			{Name: "forensicstore", Type: pluginlib.Path, Description: "forensicstore", Required: true, Argument: true},
			{Name: "timesketch", Type: pluginlib.Path, Description: "timesketch", Required: true},
			Filter,
		}
	}
	return e.parameter
}

func (e *ExportTimesketch) Output() *pluginlib.Config {
	return &pluginlib.Config{Header: []string{"message", "datetime", "timestamp_desc"}}
}

func (e *ExportTimesketch) Run(p pluginlib.Plugin, out pluginlib.LineWriter) error {
	timesketch := p.Parameter().StringValue("timesketch")
	filter := pluginlib.ExtractFilter(p.Parameter().GetStringArrayValue("filter"))
	store, teardown, err := getForensicStore(p)
	if err != nil {
		return err
	}
	defer teardown()
	return exportStore(out, store, filter, timesketch)
}

func exportStore(out pluginlib.LineWriter, store *forensicstore.ForensicStore, filter pluginlib.Filter, timesketch string) error {
	elements, err := store.Select(filter)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return nil
	}

	f, err := os.Create(timesketch)
	if err != nil {
		return err
	}
	defer f.Close()

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
				f.Write(b)          // nolint: errcheck
				f.WriteString("\n") // nolint: errcheck
			}
			return true
		})
	}

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
