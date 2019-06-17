/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/sirupsen/logrus"
	"reflect"
)

type XmasPM interface {
	GetSeq() uint32
	GetProtoOp() v10.Op
}

func (s *ImCoreService) xmas(opIn, opAck v10.Op, fnName string, ctx context.Context, req proto.Message) (res proto.Message, err error) {
	if req == nil {
		return
	}
	if rreq, ok := req.(XmasPM); ok && rreq.GetProtoOp() == opIn && rreq.GetSeq() > 0 {
		if ci, ok := s.pooledEntries[opIn]; ok {
			if fnRPC, ok := s.fwdrs[fnName]; ok {
				var r proto.Message
				r, err = fnRPC(ctx, req) // see also: service.BuildFwdr()

				rrt := ci.OutTemplate(opAck, rreq.GetSeq(), v10.Err_BACKEND_INVOKE)
				if reflect.TypeOf(r) == reflect.TypeOf(rrt) {
					res = r
				} else {
					if err != nil {
						logrus.Warnf("    [core.service] invoke backend error: %v", err)
					} else {
						logrus.Warn("    [core.service] invoke backend generic error, no futher details")
					}
					res = rrt
				}
			}

		} else {
			logrus.Warn("    [core.service] no backend or no backend api found")
			res = ci.OutTemplate(opAck, rreq.GetSeq(), v10.Err_BACKEND_API_NOTFOUND)
		}
	} // else warn

	return
}
