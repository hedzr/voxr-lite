/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	voxr_common "github.com/hedzr/voxr-common"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-common/xs/health"
	"github.com/hedzr/voxr-lite/core/impl/chat"
	"github.com/hedzr/voxr-lite/core/impl/service"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

const (
	XS_SERVER_BANNER = `
 ___ __  __    ____                 ____
|_ _|  \/  |  / ___|___  _ __ ___  / ___|  ___ _ ____   _____ _ __
 | || |\/| | | |   / _ \| '__/ _ \ \___ \ / _ \ '__\ \ / / _ \ '__|
 | || |  | | | |__| (_) | | |  __/  ___) |  __/ |   \ V /  __/ |
|___|_|  |_|  \____\___/|_|  \___| |____/ \___|_|    \_/ \___|_|
`
)

type (
	handlers struct {
		// DB *mgo.Session
		// DB    *dbi.Config
		// mgoDB *mgo.Session
	}
)

func g(entry, defaultValue string) (s string) {
	s = vxconf.GetStringR(entry, defaultValue)
	return
}

func ds(val, defVal string) string {
	if len(val) == 0 {
		return defVal
	}
	return val
}

func (h *handlers) OnGetBanner() string {
	return XS_SERVER_BANNER
}

func (h *handlers) InitRoutes(e *echo.Echo, s vxconf.CoolServer) (ready bool) {
	ready = false

	e.POST(voxr_common.GetApiPrefix()+g("server.jwt.loginUrl", "/login"), h.login)
	// e.POST(common.GetApiPrefix()+g("server.jwt.signupUrl", "/signup"), h.register)
	e.POST(voxr_common.GetApiPrefix()+g("server.jwt.refreshToken", "/refresh-token"), h.refreshToken)

	e.POST(voxr_common.GetApiPrefix()+g("server.upload.url", "/upload"), h.upload)

	e.POST(voxr_common.GetApiPrefix()+"/circle/upload", h.circleUpload)

	service.Init(func(bi *service.BuildInf) {
		e.POST(voxr_common.GetApiPrefix()+bi.Entry, service.BuildEchoHandlerFunc(bi))
	})

	health.Enable(e)

	// h.initProxy(e)

	ready = true
	return
}

func (h *handlers) InitWebSocket(e *echo.Echo, s vxconf.CoolServer) (ready bool) {
	e.GET("/v1/ws", chat.WsHandler)
	e.GET("/v1/public/ws", chat.WsPublicHandler)
	// e.GET("/v2/ws", chat.WsPublicHandler)
	ready = true
	return
}

func (h *handlers) initProxy(e *echo.Echo) {
	url1, err := url.Parse("http://localhost:8081/aaa")
	if err != nil {
		e.Logger.Fatal(err)
	}
	url2, err := url.Parse("http://localhost:8082/bbb")
	if err != nil {
		e.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL: url1,
		},
		{
			URL: url2,
		},
	}
	e.Group("/c", middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer(targets),
		Rewrite: map[string]string{
			"/c/*": "/$1",
		},
	}))
}

func (h *handlers) login(c echo.Context) (err error) {
	// defer func() { // 必须要先声明defer，否则不能捕获到panic异常
	// 	fmt.Println("c")
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err) // 这里的err其实就是panic传入的内容，55
	// 	}
	// 	fmt.Println("d")
	// }()

	in := new(v10.LoginReq)
	if err = c.Bind(in); err != nil {
		return
	}

	ch := make(chan bool)
	// scheduler.Invoke("apply", "ApplyService", "/inx.im.exporttask.ApplyService/FetchApplyByUid", &exporttask.ApplyRequest{Uid: "9xdhasjkhas"}, func(e error, input *scheduler.Input, out interface{}) {
	// 	c.JSON(200, out)
	// 	ch <- true
	// })
	logrus.Debugf("/login user.login invoking: %v, %v, %v", in.UserInfo.Nickname, in.UserInfo.Phone, in.UserInfo.Pass)
	scheduler.Invoke(api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName, service.LoginMethod, in, nil, func(e error, input *scheduler.Input, out proto.Message) {
		if e == nil {
			logrus.Debugf("/login user.login return ok: out=%v", out)
			_ = c.JSON(api.HttpOk, out)

			if r, ok := out.(*v10.Result); ok && r.Ok && len(r.Data) > 0 {
				ret := &v10.UserInfoToken{}
				err = service.DecodeBaseResult(r, ret)
				if err == nil {
					service.PF().OnLoginOk(ret)
				} else {
					logrus.Warnf("invalid grpc response: %v", *r)
					err = echo.NewHTTPError(http.StatusInternalServerError, "Please login again later")
				}
			} else {
				logrus.Warnf("invalid grpc response, expect v10.Result{}: %v", out)
				err = echo.NewHTTPError(http.StatusInternalServerError, "Please login again later")
			}
			// if r, ok := out.(*base.Result); ok {
			// 	service.RecordUserHashG(r)
			// }
		} else {
			logrus.Debugf("/login user.login return error: err=%v", e)
			err = echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}
		// if r, ok := out.(*user.UserInfoToken); ok {
		// 	logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
		// 	c.JSON(api.HttpOk, out)
		// } else {
		// 	logrus.Warnf(">> Input: %v\nhas error??? output: %v", input, out)
		// }
		ch <- true
	})

	<-ch
	return
}

func (h *handlers) refreshToken(c echo.Context) (err error) {
	// defer func() { // 必须要先声明defer，否则不能捕获到panic异常
	// 	fmt.Println("c")
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err) // 这里的err其实就是panic传入的内容，55
	// 	}
	// 	fmt.Println("d")
	// }()

	in := new(v10.UserInfoToken)
	if err = c.Bind(in); err != nil {
		return
	}

	ch := make(chan bool)
	scheduler.Invoke(api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName, service.RefreshTokenMethod, in, nil, func(e error, input *scheduler.Input, out proto.Message) {
		if e == nil {
			_ = c.JSON(api.HttpOk, out)

			if r, ok := out.(*v10.Result); ok && r.Ok && len(r.Data) > 0 {
				ret := &v10.UserInfoToken{}
				err = service.DecodeBaseResult(r, ret)
				if err == nil {
					service.PF().OnRefreshTokenOk(ret)
				} else {
					logrus.Warnf("invalid grpc response: %v", *r)
				}
			}
		} else {
			// grpc invoke failed
			logrus.Errorf("grpc invoker return failed: %v", e)
			err = echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}

		ch <- true
	})

	<-ch
	return
}

func (h *handlers) Z(c echo.Context) (err error) {
	return
}
