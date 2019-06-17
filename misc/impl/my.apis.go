/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-common/xs/health"
	"github.com/labstack/echo"
)

const (
	XS_SERVER_BANNER = `
 ___ __  __   ____
|_ _|  \/  | / ___|  ___ _ ____   _____ _ __
 | || |\/| | \___ \ / _ \ '__\ \ / / _ \ '__|
 | || |  | |  ___) |  __/ |   \ V /  __/ |
|___|_|  |_| |____/ \___|_|    \_/ \___|_|
`
)

type (
	handlers struct {
		// DB *mgo.Session
		// DB    *dbi.Config
		// mgoDB *mgo.Session
	}
)

func (h *handlers) OnGetBanner() string {
	return XS_SERVER_BANNER
}

func (h *handlers) InitRoutes(e *echo.Echo, s vxconf.CoolServer) (ready bool) {
	ready = false

	health.Enable(e)

	ready = true
	return
}

func (h *handlers) InitWebSocket(e *echo.Echo, s vxconf.CoolServer) (ready bool) {
	return
}
