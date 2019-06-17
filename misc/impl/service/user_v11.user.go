/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"context"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"github.com/sirupsen/logrus"
)

type UserServerV11 struct {
	dao dao.UserDao
}

func (s *UserServerV11) AddUser(ctx context.Context, req *v10.UserAllReq) (res *v10.UserAllReply, err error) {
	res = &v10.UserAllReply{ProtoOp: v10.Op_UserAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_UserAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	var id uint64
	model := (&models.User{}).FromProto(req.GetAur().User)
	id, err = s.dao.Add(model)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	res.Oneof = &v10.UserAllReply_Aur{Aur: &v10.AddUserReply{Id: id}}
	res.ErrorCode = v10.Err_OK
	return
}

func (s *UserServerV11) RemoveUser(ctx context.Context, req *v10.UserAllReq) (res *v10.UserAllReply, err error) {
	res = &v10.UserAllReply{ProtoOp: v10.Op_UserAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_UserAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	var (
		rows int64
		tmpl = new(models.User)
	)
	if req.GetRur().Id > 0 {
		rows, err = s.dao.RemoveById(req.GetRur().Id)
	} else if req.GetRur().User != nil {
		rows, err = s.dao.Remove(tmpl.FromProto(req.GetRur().User))
	}
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	res.Oneof = &v10.UserAllReply_Rur{Rur: &v10.RemoveUserReply{RowsAffected: rows}}
	res.ErrorCode = v10.Err_OK
	// mq.RaiseEvent(v10.GlobalEvents_EvOrgRemoved, res)
	return
}

func (s *UserServerV11) UpdateUser(ctx context.Context, req *v10.UserAllReq) (res *v10.UserAllReply, err error) {
	res = &v10.UserAllReply{ProtoOp: v10.Op_UserAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_UserAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	var (
		rows int64
		ret  *models.User
	)
	model := (&models.User{}).FromProto(req.GetUur().User)
	ret, rows, err = s.dao.Update(model)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		logrus.Errorf("ret=%v", ret)
		return
	}

	res.Oneof = &v10.UserAllReply_Uur{Uur: &v10.UpdateUserReply{RowsAffected: rows}}
	res.ErrorCode = v10.Err_OK
	// mq.RaiseEvent(v10.GlobalEvents_EvOrgUpdated, res)
	return
}

func (s *UserServerV11) ListUsers(ctx context.Context, req *v10.UserAllReq) (res *v10.UserAllReply, err error) {
	res = &v10.UserAllReply{ProtoOp: v10.Op_UserAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_UserAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	var (
		arr  []*v10.UserInfo
		ret  []*models.User
		tmpl = new(models.User)
	)
	ret, err = s.dao.List(tmpl.FromProto(req.GetLur().User))
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	for _, m := range ret {
		mx := m.ToProto()
		arr = append(arr, mx)
	}

	res.Oneof = &v10.UserAllReply_Lur{Lur: &v10.ListUsersReply{Users: arr}}
	res.ErrorCode = v10.Err_OK
	return
}

func (s *UserServerV11) GetUser(ctx context.Context, req *v10.UserAllReq) (res *v10.UserAllReply, err error) {
	res = &v10.UserAllReply{ProtoOp: v10.Op_UserAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_UserAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	var ret *models.User
	var tmpl *models.User
	tmpl = &models.User{Id: req.GetGur().Id, Mobile: req.GetGur().Mobile, Nickname: req.GetGur().Name, Email: req.GetGur().Email}
	ret, err = s.dao.Get(tmpl)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	res.Oneof = &v10.UserAllReply_Gur{Gur: &v10.GetUserReply{User: ret.ToProto()}}
	res.ErrorCode = v10.Err_OK
	return
}
