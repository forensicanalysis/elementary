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
	"log"

	"github.com/tidwall/gjson"
	goprefetch "www.velocidex.com/golang/go-prefetch"

	"github.com/forensicanalysis/elementary/plugin"
	"github.com/forensicanalysis/forensicstore"
)

func prefetch() plugin.Plugin {
	prefetchCommand := &command{
		name:  "prefetch",
		short: "Process prefetch files",
		parameter: []*plugin.Parameter{
			{Name: "forensicstore", Type: plugin.Path, Required: true, Argument: true},
			{Name: "filter", Description: "filter processed events", Type: plugin.StringArray, Required: false},
		},
		run: func(cmd plugin.Plugin) error {
			log.Printf("run prefetch")
			path := cmd.Parameter().StringValue("forensicstore")
			filtersets := cmd.Parameter().GetStringArrayValue("filter")
			return prefetchFromStore(path, plugin.ExtractFilter(filtersets), cmd)
		},
	}
	prefetchCommand.parameter = append(prefetchCommand.parameter, plugin.OutputParameter(prefetchCommand)...)

	return prefetchCommand
}

func prefetchFromStore(url string, filter plugin.Filter, cmd plugin.Plugin) error {
	store, teardown, err := forensicstore.Open(url)
	if err != nil {
		return err
	}
	defer teardown()

	for idx := range filter {
		filter[idx]["type"] = "file"
		filter[idx]["name"] = "%.pf"
	}

	if len(filter) == 0 {
		filter = plugin.Filter{{"type": "file", "name": "%.pf"}}
	}

	fileElements, err := store.Select(filter)
	if err != nil {
		return err
	}

	output := plugin.NewOutputWriterStore(cmd, store, &plugin.OutputConfig{
		Header: []string{
			"Executable",
			"FileSize",
			"Hash",
			"Version",
			"LastRunTimes",
			"FilesAccessed",
			"RunCount",
		},
	})

	for _, element := range fileElements {
		exportPath := gjson.GetBytes(element, "export_path")
		if exportPath.Exists() && exportPath.String() != "" {
			buff, err := fileToReader(store, exportPath)
			if err != nil {
				return err
			}

			prefetchInfo, err := goprefetch.LoadPrefetch(buff)
			if err != nil {
				return err
			}

			elem, err := prefetchToElement(prefetchInfo)
			if err != nil {
				return err
			}

			output.WriteLine(elem) // nolint: errcheck
		}
	}

	output.WriteFooter()
	return nil
}

func prefetchToElement(prefetchInfo *goprefetch.PrefetchInfo) (forensicstore.JSONElement, error) {
	return json.Marshal(map[string]interface{}{
		"Executable":    prefetchInfo.Executable,
		"FileSize":      prefetchInfo.FileSize,
		"Hash":          prefetchInfo.Hash,
		"Version":       prefetchInfo.Version,
		"LastRunTimes":  prefetchInfo.LastRunTimes,
		"FilesAccessed": prefetchInfo.FilesAccessed,
		"RunCount":      prefetchInfo.RunCount,
		"type":          "prefetch",
	})
}
