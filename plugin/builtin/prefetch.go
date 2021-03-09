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

	"github.com/tidwall/gjson"
	goprefetch "www.velocidex.com/golang/go-prefetch"

	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/forensicstore"
)

var _ pluginlib.Plugin = &Prefetch{}

type Prefetch struct {
	parameter pluginlib.ParameterList
}

func (p *Prefetch) Name() string {
	return "prefetch"

}
func (p *Prefetch) Short() string {
	return "Process prefetch files"
}

func (p *Prefetch) Parameter() pluginlib.ParameterList {
	if p.parameter == nil {
		p.parameter = pluginlib.ParameterList{
			{Name: "forensicstore", Type: pluginlib.Path, Description: "forensicstore", Required: true, Argument: true},
			Filter,
		}
	}
	return p.parameter
}

func (p *Prefetch) Output() *pluginlib.Config {
	return &pluginlib.Config{Header: []string{"Executable", "FileSize", "Hash", "Version", "LastRunTimes", "FilesAccessed", "RunCount"}}
}

func (p *Prefetch) Run(plg pluginlib.Plugin, out pluginlib.LineWriter) error {
	filter := pluginlib.ExtractFilter(plg.Parameter().GetStringArrayValue("filter"))
	store, teardown, err := getForensicStore(plg)
	if err != nil {
		return err
	}
	defer teardown()
	return prefetchFromStore(out, store, filter)
}

func prefetchFromStore(out pluginlib.LineWriter, store *forensicstore.ForensicStore, filter pluginlib.Filter) error {
	for idx := range filter {
		filter[idx]["type"] = "file"
		filter[idx]["name"] = "%.pf"
	}

	if len(filter) == 0 {
		filter = pluginlib.Filter{{"type": "file", "name": "%.pf"}}
	}

	fileElements, err := store.Select(filter)
	if err != nil {
		return err
	}

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

			out.WriteLine(elem) // nolint: errcheck
		}
	}
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
