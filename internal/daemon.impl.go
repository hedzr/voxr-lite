/*
 * Copyright © 2019 Hedzr Yeh.
 */

package internal

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-lite/internal/restful"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
)

type (
	daemonImpl struct {
		appTag        string
		certManager   *autocert.Manager
		mux           *http.ServeMux
		doDeregister  func()
		doStart       func()
		handlers      restful.Handlers
		buildH2Routes func(mux *http.ServeMux) (err error)
	}
)

var chStop, chDone chan struct{}

//
//
//

// NewDaemon creates an `daemon.Daemon` object
func NewDaemon(doStart func(), doDeregister func(), h restful.Handlers, buildRoutes func(mux *http.ServeMux) (err error)) daemon.Daemon {
	return &daemonImpl{
		doStart:       doStart,
		doDeregister:  doDeregister,
		handlers:      h,
		buildH2Routes: buildRoutes,
	}
}

func AppStopCh() chan struct{} {
	return chStop
}

func AppDoneCh() chan struct{} {
	return chDone
}

//
//
//

func (d *daemonImpl) OnInstall(cxt *daemon.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnInstall", cmd.GetRoot().AppName) // panic("implement me")
	return
}

func (d *daemonImpl) OnUninstall(cxt *daemon.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnUninstall", cmd.GetRoot().AppName) // panic("implement me")
	return
}

func (d *daemonImpl) OnStatus(cxt *daemon.Context, cmd *cmdr.Command, p *os.Process) (err error) {
	fmt.Printf("%s v%v\n", cmd.GetRoot().AppName, cmd.GetRoot().Version)
	fmt.Printf("PID=%v\nLOG=%v\n", cxt.PidFileName, cxt.LogFileName)
	// panic("implement me")
	return
}

func (d *daemonImpl) OnReload() {
	logrus.Debugf("%s daemon OnReload", d.appTag) // panic("implement me")
}

func (d *daemonImpl) OnStop(cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnStop", cmd.GetRoot().AppName)
	d.doDeregister()
	return
}

func (d *daemonImpl) OnRun(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}) (err error) {
	d.appTag = cmd.GetRoot().AppName
	logrus.Debugf("%s daemon OnRun, pid = %v, ppid = %v", d.appTag, os.Getpid(), os.Getppid())

	port := vxconf.GetIntR("server.port", 2300)
	if port == 0 {
		logrus.Fatal("port not defined")
	}

	if vxconf.GetBoolR("server.xs-server.enabled", true) {
		d.doStart()
		return restful.NewXsServer(cmd, args, stopCh, doneCh, d.handlers)
	}

	return restful.NewH2(cmd, args, stopCh, doneCh, port, d.buildH2Routes)
}
