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

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/forensicanalysis/forensicstore"
	forensicstoreCmd "github.com/forensicanalysis/forensicstore/cmd"
)

func main() {
	var debugLog bool

	version := ""
	version += fmt.Sprintf("\n %-30s v%d\n", "forensicstore format:", forensicstore.Version)
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, i := range info.Deps {
			if i.Path == "github.com/forensicanalysis/forensicstore" {
				version += fmt.Sprintf(" %-30s %s\n", path.Base(i.Path)+" library:", i.Version)
			}
		}
	}

	rootCmd := cobra.Command{
		Use:                "elementary",
		Version:            version,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debugLog {
				log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
				log.Println("debugLog mode enabled")
			} else {
				log.SetOutput(ioutil.Discard)
			}
		},
	}
	archiveCommand := &cobra.Command{
		Use:   "archive",
		Short: "Insert or retrieve files from the forensicstore",
	}
	archiveCommand.AddCommand(
		forensicstoreCmd.Pack(),
		forensicstoreCmd.Unpack(),
		forensicstoreCmd.Ls(),
	)
	rootCmd.AddCommand(
		run(),
		install(),
		// workflow(),
		forensicstoreCmd.Element(),
		forensicstoreCmd.Create(),
		forensicstoreCmd.Validate(),
		archiveCommand,
	)
	rootCmd.PersistentFlags().BoolVar(&debugLog, "debug", false, "show log messages")
	_ = rootCmd.PersistentFlags().MarkHidden("debug")

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
