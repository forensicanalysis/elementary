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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"

	"github.com/forensicanalysis/elementary/daggy"
	"github.com/forensicanalysis/forensicstore"
)

const appName = "elementary"

func AppDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = ""
	}
	return filepath.Join(configDir, appName, strconv.Itoa(forensicstore.Version))
}

func All() []*cobra.Command {
	cmds := []*cobra.Command{
		eventlogs(),
		export(),
		forensicStoreImport(),
		jsonImport(),
		prefetch(),
		importFile(),
		// Yara(),
		exportTimesketch(),
		bulkSearch(),
	}
	cmds = append(cmds, dockerCommands()...)
	cmds = append(cmds, scriptCommands()...)
	return cmds
}

func RequireStore(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("the following arguments are required: forensicstore")
	}
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		return fmt.Errorf("%s: %w", args[0], os.ErrNotExist)
	}
	return nil
}

func extractFilter(filtersets []string) daggy.Filter {
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

func fileToReader(store *forensicstore.ForensicStore, exportPath gjson.Result) (*bytes.Reader, error) {
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

func jsonschemaToFlags(schema JSONSchema, command *cobra.Command) error {
	for name, property := range schema.Properties {
		switch property.Type {
		case "string":
			if defaultValue, ok := property.Default.(string); ok {
				command.Flags().String(name, defaultValue, property.Description)
			} else {
				command.Flags().String(name, "", property.Description)
			}
		case "number":
			if defaultValue, ok := property.Default.(float64); ok {
				command.Flags().Float64(name, defaultValue, property.Description)
			} else {
				command.Flags().Float64(name, 0, property.Description)
			}
		case "integer":
			if defaultValue, ok := property.Default.(int64); ok {
				command.Flags().Int64(name, defaultValue, property.Description)
			} else {
				command.Flags().Int64(name, 0, property.Description)
			}
		case "boolean":
			if defaultValue, ok := property.Default.(bool); ok {
				command.Flags().Bool(name, defaultValue, property.Description)
			} else {
				command.Flags().Bool(name, false, property.Description)
			}
		}
	}
	for _, required := range schema.Required {
		err := command.MarkFlagRequired(required)
		if err != nil {
			return err
		}
	}
	return nil
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func toCommandlineArgs(flagset *pflag.FlagSet, args []string) []string {
	var cmdArgs []string
	flagset.VisitAll(func(flag *pflag.Flag) {
		value := flag.Value.String()

		endsWithSlice := strings.HasSuffix(flag.Value.Type(), "Slice")
		if endsWithSlice && strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			slice, err := readAsCSV(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"))
			if err == nil {
				for _, value := range slice {
					cmdArgs = append(cmdArgs, fmt.Sprintf("--%s=%s", flag.Name, value))
				}
				return
			}
		}
		if flag.Value.Type() == "bool" {
			b, _ := flagset.GetBool(flag.Name)
			if b {
				cmdArgs = append(cmdArgs, fmt.Sprintf("--%s", flag.Name))
			}
			return
		}
		cmdArgs = append(cmdArgs, fmt.Sprintf("--%s=%s", flag.Name, value))
	})
	cmdArgs = append(cmdArgs, args...)
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
