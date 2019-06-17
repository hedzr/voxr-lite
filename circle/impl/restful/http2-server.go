/*
 * Copyright © 2019 Hedzr Yeh.
 */

package restful

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"path"
	"time"
)

func (d *daemonImpl) OnHttp2ServerRun(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}, port int) (err error) {
	d.mux = http.NewServeMux()
	err = d.buildRoutes(d.mux)
	if err != nil {
		return
	}

	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   d.mux, // http.HandlerFunc(d.handle),
		TLSConfig: d.checkAndEnableAutoCert(),
	}

	d.enableGracefulShutdown(srv, stopCh, doneCh)

	// TODO server push, ...
	// https://posener.github.io/http2/

	go func() {
		var e error
		// Start the server with TLS, since we are running HTTP/2 it must be
		// run with TLS.
		// Exactly how you would run an HTTP/1.1 server with TLS connection.
		if srv.TLSConfig.GetCertificate == nil {
			logrus.Printf("Serving on https://0.0.0.0:%d ...", port)
			certFile, keyFile := d.findLocalCertFiles()
			e = srv.ListenAndServeTLS(certFile, keyFile)
		} else {
			logrus.Printf("Serving on https://0.0.0.0:%d with autocert...", port)
			e = srv.ListenAndServeTLS("", "")
		}
		if e != nil {
			logrus.Error("restful server stopped failed: ", e)
		}
	}()

	return
}

func (d *daemonImpl) selectLocalCertDir() (dir string) {
	if cmdr.FileExists("./ci") {
		return os.ExpandEnv(defaultCertDirs[0]) // return ci/etc/xxxx/assets
	}
	return os.ExpandEnv(defaultCertDirs[1]) // return /etc/xxxx/assets
}

func (d *daemonImpl) findLocalCertFiles() (certFile, keyFile string) {
	for _, dir := range defaultCertDirs {
		dir = os.ExpandEnv(dir)
		if cmdr.FileExists(dir) {
			certFile = path.Join(dir, defaultCert)
			keyFile = path.Join(dir, defaultCertKey)
			if cmdr.FileExists(certFile) && cmdr.FileExists(keyFile) {
				return
			}
		}
	}
	return
}

func (d *daemonImpl) domains(topDomains ...string) (domainList []string) {
	for _, top := range vxconf.GetStringSliceR("server.domains", topDomains) {
		domainList = append(domainList, top)
		for _, s := range []string{"aurora", "api", "home", "res"} {
			domainList = append(domainList, fmt.Sprintf("%s.%s", s, top))
		}
	}
	return
}

func (d *daemonImpl) checkAndEnableAutoCert() (tlsConfig *tls.Config) {
	tlsConfig = &tls.Config{}

	if vxconf.GetBoolR("server.autocert.enabled", vxconf.IsProd()) {
		logrus.Debugf("...autocert enabled")
		d.certManager = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(d.domains()...), // 测试时使用的域名：example.com
			Cache:      autocert.DirCache("ci/certs"),
		}
		go func() {
			if err := http.ListenAndServe(":80", d.certManager.HTTPHandler(nil)); err != nil {
				logrus.Fatal("autocert tool listening on :80 failed.", err)
			}
		}()
		tlsConfig.GetCertificate = d.certManager.GetCertificate
	}

	return
}

func (d *daemonImpl) enableGracefulShutdown(srv *http.Server, stopCh, doneCh chan struct{}) {

	go func() {
		for {
			select {
			case <-stopCh:
				logrus.Debugf("...shutdown going on.")
				ctx, cancelFunc := context.WithTimeout(context.TODO(), 8*time.Second)
				defer cancelFunc()
				if err := srv.Shutdown(ctx); err != nil {
					logrus.Error("Shutdown failed: ", err)
				} else {
					logrus.Debugf("Shutdown ok.")
				}
				<-doneCh
				return
			}
		}
	}()

}

var (
	defaultCertDirs = []string{
		"ci/etc/$APPNAME/assets",
		"/etc/$APPNAME/assets",
		"$HOME/.$APPNAME/assets",
	}
)

const (
	idleTimeout         = 5 * time.Minute
	activeTimeout       = 10 * time.Minute
	maxIdleConns        = 1000
	maxIdleConnsPerHost = 100

	defaultCert    = "server.cert"
	defaultCertKey = "server.key"
)
