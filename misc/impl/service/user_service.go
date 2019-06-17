/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

// tokenBody是未加密的，payload是tokenBody加密base64后的
// token = payload.sign

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-api/util"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"strconv"
)

// const (
// 	secret = "A58oj3c8CsxdVpod"
// 	tokenDuration = 1 * 24 * 3600 * 1000 //一天（毫秒）
// 	tokenBodyLinkSymbol = "::"//连接符
// )

type TokenBody struct {
	Expire int64
	Uid    string
}

type AuthServer struct{}

func (s *AuthServer) Login(ctx context.Context, req *v10.LoginReq) (res *v10.Result, err error) {
	var (
		ui    *v10.UserInfo
		token string
		data  *any.Any
	)
	ui, err = dao.UserLoginUsePass(req.UserInfo)
	if err != nil {
		logrus.Errorf("login failed: %v", err)
		return
	}

	token, err = redis_op.JwtSign(strconv.FormatInt(ui.Id, 10), req.Device.Unique)
	if err != nil {
		logrus.Errorf("login failed: %v", err)
		return
	}
	logrus.Debugf("token maden: %v", token)

	data, err = ptypes.MarshalAny(&v10.UserInfoToken{UserInfo: ui, Token: token, DeviceId: req.Device.Unique})
	if err == nil {
		_, err = dao.SaveOrUpdateDevice(ui.Id, req.Device)
		if err != nil {
			logrus.Errorf("login failed (can't save device info): %v", err)
			return
		}
		res = &v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}}
		logrus.Debugf("Login ok: %v", res)
	} else {
		logrus.Errorf("Err: %v", err)
	}
	return
}

// func (s *AuthServer) Login1(ctx context.Context, lr *v10.LoginReq) (*v10.Result, error) {
// 	var result = v10.Result{Ok: false, Msg: "err", ErrCode: -1, Count: 0, Data: nil};
// 	user_login_res := dao.UserLoginUsePass(lr.UserInfo)
// 	fmt.Println("user_login_res:", user_login_res)
// 	if user_login_res == nil {
// 		result = v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 0, Data: nil};
// 		return &result, nil
// 	}
// 	uid := user_login_res.Id
// 	deviceId := lr.Device.Unique
// 	token, err := util.JwtSign(strconv.FormatInt(uid, 10), deviceId)
// 	fmt.Println("token:" + token)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	data, err := ptypes.MarshalAny(&v10.UserInfoToken{UserInfo: user_login_res, Token: token, DeviceId: lr.Device.Unique})
//
// 	if err == nil {
// 		_, err := dao.SaveDevice(user_login_res.Id, lr.Device)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result = v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}};
// 		logrus.Info("Login:", result)
// 		return &result, nil
// 	}
// 	return &result, err
// }

func (s *AuthServer) Register(ctx context.Context, req *v10.LoginReq) (res *v10.Result, err error) {
	var (
		ui    *v10.UserInfo
		token string
		data  *any.Any
	)
	res = util.MyBaseResult{}.FailResult()

	ui, err = dao.UserRegister(req.UserInfo)
	if err != nil {
		logrus.Errorf("register failed: %v", err)
		return
	}

	token, err = redis_op.JwtSign(strconv.FormatInt(ui.Id, 10), req.Device.Unique)
	if err != nil {
		logrus.Errorf("login failed: %v", err)
		return
	}
	logrus.Debugf("token maden: %v", token)

	data, err = ptypes.MarshalAny(&v10.UserInfoToken{UserInfo: ui, Token: token, DeviceId: req.Device.Unique})
	if err == nil {
		_, err = dao.SaveOrUpdateDevice(ui.Id, req.Device)
		if err != nil {
			logrus.Errorf("login failed (can't save device info): %v", err)
			return
		}
		res = &v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}}
		logrus.Debugf("Register ok: %v", res)
	} else {
		logrus.Errorf("Err: %v", err)
	}
	return
}

// func (s *AuthServer) Register1(ctx context.Context, lr *v10.LoginReq) (*v10.Result, error) {
// 	var result = v10.Result{Ok: false, Msg: "err", ErrCode: -1, Count: 0, Data: nil};
// 	user_register_res, err := dao.UserRegister(lr.UserInfo)
// 	fmt.Println("user_register_res:", user_register_res)
// 	if user_register_res == nil {
// 		result = v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 0, Data: nil};
// 		return &result, nil
// 	}
// 	uid := user_register_res.Id
// 	deviceId := lr.Device.Unique
// 	token, err := util.JwtSign(strconv.FormatInt(uid, 10), deviceId)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	data, err := ptypes.MarshalAny(&v10.UserInfoToken{UserInfo: user_register_res, Token: token, DeviceId: lr.Device.Unique})
//
// 	if err == nil {
// 		dao.SaveDevice(user_register_res.Id, lr.Device)
// 		result = v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}};
// 		logrus.Info("Register:", result)
// 		return &result, nil
// 	}
// 	return &result, err
// }

func (s *AuthServer) RefreshUserInfo(ctx context.Context, req *v10.UserInfoToken) (res *v10.Result, err error) {
	var (
		ui   *v10.UserInfo
		data *any.Any
	)
	res = util.MyBaseResult{}.FailResult()

	ui, err = dao.GetUserInfoByPhone(req.UserInfo.Phone)
	if err != nil {
		logrus.Errorf("register failed: %v", err)
		return
	}

	req.UserInfo = ui
	data, err = ptypes.MarshalAny(req)
	if err == nil {
		res = &v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}}
		logrus.Debugf("RefreshUserInfo ok: %v", res)
	} else {
		logrus.Errorf("Err: %v", err)
	}
	return
}

// func (s *AuthServer) RefreshUserInfo1(ctx context.Context, ut *v10.UserInfoToken) (*v10.Result, error) {
// 	var result = v10.Result{Ok: false, Msg: "err", ErrCode: -1, Count: 0, Data: nil};
// 	userInfo, err := dao.GetUserInfoByPhone(ut.UserInfo.Phone)
// 	ut.UserInfo = userInfo
//
// 	data, err := ptypes.MarshalAny(ut)
//
// 	if err == nil {
// 		result = v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}};
// 		return &result, nil
// 	}
//
// 	return &result, err
// }

func (s *AuthServer) UpdateUserInfo(ctx context.Context, req *v10.UserInfoToken) (res *v10.Result, err error) {
	var (
		ok   bool
		obj  *models.User
		data *any.Any
	)
	res = util.MyBaseResult{}.FailResult()

	ok, err, obj = dao.UpdateUserInfo(req.UserInfo)
	if err != nil {
		logrus.Errorf("UpdateUserInfo failed: %v | ok = %v", err, ok)
		return
	}

	req.UserInfo = obj.ToProto()
	data, err = ptypes.MarshalAny(req)
	if err == nil {
		res = &v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}}
		logrus.Debugf("UpdateUserInfo ok: %v", res)
	} else {
		logrus.Errorf("Err: %v", err)
	}
	return
}

// func (s *AuthServer) UpdateUserInfo1(ctx context.Context, ut *v10.UserInfoToken) (*v10.Result, error) {
// 	var result = v10.Result{Ok: false, Msg: "err", ErrCode: -1, Count: 0, Data: nil};
// 	dao.UpdateUserInfo(ut.UserInfo)
//
// 	data, err := ptypes.MarshalAny(ut)
//
// 	if err == nil {
// 		result = v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}};
// 		return &result, nil
// 	}
//
// 	return &result, err
// }

func (s *AuthServer) RefreshToken(ctx context.Context, req *v10.UserInfoToken) (res *v10.Result, err error) {
	var (
		ui    *v10.UserInfo
		valid bool
		token string
		data  *any.Any
	)
	res = util.MyBaseResult{}.FailResult()
	token = req.Token

	valid, err = redis_op.JwtVerifyToken(token)
	if err != nil {
		logrus.Errorf("RefreshToken failed: %v", err)
		return
	}
	if !valid {
		err = exception.UnwrapErr(exception.TokenErr)
		logrus.Errorf("RefreshToken failed: %v", err)
		return
	}

	token, err = redis_op.JwtSign(strconv.FormatInt(req.UserInfo.Id, 10), req.DeviceId)
	if err != nil {
		logrus.Errorf("RefreshToken failed: %v", err)
		return
	}

	ui, err = dao.GetUserInfoByAutoId(req.UserInfo.Id)
	if err != nil {
		logrus.Errorf("RefreshToken failed: %v", err)
		return
	}

	req.UserInfo = ui
	req.Token = token
	data, err = ptypes.MarshalAny(req)
	if err == nil {
		res = &v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}}
		logrus.Debugf("RefreshUserInfo ok: %v", res)
	} else {
		logrus.Errorf("Err: %v", err)
	}

	return
}

// func (s *AuthServer) RefreshToken1(ctx context.Context, ut *v10.UserInfoToken) (*v10.Result, error) {
// 	var result = util.MyBaseResult{}.FailResult()
//
// 	var token = ut.Token
//
// 	isOk, err := util.JwtVerifyToken(token)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !isOk {
// 		return nil, excep.UnwrapErr(excep.TokenErr)
// 	}
//
// 	token, err = util.JwtSign(strconv.FormatInt(ut.UserInfo.Id, 10), ut.DeviceId)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	userInfo, _ := dao.GetUserInfoByAutoId(ut.UserInfo.Id)
//
// 	data, err := ptypes.MarshalAny(&v10.UserInfoToken{UserInfo: userInfo, Token: token})
//
// 	if err == nil {
//
// 		result = util.MyBaseResult{}.SuccessResult([]*any.Any{data})
// 		return result, nil
// 	}
//
// 	return result, err
// }

func (s *AuthServer) ValidateToken(ctx context.Context, req *v10.UserInfoToken) (res *v10.Result, err error) {
	var (
		valid bool
		token string
		data  *any.Any
	)
	res = util.MyBaseResult{}.FailResult()
	validate := &v10.TokenValidate{Ok: false}
	token = req.Token

	valid, err = redis_op.JwtVerifyToken(token)
	if err != nil || !valid {
		data, err = ptypes.MarshalAny(validate)
		res.Data = []*any.Any{data}
		res.Count = 1
		if err == nil {
			err = exception.UnwrapErr(exception.TokenErr)
		}
		logrus.Errorf("ValidateToken failed: %v", err)
		return
	}

	validate.Ok = valid
	data, err = ptypes.MarshalAny(validate)
	if err == nil {
		res = &v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}}
		logrus.Debugf("ValidateToken ok: %v", res)
	} else {
		logrus.Errorf("Err: %v", err)
	}

	return
}

// func (s *AuthServer) ValidateToken1(ctx context.Context, ut *v10.UserInfoToken) (*v10.Result, error) {
// 	var result = util.MyBaseResult{}.FailResult()
//
// 	token := ut.Token
// 	isOk, err := util.JwtVerifyToken(token)
// 	if err != nil {
// 		return nil, err
// 	}
// 	validate := v10.TokenValidate{}
// 	validate.Ok = isOk
//
// 	data, err := ptypes.MarshalAny(&validate)
//
// 	if err == nil {
// 		result = util.MyBaseResult{}.SuccessResult([]*any.Any{data})
// 		return result, nil
// 	}
//
// 	return result, err
// }

func (s *AuthServer) ValidateTokenString(ctx context.Context, tokenT *v10.Token) (res *v10.Result, err error) {
	var (
		valid bool
		token string
		data  *any.Any
	)
	res = util.MyBaseResult{}.FailResult()
	validate := &v10.TokenValidate{Ok: false}
	token = tokenT.Token

	valid, err = redis_op.JwtVerifyToken(token)
	if err != nil || !valid {
		data, err = ptypes.MarshalAny(validate)
		res.Data = []*any.Any{data}
		res.Count = 1
		if err == nil {
			err = exception.UnwrapErr(exception.TokenErr)
		}
		logrus.Errorf("ValidateTokenString failed: %v", err)
		return
	}

	validate.Ok = valid
	data, err = ptypes.MarshalAny(validate)
	if err == nil {
		res = &v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}}
		logrus.Debugf("ValidateTokenString ok: %v", res)
	} else {
		logrus.Errorf("Err: %v", err)
	}

	return
}

// func (s *AuthServer) ValidateTokenString1(ctx context.Context, token *v10.Token) (*v10.Result, error) {
// 	var result = v10.Result{Ok: false, Msg: "err", ErrCode: -1, Count: 0, Data: nil};
// 	tokenStr := token.Token
// 	isOk, err := util.JwtVerifyToken(tokenStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	validate := v10.TokenValidate{}
// 	validate.Ok = isOk
//
// 	data, err := ptypes.MarshalAny(&validate)
//
// 	if err == nil {
// 		result = v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data}};
// 		return &result, nil
// 	}
// 	return &result, err
// }
