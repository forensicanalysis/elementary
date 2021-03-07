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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"

	"github.com/forensicanalysis/elementary/daggy"
	"github.com/forensicanalysis/elementary/plugin/meta"
	"github.com/forensicanalysis/forensicstore"
)

// workflow is a subcommand to run a forensic workflow.
func workflow() *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "run a workflow",
		Long: `process can run parallel workflows locally. Those workflows are a directed acyclic graph of tasks.
Those tasks can be defined to be run on the system itself or in a containerized way.`,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("the following arguments are required: forensicstore")
			}
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				return fmt.Errorf("%s: %w", args[0], os.ErrNotExist)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse workflow yaml
			workflowFile, _ := cmd.Flags().GetString("file")
			if _, err := os.Stat(workflowFile); os.IsNotExist(err) {
				log.Fatal(err, workflowFile)
			}

			mcp := meta.PluginProvider{}
			engine := daggy.New(mcp.List())
			workflow, err := daggy.Parse(workflowFile)
			if err != nil {
				return err
			}

			if err := insertTasks(args[0], workflow); err != nil {
				return err
			}

			return engine.Run(workflow, args[0])
		},
	}
	workflowCmd.Flags().StringP("file", "f", "", "workflow definition file")
	_ = workflowCmd.MarkFlagRequired("file")
	return workflowCmd
}

func insertTasks(storeURL string, workflow *daggy.Workflow) error {
	store, teardown, err := forensicstore.Open(storeURL)
	if err != nil {
		return err
	}
	defer teardown()
	for _, task := range workflow.Tasks {
		jsonTask, err := json.Marshal(task)
		if err != nil {
			return err
		}
		jsonTask, err = sjson.SetBytes(jsonTask, "type", "task")
		if err != nil {
			return err
		}
		_, err = store.Insert(jsonTask)
		if err != nil {
			return err
		}
	}
	return nil
}
