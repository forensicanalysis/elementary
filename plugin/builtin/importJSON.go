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
	"errors"
	"os"

	"github.com/tidwall/gjson"

	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/forensicstore"
)

var _ pluginlib.Plugin = &JSONImport{}

type JSONImport struct {
	parameter pluginlib.ParameterList
}

func (j *JSONImport) Name() string {
	return "import-json"
}

func (j *JSONImport) Short() string {
	return "Import json files"
}

func (j *JSONImport) Parameter() pluginlib.ParameterList {
	if j.parameter == nil {
		j.parameter = pluginlib.ParameterList{
			{Name: "forensicstore", Type: pluginlib.Path, Description: "forensicstore", Required: true, Argument: true},
			{Name: "file", Description: "file to import", Type: pluginlib.Path, Required: true},
			Filter,
		}
	}
	return j.parameter
}

func (j *JSONImport) Output() *pluginlib.Config {
	return nil
}

func (j *JSONImport) Run(p pluginlib.Plugin, _ pluginlib.LineWriter) error {
	file := p.Parameter().StringValue("file")
	filter := pluginlib.ExtractFilter(p.Parameter().GetStringArrayValue("filter"))
	store, teardown, err := getForensicStore(p)
	if err != nil {
		return err
	}
	defer teardown()

	b, err := os.ReadFile(file) // #nosec
	if err != nil {
		return err
	}

	topLevel := gjson.GetBytes(b, "@this")
	if !topLevel.IsArray() {
		return errors.New("imported json must have a top level array containing objects")
	}

	topLevel.ForEach(func(_, element gjson.Result) bool {
		elementType := element.Get("type")
		if elementType.Exists() && filter.Match(forensicstore.JSONElement(element.Raw)) {
			_, err = store.Insert(forensicstore.JSONElement(element.Raw))
			if err != nil {
				return false
			}
		}
		return true
	})

	return nil
}
