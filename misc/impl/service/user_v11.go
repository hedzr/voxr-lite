/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"context"
	"github.com/hedzr/voxr-api/api/v10"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"github.com/sirupsen/logrus"
	"strconv"
)

type AuthServerV11 struct {
	UserServerV11
	AuthServer
	FriendServer
	UserContactServer
}

func NewAuthServerV11() *AuthServerV11 {
	return &AuthServerV11{}
}

func (s *AuthServer) LoginV11(ctx context.Context, req *v10.AuthReq) (res *v10.AuthReply, err error) {
	var (
		r     *v10.LoginReq
		ui    *v10.UserInfo
		token string
	)

	res = &v10.AuthReply{ProtoOp: v10.Op_AuthAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	r = req.GetReq()
	if r == nil {
		err = exception.New(exception.InvalidParams)
		return
	}

	ui, err = dao.UserLoginUsePass(r.UserInfo)
	if err != nil {
		logrus.Errorf("login failed: %v", err)
		return
	}

	token, err = redis_op.JwtSign(strconv.FormatInt(ui.Id, 10), r.Device.Unique)
	if err != nil {
		logrus.Errorf("login failed: %v", err)
		return
	}
	logrus.Debugf("token made: %v", token)

	_, err = dao.SaveOrUpdateDevice(ui.Id, r.Device)
	if err == nil {
		res.ErrorCode = v10.Err_OK
		res.Oneof = &v10.AuthReply_Uit{Uit: &v10.UserInfoToken{UserInfo: ui, Token: token, DeviceId: r.Device.Unique}}
		logrus.Debugf("Login v11 ok: %v", res)
	}
	return
}

func (s *AuthServer) RegisterV11(ctx context.Context, req *v10.AuthReq) (res *v10.AuthReply, err error) {
	var (
		r     *v10.LoginReq
		ui    *v10.UserInfo
		token string
	)

	res = &v10.AuthReply{ProtoOp: v10.Op_AuthAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	r = req.GetReq()
	if r == nil {
		err = exception.New(exception.InvalidParams)
		return
	}

	ui, err = dao.UserRegister(r.UserInfo)
	if err != nil {
		logrus.Errorf("register failed: %v", err)
		return
	}

	token, err = redis_op.JwtSign(strconv.FormatInt(ui.Id, 10), r.Device.Unique)
	if err != nil {
		logrus.Errorf("login failed: %v", err)
		return
	}
	logrus.Debugf("token made: %v", token)

	_, err = dao.SaveOrUpdateDevice(ui.Id, r.Device)
	if err == nil {
		res.ErrorCode = v10.Err_OK
		res.Oneof = &v10.AuthReply_Uit{Uit: &v10.UserInfoToken{UserInfo: ui, Token: token, DeviceId: r.Device.Unique}}
		logrus.Debugf("Register v11 ok: %v", res)
	}
	return
}

func (s *AuthServer) RefreshTokenV11(ctx context.Context, req *v10.AuthReq) (res *v10.AuthReply, err error) {
	var (
		r     *v10.UserInfoToken
		ui    *v10.UserInfo
		valid bool
		token string
	)

	res = &v10.AuthReply{ProtoOp: v10.Op_AuthAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	r = req.GetUit()
	if r == nil {
		err = exception.New(exception.InvalidParams)
		return
	}

	token = r.Token

	valid, err = redis_op.JwtVerifyToken(token)
	if err != nil {
		logrus.Errorf("RefreshToken failed: %v", err)
		return
	}
	if !valid {
		err = exception.New(exception.TokenErr)
		logrus.Errorf("RefreshToken failed: %v", err)
		return
	}

	token, err = redis_op.JwtSign(strconv.FormatInt(r.UserInfo.Id, 10), r.DeviceId)
	if err != nil {
		logrus.Errorf("RefreshToken failed: %v", err)
		return
	}

	if err == nil {
		res.ErrorCode = v10.Err_OK
		res.Oneof = &v10.AuthReply_Uit{Uit: &v10.UserInfoToken{UserInfo: ui, Token: token, DeviceId: r.DeviceId}}
		logrus.Debugf("RefreshToken v11 ok: %v", res)
	}

	return
}

func (s *AuthServer) AddRole(ctx context.Context, req *v10.RoleAllReq) (res *v10.RoleAllReply, err error) {
	return
}

func (s *AuthServer) RemoveRole(ctx context.Context, req *v10.RoleAllReq) (res *v10.RoleAllReply, err error) {
	return
}

func (s *AuthServer) GetRole(ctx context.Context, req *v10.RoleAllReq) (res *v10.RoleAllReply, err error) {
	return
}

func (s *AuthServer) UpdateRole(ctx context.Context, req *v10.RoleAllReq) (res *v10.RoleAllReply, err error) {
	return
}

func (s *AuthServer) ListRoles(ctx context.Context, req *v10.RoleAllReq) (res *v10.RoleAllReply, err error) {
	return
}

func (s *AuthServer) GetUserRoles(ctx context.Context, req *v10.RoleAllReq) (res *v10.RoleAllReply, err error) {
	return
}

func (s *AuthServer) SetUserRoles(ctx context.Context, req *v10.RoleAllReq) (res *v10.RoleAllReply, err error) {
	return
}
