package plugin

import "github.com/spf13/cobra"

type Provider interface {
	List() []Plugin
}

func ToCobra(provider Provider) []*cobra.Command {
	var cobraCommands []*cobra.Command
	for _, cmd := range provider.List() {
		cmd := cmd
		cobraCommands = append(cobraCommands, &cobra.Command{
			Use:   cmd.Name(),
			Short: cmd.Short(),
			RunE:  func(c *cobra.Command, args []string) error { return cmd.Run(cmd) },
		})
	}

	return cobraCommands
}
