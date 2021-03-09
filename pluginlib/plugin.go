package pluginlib

import (
	"github.com/spf13/cobra"
)

type Config struct {
	Header []string `json:"header,omitempty"`
}

type Plugin interface {
	Name() string
	Short() string
	Parameter() ParameterList
	Output() *Config
	Run(Plugin, LineWriter) error
}

func ToCobra(plugins []Plugin) []*cobra.Command {
	var cobraCommands []*cobra.Command
	for _, cmd := range plugins {
		cmd := cmd
		cobraCommand := &cobra.Command{
			Use:   cmd.Name(),
			Short: cmd.Short(),
			RunE: func(c *cobra.Command, args []string) error {
				return cmd.Run(cmd, nil) // TODO
			},
		}
		// TODO add parameter => flags
		cobraCommands = append(cobraCommands, cobraCommand)
	}

	return cobraCommands
}
