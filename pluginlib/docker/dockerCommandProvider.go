package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forensicanalysis/elementary/pluginlib"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type PluginProvider struct {
	Prefix string
	Images []string
}

func (d *PluginProvider) List() []pluginlib.Plugin {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second) // TODO: adjust time
	defer cancel()

	cli, err := client.NewEnvClient()
	if err != nil {
		return nil
	}

	options := types.ImageListOptions{All: true}
	imageSummaries, err := cli.ImageList(timeoutCtx, options)
	if err != nil {
		log.Printf("docker plugins disabled: %s", err)
		return nil
	}

	var cmds []pluginlib.Plugin
	commandNames := map[string]bool{}
	for _, imageSummary := range imageSummaries {
		for _, dockerImage := range imageSummary.RepoTags {
			name, err := commandName(d.Prefix, dockerImage)
			if err != nil {
				continue
			}

			cmd := newCommand(name, dockerImage, imageSummary.Labels)
			cmds = append(cmds, cmd)
			commandNames[name] = true
		}
	}
	for _, dockerImage := range d.Images {
		name, err := commandName(d.Prefix, dockerImage)
		if err != nil {
			continue
		}
		if _, ok := commandNames[name]; !ok {
			labels := map[string]string{
				"short": fmt.Sprintf("Use '%s install -f' to download", filepath.Base(os.Args[0])),
			}
			cmds = append(cmds, newCommand(name, dockerImage, labels))
		}
	}

	return cmds
}
