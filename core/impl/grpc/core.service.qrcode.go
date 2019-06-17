/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/skip2/go-qrcode"
)

func (s *ImCoreService) GenerateQrCodeX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.GenerateQrCode(ctx, req.(*v10.GenerateQrCodeReq))
}

func (s *ImCoreService) GenerateQrCode(ctx context.Context, req *v10.GenerateQrCodeReq) (res *v10.GenerateQrCodeReply, err error) {
	var blob []byte
	if blob, err = qrcode.Encode(req.Content, qrcode.RecoveryLevel(int(req.Level)), int(req.Size)); err == nil {
		res = &v10.GenerateQrCodeReply{ProtoOp: v10.Op_GenQrCode, Seq: req.Seq, ErrorCode: v10.Err_OK, Blob: blob}
	} else {
		res = &v10.GenerateQrCodeReply{ProtoOp: v10.Op_GenQrCodeAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_ERROR}
	}
	return
}
