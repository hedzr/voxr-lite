/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-lite/core/impl/service"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

func (s *ImCoreService) RegisterV11X(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_RegisterV11, v10.Op_RegisterV11Ack, "RegisterV11", ctx, req)
}

func (s *ImCoreService) RegisterV11(ctx context.Context, req *v10.AuthReq) (res *v10.AuthReply, err error) {
	r, e := s.RegisterV11X(ctx, req)
	res = r.(*v10.AuthReply)
	err = e
	return
}

func (s *ImCoreService) LoginV11X(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_LoginV11, v10.Op_LoginV11Ack, "LoginV11", ctx, req)
}

func (s *ImCoreService) LoginV11(ctx context.Context, req *v10.AuthReq) (res *v10.AuthReply, err error) {
	r, e := s.LoginV11X(ctx, req)
	res = r.(*v10.AuthReply)
	err = e
	return
}

func (s *ImCoreService) RefreshTokenV11X(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_RefreshTokenV11, v10.Op_RefreshTokenV11Ack, "RefreshTokenV11", ctx, req)
}

func (s *ImCoreService) RefreshTokenV11(ctx context.Context, req *v10.AuthReq) (res *v10.AuthReply, err error) {
	r, e := s.RefreshTokenV11X(ctx, req)
	res = r.(*v10.AuthReply)
	err = e
	return
}

//

func (s *ImCoreService) UserOperateX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	var fn string
	if r, ok := req.(*v10.UserAllReq); ok {
		if r.GetAur() != nil {
			fn = "AddUser"
		} else if r.GetGur() != nil {
			fn = "GetUser"
		} else if r.GetLur() != nil {
			fn = "ListUsers"
		} else if r.GetRur() != nil {
			fn = "RemoveUser"
		} else if r.GetUur() != nil {
			fn = "UpdateUser"
		}
		return s.xmas(v10.Op_UserAll, v10.Op_UserAllAck, fn, ctx, req)
	}
	res = &v10.UserAllReply{ProtoOp: v10.Op_UserAllAck, Seq: 0, ErrorCode: 1001}
	return
}

func (s *ImCoreService) UserOperate(ctx context.Context, req *v10.UserAllReq) (res *v10.UserAllReply, err error) {
	r, e := s.UserOperateX(ctx, req)
	res = r.(*v10.UserAllReply)
	err = e
	return
}

//

func (s *ImCoreService) Login(ctx context.Context, request *v10.AuthReq) (res *v10.AuthReply, err error) {
	ch := make(chan bool)

	req := request.GetReq()
	res = &v10.AuthReply{ProtoOp: v10.Op_AuthAck, Oneof: &v10.AuthReply_Uit{Uit: &v10.UserInfoToken{}}}
	logrus.Debugf("    user.login invoking: %v, %v", req.UserInfo.Nickname, req.UserInfo.Pass)
	scheduler.Invoke(api.GrpcAuth, api.GrpcAuthPackageName, userActionName, loginMethod, req, nil, func(e error, input *scheduler.Input, out proto.Message) {
		if e == nil {
			logrus.Debugf("user.login return ok: %v", out)
			if r, ok := out.(*v10.Result); ok && r.Ok && r.Count == 1 {
				// RecordUserHashG(r)
				if err := ptypes.UnmarshalAny(r.Data[0], res.GetUit()); err == nil {
					service.PF().OnLoginOk(res.GetUit())
					// recordUserHash(res.GetUit())
				} else {
					logrus.Warnf("cannot decode to %v: %v", reflect.TypeOf(res), r)
				}
			}
		} else {
			logrus.Debugf("user.login return error: %v", e)
			res = nil // ErrorBaseResult(http.StatusUnauthorized, "Please provide valid credentials")
			err = echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}
		ch <- true
	})

	<-ch
	return
}

func (s *ImCoreService) RefreshUserInfo(ctx context.Context, request *v10.AuthReq) (res *v10.AuthReply, err error) {
	if fnRPC, ok := s.fwdrs["RefreshUserInfo"]; ok {
		req := request.GetReq()
		res = &v10.AuthReply{ProtoOp: v10.Op_AuthAck, Oneof: &v10.AuthReply_Uit{Uit: &v10.UserInfoToken{}}}
		var r proto.Message
		r, err = fnRPC(ctx, req)
		if rr, ok := r.(*v10.UserInfoToken); ok {
			res.Oneof = &v10.AuthReply_Uit{rr}
		}
	} else {
		logrus.Error(CANNOT_FOUND_FWDRS)
	}
	return
}

func (s *ImCoreService) UpdateUserInfo(ctx context.Context, request *v10.AuthReq) (res *v10.AuthReply, err error) {
	if fnRPC, ok := s.fwdrs["UpdateUserInfo"]; ok {
		req := request.GetReq()
		res = &v10.AuthReply{ProtoOp: v10.Op_AuthAck, Oneof: &v10.AuthReply_Uit{Uit: &v10.UserInfoToken{}}}
		var r proto.Message
		r, err = fnRPC(ctx, req)
		if rr, ok := r.(*v10.UserInfoToken); ok {
			res.Oneof = &v10.AuthReply_Uit{rr}
		}
	} else {
		logrus.Error(CANNOT_FOUND_FWDRS)
	}
	return
}

func (s *ImCoreService) RefreshToken(ctx context.Context, request *v10.AuthReq) (res *v10.AuthReply, err error) {
	if fnRPC, ok := s.fwdrs["RefreshToken"]; ok {
		req := request.GetReq()
		res = &v10.AuthReply{ProtoOp: v10.Op_AuthAck, Oneof: &v10.AuthReply_Uit{Uit: &v10.UserInfoToken{}}}
		var r proto.Message
		r, err = fnRPC(ctx, req)
		if rr, ok := r.(*v10.UserInfoToken); ok {
			res.Oneof = &v10.AuthReply_Uit{rr}
		}
	} else {
		logrus.Error(CANNOT_FOUND_FWDRS)
	}
	return
}

func (s *ImCoreService) ValidateToken(ctx context.Context, request *v10.AuthReq) (res *v10.AuthReply, err error) {
	if fnRPC, ok := s.fwdrs["ValidateToken"]; ok {
		req := request.GetReq()
		res = &v10.AuthReply{ProtoOp: v10.Op_AuthAck, Oneof: &v10.AuthReply_Uit{Uit: &v10.UserInfoToken{}}}
		var r proto.Message
		r, err = fnRPC(ctx, req)
		if rr, ok := r.(*v10.UserInfoToken); ok {
			res.Oneof = &v10.AuthReply_Uit{rr}
		}
	} else {
		logrus.Error(CANNOT_FOUND_FWDRS)
	}
	return
}

// func (s *ImCoreService) GenerateQrCode(ctx context.Context, request *core.AuthReq) (res *core.AuthReply, err error) {
// 	return
// }

const (
	CANNOT_FOUND_FWDRS = "Something was wrong. NO fwdrs definition can be FOUND."
)
