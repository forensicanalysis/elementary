package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/gorilla/mux"

	"github.com/forensicanalysis/elementary/server"
	"github.com/forensicanalysis/forensicstore"
)

type resW struct {
	*bytes.Buffer
	header http.Header
}

func NewResW() *resW {
	return &resW{
		Buffer: &bytes.Buffer{},
		header: http.Header{},
	}
}

func (r resW) Header() http.Header        { return r.header }
func (r resW) WriteHeader(statusCode int) {}

type closingBuffer struct{ *bytes.Buffer }

func (c closingBuffer) Close() error { return nil }

func open(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	storePath := strings.Trim(string(m.Payload), "\"")

	switch m.Name {
	case "new":
		_, teardown, err := forensicstore.New(storePath)
		if err != nil {
			return nil, err
		}
		defer teardown()
		return nil, storeWindow(storePath)
	// case "image":
	// 	_, teardown, err := forensicstore.New(storePath)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer teardown()
	// 	cmd := commands.DockerCommand(
	// 		"import-image",
	// 		"forensicanalysis/elementary-import-image:v0.3.5",
	// 		map[string]string{"mounts": "input-dir,artifacts-dir"},
	// 	)
	// 	cmd.Flags().Set("input-file", path.Dir(storePath))
	// 	cmd.Flags().Set("input-dir", path.Base(storePath))
	// 	cmd.Run(cmd, []string{path.Join(path.Dir(storePath), "elementary.forensicstore")})
	// 	return nil, storeWindow(storePath)
	case "open":
		return nil, storeWindow(storePath)
	}
	return
}

type save struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

// handleMessages handles messages
func handleStoreMessages(storeURL string) func(_ *astilectron.Window, m MessageIn) (payload interface{}, err error) {
	return func(_ *astilectron.Window, m MessageIn) (payload interface{}, err error) {
		log.Println("using store", storeURL)

		if m.Name == "save" {
			var s save
			err = json.Unmarshal(m.Payload, &s)
			if err != nil {
				return nil, err
			}
			store, teardown, err := forensicstore.Open(storeURL)
			if err != nil {
				return nil, err
			}
			defer teardown()
			srcFile, loadTeardown, err := store.LoadFile(s.Src)
			if err != nil {
				return nil, err
			}
			defer loadTeardown()

			destFile, err := os.Create(s.Dest)
			if err != nil {
				return nil, err
			}

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return nil, err
			}
			return nil, nil
		}

		router := server.Router(server.Commands(), func() []string { return []string{storeURL} })

		u, err := url.Parse("/api" + m.Name)
		if err != nil {
			return nil, err
		}
		buffer := closingBuffer{bytes.NewBuffer(m.Payload)}
		req := &http.Request{Method: m.Method, URL: u, Body: buffer}

		var match mux.RouteMatch
		if router.Match(req, &match) {
			log.Println("MATCHING ROUTE FOUND", u, match.Handler, match.Vars, storeURL)
			res := NewResW()
			match.Handler.ServeHTTP(res, req)
			var payload interface{}
			err = json.Unmarshal(res.Bytes(), &payload)
			if err != nil {
				return res.Bytes(), nil
			}
			return payload, nil
		} else {
			log.Println("NO MATCHING ROUTE FOUND", u, match.MatchErr)
		}

		return payload, err
	}
}
