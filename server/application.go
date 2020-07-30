// Copyright (c) 2020 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

package server

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Command struct {
	Name       string
	Route      string
	Method     string
	SetupFlags func(f *pflag.FlagSet)
	Handler    func(io.Writer, io.Reader, *pflag.FlagSet) error
}

func Application(name string, width, height int, static http.FileSystem, uiRoot bool, cmds ...*Command) *cobra.Command {
	rootCmd := &cobra.Command{Use: name}

	rootCmd.AddCommand(Commandline(cmds...)...)
	rootCmd.AddCommand(ServeCommand(static, cmds...))

	// parse config
	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+name+".yaml)")
	initConfig(name, cfgFile)

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig(name, cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".view" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("." + name)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setupFlags(command *Command, flagset *pflag.FlagSet, formatDefault string) {
	flagset.String("format", formatDefault, "Set output format")
	if command.SetupFlags != nil {
		command.SetupFlags(flagset)
	}
}
