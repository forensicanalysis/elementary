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
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/forensicanalysis/elementary/pluginlib"
)

var _ pluginlib.Plugin = &command{}

type command struct {
	ScriptName      string               `json:"name,omitempty"`
	ScriptShort     string               `json:"short,omitempty"`
	ScriptArguments pluginlib.JSONSchema `json:"arguments,omitempty"`
	ScriptOutput    *pluginlib.Config    `json:"output,omitempty"`

	parameter pluginlib.ParameterList
	run       func(pluginlib.Plugin, io.Writer) error
}

func newCommand(path string) pluginlib.Plugin {
	scriptCommand := &command{}

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
	scriptCommand.run = func(cmd pluginlib.Plugin, out io.Writer) error {
		shellCommand := strings.Join(append(
			[]string{`"` + filepath.ToSlash(path) + `"`},
			scriptCommand.Parameter().ToCommandlineArgs()...,
		), " ")

		if strings.HasSuffix(path, ".py") {
			name, err := exec.LookPath("python3")
			if err == nil {
				shellCommand = name + " " + shellCommand
			}
		}

		log.Println("sh", "-c", shellCommand)
		script := exec.Command("sh", "-c", shellCommand) // #nosec

		script.Stdout = out
		script.Stderr = log.Writer()
		err := script.Run()
		if err != nil {
			return fmt.Errorf("%s script failed with %w", scriptCommand.ScriptName, err)
		}

		return nil
	}
	scriptCommand.parameter = append(scriptCommand.parameter, pluginlib.JsonschemaToParameter(scriptCommand.ScriptArguments)...)

	return scriptCommand
}

func (s *command) Name() string {
	return s.ScriptName
}

func (s *command) Short() string {
	return s.ScriptShort
}

func (s *command) Parameter() pluginlib.ParameterList {
	return s.parameter
}

func (s *command) Output() *pluginlib.Config {
	return s.ScriptOutput
}

func (s *command) Run(c pluginlib.Plugin, writer pluginlib.LineWriter) error {
	lbw := &pluginlib.LineWriterBuffer{Writer: writer}
	defer lbw.WriteFooter()
	return s.run(c, lbw)
}
