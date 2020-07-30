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

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/forensicanalysis/elementary/commands"
)

type Task struct {
	Name        string
	Description string
	Schema      *commands.JSONSchema
}

func ListTasks() *Command {
	return &Command{
		Name:   "listTasks",
		Route:  "/tasks",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			// f.String("directory", "/", "current directory")
			// f.String("type", "file", "item type")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			var children []Task
			for _, command := range commands.All() {
				if command.Annotations != nil {
					if properties, ok := command.Annotations["plugin_property_flags"]; ok {
						if strings.Contains(properties, "ex") {
							continue
						}
					}
				}

				schema := flagsToSchema(command.Flags())
				children = append(children, Task{
					Name:        command.Name(),
					Description: command.Short,
					Schema:      &schema,
				})
			}

			return PrintAny(w, children)
		},
	}
}

func RunTask() *Command {
	return &Command{
		Name:   "run",
		Route:  "/run",
		Method: http.MethodPost,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("name", "", "command name")
		},
		Handler: func(w io.Writer, r io.Reader, flags *pflag.FlagSet) error {
			name, err := flags.GetString("name")
			if err != nil {
				return err
			}

			var plugin *cobra.Command
			for _, command := range commands.All() {
				if command.Name() == name {
					plugin = command
				}
			}

			if plugin == nil || plugin.RunE == nil {
				return fmt.Errorf("plugin %s cannot be run", name)
			}

			var arguments map[string]interface{}
			b, err := ioutil.ReadAll(r)
			if err != nil {
				return err
			}
			err = json.Unmarshal(b, &arguments)
			if err != nil {
				return err
			}

			plugin.SetOut(w)
			plugin.Flags().AddFlagSet(flags)
			for name, arg := range arguments {
				fmt.Println(name, arg)
				err = plugin.Flags().Set(name, fmt.Sprint(arg))
				if err != nil {
					return err
				}
			}
			err = plugin.Flags().Set("format", "json")
			if err != nil {
				return err
			}
			err = plugin.Flags().Set("add-to-store", "true")
			if err != nil {
				return err
			}
			return plugin.RunE(plugin, flags.Args())
		},
	}
}

func flagsToSchema(flags *pflag.FlagSet) commands.JSONSchema {
	schema := commands.JSONSchema{
		Properties: map[string]commands.Property{},
		Required:   []string{},
	}

	flags.VisitAll(func(f *pflag.Flag) {
		typeMapping := map[string]string{
			"string": "string",
			"int":    "integer",
			"bool":   "boolean",
			"float":  "number",
		}

		property := commands.Property{
			Type:        typeMapping[f.Value.Type()],
			Description: f.Usage,
		}

		if f.DefValue != "" {
			var defaultValue interface{}
			var err error
			switch f.Value.Type() {
			case "string":
				property.Default = f.DefValue
			case "int":
				defaultValue, err = strconv.ParseInt(f.DefValue, 10, 64)
			case "bool":
				defaultValue, err = strconv.ParseBool(f.DefValue)
			case "float":
				defaultValue, err = strconv.ParseFloat(f.DefValue, 64)
			}
			if err == nil {
				property.Default = defaultValue
			}
		}
		schema.Properties[f.Name] = property
		if _, ok := f.Annotations[cobra.BashCompOneRequiredFlag]; ok {
			schema.Required = append(schema.Required, f.Name)
		}
	})
	return schema
}
