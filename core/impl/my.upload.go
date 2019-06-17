/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-common/tool"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-common/xs/mjwt"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func (h *handlers) upload(c echo.Context) (err error) {
	// Read form fields
	// name := c.FormValue("name")
	// email := c.FormValue("email")

	// -----------
	// Read file
	// -----------

	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			logrus.Errorln("Websocket client closing error:", err)
		}
	}()

	// Source
	file, err := c.FormFile(vxconf.GetStringR("server.upload.formField", "file"))
	if err != nil {
		return
	}
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(path.Join(h.UploadDir(c), file.Filename))
	if err != nil {
		return
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "msg": "OK"})
	// return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with fields name=%s and email=%s.</p>", file.Filename, name, email))
}

func (h *handlers) UploadDir(c echo.Context) (dir string) {
	base := vxconf.GetStringR("server.upload.base", "/var/lib/$APPNAME/public/uploads")
	if strings.HasPrefix(base, "./") {
		if cwd, err := os.Getwd(); err == nil {
			base = path.Join(cwd, base)
		} else {
			dir = "/tmp"
			return
		}
	} else {
		base = os.ExpandEnv(base)
	}

	token, valid, err := redis_op.JwtExtract(c)
	if err != nil {
		logrus.Warnf("err: %v", err)
	}
	if !valid {
		logrus.Warn("err: invaild: invalid jwt token")
	}

	if imtk, ok := token.Claims.(*mjwt.ImClaims); ok {
		now := time.Now()
		dir = fmt.Sprintf("%v/%v/%v/%v/%v", base, now.Year(), now.Month(), imtk.Id, imtk.DeviceId)
	} else {
		logrus.Warnf("err: invalid: token.Claims is not ImClaims")
	}
	return
}

func (h *handlers) circleUpload(c echo.Context) (err error) {
	// Read form fields
	// name := c.FormValue("name")
	// email := c.FormValue("email")

	// -----------
	// Read file
	// -----------

	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			logrus.Errorln("Websocket client closing error:", err)
		}
	}()

	var (
		mime, filename, size, seq, userId, circleId string
		seqNo, uid, sizeInt, circleIdInt            uint64
		form                                        *multipart.Form
		file                                        *multipart.FileHeader
		files                                       []*multipart.FileHeader
		src                                         multipart.File
		urls                                        []string
	)

	circleId = c.FormValue("circleId")
	userId = c.FormValue("userId")
	mime = c.FormValue("mime")
	filename = c.FormValue("filename")
	size = c.FormValue("size")
	seq = c.FormValue("seq")
	seqNo, err = strconv.ParseUint(seq, 10, 32)
	uid, err = strconv.ParseUint(userId, 10, 64)
	circleIdInt, err = strconv.ParseUint(circleId, 10, 64)
	sizeInt, err = strconv.ParseUint(size, 10, 64)

	form, err = c.MultipartForm()
	if err != nil {
		return err
	}

	files = form.File[vxconf.GetStringR("server.upload.formField", "file")]
	if files == nil {
		files = form.File[vxconf.GetStringR("server.upload.formField", "file")+"[]"]
	}

	for _, file = range files {
		// file, err = c.FormFile(vxconf.GetStringR("server.upload.formField", "file))
		// if err != nil {
		// 	return
		// }
		src, err = file.Open()
		if err != nil {
			return
		}
		defer src.Close()

		mime = ds(file.Header.Get("Content-Type"), mime)
		// size = g(file.Header.Get("Size"), size)
		// sizeInt, err = strconv.ParseUint(size, 10, 64)
		sizeInt = uint64(file.Size)
		filename = file.Filename

		if sizeInt == 0 {
			err = echo.NewHTTPError(http.StatusBadRequest, "size == 0")
			return
		}
		if uid == 0 {
			err = echo.NewHTTPError(http.StatusBadRequest, "userId == 0")
			return
		}
		if seqNo == 0 {
			err = echo.NewHTTPError(http.StatusBadRequest, "seq == 0")
			return
		}
		if len(filename) == 0 {
			err = echo.NewHTTPError(http.StatusBadRequest, "filename == ''")
			return
		}
		if len(mime) == 0 {
			err = echo.NewHTTPError(http.StatusBadRequest, "mime == ''")
			return
		}

		urlPrefix := vxconf.GetStringR("server.upload.url", "/public/uploads")
		root := vxconf.GetStringR("server.upload.base", "/var/lib/$APPNAME/public/uploads")
		if strings.Contains(root, "{{.AppName}}") {
			root = strings.Replace(root, "{{.AppName}}", vxconf.GetStringR("server.serviceName", "voxr-lite"), -1)
		}
		if !tool.FileExists(root) {
			logrus.Warnf("`static` warn not exists: root=%s", root)

			root = vxconf.GetStringR("server.static.root", "/var/lib/$APPNAME/public")
			if strings.Contains(root, "{{.AppName}}") {
				root = strings.Replace(root, "{{.AppName}}", vxconf.GetStringR("server.serviceName", "voxr-lite"), -1)
			}

			if !tool.FileExists(root) {
				logrus.Warnf("`static` warn not exists: root=%s", root)
				root = path.Join(tool.GetCurrentDir(), vxconf.GetStringR("server.static.urlPrefix", "/public"), path.Base(vxconf.GetStringR("server.upload.base", "/var/lib/$APPNAME/public/uploads"))) // 用和 urlPrefix 相同的名字
			} else {
				root = path.Join(root, urlPrefix)
			}
		}
		if root[0] != '/' {
			root = fmt.Sprintf("%v/%v", tool.GetCurrentDir(), root)
		}

		immPart := time.Now().Format("2006/09")
		filePath := fmt.Sprintf("%v/%s/%v", root, immPart, filename)
		if filePath, err = saveBlob(src, filePath); err != nil {
			logrus.Errorf("CANT'T save to '%s': %v", filePath, err)
			return
		}

		url := fmt.Sprintf("%v/%v", urlPrefix, strings.Trim(filePath, root))
		logrus.Debugf("<upload> url=%v, filePath=%v", url, filePath)
		urls = append(urls, url)

		// invoke `vx-circle` uploadImage API, save the record into DB.
		// NOTE: bug exists, but it's a temporary case.
		scheduler.Invoke(api.GrpcCircle, api.GrpcCirclePackageName, api.ImCircleActionName,
			"UploadImage", // &v10.CircleAllReq{
			// ProtoOp: v10.Op_CirclesAll, Seq: uint32(seqNo),
			// Oneof: &v10.CircleAllReq_Upcr{Upcr:
			&v10.UploadCircleReq{
				ProtoOp: v10.Op_CirclesAll, Seq: uint32(seqNo), UserId: uid, CircleId: circleIdInt,
				FileName: path.Base(filePath), Mime: mime, Size: sizeInt,
				Blob: getFileContentBlob(filePath),
				// },
				// },
			}, &v10.UploadCircleReply{}, func(e error, input *scheduler.Input, out proto.Message) {

				// 回调的处理过程，当前是被忽略的，不处理 vx-circle grpc 失败的情况，也不处理文件的最终url不一致的情况。

				if e == nil {

					// trigger the post-process function:
					if r, ok := out.(*v10.UploadCircleReply); ok {
						logrus.Infof("vx-circle return the url: %v, vx-core generated url: %v", r.ImageUrl, url)
						if url != r.ImageUrl {
							logrus.Warnf("不一致的url，刚刚上传的图片，出现了问题！url=%v, ret=%v", url, r.ImageUrl)
						}
					}

					// err = c.JSON(api.HttpOk, out)

				} else {
					// grpc invoke failed
					logrus.Errorf("grpc invoker return failed: %v", e)
					// err = echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
				}
			})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "msg": "OK", "url": urls[0], "urls": urls})
	// return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with fields name=%s and email=%s.</p>", file.Filename, name, email))
}

func exists(pathname string) (size int64, yes bool) {
	if fi, err := os.Stat(pathname); err != nil {
		return -1, false
	} else if fi.Size() == 0 {
		return 0, yes
	} else {
		return fi.Size(), true
	}
}

func getFileContentBlob(filePath string) (blob []byte) {
	var err error
	if blob, err = ioutil.ReadFile(filePath); err != nil {
		logrus.Warnf("cannot read the uploaded file: %v", filePath)
	}
	return
}

func saveBlob(src multipart.File, pathname string) (filePath string, err error) {
	dir := path.Dir(pathname)
	f := path.Base(pathname)
	e := path.Ext(pathname)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	if len(e) > 0 {
		f = f[0 : len(f)-len(e)]
	}
	f = path.Join(dir, f)

	var i = 1
	for size, exist := exists(pathname); exist; i++ {
		if size > 0 {
			pathname = fmt.Sprintf("%v.%05d%v", f, i, e)
			size, exist = exists(pathname)
		} else {
			break
		}
	}

	var dst *os.File
	if dst, err = os.Create(pathname); err != nil {
		return
	}
	defer dst.Close()

	// w := bufio.NewWriter(f)
	// if _, err = w.Write(blob); err != nil {
	// 	return
	// }

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return
	}

	filePath = pathname
	return
}
