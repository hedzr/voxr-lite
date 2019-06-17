/*
 * Copyright © 2019 Hedzr Yeh.
 */

package restful

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
)

type (
	daemonImpl struct {
		appTag      string
		certManager *autocert.Manager
		mux         *http.ServeMux
	}
)

// NewDaemon creates an `daemon.Daemon` object
func NewDaemon() daemon.Daemon {
	return &daemonImpl{}
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
	doDeregister()
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
		doStart()
		return d.OnXsServerRun(cmd, args, stopCh, doneCh)
	}

	return d.OnHttp2ServerRun(cmd, args, stopCh, doneCh, port)
}
