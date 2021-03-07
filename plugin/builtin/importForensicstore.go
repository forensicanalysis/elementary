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

	"github.com/forensicanalysis/elementary/plugin"
	"github.com/forensicanalysis/forensicstore"
)

func forensicStoreImport() plugin.Plugin {
	return &command{
		name:      "import-forensicstore",
		short:     "Import forensicstore files",
		parameter: []*plugin.Parameter{ForensicStore, {Name: "file", Type: plugin.Path, Required: true}, Filter},
		run: func(cmd plugin.Plugin) error {
			path := cmd.Parameter().StringValue("forensicstore")
			file := cmd.Parameter().StringValue("file")
			filter := plugin.ExtractFilter(cmd.Parameter().GetStringArrayValue("filter"))
			return singleImport(path, file, filter)
		},
		annotations: []plugin.Annotation{plugin.Di, plugin.Importer},
	}
}

func singleImport(url string, file string, filter plugin.Filter) error {
	store, teardown, err := forensicstore.Open(url)
	if err != nil {
		return err
	}
	defer teardown()

	err = merge(store, file, filter)
	if err != nil {
		return err
	}
	return nil
}

// Merge merges another JSONLite into this one.
func merge(db *forensicstore.ForensicStore, url string, filter plugin.Filter) (err error) {
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
