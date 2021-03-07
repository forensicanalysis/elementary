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

package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/forensicanalysis/elementary/plugin"
	output "github.com/forensicanalysis/elementary/plugin/output"
)

var _ plugin.Plugin = &command{}

type command struct {
	name        string
	short       string
	parameter   plugin.ParameterList
	run         func(plugin.Plugin) error
	annotations []plugin.Annotation
}

func newCommand(name, image string, labels map[string]string) plugin.Plugin {
	dockerCmd := &command{
		name:      name,
		short:     "(docker: " + image + ")",
		parameter: []*plugin.Parameter{output.File, output.Format},
		run: func(cmd plugin.Plugin) error {
			log.Println("run", cmd.Name())

			mounts := parseMounts(cmd)

			path := cmd.Parameter().StringValue("output")
			format := cmd.Parameter().StringValue("format")
			out := output.New(path, format, nil)
			defer out.WriteFooter()

			args := cmd.Parameter().ToCommandlineArgs()
			return dockerCreate(image, args, mounts, out)
		},
	}

	if short, ok := labels["short"]; ok {
		dockerCmd.short = short + " (docker: " + image + ")"
	}

	addOutput := true
	if properties, ok := labels["properties"]; ok {
		if strings.Contains(properties, string(plugin.Di)) { // TODO: use constant
			addOutput = false
		}
		// dockerCmd.annotations = append(dockerCmd.annotations, properties) TODO
	}

	dockerCmd.parameter = append(dockerCmd.parameter, getLabelParameter(labels)...)
	if addOutput {
		dockerCmd.parameter = append(dockerCmd.parameter, output.File, output.Format)
	}

	return dockerCmd
}

func (s *command) Name() string {
	return s.name
}

func (s *command) Short() string {
	return s.short
}

func (s *command) Parameter() plugin.ParameterList {
	return s.parameter
}

func (s *command) Run(c plugin.Plugin) error {
	return s.run(c)
}

func (s *command) Annotations() []plugin.Annotation {
	return s.annotations
}

func parseMounts(cmd plugin.Plugin) map[string]string {
	mounts := map[string]string{}

	// TODO: check if application id == "eldr"

	for _, parameter := range cmd.Parameter() {
		if parameter.Type == plugin.Path {
			mountPointValue := parameter.StringValue()
			if mountPointValue == "" {
				continue
			}
			abs, err := filepath.Abs(mountPointValue)
			if err != nil {
				continue
			}
			mounts[abs] = parameter.Name
		}
		if parameter.Type == plugin.PathArray {
			// TODO
		}

	}
	return mounts
}

func getLabelParameter(labels map[string]string) []*plugin.Parameter {
	var parameters []*plugin.Parameter
	if use, ok := labels["arguments"]; ok {
		var schema plugin.JSONSchema
		err := json.Unmarshal([]byte(use), &schema)
		if err != nil {
			log.Println(err)
		} else {
			parameters = append(parameters, plugin.JsonschemaToParameter(schema)...)
		}
	}
	return parameters
}

func commandName(prefix, image string) (string, error) {
	idx := strings.LastIndex(image, "/")
	if strings.HasPrefix(image[idx+1:], prefix+"-") {
		name := image[idx+len(prefix)+2:]
		parts := strings.Split(name, ":")
		name = parts[0]
		return name, nil
	}
	return "", errors.New("no plugin")
}

func dockerCreate(image string, args []string, mountDirs map[string]string, w io.Writer) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	mounts, err := getMounts(mountDirs)
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{Image: image, Cmd: args, Tty: true, WorkingDir: "/elementary"},
		&container.HostConfig{Binds: mounts},
		nil,
		"",
	)
	if err != nil {
		return err
	}

	log.Println("start docker container")
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	defer cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}) // nolint: errcheck

	go streamLogs(ctx, cli, resp.ID, w, true, false)            // nolint: errcheck
	go streamLogs(ctx, cli, resp.ID, log.Writer(), false, true) // nolint: errcheck

	log.Println("wait for docker container")
	statusCode, err := cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		return err
	}
	if statusCode != 0 {
		return fmt.Errorf("container returned status code %d", statusCode)
	}
	return nil
}

func streamLogs(ctx context.Context, cli *client.Client, id string, w io.Writer, stdout, stderr bool) error {
	options := types.ContainerLogsOptions{ShowStderr: stderr, ShowStdout: stdout, Follow: true}
	out, err := cli.ContainerLogs(ctx, id, options)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(w, out)
	return err
}

func getMounts(mountDirs map[string]string) ([]string, error) {
	for localDir := range mountDirs {
		// create directory if not exists
		_, err := os.Open(localDir) // #nosec
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s does not exist", localDir)
		} else if err != nil {
			return nil, err
		}
	}
	for localDir := range mountDirs {
		if localDir[1] == ':' {
			mountDirs["/"+strings.ToLower(string(localDir[0]))+filepath.ToSlash(localDir[2:])] = mountDirs[localDir]
			delete(mountDirs, localDir)
		}
	}

	var mounts []string
	for localDir, containerDir := range mountDirs {
		// mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: localDir, Target: "/" + containerDir})
		mounts = append(mounts, localDir+":/elementary/"+containerDir)
	}
	return mounts, nil
}
