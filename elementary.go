package elementary

import (
	"embed"
	"os"
	"path/filepath"
	"strconv"

	"github.com/forensicanalysis/elementary/plugin"
	"github.com/forensicanalysis/elementary/plugin/meta"
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

//go:embed scripts
var scripts embed.FS

func PluginProvider() plugin.Provider {
	return &meta.PluginProvider{Name: Name(), Dir: AppDir(), Images: Images(), Scripts: scripts}
}
