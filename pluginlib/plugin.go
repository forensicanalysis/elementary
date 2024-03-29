package pluginlib

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Config struct {
	Header []string `json:"header,omitempty"`
}

type Plugin interface {
	Name() string
	Short() string
	Parameter() ParameterList
	Output() *Config
	Run(Plugin, LineWriter) error
}

type SimpleLineWriter struct{}

func (s SimpleLineWriter) WriteLine(bytes []byte) {
	fmt.Println(string(bytes))
}

func ToCobra(plugins []Plugin) []*cobra.Command {
	var cobraCommands []*cobra.Command
	for _, plugin := range plugins {
		plgn := plugin
		cobraCommand := &cobra.Command{
			Use:   plgn.Name(),
			Short: plgn.Short(),
			RunE: func(c *cobra.Command, args []string) error {
				err := setParameterValues(plgn.Parameter(), c.Flags(), args)
				if err != nil {
					return err
				}
				return plgn.Run(plgn, &SimpleLineWriter{})
			},
		}
		for _, parameter := range plgn.Parameter() {
			if parameter.Argument {
				cobraCommand.Use += " <" + parameter.Name + ">"
				continue
			}
			switch parameter.Type {
			case String, Path:
				cobraCommand.Flags().String(parameter.Name, parameter.StringValue(), parameter.Description)
			case StringArray, PathArray:
				cobraCommand.Flags().StringArray(parameter.Name, parameter.StringArray(), parameter.Description)
			case Bool:
				cobraCommand.Flags().Bool(parameter.Name, parameter.BoolValue(), parameter.Description)
			default:
				log.Printf("unknown parameter type %v", parameter.Type)
			}
			if parameter.Required {
				_ = cobraCommand.MarkFlagRequired(parameter.Name)
			}
		}
		cobraCommands = append(cobraCommands, cobraCommand)
	}

	return cobraCommands
}

func setParameterValues(parameterList ParameterList, flags *pflag.FlagSet, args []string) error {
	flags.VisitAll(func(flag *pflag.Flag) {
		var value interface{}
		var err error
		switch flag.Value.Type() {
		case "stringArray":
			value, err = flags.GetStringArray(flag.Name)
		case "string":
			value, err = flags.GetString(flag.Name)
		case "bool":
			value, err = flags.GetBool(flag.Name)
		}
		if err != nil {
			log.Println(err)
		}
		_, err = parameterList.Get(flag.Name)
		if err != nil {
			return
		}
		parameterList.Set(flag.Name, value)
	})

	i := 0
	for _, parameter := range parameterList {
		if parameter.Argument {
			if i+1 > len(args) {
				return fmt.Errorf("missing argument %s", parameter.Name)
			}
			parameterList.Set(parameter.Name, args[i])
			i++
		}
	}
	return nil
}
