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
	"embed"
	"github.com/forensicanalysis/elementary/commands/meta"
	"github.com/forensicanalysis/forensicstore"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/forensicanalysis/elementary/cmd/elementary-server/server"
)

//go:embed dist
var static embed.FS

const appName = "elementary"

func main() {
	sub, err := fs.Sub(static, "dist")
	if err != nil {
		log.Fatal(err)
	}

	mcp := &meta.CommandProvider{Name: appName, Dir: appDir()}
	rootCmd := server.Application("fstore", http.FS(sub), server.Commands(mcp)...)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func appDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = ""
	}
	return filepath.Join(configDir, appName, strconv.Itoa(forensicstore.Version))
}
