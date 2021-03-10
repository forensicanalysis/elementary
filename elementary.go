package elementary

import (
	"embed"
	"os"
	"path/filepath"
	"strconv"

	"github.com/forensicanalysis/elementary/pluginlib/output"

	"github.com/forensicanalysis/elementary/plugin/builtin"
	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/elementary/pluginlib/meta"
	"github.com/forensicanalysis/forensicstore"
)

func Name() string {
	return "elementary"
}

func AppDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = ""
	}
	return filepath.Join(configDir, Name(), strconv.Itoa(forensicstore.Version))
}

func Images() []string {
	return []string{
		"docker.io/forensicanalysis/elementary-shimcache:v0.3.6",
		"docker.io/forensicanalysis/elementary-plaso:v0.3.6",
		"docker.io/forensicanalysis/elementary-import-image:v0.3.6",
		"docker.io/forensicanalysis/elementary-yara:v0.3.6",
		"docker.io/forensicanalysis/elementary-sigma:v0.3.6",
	}
}

//go:embed plugin/scripts
var scripts embed.FS

func NewPluginProvider() pluginlib.Provider {
	return &PluginProvider{Name: Name(), Dir: AppDir(), Images: Images(), Scripts: scripts}
}

type PluginProvider struct {
	Name    string
	Dir     string
	Images  []string
	Scripts embed.FS
}

func (cp *PluginProvider) List() []pluginlib.Plugin {
	mpp := meta.PluginProvider{
		Name:    cp.Name,
		Dir:     cp.Dir,
		Images:  cp.Images,
		Scripts: cp.Scripts,
		Plugins: builtin.List(),
	}

	return storeOutputLayer(mpp.List())
}

func storeOutputLayer(plugins []pluginlib.Plugin) []pluginlib.Plugin {
	var layerd []pluginlib.Plugin
	for _, p := range plugins {
		layerd = append(layerd,
			&pluginlib.LoggerOutputPlugin{
				Internal: &StoreOutputPlugin{
					Internal: &output.FormatOutputPlugin{
						Internal: p,
					},
				},
			},
		)
	}
	return layerd
}
