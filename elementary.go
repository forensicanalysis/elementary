package elementary

import (
	"os"
	"path/filepath"
	"strconv"

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
