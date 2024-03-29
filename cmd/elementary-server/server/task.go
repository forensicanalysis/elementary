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

	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/elementary/pluginlib/output"

	"github.com/spf13/pflag"
)

type Task struct {
	Name        string
	Description string
	Schema      *pluginlib.JSONSchema
}

func ListTasks(mcp pluginlib.Provider) *Command {
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
			for _, command := range mcp.List() {
				schema := parameterToSchema(command.Parameter())
				children = append(children, Task{
					Name:        command.Name(),
					Description: command.Short(),
					Schema:      &schema,
				})
			}

			return PrintAny(w, children)
		},
	}
}

func RunTask(cp pluginlib.Provider) *Command {
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

			var p pluginlib.Plugin
			for _, command := range cp.List() {
				if command.Name() == name {
					p = command
				}
			}

			if p == nil {
				return fmt.Errorf("p %s cannot be run", name)
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

			// p.SetOut(w) TODO
			// p.Flags().AddFlagSet(flags) TODO
			for name, arg := range arguments {
				fmt.Println(name, arg)
				p.Parameter().Set(name, fmt.Sprint(arg))
			}

			return p.Run(p, output.NewJSONOutput(w)) // TODO also to store?
		},
	}
}

func parameterToSchema(parameters pluginlib.ParameterList) pluginlib.JSONSchema {
	schema := pluginlib.JSONSchema{
		Properties: map[string]pluginlib.Property{},
		Required:   []string{},
	}

	for _, parameter := range parameters {
		typeMapping := map[pluginlib.ParameterType]string{
			pluginlib.String: "string",
			pluginlib.Bool:   "boolean",
		}

		schema.Properties[parameter.Name] = pluginlib.Property{
			Type:        typeMapping[parameter.Type],
			Description: parameter.Description,
			Default:     parameter.Value,
		}
		if parameter.Required {
			schema.Required = append(schema.Required, parameter.Name)
		}
	}
	return schema
}
