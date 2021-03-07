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
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/forensicanalysis/elementary/daggy"
	"github.com/forensicanalysis/forensicstore"
)

func ExtractFilter(filtersets []string) daggy.Filter {
	filter := daggy.Filter{}
	for _, filterset := range filtersets {
		filterelement := map[string]string{}
		for _, kv := range strings.Split(filterset, ",") {
			kvl := strings.SplitN(kv, "=", 2)
			if len(kvl) == 2 { //nolint: gomnd
				filterelement[kvl[0]] = kvl[1]
			}
		}

		filter = append(filter, filterelement)
	}
	return filter
}

func FileToReader(store *forensicstore.ForensicStore, exportPath gjson.Result) (*bytes.Reader, error) {
	file, teardown, err := store.LoadFile(exportPath.String())
	if err != nil {
		return nil, err
	}
	defer teardown()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func ToCommandlineArgs(list daggy.ParameterList) []string {
	var cmdArgs []string
	for _, p := range list {

		if p.Argument {
			cmdArgs = append(cmdArgs, p.Name)
		}

		value := fmt.Sprint(p.Value)
		if p.Type.IsList() && strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			slice, err := readAsCSV(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"))
			if err == nil {
				for _, value := range slice {
					cmdArgs = append(cmdArgs, fmt.Sprintf("--%s=%s", p.Name, value))
				}
				continue
			}
		}
		if p.Type == daggy.Bool {
			if p.BoolValue() {
				cmdArgs = append(cmdArgs, fmt.Sprintf("--%s", p.Name))
			}
			continue
		}
		cmdArgs = append(cmdArgs, fmt.Sprintf("--%s=%s", p.Name, value))
	}
	return cmdArgs
}

type Property struct {
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

type JSONSchema struct {
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

func JsonschemaToParameter(schema JSONSchema) []*daggy.Parameter {
	var parameters []*daggy.Parameter
	for name, property := range schema.Properties {
		p := &daggy.Parameter{Name: name, Description: property.Description}
		switch property.Type {
		case "string":
			p.Type = daggy.String
			if defaultValue, ok := property.Default.(string); ok {
				p.Value = defaultValue
			} else {
				p.Value = ""
			}
		case "boolean":
			p.Type = daggy.Bool
			if defaultValue, ok := property.Default.(bool); ok {
				p.Value = defaultValue
			} else {
				p.Value = ""
			}
		default:
			panic(fmt.Sprintf("unknown jsonschema type %s", property.Type))
		}
		if contains(schema.Required, name) {
			p.Required = true
		}
		parameters = append(parameters, p)
	}
	return nil
}

func contains(list []string, elem string) bool {
	for _, i := range list {
		if i == elem {
			return true
		}
	}
	return false
}
