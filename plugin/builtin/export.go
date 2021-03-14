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

import "github.com/forensicanalysis/elementary/pluginlib"

var _ pluginlib.Plugin = &Export{}

type Export struct{}

func (e *Export) Name() string {
	return "export"
}

func (e *Export) Short() string {
	return "Export selected elements"
}

func (e *Export) Parameter() pluginlib.ParameterList {
	return []*pluginlib.Parameter{Filter}
}

func (e *Export) Output() *pluginlib.Config {
	return nil
}

func (e *Export) Run(p pluginlib.Plugin, out pluginlib.LineWriter) error {
	filter := pluginlib.ExtractFilter(p.Parameter().GetStringArrayValue("filter"))

	store, teardown, err := getForensicStore(p)
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

	/*
		var header []string
		gjson.GetBytes(elements[0], "@this").ForEach(func(key, _ gjson.Result) bool {
			header = append(header, key.String())
			return true
		})
		out.SetConfig(&output.Config{Header: header})
	*/

	for _, element := range elements {
		out.WriteLine(element) // nolint: errcheck
	}
	return nil
}
