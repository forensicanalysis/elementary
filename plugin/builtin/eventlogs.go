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

	"github.com/Velocidex/ordereddict"
	"github.com/tidwall/gjson"
	"www.velocidex.com/golang/evtx"

	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/forensicstore"
)

var _ pluginlib.Plugin = &Eventlogs{}

type Eventlogs struct {
	parameter pluginlib.ParameterList
}

func (e *Eventlogs) Name() string {
	return "eventlogs"
}

func (e *Eventlogs) Short() string {
	return "Process eventlogs into single events"
}

func (e *Eventlogs) Parameter() pluginlib.ParameterList {
	if e.parameter == nil {
		e.parameter = pluginlib.ParameterList{
			{Name: "forensicstore", Type: pluginlib.Path, Description: "forensicstore", Required: true, Argument: true},
			Filter,
		}
	}
	return e.parameter
}

func (e *Eventlogs) Output() *pluginlib.Config {
	return &pluginlib.Config{Header: []string{
		"System.Computer",
		"System.TimeCreated.SystemTime",
		"System.EventRecordID",
		"System.EventID.Value",
		"System.Level",
		"System.Channel",
		"System.Provider.Name",
	}}
}

func (e *Eventlogs) Run(p pluginlib.Plugin, out pluginlib.LineWriter) error {
	store, teardown, err := getForensicStore(p)
	if err != nil {
		return err
	}
	defer teardown()

	filter := pluginlib.ExtractFilter(p.Parameter().GetStringArrayValue("filter"))
	return eventlogsFromStore(out, store, filter)
}

func eventlogsFromStore(out pluginlib.LineWriter, store *forensicstore.ForensicStore, filter pluginlib.Filter) error {
	for idx := range filter {
		filter[idx]["type"] = "file"
		filter[idx]["name"] = "%.evtx"
	}

	if len(filter) == 0 {
		filter = pluginlib.Filter{{"type": "file", "name": "%.evtx"}}
	}

	fileElements, err := store.Select(filter)
	if err != nil {
		return err
	}

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
				out.WriteLine(event) // nolint: errcheck
			}
		}
	}
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
