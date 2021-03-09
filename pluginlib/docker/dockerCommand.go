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

	"github.com/forensicanalysis/elementary/pluginlib"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type command struct {
	name      string
	short     string
	parameter pluginlib.ParameterList
	run       func(pluginlib.Plugin, io.Writer) error
	output    *pluginlib.Config
}

func newCommand(name, image string, labels map[string]string) pluginlib.Plugin {
	dockerCmd := &command{
		name:  name,
		short: "(docker: " + image + ")",
		run: func(cmd pluginlib.Plugin, writer io.Writer) error {
			mounts := parseMounts(cmd)
			args := cmd.Parameter().ToCommandlineArgs()
			return dockerCreate(image, args, mounts, writer)
		},
	}

	if short, ok := labels["short"]; ok {
		dockerCmd.short = short + " (docker: " + image + ")"
	}

	if headers, ok := labels["headers"]; ok {
		dockerCmd.output = &pluginlib.Config{Header: strings.Split(headers, ",")}
	}

	dockerCmd.parameter = append(dockerCmd.parameter, getLabelParameter(labels)...)

	return dockerCmd
}

func (s *command) Name() string {
	return s.name
}

func (s *command) Short() string {
	return s.short
}

func (s *command) Parameter() pluginlib.ParameterList {
	return s.parameter
}

func (s *command) Output() *pluginlib.Config {
	return s.output
}

func (s *command) Run(c pluginlib.Plugin, writer pluginlib.LineWriter) error {
	lbw := &pluginlib.LineWriterBuffer{Writer: writer}
	defer lbw.WriteFooter()
	return s.run(c, lbw)
}

func parseMounts(cmd pluginlib.Plugin) map[string]string {
	mounts := map[string]string{}

	// TODO: check if application id == "eldr"

	for _, parameter := range cmd.Parameter() {
		if parameter.Type == pluginlib.Path {
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
		if parameter.Type == pluginlib.PathArray {
			// TODO
		}

	}
	return mounts
}

func getLabelParameter(labels map[string]string) []*pluginlib.Parameter {
	var parameters []*pluginlib.Parameter
	if use, ok := labels["arguments"]; ok {
		var schema pluginlib.JSONSchema
		err := json.Unmarshal([]byte(use), &schema)
		if err != nil {
			log.Println(err)
		} else {
			parameters = append(parameters, pluginlib.JsonschemaToParameter(schema)...)
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
