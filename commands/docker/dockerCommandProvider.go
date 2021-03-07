package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/forensicanalysis/elementary/daggy"
)

var _ daggy.CommandProvider = &CommandProvider{}

type CommandProvider struct {
	Prefix string
}

func (d *CommandProvider) List() []daggy.Command {
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

	var cmds []daggy.Command
	commandNames := map[string]bool{}
	for _, imageSummary := range imageSummaries {
		for _, dockerImage := range imageSummary.RepoTags {
			name, err := commandName(d.Prefix, dockerImage)
			if err != nil {
				continue
			}

			cmd := NewDockerCommand(name, dockerImage, imageSummary.Labels)
			cmds = append(cmds, cmd)
			commandNames[name] = true
		}
	}
	for _, dockerImage := range Images() {
		name, err := commandName(d.Prefix, dockerImage)
		if err != nil {
			continue
		}
		if _, ok := commandNames[name]; !ok {
			labels := map[string]string{"short": fmt.Sprintf("Use '%s install -f' to download", os.Args[0])}
			cmds = append(cmds, NewDockerCommand(name, dockerImage, labels))
		}
	}

	return cmds
}
