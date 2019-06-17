/*
 * Copyright © 2019 Hedzr Yeh.
 */

package restful

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func (d *H2) buildRoutes(mux *http.ServeMux) (err error) {
	err = d.onBuildRoutes(mux)
	return
}

func (d *H2) handle(w http.ResponseWriter, r *http.Request) {
	// Log the request protocol
	logrus.Printf("Got connection: %s", r.Proto)
	// Send a message back to the client
	_, _ = w.Write([]byte("Hello"))
}
