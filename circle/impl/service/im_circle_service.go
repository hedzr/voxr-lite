/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/tool"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-lite/circle/impl/dao"
	"github.com/hedzr/voxr-lite/circle/impl/models"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

type ImCircleService struct {
	dao *dao.CircleDao
	// daoMember *dao.MemberDao
}

func NewImCircleService() *ImCircleService {
	return &ImCircleService{dao.NewCircleDao()}
}

func (s *ImCircleService) List(ctx context.Context, req *v10.ListCirclesReq) (res *v10.ListCirclesReply, err error) {
	res = &v10.ListCirclesReply{ProtoOp: v10.Op_CirclesAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_CirclesAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	if req.UserId == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}

	var ret []*models.Circle
	var qry = ""
	var args interface{}
	if req.Anyone {
		// no limit
	} else if req.IncludesFriends {
		qry = "user_id in (?)"
		// find user's friend...
		if friends, e := s.dao.SearchUserFriends(req.UserId); e != nil {
			err = e
		} else {
			args = append(friends, req.UserId)
		}
	} else {
		qry = "user_id=?"
		args = req.UserId
	}

	if ret, err = s.dao.List(req.Limit, req.Start, req.OrderBy, qry, args); err != nil {
		return
	}

	res.ErrorCode = v10.Err_OK
	res.Circles = make([]*v10.Circle, len(ret))
	for ix, r := range ret {
		res.Circles[ix] = r.ToProto()
	}
	return
}

func (s *ImCircleService) Send(ctx context.Context, req *v10.SendCircleReq) (res *v10.SendCircleReply, err error) {
	res = &v10.SendCircleReply{ProtoOp: v10.Op_CirclesAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_CirclesAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	model := (&models.Circle{}).FromProto(req.Circle)
	err = s.dao.Add(model, req.ParentId)
	return
}

func (s *ImCircleService) UploadImage(ctx context.Context, req *v10.UploadCircleReq) (res *v10.UploadCircleReply, err error) {
	res = &v10.UploadCircleReply{ProtoOp: v10.Op_CirclesAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_CirclesAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	filename := req.FileName
	// immPart := time.Now().Format("2006/09")
	// fileName := fmt.Sprintf("%v/%s/%v", root, immPart, req.FileName)
	// if fileName, err = saveBlob(req.Blob, fileName); err != nil {
	// 	logrus.Errorf("CANT'T save to '%s': %v", fileName, err)
	// 	return
	// }

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
	if filePath, err = saveBlob(req.Blob, filePath); err != nil {
		logrus.Errorf("CANT'T save to '%s': %v", filePath, err)
		return
	}

	url := fmt.Sprintf("%v/%v", urlPrefix, strings.Trim(filePath, root))
	logrus.Debugf("<upload> url=%v, filePath=%v", url, filePath)
	// urls = append(urls, url)

	res.ErrorCode = v10.Err_OK
	res.CircleId = req.CircleId
	// url := fmt.Sprintf("%v/%v/%v", urlPrefix, immPart, strings.Trim(fileName, root))
	// logrus.Debugf("uploaded file's url: %s", res.ImageUrl)
	res.ImageUrl = url

	if err = s.dao.SaveImage(&models.CircleImage{CircleId: req.CircleId, UserId: req.UserId,
		BaseName:  path.Base(filename),
		Mime:      req.Mime,
		Size:      int64(req.Size),
		LocalPath: filename,
		Url:       res.ImageUrl,
	}); err != nil {
		return
	}

	if req.CircleId > 0 {
		if err = s.dao.Where("id=?", req.CircleId).Update("image_url", res.ImageUrl).Error; err != nil {
			return
		}
	}

	logrus.Debugf("image uploaded ok. url = %v | f=%v", url, req.FileName)
	return
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

func saveBlob(blob []byte, pathname string) (filePath string, err error) {
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

	var wf *os.File
	if wf, err = os.Create(pathname); err != nil {
		return
	}
	defer wf.Close()

	w := bufio.NewWriter(wf)
	if _, err = w.Write(blob); err != nil {
		return
	}

	filePath = pathname
	return
}

func (s *ImCircleService) Remove(ctx context.Context, req *v10.RemoveCircleReq) (res *v10.RemoveCircleReply, err error) {
	res = &v10.RemoveCircleReply{ProtoOp: v10.Op_CirclesAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_CirclesAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	if err = s.dao.RemoveById(req.CircleId); err != nil {
		return
	}

	res.ErrorCode = v10.Err_OK
	res.RowsAffected = 1
	return
}

func (s *ImCircleService) Update(ctx context.Context, req *v10.UpdateCircleReq) (res *v10.UpdateCircleReply, err error) {
	res = &v10.UpdateCircleReply{ProtoOp: v10.Op_CirclesAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_CirclesAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	model := (&models.Circle{}).FromProto(req.Circle)
	if err = s.dao.Update(model); err != nil {
		return
	}

	res.ErrorCode = v10.Err_OK
	res.RowsAffected = 1
	return
}

func (s *ImCircleService) Get(ctx context.Context, req *v10.GetCircleReq) (res *v10.GetCircleReply, err error) {
	res = &v10.GetCircleReply{ProtoOp: v10.Op_CirclesAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_CirclesAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	var model *models.Circle
	if model, err = s.dao.GetById(req.CircleId); err != nil {
		return
	}

	res.ErrorCode = v10.Err_OK
	res.Circle = model.ToProto()
	return
}
