package main

import (
	"github.com/forensicanalysis/elementary"
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
	internal pluginlib.Plugin
}

func (s *StoreOutputPlugin) Name() string {
	return s.internal.Name()
}

func (s *StoreOutputPlugin) Short() string {
	return s.internal.Short()
}

func (s *StoreOutputPlugin) Parameter() pluginlib.ParameterList {
	return append(s.internal.Parameter(), AddToStoreParameter, ForensicStoreParameter)
}

func (s *StoreOutputPlugin) Output() *pluginlib.Config {
	return s.internal.Output()
}

func (s *StoreOutputPlugin) Run(p pluginlib.Plugin, writer pluginlib.LineWriter) error {
	path := p.Parameter().StringValue("forensicstore")
	store, teardown, err := forensicstore.New(path)
	if err != nil {
		return err
	}
	defer teardown()

	if p.Parameter().BoolValue("add-to-store") {
		lw := elementary.NewForensicStoreOutput(store)
		writer = &pluginlib.MultiLineWriter{LineWriter: []pluginlib.LineWriter{writer, lw}}
		defer lw.WriteFooter()
	}
	return s.internal.Run(p, writer)
}
