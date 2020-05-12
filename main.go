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
	workflowCmd "github.com/forensicanalysis/forensicworkflows/cmd"
)

func main() {
	var debugLog bool

	version := ""
	version += fmt.Sprintf("\n %-30s v%d\n", "forensicstore format:", forensicstore.Version)
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, i := range info.Deps {
			if i.Path == "github.com/forensicanalysis/forensicstore" || i.Path == "github.com/forensicanalysis/forensicworkflows" {
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
		workflowCmd.Run(),
		workflowCmd.Install(),
		workflowCmd.Workflow(),
		forensicstoreCmd.Element(),
		forensicstoreCmd.Create(),
		forensicstoreCmd.Validate(),
		archiveCommand,
	)
	rootCmd.PersistentFlags().BoolVar(&debugLog, "debug", false, "show log messages")
	rootCmd.PersistentFlags().MarkHidden("debug")

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
