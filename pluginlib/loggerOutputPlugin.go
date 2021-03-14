package pluginlib

import (
	"log"
)

type LoggerOutputPlugin struct {
	Internal Plugin
}

func (s *LoggerOutputPlugin) Name() string {
	return s.Internal.Name()
}

func (s *LoggerOutputPlugin) Short() string {
	return s.Internal.Short()
}

func (s *LoggerOutputPlugin) Parameter() ParameterList {
	return s.Internal.Parameter()
}

func (s *LoggerOutputPlugin) Output() *Config {
	return s.Internal.Output()
}

func (s *LoggerOutputPlugin) Run(p Plugin, w LineWriter) error {
	log.Printf("run %s\n", p.Name())
	return s.Internal.Run(p, w)
}
