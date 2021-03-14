package meta

import (
	"embed"
	"path/filepath"

	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/elementary/pluginlib/docker"
	"github.com/forensicanalysis/elementary/pluginlib/script"
)

type PluginProvider struct {
	Name    string
	Dir     string
	Images  []string
	Scripts embed.FS
	Plugins []pluginlib.Plugin
}

func (cp *PluginProvider) List() []pluginlib.Plugin {
	scriptPluginProvider := script.PluginProvider{
		Prefix: cp.Name, Dir: filepath.Join(cp.Dir, "scripts"), Scripts: cp.Scripts,
	}
	dockerPluginProvider := docker.PluginProvider{Prefix: cp.Name, Images: cp.Images}

	l := scriptPluginProvider.List()
	l = append(l, dockerPluginProvider.List()...)
	l = append(l, cp.Plugins...)
	return l
}
