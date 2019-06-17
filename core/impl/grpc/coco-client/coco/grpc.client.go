/*
 * Copyright © 2019 Hedzr Yeh.
 */

package coco

import (
	"github.com/hedzr/voxr-api/api/v10"
	"google.golang.org/grpc"
)

type (
	GrpcClient struct {
		conn              *grpc.ClientConn
		imCoreClient      v10.ImCoreClient
		callback          func(c *GrpcClient, uit *v10.UserInfoToken)
		closeAfterQueried bool
	}
)

func (c *GrpcClient) RequestClose() {
	grpcHub.CloseClient(c)
}

func (c *GrpcClient) DoClose() {
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	if c.imCoreClient != nil {
		c.imCoreClient = nil
	}
	c.callback = nil
}

func (c *GrpcClient) Query(name string) {
	// Contact the server and print out its response.
	// grpcHub.querying <- []byte(name)
	grpcHub.QueryForClient(c, name)
}
