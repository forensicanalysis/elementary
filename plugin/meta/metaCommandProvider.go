package meta

import (
	"path/filepath"

	"github.com/forensicanalysis/elementary/plugin"

	"github.com/forensicanalysis/elementary/plugin/builtin"
	"github.com/forensicanalysis/elementary/plugin/docker"
	"github.com/forensicanalysis/elementary/plugin/script"
)

type CommandProvider struct {
	Name string
	Dir  string
}

func (cp *CommandProvider) List() []plugin.Plugin {
	scp := script.CommandProvider{Prefix: cp.Name, Dir: filepath.Join(cp.Dir, "scripts")}
	dcp := docker.CommandProvider{Prefix: cp.Name}
	bicp := builtin.CommandProvider{}

	l := append(scp.List(), dcp.List()...)
	return append(l, bicp.List()...)
}
