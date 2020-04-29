package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime/debug"

	"github.com/spf13/cobra"

	forensicstoreCmd "github.com/forensicanalysis/forensicstore/cmd"
	workflowCmd "github.com/forensicanalysis/forensicworkflows/cmd"
)

func main() {
	var debugLog bool

	rootCmd := cobra.Command{
		Use:                "elementary",
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
	rootCmd.AddCommand(
		workflowCmd.Run(),
		workflowCmd.Install(),
		workflowCmd.Workflow(),
		forensicstoreCmd.Item(),
		forensicstoreCmd.Create(),
		forensicstoreCmd.Validate(),
	)
	rootCmd.PersistentFlags().BoolVar(&debugLog, "debug", false, "show log messages")

	info, ok := debug.ReadBuildInfo()
	if ok {
		fmt.Println(rootCmd.Name())
		for _, i := range info.Deps {
			if i.Path == "github.com/forensicanalysis/forensicstore" || i.Path == "github.com/forensicanalysis/forensicworkflows" {
				fmt.Printf(" %-20s %s\n", path.Base(i.Path)+":", i.Version)
			}
		}
	}

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
