/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
)

func (s *ImCoreService) VerifyIdCardX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_VerifyIdCard, v10.Op_VerifyIdCardAck, "VerifyIdCard", ctx, req)
}

func (s *ImCoreService) VerifyIdCard(ctx context.Context, req *v10.Empty) (res *v10.Empty, err error) {
	r, e := s.VerifyIdCardX(ctx, req)
	res = r.(*v10.Empty)
	err = e
	return
}

func (s *ImCoreService) VerifyMobileNumberX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_VerifyMobileNumber, v10.Op_VerifyMobileNumberAck, "VerifyMobileNumber", ctx, req)
}

func (s *ImCoreService) VerifyMobileNumber(ctx context.Context, req *v10.Empty) (res *v10.Empty, err error) {
	r, e := s.VerifyMobileNumberX(ctx, req)
	res = r.(*v10.Empty)
	err = e
	return
}

func (s *ImCoreService) UploadIdCardX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_Unknown, v10.Op_SendSMSAck, "UploadIdCard", ctx, req)
}

func (s *ImCoreService) UploadIdCard(ctx context.Context, req *v10.Empty) (res *v10.Empty, err error) {
	r, e := s.UploadIdCardX(ctx, req)
	res = r.(*v10.Empty)
	err = e
	return
}

func (s *ImCoreService) UploadVerifyAttachmentsX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_Unknown, v10.Op_SendSMSAck, "UploadVerifyAttachments", ctx, req)
}

func (s *ImCoreService) UploadVerifyAttachments(ctx context.Context, req *v10.Empty) (res *v10.Empty, err error) {
	r, e := s.UploadVerifyAttachmentsX(ctx, req)
	res = r.(*v10.Empty)
	err = e
	return
}
