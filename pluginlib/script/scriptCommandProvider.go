package script

import (
	"embed"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/forensicanalysis/elementary/pluginlib"
)

var _ pluginlib.Provider = &PluginProvider{}

type PluginProvider struct {
	Scripts embed.FS
	Prefix  string
	Dir     string
}

func (s *PluginProvider) List() []pluginlib.Plugin {
	infos, err := ioutil.ReadDir(s.Dir)
	if err != nil {
		log.Printf("script plugins disabled: %s, ", err)
		return nil
	}

	var cmds []pluginlib.Plugin
	for _, info := range infos {
		validName := strings.HasPrefix(info.Name(), s.Prefix+"-") && !strings.HasSuffix(info.Name(), ".json")
		if info.Mode().IsRegular() && validName {
			cmds = append(cmds, newCommand(filepath.Join(s.Dir, info.Name())))
		}
	}
	return cmds
}
