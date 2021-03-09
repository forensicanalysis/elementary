package main

import (
	"log"

	"github.com/forensicanalysis/elementary/pluginlib"
)

type LoggerOutputPlugin struct {
	internal pluginlib.Plugin
}

func (s *LoggerOutputPlugin) Name() string {
	return s.internal.Name()
}

func (s *LoggerOutputPlugin) Short() string {
	return s.internal.Short()
}

func (s *LoggerOutputPlugin) Parameter() pluginlib.ParameterList {
	return s.internal.Parameter()
}

func (s *LoggerOutputPlugin) Output() *pluginlib.Config {
	return s.internal.Output()
}

func (s *LoggerOutputPlugin) Run(p pluginlib.Plugin, w pluginlib.LineWriter) error {
	log.Printf("run %s\n", p.Name())
	return s.internal.Run(p, w)
}
