/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"io"
	"net/http"
)

func buildH2Routes(mux *http.ServeMux) (err error) {
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/", echoHandler)
	return
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Hello, world!\n")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, r.URL.Path)
}
