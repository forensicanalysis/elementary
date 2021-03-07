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

import (
	"crypto/md5"  // #nosec
	"crypto/sha1" // #nosec
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/forensicanalysis/elementary/plugin"
	"github.com/forensicanalysis/forensicstore"
)

func importFile() plugin.Plugin {
	return &command{
		name:      "import-file",
		short:     "Import files",
		parameter: []*plugin.Parameter{ForensicStore, {Name: "file", Description: "file to import", Type: plugin.PathArray, Required: true}},
		run: func(cmd plugin.Plugin) error {
			path := cmd.Parameter().StringValue("forensicstore")
			files := cmd.Parameter().GetStringArrayValue("file")
			return singleFileImport(path, files)
		},
		annotations: []plugin.Annotation{plugin.Di, plugin.Importer},
	}
}

func singleFileImport(url string, files []string) error {
	store, teardown, err := forensicstore.Open(url)
	if err != nil {
		return err
	}
	defer teardown()

	for _, filePath := range files {
		err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			return insertFile(store, path)
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func insertFile(store *forensicstore.ForensicStore, srcpath string) error {
	file := forensicstore.NewFile()
	file.Name = filepath.Base(srcpath)

	dstpath, storeFile, teardown, err := store.StoreFile(srcpath)
	if err != nil {
		return fmt.Errorf("error storing file: %w", err)
	}
	defer teardown()

	srcFile, err := os.Open(srcpath) // #nosec
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer srcFile.Close()

	size, hashes, err := hashCopy(storeFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	file.Size = float64(size)
	file.ExportPath = filepath.ToSlash(dstpath)
	file.Hashes = hashes

	_, err = store.InsertStruct(file)
	return err
}

func hashCopy(dst io.Writer, src io.Reader) (int64, map[string]interface{}, error) {
	md5hash, sha1hash, sha256hash := md5.New(), sha1.New(), sha256.New() // #nosec
	size, err := io.Copy(io.MultiWriter(dst, sha1hash, md5hash, sha256hash), src)
	if err != nil {
		return 0, nil, err
	}
	return size, map[string]interface{}{
		"MD5":     fmt.Sprintf("%x", md5hash.Sum(nil)),
		"SHA-1":   fmt.Sprintf("%x", sha1hash.Sum(nil)),
		"SHA-256": fmt.Sprintf("%x", sha256hash.Sum(nil)),
	}, nil
}
