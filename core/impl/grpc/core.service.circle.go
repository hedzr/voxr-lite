/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
)

func (s *ImCoreService) CircleOperateX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	var fn string
	if r, ok := req.(*v10.CircleAllReq); ok {
		if r.GetScr() != nil {
			fn = "SendCircle"
		} else if r.GetUpcr() != nil {
			fn = "UploadImage"
		} else if r.GetLcr() != nil {
			fn = "ListCircles"
		} else if r.GetAcr() != nil {
			fn = "AddCircle" // never used, replaced with SendCircle
		} else if r.GetGcr() != nil {
			fn = "GetCircle"
		} else if r.GetRcr() != nil {
			fn = "RemoveCircle"
		} else if r.GetUcr() != nil {
			fn = "UpdateCircle"
		}
		return s.xmas(v10.Op_CirclesAll, v10.Op_CirclesAllAck, fn, ctx, req)
	}
	res = &v10.CircleAllReply{ProtoOp: v10.Op_CirclesAllAck, Seq: 0, ErrorCode: 1001}
	return
}

func (s *ImCoreService) CircleOperate(ctx context.Context, req *v10.CircleAllReq) (res *v10.CircleAllReply, err error) {
	r, e := s.CircleOperateX(ctx, req)
	res = r.(*v10.CircleAllReply)
	err = e
	return
}
