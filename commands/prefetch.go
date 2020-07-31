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
	"encoding/json"
	"log"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	goprefetch "www.velocidex.com/golang/go-prefetch"

	"github.com/forensicanalysis/elementary/daggy"
	"github.com/forensicanalysis/forensicstore"
)

func prefetch() *cobra.Command {
	var filtersets []string
	prefetchCommand := &cobra.Command{
		Use:   "prefetch <forensicstore>",
		Short: "Process prefetch files",
		Args:  RequireStore,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("run prefetch %s", args)
			return prefetchFromStore(args[0], extractFilter(filtersets), cmd)
		},
	}
	addOutputFlags(prefetchCommand)
	prefetchCommand.Flags().StringArrayVar(&filtersets, "filter", nil, "filter processed events")
	return prefetchCommand
}

func prefetchFromStore(url string, filter daggy.Filter, cmd *cobra.Command) error {
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
		filter = daggy.Filter{{"type": "file", "name": "%.pf"}}
	}

	fileElements, err := store.Select(filter)
	if err != nil {
		return err
	}

	output := newOutputWriterStore(cmd, store, &outputConfig{
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

			output.writeLine(elem) // nolint: errcheck
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
