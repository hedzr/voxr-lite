/*
 * Copyright © 2019 Hedzr Yeh.
 */

package restful

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (d *daemonImpl) buildRoutes(mux *http.ServeMux) (err error) {
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/", echoHandler)
	return
}

func (d *daemonImpl) handle(w http.ResponseWriter, r *http.Request) {
	// Log the request protocol
	logrus.Printf("Got connection: %s", r.Proto)
	// Send a message back to the client
	_, _ = w.Write([]byte("Hello"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Hello, world!\n")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, r.URL.Path)
}
