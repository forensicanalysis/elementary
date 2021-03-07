package builtin

import (
	"github.com/forensicanalysis/elementary/plugin"
)

var _ plugin.Provider = &PluginProvider{}

type PluginProvider struct{}

func (d *PluginProvider) List() []plugin.Plugin {
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
