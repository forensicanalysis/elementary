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
	"context"
	"github.com/forensicanalysis/elementary"
	"github.com/forensicanalysis/elementary/commands/docker"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// install required assets.
func install() *cobra.Command {
	var force bool
	var dockerUser, dockerPassword, dockerServer string
	cmd := &cobra.Command{
		Use:          "install",
		Short:        "Setup required assets",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var auth types.AuthConfig
			auth.Username = dockerUser
			auth.Password = dockerPassword
			auth.ServerAddress = dockerServer

			if force {
				setup(&auth, true)
			}
			return nil // fmt.Errorf("%s already exists, use --force to recreate", appDir)
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "workflow definition file")
	cmd.Flags().StringVar(&dockerUser, "docker-user", "", "docker registry username")
	cmd.Flags().StringVar(&dockerPassword, "docker-password", "", "docker registry password")
	cmd.Flags().StringVar(&dockerServer, "docker-server", "", "docker registry server")
	return cmd
}

func ensureSetup() {
	_, err := os.UserConfigDir()
	if err != nil {
		log.Printf("config dir not found: %s, using current directory", err)
	}
	appDir := elementary.AppDir()
	info, err := os.Stat(appDir)
	if os.IsNotExist(err) {
		setup(nil, false)
		return
	}
	if err != nil {
		log.Println(err)
	}
	if !info.IsDir() {
		log.Printf("%s is not a directory", appDir)
	}
}

func setup(auth *types.AuthConfig, pull bool) {
	appDir := elementary.AppDir()

	// unpack scripts
	err := unpack(appDir)
	if err != nil {
		log.Println("error unpacking scripts:", err)
	}

	// install python requirements
	pipPath, err := exec.LookPath("pip3")
	if err != nil {
		pipPath, err = exec.LookPath("pip")
		if err != nil {
			log.Println("pip is not installed")
			pipPath = ""
		}
	}
	if pipPath != "" {
		log.Println(pipPath, "install",
			"--target", filepath.Join(appDir, "scripts"),
			"-r", filepath.Join(appDir, "requirements.txt"))
		pip := exec.Command(pipPath, "install",
			"--target", filepath.Join(appDir, "scripts"),
			"-r", filepath.Join(appDir, "requirements.txt")) // #nosec
		err := pip.Run()
		if err != nil {
			log.Println("error installing python requirements:", err)
		}
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println("error setting up docker client:", err)
		return
	}

	if pull {
		pullImages(ctx, cli, auth)
	}
}

func pullImages(ctx context.Context, cli *client.Client, auth *types.AuthConfig) { // nolint: unused
	// remove old images
	options := types.ImageListOptions{All: true}
	imageSummaries, err := cli.ImageList(ctx, options)
	if err != nil {
		log.Printf("could not list images: %s", err)
	}
	for _, imageSummary := range imageSummaries {
		for _, dockerImage := range imageSummary.RepoTags {
			isElementary := strings.HasPrefix(dockerImage, "forensicanalysis/elementary-")
			if isElementary && !contains(docker.Images(), dockerImage) {
				_, _ = cli.ImageRemove(ctx, dockerImage, types.ImageRemoveOptions{Force: true})
			}
		}
	}

	// pull docker images
	for _, image := range docker.Images() {
		log.Println("pull docker image", image)
		err = pullImage(ctx, cli, image, auth)
		if err != nil {
			log.Println("error pulling docker images:", err)
		}
	}
}

func contains(l []string, s string) bool {
	for _, e := range l {
		if e == s {
			return true
		}
	}
	return false
}

func unpack(appDir string) (err error) {
	return RestoreAssets(appDir, "")
}

func pullImage(ctx context.Context, cli *client.Client, image string, auth *types.AuthConfig) error { // nolint: unused
	if auth != nil {
		body, err := cli.RegistryLogin(ctx, *auth)
		if err != nil {
			return err
		}
		log.Println("login", body)
	}

	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = io.Copy(log.Writer(), reader)
	if err != nil {
		return err
	}
	return nil
}

