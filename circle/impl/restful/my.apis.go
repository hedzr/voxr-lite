/*
 * Copyright © 2019 Hedzr Yeh.
 */

package restful

import (
	"encoding/base64"
	"fmt"
	voxr_common "github.com/hedzr/voxr-common"
	"github.com/hedzr/voxr-common/tool"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-common/xs/health"
	"github.com/hedzr/voxr-lite/circle/impl/dao"
	"github.com/hedzr/voxr-lite/circle/impl/models"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"path"
	"strconv"
	"strings"
)

type (
	Handlers struct {
		// DB *mgo.Session
		// DB    *dbi.Config
		// mgoDB *mgo.Session
	}
)

func (h *Handlers) Init(e *echo.Echo, s *myService) (ready bool) {
	ready = false

	health.Enable(e)

	var base = vxconf.GetStringR("server.upload.url", "/upload") + "/:idOrBase64Path"
	e.GET(base, uploadedImageFunc)
	e.GET(voxr_common.GetApiPrefix()+base, uploadedImageFunc)

	ready = true
	return
}

var circleDao = dao.NewCircleDao()

func uploadedImageFunc(c echo.Context) (err error) {
	var (
		str       = c.Param("idOrBase64Path")
		localPath string
		name      string
		id        uint64
		model     *models.CircleImage
	)

	localPath = fmt.Sprintf("%v/images/no.png", vxconf.GetStringR("server.static.root", "/var/lib/$APPNAME/public"))
	if !tool.FileExists(localPath) {
		localPath = fmt.Sprintf("%v/public/images/no.png", tool.GetCurrentDir())
	}
	name = "notfound"

	if id, err = strconv.ParseUint(str, 10, 64); err == nil {
		if model, err = circleDao.GetImageById(id); err != nil {
			return c.Inline("./images/no.png", "notfound")
		}
		localPath = model.LocalPath
		name = model.BaseName
	} else {
		var b []byte
		if b, err = base64.StdEncoding.DecodeString(str); err != nil {
			if tool.FileExists(str) {
				return c.Inline(str, path.Base(str))
			} else {
				// return c.Inline("./images/no.png", "notfound")
			}
		}
		x := strings.Trim(string(b), "\n")
		name = path.Base(x)
		if !tool.FileExists(x) {
			y := fmt.Sprintf("%v/%v", vxconf.GetStringR("server.static.root", "/var/lib/$APPNAME/public"), x)
			if !tool.FileExists(y) {
				y = fmt.Sprintf("%v/%v/%v", tool.GetCurrentDir(), "public", x)
				if tool.FileExists(y) {
					localPath = y
				}
			} else {
				localPath = y
			}
		} else {
			localPath = x
		}
	}
	// err = c.JSON(api.HttpOk, &health)
	logrus.Debugf("pwd=%v, localPath=%v", tool.GetCurrentDir(), localPath)
	return c.Inline(localPath, name)
}
