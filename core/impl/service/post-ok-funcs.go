/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hedzr/voxr-api/api/v10"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type (
	PostFunctions struct{}
)

var (
	pf = &PostFunctions{}
)

func PF() *PostFunctions {
	return pf
}

//

func (s *PostFunctions) OnLoginV11Ok(ret *v10.AuthReply) {
	recordUserHash(ret.GetUit())
}

func (s *PostFunctions) OnRefreshTokenV11Ok(ret *v10.AuthReply) {
	recordUserHash(ret.GetUit())
}

func (s *PostFunctions) OnRegisterV11Ok(ret *v10.AuthReply) {
	// recordUserHash(ret.GetUit())
}

func (s *PostFunctions) OnLoginOk(ret *v10.UserInfoToken) {
	recordUserHash(ret)
}

func (s *PostFunctions) OnRefreshTokenOk(ret *v10.UserInfoToken) {
	recordUserHash(ret)
}

//

func (s *PostFunctions) OnX(ret proto.Message) {
	//
}

//

//

//

//

func recordUserHash(r *v10.UserInfoToken) {
	var exp time.Duration = vxconf.GetDurationR("server.websocket.client-expiration", 60*time.Second)
	exp += time.Second
	redis_op.PutUserHash(uint64(r.UserInfo.Id), r.DeviceId /*uint64(r.UserInfo.DeviceID)*/, r.Token, exp)
}

func RefreshUserHash(uid uint64, did string, token string) {
	var exp time.Duration = vxconf.GetDurationR("server.websocket.client-expiration", 60*time.Second)
	exp += time.Second
	redis_op.PutUserHash(uid, did, token, exp)
}

func RemoveUserHash(uid uint64, did string, token string) {
	redis_op.DelUserHash(uid, did, token)
}

func DecodeBaseResult(r *v10.Result, ret proto.Message) (err error) {
	if r.Ok == false || r.Count != 1 || r.Data == nil || len(r.Data) < 1 {
		err = errors.New("no data")
	}

	if err = ptypes.UnmarshalAny(r.Data[0], ret); err == nil {
		return
	} else {
		err = errors.New(fmt.Sprintf("decode to %v failed from %v", reflect.TypeOf(ret), r))
		logrus.Warnf("cannot decode to %v: %v", reflect.TypeOf(ret), r)
	}
	return
}

// func RecordUserHashG(r *base.Result) {
// 	if r.Ok == false || r.Count != 1 {
// 		return
// 	}
//
// 	ret := &user.UserInfoToken{}
// 	if err := ptypes.UnmarshalAny(r.Data[0], ret); err == nil {
// 		recordUserHash(ret)
// 	} else {
// 		logrus.Warnf("cannot decode to %v: %v", reflect.TypeOf(ret), r)
// 	}
// }
//
// func recordUserHash(r *user.UserInfoToken) {
// 	redis_op.PutUserHash(uint64(r.UserInfo.Id), r.DeviceId, r.Token)
// }

// func RecordUserHashG(r *base.Result) {
// 	if r.Ok == false || r.Count != 1 || r.Data == nil || len(r.Data) < 1 {
// 		return
// 	}
//
// 	ret := &user.UserInfoToken{}
// 	if err := ptypes.UnmarshalAny(r.Data[0], ret); err == nil {
// 		pf.OnLoginOk(ret)
// 	} else {
// 		logrus.Warnf("cannot decode to %v: %v", reflect.TypeOf(ret), r)
// 	}
// }
