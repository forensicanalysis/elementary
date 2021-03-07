package builtin

import (
	"github.com/forensicanalysis/elementary/plugin"
)

var _ plugin.Provider = &CommandProvider{}

type CommandProvider struct{}

func (d *CommandProvider) List() []plugin.Plugin {
	return []plugin.Plugin{
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
