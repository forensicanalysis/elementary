package elementary

import (
	"embed"

	"github.com/forensicanalysis/elementary/plugin/builtin"
	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/elementary/pluginlib/meta"
)

type ElementaryPluginProvider struct {
	Name    string
	Dir     string
	Images  []string
	Scripts embed.FS
}

func (cp *ElementaryPluginProvider) List() []pluginlib.Plugin {
	mpp := meta.PluginProvider{
		Name:    cp.Name,
		Dir:     cp.Dir,
		Images:  cp.Images,
		Scripts: cp.Scripts,
		Plugins: builtin.List(),
	}

	return mpp.List()
}
