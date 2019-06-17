/*
 * Copyright © 2019 Hedzr Yeh.
 */

package coco

import (
	"context"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"time"
)

func (h *ClientHub) start() {
	logrus.Println("=== gRPC client hub: grpc client hub started.")
	ticker := time.NewTicker(8 * time.Second)
	defer func() {
		ticker.Stop()
		logrus.Infof("=== gRPC client hub: grpc client hub exited.")
	}()
	for {
		h.exited = false
		select {
		case exit := <-h.exitCh:
			if exit {
				logrus.Infof("=== gRPC client hub: grpc client hub exiting.")
				for k, ok := range h.clients {
					if ok && k.conn != nil {
						k.conn.Close()
					}
				}
				h.exited = exit
				return
			}

		case tm := <-ticker.C:
			if vxconf.GetBoolR("server.deps-debug", false) {
				logrus.Debugf("=== gRPC client hub: %v clients | %v", len(h.clients), tm)
			}

		case c := <-h.newCocoClientAdding:
			h.clients[c] = true                                               // flag the client conn object is true inside map `clients`
			logrus.Println("=== new gRPC client session created.", h.clients) // and log it for debugging

		case c := <-h.destroyClient:
			logrus.Println("=== destroying gRPC client: ", c)
			delete(h.clients, c)
			c.DoClose()

		case name := <-h.querying:
			// NOTE h.querying 支持多端点并发分发
			// NOTE 目前，多个clients 触发 h.querying 时，导致这里会重复对全部grpc连接进行请求
			// 直到请求结束和该 client 被 defer 所关闭（see also coco.ClientSend, coco.autoClose）和撤销为止。
			// 由于 coco 只是被用于内部测试，因此这个问题无需被修复。

			// logrus.Println(name)
			// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*len(h.clients))*time.Second)
			// defer cancel()
			ctx := context.Background()
			for c, ok := range h.clients {
				if ok && c.conn != nil {
					r, err := c.imCoreClient.Login(ctx, &v10.AuthReq{
						Oneof: &v10.AuthReq_Req{Req: &DemoLoginReq},
					})
					if err != nil {
						logrus.Warnf("=== gRPC client hub: could not greet: %v / %v", err, name)
						// } else if uit r.GetUit()Oneof.(*user.UserInfoToken) {
						// 	logrus.Warnf("could not greet - server return failed: code=%d, %v", r.ErrCode, r.Msg)
					} else {
						uit := r.GetUit()
						logrus.Printf("=== gRPC client hub: core.login return: %v", uit)
						if c.callback != nil {
							c.callback(c, uit)
						}
						if c.closeAfterQueried {
							go func() {
								time.Sleep(2 * time.Second)
								h.CloseClient(c)
							}()
						}
					}
				}
			}

		}
	}
}

func (h *ClientHub) CloseClient(client *GrpcClient) {
	if enabled, ok := h.clients[client]; ok && enabled {
		h.destroyClient <- client
	}
}

func (h *ClientHub) QueryForClient(client *GrpcClient, s string) {
	h.querying <- []byte(s)
}
