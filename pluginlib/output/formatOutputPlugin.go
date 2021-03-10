package output

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/forensicanalysis/elementary/pluginlib"
)

var (
	OutputParameter = &pluginlib.Parameter{
		Name:        "output",
		Description: "choose an output file",
		Type:        pluginlib.Path,
		Value:       "",
		Required:    false,
	}
	FormatParameter = &pluginlib.Parameter{
		Name:        "format",
		Description: "choose output format [csv, jsonl, table, json, none]",
		Type:        pluginlib.String,
		Value:       "jsonl",
		Required:    false,
	}
)

type FormatOutputPlugin struct {
	Internal pluginlib.Plugin
}

func (s *FormatOutputPlugin) Name() string {
	return s.Internal.Name()
}

func (s *FormatOutputPlugin) Short() string {
	return s.Internal.Short()
}

func (s *FormatOutputPlugin) Parameter() pluginlib.ParameterList {
	return append(s.Internal.Parameter(), OutputParameter, FormatParameter)
}

func (s *FormatOutputPlugin) Output() *pluginlib.Config {
	return s.Internal.Output()
}

func (s *FormatOutputPlugin) Run(p pluginlib.Plugin, _ pluginlib.LineWriter) error {
	path := p.Parameter().StringValue("output")
	format := p.Parameter().StringValue("format")

	var dest io.Writer
	if path != "" {
		f, err := os.Create(path)
		if err != nil {
			log.Println(err)
		}
		dest = f
	} else {
		dest = os.Stdout
	}

	var w pluginlib.LineWriter
	switch format {
	case "table":
		if p.Output() == nil {
			return fmt.Errorf("%s does not support table output", p.Name())
		}
		o := NewTableOutput(dest, p.Output().Header)
		w = o
		defer o.WriteFooter()
	case "csv":
		if p.Output() == nil {
			return fmt.Errorf("%s does not support csv output", p.Name())
		}
		o := NewCSVOutput(dest, p.Output().Header)
		w = o
		defer o.WriteFooter()
	case "jsonl":
		w = NewJsonlOutput(dest)
	case "json":
		o := NewJSONOutput(dest)
		w = o
		defer o.WriteFooter()
	case "none":
	default:
		return fmt.Errorf("unknown output format %s", format)
	}

	return s.Internal.Run(p, w)
}
