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

package script

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/forensicanalysis/elementary/commands"
	"github.com/forensicanalysis/elementary/daggy"
)



var _ daggy.Command = &Command{}

type Command struct {
	ScriptName        string              `json:"name,omitempty"`
	ScriptShort       string              `json:"short,omitempty"`
	Arguments         commands.JSONSchema `json:"arguments,omitempty"`
	ScriptAnnotations []daggy.Annotation  `json:"annotations,omitempty"`
	parameter         daggy.ParameterList
	run               func(daggy.Command) error
}

func (s *Command) Name() string {
	return s.ScriptName
}

func (s *Command) Short() string {
	return s.ScriptShort
}

func (s *Command) Parameter() daggy.ParameterList {
	return s.parameter
}

func (s *Command) Run(c daggy.Command) error {
	return s.run(c)
}

func (s *Command) Annotations() []daggy.Annotation {
	return s.ScriptAnnotations
}

func NewScriptCommand(path string) daggy.Command {
	scriptCommand := &Command{}

	out, err := ioutil.ReadFile(path + ".info") // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(path + ".info does not exist")
		} else {
			log.Println(path, err)
		}
	} else {
		err = json.Unmarshal(out, &scriptCommand)
		if err != nil {
			log.Println(err)
		}
	}

	if scriptCommand.ScriptName == "" {
		scriptCommand.ScriptName = filepath.Base(path)
	}
	scriptCommand.ScriptShort += " (script)"
	scriptCommand.run = func(cmd daggy.Command) error {
		log.Printf("run %s", scriptCommand.Name())
		shellCommand := strings.Join(append(
			[]string{`"` + filepath.ToSlash(path) + `"`},
			commands.ToCommandlineArgs(scriptCommand.Parameter())...,
		), " ")

		if strings.HasSuffix(path, ".py") {
			name, err := exec.LookPath("python3")
			if err == nil {
				shellCommand = name + " " + shellCommand
			}
		}

		log.Println("sh", "-c", shellCommand)
		script := exec.Command("sh", "-c", shellCommand) // #nosec

		path := cmd.Parameter().StringValue("forensicstore")
		output, teardown := commands.NewOutputWriterURL(scriptCommand, path)
		defer teardown()

		script.Stdout = output
		script.Stderr = log.Writer()
		err := script.Run()
		if err != nil {
			return fmt.Errorf("%s script failed with %w", scriptCommand.ScriptName, err)
		}

		output.WriteFooter()
		return nil
	}
	scriptCommand.parameter = append(scriptCommand.parameter, commands.JsonschemaToParameter(scriptCommand.Arguments)...)
	scriptCommand.parameter = append(scriptCommand.parameter, commands.OutputParameter(scriptCommand)...)

	return scriptCommand
}