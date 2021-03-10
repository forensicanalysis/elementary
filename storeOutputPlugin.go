package elementary

import (
	"github.com/forensicanalysis/elementary/pluginlib"
	"github.com/forensicanalysis/forensicstore"
)

var (
	AddToStoreParameter = &pluginlib.Parameter{
		Name:        "add-to-store",
		Description: "output to store",
		Type:        pluginlib.Bool,
		Value:       false,
		Required:    false,
	}
	ForensicStoreParameter = &pluginlib.Parameter{
		Name:     "forensicstore",
		Type:     pluginlib.Path,
		Required: true,
		Argument: true,
	}
)

type StoreOutputPlugin struct {
	Internal pluginlib.Plugin
}

func (s *StoreOutputPlugin) Name() string {
	return s.Internal.Name()
}

func (s *StoreOutputPlugin) Short() string {
	return s.Internal.Short()
}

func (s *StoreOutputPlugin) Parameter() pluginlib.ParameterList {
	pl := append(s.Internal.Parameter(), AddToStoreParameter)
	_, err := s.Internal.Parameter().Get("forensicstore")
	if err != nil {
		return append(pl, ForensicStoreParameter)
	}
	return pl
}

func (s *StoreOutputPlugin) Output() *pluginlib.Config {
	return s.Internal.Output()
}

func (s *StoreOutputPlugin) Run(p pluginlib.Plugin, writer pluginlib.LineWriter) error {
	path := p.Parameter().StringValue("forensicstore")
	store, teardown, err := forensicstore.New(path)
	if err != nil {
		return err
	}
	defer teardown()

	if p.Parameter().BoolValue("add-to-store") {
		lw := NewForensicStoreOutput(store)
		writer = &pluginlib.MultiLineWriter{LineWriter: []pluginlib.LineWriter{writer, lw}}
		defer lw.WriteFooter()
	}
	return s.Internal.Run(p, writer)
}
