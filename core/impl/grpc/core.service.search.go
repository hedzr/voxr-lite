/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
)

func (s *ImCoreService) SearchGlobalX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.SearchGlobal(ctx, req.(*v10.SearchGlobalReq))
}

func (s *ImCoreService) SearchGlobal(ctx context.Context, req *v10.SearchGlobalReq) (res *v10.SearchGlobalReply, err error) {
	return
}
