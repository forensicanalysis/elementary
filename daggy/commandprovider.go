package daggy

import "github.com/spf13/cobra"

type CommandProvider interface {
	List() []Command
}

func ToCobra(commandProvider CommandProvider) []*cobra.Command {
	var cobraCommands []*cobra.Command
	for _, cmd := range commandProvider.List() {
		cobraCommand := &cobra.Command{
			Use:   cmd.Name(),
			Short: cmd.Short(),
		}
		cobraCommand.RunE = func(c *cobra.Command, args []string) error {
			return cmd.Run(cmd)
		}
		cobraCommands = append(cobraCommands, cobraCommand)
	}

	return cobraCommands
}
