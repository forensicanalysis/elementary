package meta

import (
	"embed"
	"path/filepath"

	"github.com/forensicanalysis/elementary/plugin"

	"github.com/forensicanalysis/elementary/plugin/builtin"
	"github.com/forensicanalysis/elementary/plugin/docker"
	"github.com/forensicanalysis/elementary/plugin/script"
)

type PluginProvider struct {
	Name    string
	Dir     string
	Images  []string
	Scripts embed.FS
}

func (cp *PluginProvider) List() []plugin.Plugin {
	builtinPluginProvider := builtin.PluginProvider{}
	scriptPluginProvider := script.PluginProvider{Prefix: cp.Name, Dir: filepath.Join(cp.Dir, "scripts"), Scripts: cp.Scripts}
	dockerPluginProvider := docker.PluginProvider{Prefix: cp.Name, Images: cp.Images}

	l := builtinPluginProvider.List()
	l = append(l, scriptPluginProvider.List()...)
	l = append(l, dockerPluginProvider.List()...)
	return l
}
