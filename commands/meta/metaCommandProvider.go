package meta

import (
	"path/filepath"

	"github.com/forensicanalysis/elementary/commands/builtin"
	"github.com/forensicanalysis/elementary/commands/docker"
	"github.com/forensicanalysis/elementary/commands/script"
	"github.com/forensicanalysis/elementary/daggy"
)

type CommandProvider struct {
	Name string
	Dir  string
}

func (cp *CommandProvider) List() []daggy.Command {
	scp := script.CommandProvider{Prefix: cp.Name, Dir: filepath.Join(cp.Dir, "scripts")}
	dcp := docker.CommandProvider{Prefix: cp.Name}
	bicp := builtin.BuiltInCommandProvider{}

	l := append(scp.List(), dcp.List()...)
	return append(l, bicp.List()...)
}
