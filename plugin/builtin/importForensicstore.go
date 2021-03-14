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
	"fmt"
	"io"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/forensicstore"
)

var _ pluginlib.Plugin = &ImportForensicstore{}

type ImportForensicstore struct {
	parameter pluginlib.ParameterList
}

func (i *ImportForensicstore) Name() string {
	return "import-forensicstore"
}

func (i *ImportForensicstore) Short() string {
	return "Import forensicstore files"
}

func (i *ImportForensicstore) Parameter() pluginlib.ParameterList {
	if i.parameter == nil {
		i.parameter = pluginlib.ParameterList{
			{Name: "forensicstore", Type: pluginlib.Path, Description: "forensicstore", Required: true, Argument: true},
			{Name: "file", Description: "file to import", Type: pluginlib.Path, Required: true},
			Filter,
		}
	}
	return i.parameter
}

func (i *ImportForensicstore) Output() *pluginlib.Config {
	return nil
}

func (i *ImportForensicstore) Run(p pluginlib.Plugin, _ pluginlib.LineWriter) error {
	file := p.Parameter().StringValue("file")
	filter := pluginlib.ExtractFilter(p.Parameter().GetStringArrayValue("filter"))
	store, teardown, err := getForensicStore(p)
	if err != nil {
		return err
	}
	defer teardown()
	return merge(store, file, filter)
}

// Merge merges another JSONLite into this one.
func merge(db *forensicstore.ForensicStore, url string, filter pluginlib.Filter) (err error) {
	// TODO: import elements with "_path" on sublevel"…
	// TODO: import does not need to unflatten and flatten

	importStore, teardown, err := forensicstore.Open(url)
	if err != nil {
		return err
	}
	defer teardown()

	elements, err := importStore.All()
	if err != nil {
		return err
	}

	for _, element := range elements {
		element := element
		if !filter.Match(element) {
			continue
		}

		var ferr error
		r := gjson.GetBytes(element, "@this")
		r.ForEach(func(field, value gjson.Result) bool {
			if strings.HasSuffix(field.String(), "_path") {
				dstPath, writer, teardownStoreFile, err := db.StoreFile(value.String())
				if err != nil {
					ferr = fmt.Errorf("could not store file: %w", err)
					return false
				}
				reader, teardownLoadFile, err := importStore.LoadFile(value.String())
				if err != nil {
					ferr = fmt.Errorf("could not load file: %w", err)
					return false
				}
				_, err = io.Copy(writer, reader)
				_ = teardownLoadFile()
				_ = teardownStoreFile()
				if err != nil {
					ferr = err
					return false
				}

				element, err = sjson.SetBytes(element, field.String(), dstPath)
				if err != nil {
					ferr = err
					return false
				}
			}
			return true
		})
		if ferr != nil {
			return ferr
		}

		_, err = db.Insert(element)
		if err != nil {
			return err
		}
	}
	return err
}
