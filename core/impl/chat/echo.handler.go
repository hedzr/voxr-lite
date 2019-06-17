/*
 * Copyright © 2019 Hedzr Yeh.
 */

package chat

import (
	"github.com/gorilla/websocket"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-common/xs/mjwt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    int(maxMessageSize) * 2,
	WriteBufferSize:   int(maxMessageSize) * 2,
	EnableCompression: true,
}

func WsPublicHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	// defer conn.Close()
	var userAgent string = c.Request().UserAgent()
	var cookies = c.Request().Cookies()
	// var token = c.Request().
	// uid, did := extractUserId(token)
	// token := JwtExtractor(c.Request())

	NewWsClient(conn, 0, "", "", userAgent, cookies, true)

	return nil
}

func WsHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	var (
		uid   uint64
		did   string
		ok    bool
		token string
		// err error
	)
	if uid, ok = c.Get("uid").(uint64); !ok {
		logrus.Warnf("[ws-upgrading] expect user-id is an uint64 number, but it doesn't.")
		return mjwt.ErrJWTMissing
	}
	if did, ok = c.Get("did").(string); !ok {
		logrus.Warnf("[ws-upgrading] expect device-id is a string, but it doesn't.")
		return mjwt.ErrJWTMissing
	}
	if token, ok = c.Get("token").(string); !ok {
		logrus.Warnf("[ws-upgrading] expect token is a string, but it doesn't.")
		return mjwt.ErrJWTMissing
	}

	// defer conn.Close()
	var userAgent string = c.Request().UserAgent()
	var cookies = c.Request().Cookies()
	// var token = c.Request().
	// uid, did := extractUserId(token)
	// token := JwtExtractor(c.Request())

	NewWsClient(conn, uid, did, token, userAgent, cookies, true)

	return nil
}

// never used
func extractUserId(token string) (uid uint64, did string) {
	tk, valid, err := redis_op.JwtDecodeToken(token)
	if err == nil && valid && tk != nil {
		if imtk, ok := tk.Claims.(*mjwt.ImClaims); ok {
			uid, err = strconv.ParseUint(imtk.Id, 10, 64)
			did = imtk.DeviceId
		}
	}
	return
}
