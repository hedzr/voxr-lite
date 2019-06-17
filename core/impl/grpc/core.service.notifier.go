/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
)

func (s *ImCoreService) SendSMSX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_SendSMS, v10.Op_SendSMSAck, "SendSMS", ctx, req)
}

func (s *ImCoreService) SendSMS(ctx context.Context, req *v10.Empty) (res *v10.Empty, err error) {
	r, e := s.SendSMSX(ctx, req)
	res = r.(*v10.Empty)
	err = e
	return
}

func (s *ImCoreService) SendMailX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_SendMail, v10.Op_SendMailAck, "SendMail", ctx, req)
}

func (s *ImCoreService) SendMail(ctx context.Context, req *v10.Empty) (res *v10.Empty, err error) {
	r, e := s.SendMailX(ctx, req)
	res = r.(*v10.Empty)
	err = e
	return
}
