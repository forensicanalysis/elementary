package script

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/forensicanalysis/elementary/plugin"
)

var _ plugin.Provider = &CommandProvider{}

type CommandProvider struct {
	Prefix string
	Dir    string
}

func (s *CommandProvider) List() []plugin.Plugin {
	infos, err := ioutil.ReadDir(s.Dir)
	if err != nil {
		log.Printf("script plugins disabled: %s, ", err)
		return nil
	}

	var cmds []plugin.Plugin
	for _, info := range infos {
		validName := strings.HasPrefix(info.Name(), s.Prefix+"-") && !strings.HasSuffix(info.Name(), ".info")
		if info.Mode().IsRegular() && validName {
			cmds = append(cmds, newCommand(filepath.Join(s.Dir, info.Name())))
		}
	}
	return cmds
}
