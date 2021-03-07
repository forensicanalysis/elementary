package builtin

import (
	"bytes"
	"io/ioutil"

	"github.com/forensicanalysis/elementary/plugin"

	"github.com/tidwall/gjson"

	"github.com/forensicanalysis/forensicstore"
)

var _ plugin.Plugin = &command{}

type command struct {
	name        string
	short       string
	parameter   plugin.ParameterList
	run         func(plugin.Plugin) error
	annotations []plugin.Annotation
}

func (cmd *command) Name() string {
	return cmd.name
}

func (cmd *command) Short() string {
	return cmd.short
}

func (cmd *command) Parameter() plugin.ParameterList {
	return cmd.parameter
}

func (cmd *command) Run(c plugin.Plugin) error {
	return cmd.run(c)
}

func (cmd *command) Annotations() []plugin.Annotation {
	return cmd.annotations
}

func fileToReader(store *forensicstore.ForensicStore, exportPath gjson.Result) (*bytes.Reader, error) {
	file, teardown, err := store.LoadFile(exportPath.String())
	if err != nil {
		return nil, err
	}
	defer teardown()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
