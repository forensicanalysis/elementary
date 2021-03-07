package builtin

import (
	"github.com/forensicanalysis/elementary/daggy"
)

var _ daggy.CommandProvider = &BuiltInCommandProvider{}

type BuiltInCommandProvider struct{}

func (d *BuiltInCommandProvider) List() []daggy.Command {
	return []daggy.Command{
		eventlogs(),
		export(),
		forensicStoreImport(),
		jsonImport(),
		prefetch(),
		importFile(),
		exportTimesketch(),
		bulkSearch(),
	}
}

var _ daggy.Command = &BuiltInCommand{}

type BuiltInCommand struct {
	name        string
	short       string
	parameter   daggy.ParameterList
	run         func(daggy.Command) error
	annotations []daggy.Annotation
}

func (cmd *BuiltInCommand) Name() string {
	return cmd.name
}

func (cmd *BuiltInCommand) Short() string {
	return cmd.short
}

func (cmd *BuiltInCommand) Parameter() daggy.ParameterList {
	return cmd.parameter
}

func (cmd *BuiltInCommand) Run(c daggy.Command) error {
	return cmd.run(c)
}

func (cmd *BuiltInCommand) Annotations() []daggy.Annotation {
	return cmd.annotations
}
