/*
 * Copyright © 2019 Hedzr Yeh.
 */

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"path"
	"time"
)

type (
	daemonImpl struct {
		appTag      string
		certManager *autocert.Manager
		mux         *http.ServeMux
	}
)

//
//
//

// NewDaemon creates an `daemon.Daemon` object
func NewDaemon() daemon.Daemon {
	return &daemonImpl{}
}

func OnBuildCmd(root *cmdr.RootCommand) {
	_ = os.Setenv("APPNAME", root.AppName)

	cmdr.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {

		// app.server.port
		if cmd := cmdr.FindSubCommand("server", &root.Command); cmd != nil {

			// logrus.Debugf("`server` command found")

			opt := cmdr.NewCmdFrom(cmd)

			opt.NewFlag(cmdr.OptFlagTypeInt).
				Titles("p", "port").
				Description("the port to listen.", "").
				Group("").
				DefaultValue(2913, "PORT")

		}
	})
}

//
//
//

func (d *daemonImpl) OnInstall(cxt *daemon.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnInstall", cmd.GetRoot().AppName)
	return
	// panic("implement me")
}

func (d *daemonImpl) OnUninstall(cxt *daemon.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnUninstall", cmd.GetRoot().AppName)
	return
	// panic("implement me")
}

func (d *daemonImpl) OnStatus(cxt *daemon.Context, cmd *cmdr.Command, p *os.Process) (err error) {
	fmt.Printf("%s v%v\n", cmd.GetRoot().AppName, cmd.GetRoot().Version)
	fmt.Printf("PID=%v\nLOG=%v\n", cxt.PidFileName, cxt.LogFileName)
	return
}

func (d *daemonImpl) OnReload() {
	logrus.Debugf("%s daemon OnReload", d.appTag)
}

func (d *daemonImpl) OnStop(cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnStop", cmd.GetRoot().AppName)
	return
}

func (d *daemonImpl) OnRun(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}) (err error) {
	d.appTag = cmd.GetRoot().AppName
	logrus.Debugf("%s daemon OnRun, pid = %v, ppid = %v", d.appTag, os.Getpid(), os.Getppid())

	port := vxconf.GetIntR("server.port", 2300)
	if port == 0 {
		logrus.Fatal("port not defined")
	}

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
		// Start the server with TLS, since we are running HTTP/2 it must be
		// run with TLS.
		// Exactly how you would run an HTTP/1.1 server with TLS connection.
		if srv.TLSConfig.GetCertificate == nil {
			logrus.Printf("Serving on https://0.0.0.0:%d ...", port)
			certFile, keyFile := d.findLocalCertFiles()
			logrus.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
		} else {
			logrus.Printf("Serving on https://0.0.0.0:%d with autocert...", port)
			logrus.Fatal(srv.ListenAndServeTLS("", ""))
		}
	}()

	// go worker(stopCh, doneCh)
	return
}

func (d *daemonImpl) worker(stopCh, doneCh chan struct{}) {
LOOP:
	for {
		time.Sleep(3 * time.Second) // this is work to be done by worker.
		select {
		case <-stopCh:
			break LOOP
		default:
			logrus.Debugf("%s running at %d", d.appTag, os.Getpid())
		}
	}
	doneCh <- struct{}{}
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
			HostPolicy: autocert.HostWhitelist(d.domains("example.com")...), // 测试时使用的域名：example.com
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
