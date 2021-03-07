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
	"github.com/tidwall/gjson"

	"github.com/forensicanalysis/elementary/plugin"
	"github.com/forensicanalysis/forensicstore"
)

func export() plugin.Plugin {
	outputCommand := &command{
		name:  "export",
		short: "Export selected elements",
		parameter: []*plugin.Parameter{
			{Name: "forensicstore", Type: plugin.Path, Required: true, Argument: true},
			{Name: "filter", Description: "filter processed events", Type: plugin.StringArray, Required: false},
		},
		run: func(cmd plugin.Plugin) error {
			filter := plugin.ExtractFilter(cmd.Parameter().GetStringArrayValue("filter"))

			path := cmd.Parameter().StringValue("forensicstore")
			store, teardown, err := forensicstore.Open(path)
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

			var header []string
			gjson.GetBytes(elements[0], "@this").ForEach(func(key, _ gjson.Result) bool {
				header = append(header, key.String())
				return true
			})
			output := plugin.NewOutputWriterStore(cmd, store, &plugin.OutputConfig{
				Header: header,
			})
			for _, element := range elements {
				output.Write(element) // nolint: errcheck
			}
			output.WriteFooter()
			return nil
		},
		annotations: []plugin.Annotation{plugin.Exporter},
	}
	outputCommand.parameter = append(outputCommand.parameter, plugin.OutputParameter(outputCommand)...)
	return outputCommand
}
