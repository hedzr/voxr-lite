/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hedzr/voxr-api/api/v10"
)

func BaseResult(ok bool, count int32, msg string, errCode int32, data []*any.Any) (res *v10.Result) {
	res = &v10.Result{Ok: ok, Count: count, Msg: msg, ErrCode: errCode, Data: data}
	return
}

func ErrorBaseResult(errCode int, msg string) (res *v10.Result) {
	res = &v10.Result{Msg: msg, ErrCode: int32(errCode)}
	return
}

func wrapIn(data proto.Message) *any.Any {
	a, err := ptypes.MarshalAny(data)
	if err == nil {
		return a
	}
	return nil
}

func wrapOut(a *any.Any, data proto.Message) (err error) {
	err = ptypes.UnmarshalAny(a, data)
	return
}
