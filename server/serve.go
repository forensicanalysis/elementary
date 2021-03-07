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
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"net/http"
	"regexp"
)

func setupResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	)
}

func logger() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Println(req.URL.String())

			next.ServeHTTP(w, req)
		})
	}
}

func ServeCommand(staticPath http.FileSystem, commands ...*Command) *cobra.Command {
	serveCmd := &cobra.Command{Use: "serve", Short: "Start the api server"}
	router := Router(commands, serveCmd.Flags().Args)
	router.Use(logger())
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(staticPath)))

	var port int
	serveCmd.Run = func(_ *cobra.Command, _ []string) {
		http.Handle("/", router)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port")
	return serveCmd
}

func Router(commands []*Command, args func() []string) *mux.Router {
	router := mux.NewRouter()

	re := regexp.MustCompile(`^(\w+)(?:\[(.+)\])?$`)

	for _, commandP := range commands {
		command := *commandP
		router.HandleFunc("/api"+command.Route, func(w http.ResponseWriter, r *http.Request) {
			setupResponse(&w)
			if r.Method == "OPTIONS" {
				return
			}

			var argsString []string
			for key, values := range r.URL.Query() {
				submatches := re.FindAllStringSubmatch(key, -1)
				group := ""
				if submatches[0][2] != "" {
					group = submatches[0][2] + ":"
				}
				for _, value := range values {
					argsString = append(argsString, "--"+submatches[0][1], group+value)
				}
			}
			argsString = append(argsString, args()...)
			flagset := &pflag.FlagSet{}
			setupFlags(&command, flagset, "json")
			if err := flagset.Parse(argsString); err != nil {
				fmt.Fprintf(w, "{'error':'%s'}", err)
				return
			}

			if err := command.Handler(w, r.Body, flagset); err != nil {
				fmt.Fprintf(w, "{'error':'%s'}", err)
				// w.WriteHeader(http.StatusInternalServerError)
			}
		}).Methods(command.Method, http.MethodOptions)
	}

	return router
}
