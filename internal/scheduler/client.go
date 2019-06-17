/*
 * Copyright © 2019 Hedzr Yeh.
 */

package scheduler

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/dc"
	balance2 "github.com/hedzr/voxr-lite/internal/scheduler/balance"
	"github.com/sirupsen/logrus"
)

type (
	GrpcClient struct {
		// Conn  *grpc.ClientConn
		Name      string
		Peers     []*balance2.ServicePeer
		invoking  chan *Input
		exiting   chan bool
		exited    bool
		depRecord *DepRecord
		balancer  balance2.Balancer
	}
)

// func init() {
// 	r2.RegisterResolver("voxr/services")
// }

func newClientInternal(serviceName string, dr *DepRecord) (client *GrpcClient) {
	client = lookupClients(serviceName, dr)
	go client.writePump()
	return
}

func (c *GrpcClient) Send(input *Input) {
	input.client = c
	c.invoking <- input
}

func (c *GrpcClient) reset_() {
	if c.exited {
		return
	}

	// do sth
}

func (c *GrpcClient) writePump() {
	c.exited = false
	defer func() {
		logrus.Debugf("GrpcClient.writePump - stopped.")
	}()
	for {
		select {
		case exit := <-c.exiting:
			if exit {
				logrus.Infof("GrpcClient.writePump - grpc client write pump exiting.")
				// ticker.Stop()
				c.reset_()
				c.exited = true
				return
			}

		case input := <-c.invoking:
			// invoke_(input, grpc.EmptyCallOption{})
			shortDebugf("    gRPC invoking: %v", input)
			invoke_nolock_(input)

		}
	}

}

func shortDebugf(fmt string, in *Input) {
	i := &Input{}
	_ = dc.StandardCopier.Copy(i, in)
	if c, ok := i.In.(*v10.CircleAllReq); ok {
		if c.GetUpcr() != nil {
			c.GetUpcr().Blob = nil
			logrus.Debugf("    gRPC invoking: %v", c)
			return
		}
	} else if c, ok := i.In.(*v10.UploadCircleReq); ok {
		c.Blob = nil
		logrus.Debugf("    gRPC invoking: %v", c)
		return
	}
	logrus.Debugf("    gRPC invoking: %v", i)
}
