/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
)

type ApplyService struct {
}

func (*ApplyService) FetchApplyByUid(ctx context.Context, req *v10.ApplyRequest) (*v10.Result, error) {
	if len(req.Uid) > 0 {
		val := &v10.ApplyResponse{Result: fmt.Sprintf("Hello, %v", req.Uid)}
		a, err := ptypes.MarshalAny(val)
		if err == nil {
			return &v10.Result{Ok: true, Count: 1, Msg: "OK", ErrCode: api.OK, Data: []*any.Any{a}}, nil
			// return &pb.ApplyResponse{Result: fmt.Sprintf("Hello, %v", req.Uid)}, nil
		}
	}
	return &v10.Result{Msg: "ApplyService empty response package.", ErrCode: api.ErrorUnknown}, nil
	// return &pb.ApplyResponse{Result: ""}, errors.New("ApplyService empty response package.")
}
