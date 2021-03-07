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
	"io"
	"log"

	"github.com/Velocidex/ordereddict"
	"github.com/tidwall/gjson"
	"www.velocidex.com/golang/evtx"

	"github.com/forensicanalysis/elementary/plugin"
	"github.com/forensicanalysis/forensicstore"
)

func eventlogs() plugin.Plugin {
	eventlogsCmd := &command{
		name:  "eventlogs",
		short: "Process eventlogs into single events",
		parameter: []*plugin.Parameter{
			{Name: "forensicstore", Type: plugin.Path, Required: true, Argument: true},
			{Name: "filter", Description: "filter processed events", Type: plugin.StringArray, Required: false},
		},
		run: func(cmd plugin.Plugin) error {
			log.Printf("run eventlogs")
			filtersets := cmd.Parameter().GetStringArrayValue("filter")
			return eventlogsFromStore(
				cmd.Parameter().StringValue("forensicstore"),
				plugin.ExtractFilter(filtersets),
				cmd,
			)
		},
	}
	eventlogsCmd.parameter = append(eventlogsCmd.parameter, plugin.OutputParameter(eventlogsCmd)...)
	return eventlogsCmd
}

func eventlogsFromStore(url string, filter plugin.Filter, cmd plugin.Plugin) error {
	store, teardown, err := forensicstore.Open(url)
	if err != nil {
		return err
	}
	defer teardown()

	for idx := range filter {
		filter[idx]["type"] = "file"
		filter[idx]["name"] = "%.evtx"
	}

	if len(filter) == 0 {
		filter = plugin.Filter{{"type": "file", "name": "%.evtx"}}
	}

	fileElements, err := store.Select(filter)
	if err != nil {
		return err
	}

	output := plugin.NewOutputWriterStore(cmd, store, &plugin.OutputConfig{
		Header: []string{
			"System.Computer",
			"System.TimeCreated.SystemTime",
			"System.EventRecordID",
			"System.EventID.Value",
			"System.Level",
			"System.Channel",
			"System.Provider.Name",
		},
	})
	for _, element := range fileElements {
		exportPath := gjson.GetBytes(element, "export_path")
		if exportPath.Exists() && exportPath.String() != "" {
			r, err := fileToReader(store, exportPath)
			if err != nil {
				return err
			}

			events, err := getEvents(exportPath.String(), r)
			if err != nil {
				return err
			}

			for _, event := range events {
				output.WriteLine(event) // nolint: errcheck
			}
		}
	}
	output.WriteFooter()
	return nil
}

func getEvents(originPath string, file io.ReadSeeker) ([]forensicstore.JSONElement, error) {
	var elements []forensicstore.JSONElement

	chunks, err := evtx.GetChunks(file)
	if err != nil && err.Error() == "Unsupported EVTX version." {
		evtxVersionError, _ := json.Marshal(map[string]interface{}{
			"origin": map[string]string{"path": originPath},
			"type":   "eventlog",
			"errors": []string{err.Error()},
		})
		return []forensicstore.JSONElement{evtxVersionError}, nil
	} else if err != nil {
		return nil, err
	}

	for _, chunk := range chunks {
		records, err := chunk.Parse(int(chunk.Header.FirstEventRecID))
		if err != nil {
			return nil, err
		}

		for _, i := range records {
			eventMap, ok := i.Event.(*ordereddict.Dict)
			if ok {
				event, ok := ordereddict.GetMap(eventMap, "Event")
				if !ok {
					continue
				}

				event.Set("type", "eventlog")
				event.Set("origin", map[string]string{"path": originPath})
				// self.maybeExpandMessage(event)

				serialized, err := json.Marshal(event)
				if err != nil {
					return nil, err
				}

				elements = append(elements, serialized)
			}
		}
	}

	return elements, nil
}
