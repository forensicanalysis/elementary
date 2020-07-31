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
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/spf13/pflag"

	"github.com/forensicanalysis/forensicstore"
)

func LoadFile() *Command {
	return &Command{
		Name:   "loadFile",
		Route:  "/file",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("path", "", "path")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			p, err := flags.GetString("path")
			if err != nil {
				return err
			}
			if p == "" {
				return errors.New("path must be set")
			}

			storeName := flags.Args()[0]
			store, teardown, err := forensicstore.Open(storeName)
			if err != nil {
				return err
			}
			defer teardown()

			f, teardownFile, err := store.LoadFile(p)
			if err != nil {
				return err
			}
			defer teardownFile()
			_, err = io.Copy(w, f)
			return err
		},
	}
}

func Files() *Command {
	return &Command{
		Name:   "listFiles",
		Route:  "/files",
		Method: http.MethodGet,
		SetupFlags: func(f *pflag.FlagSet) {
			f.String("path", "", "path")
		},
		Handler: func(w io.Writer, _ io.Reader, flags *pflag.FlagSet) error {
			p, err := flags.GetString("path")
			if err != nil {
				return err
			}
			if p == "" {
				p = "/"
			}

			storeName := flags.Args()[0]
			store, teardown, err := forensicstore.Open(storeName)
			if err != nil {
				return err
			}
			defer teardown()

			root, err := store.Fs.Open(p)
			if err != nil {
				return err
			}
			infos, err := root.Readdir(0)
			if err != nil {
				return err
			}

			var children []forensicstore.JSONElement
			for _, info := range infos {
				b, _ := json.Marshal(struct {
					Name    string    `json:"name"`
					Path    string    `json:"path"`
					Size    int64     `json:"size"`
					Dir     bool      `json:"dir"`
					ModTime time.Time `json:"mtime"`
				}{
					info.Name(),
					path.Join(p, info.Name()),
					info.Size(),
					info.IsDir(),
					info.ModTime(),
				})
				children = append(children, b)
			}

			return PrintJSONList(w, int64(len(children)), children)
		},
	}
}
