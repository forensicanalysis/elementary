package output

import "github.com/forensicanalysis/elementary/plugin"

var File = &plugin.Parameter{
	Name:        "output",
	Description: "choose an output file",
	Type:        plugin.Path,
	Value:       "",
	Required:    false,
}

var Format = &plugin.Parameter{Name: "format",
	Description: "choose output format [csv, jsonl, table, json, none]",
	Type:        plugin.String,
	Value:       "table",
	Required:    false,
}
