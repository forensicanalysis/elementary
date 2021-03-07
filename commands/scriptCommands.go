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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func scriptCommands() []*cobra.Command {
	scriptDir := filepath.Join(AppDir(), "scripts")

	infos, err := ioutil.ReadDir(scriptDir)
	if err != nil {
		log.Printf("script plugins disabled: %s, ", err)
		return nil
	}

	var commands []*cobra.Command
	for _, info := range infos {
		validName := strings.HasPrefix(info.Name(), appName+"-") && !strings.HasSuffix(info.Name(), ".info")
		if info.Mode().IsRegular() && validName {
			commands = append(commands, scriptCommand(filepath.Join(scriptDir, info.Name())))
		}
	}
	return commands
}

type commandTemplate struct {
	*cobra.Command
	Arguments JSONSchema `json:"arguments,omitempty"`
}

func scriptCommand(path string) *cobra.Command {
	cmd := commandTemplate{}

	out, err := ioutil.ReadFile(path + ".info") // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(path + ".info does not exist")
		} else {
			log.Println(path, err)
		}
	} else {
		err = json.Unmarshal(out, &cmd)
		if err != nil {
			log.Println(err)
		}
	}

	if cmd.Use == "" {
		cmd.Use = filepath.Base(path)
	}
	cmd.Short += " (script)"
	cmd.Args = RequireStore
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		log.Printf("run %s %s", cmd.Name(), args[0])
		shellCommand := strings.Join(append(
			[]string{`"` + filepath.ToSlash(path) + `"`},
			toCommandlineArgs(cmd.Flags(), []string{filepath.ToSlash(args[0])})...,
		), " ")

		if strings.HasSuffix(path, ".py") {
			name, err := exec.LookPath("python3")
			if err == nil {
				shellCommand = name + " " + shellCommand
			}
		}

		log.Println("sh", "-c", shellCommand)
		script := exec.Command("sh", "-c", shellCommand) // #nosec

		output, teardown := newOutputWriterURL(cmd, args[0])
		defer teardown()

		script.Stdout = output
		script.Stderr = log.Writer()
		err := script.Run()
		if err != nil {
			return fmt.Errorf("%s script failed with %w", cmd.Use, err)
		}

		output.WriteFooter()
		return nil
	}
	err = jsonschemaToFlags(cmd.Arguments, cmd.Command)
	if err != nil {
		log.Println(err)
	}
	addOutputFlags(cmd.Command)
	return cmd.Command
}
