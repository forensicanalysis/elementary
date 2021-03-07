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

package builtin

/*
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hillu/go-yara"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/forensicanalysis/forensicstore"
)

func Yara() daggy.Command {
	var rulesPath string
	prefetchCommand := &BuiltInCommand{
		Use:   "yara <forensicstore>",
		Short: "Process prefetch files",
		Args:  RequireStore,
		run: func(cmd daggy.Command, args []string) error {
			log.Printf("run yara %s", args)

			c, err := yara.NewCompiler()
			if err != nil {
				return err
			}

			f, err := os.Open(rulesPath) // #nosec
			if err != nil {
				return err
			}
			err = c.AddFile(f, "default")
			if err != nil {
				return err
			}
			err = f.Close()
			if err != nil {
				return err
			}

			rules, err := c.GetRules()
			if err != nil {
				return err
			}

			for _, url := range args {
				err = yaraStore(url, rules, cmd)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	addOutputFlags(prefetchCommand)
	prefetchCommand.Flags().StringVar(&rulesPath, "rules", "", "yara rule directory")
	return prefetchCommand
}

func yaraStore(url string, rules *yara.Rules, cmd daggy.Command) error {
	store, teardown, err := forensicstore.Open(url)
	if err != nil {
		return err
	}
	defer teardown()

	var elements []forensicstore.JSONElement

	err = afero.Walk(store.Fs, "/", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		f, err := store.LoadFile(filepath.ToSlash(path))
		if err != nil {
			return err
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		matches, err := rules.ScanMem(b, 0, time.Second*10) // nolint: gomnd
		if err != nil {
			return err
		}
		for _, match := range matches {
			b, err := json.Marshal(struct {
				Type string                 `json:"type"`
				Rule string                 `json:"rule"`
				Meta map[string]interface{} `json:"meta"`
			}{
				Type: "yara",
				Rule: match.Rule,
				Meta: match.Meta,
			})
			if err != nil {
				return err
			}
			elements = append(elements, b)
		}
		return nil
	})
	if err != nil {
		return err
	}

	config := &outputConfig{
		Header: []string{
			"Rule",
		},
	}
	printElements(cmd, config, elements, store)
	return nil
}
*/
