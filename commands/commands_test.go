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

package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

func setup(subdirs ...string) (string, error) {
	tempDir, err := ioutil.TempDir("", "forensicstoreprocesstest")
	if err != nil {
		return "", err
	}

	if len(subdirs) == 0 {
		err = copy.Copy(filepath.Join("..", "test", "data"), tempDir)
	} else {
		for _, subdir := range subdirs {
			err = copy.Copy(filepath.Join("..", "test", "data", subdir), filepath.Join(tempDir, subdir))
		}
	}
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func cleanup(folders ...string) (err error) {
	for _, folder := range folders {
		err := os.RemoveAll(folder)
		if err != nil {
			return err
		}
	}
	return nil
}
