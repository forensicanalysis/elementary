package builtin

import "github.com/forensicanalysis/elementary/pluginlib"

func List() []pluginlib.Plugin {
	return []pluginlib.Plugin{
		&Eventlogs{},
		&Export{},
		&ImportForensicstore{},
		&JSONImport{},
		&Prefetch{},
		&ImportFile{},
		&ExportTimesketch{},
		&BulkSearch{},
	}
}
