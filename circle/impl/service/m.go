/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hedzr/voxr-api/api/v10"
)

func M(data proto.Message) (res *v10.Result, err error) {
	var data1 *any.Any
	data1, err = ptypes.MarshalAny(data)
	if err == nil {
		res = &v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: 1, Data: []*any.Any{data1}}
	}
	return
}

func ML(dataList []proto.Message) (res *v10.Result, err error) {
	var dataA []*any.Any
	for _, data := range dataList {
		var data1 *any.Any
		data1, err = ptypes.MarshalAny(data)
		if err == nil {
			dataA = append(dataA, data1)
		} else {
			break
		}
	}
	if err == nil {
		res = &v10.Result{Ok: true, Msg: "success", ErrCode: 0, Count: int32(len(dataList)), Data: dataA}
	}
	return
}
