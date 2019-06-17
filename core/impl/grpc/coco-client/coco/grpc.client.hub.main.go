/*
 * Copyright © 2019 Hedzr Yeh.
 */

package coco

// coco-client:
//   a psuedo ws client for login request
//
// while a ws text message ('u:XXX') sent to vx-core ws service,
// coco-client will make a gRPC request to vx-core grpc service 'Login'
//

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/tool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

const (
	defaultName = "world"
)

type (
	ClientHub struct {
		clients             map[*GrpcClient]bool
		exitCh              chan bool
		exited              bool
		newCocoClientAdding chan *GrpcClient
		destroyClient       chan *GrpcClient
		querying            chan []byte
	}
)

var (
	grpcHub = &ClientHub{
		clients:             make(map[*GrpcClient]bool),
		exitCh:              make(chan bool),
		exited:              true,
		newCocoClientAdding: make(chan *GrpcClient),
		destroyClient:       make(chan *GrpcClient),
		querying:            make(chan []byte),
	}
)

func GrpcStartClient() {
	go grpcHub.start()
}

func GrpcStopClient() {
	if !grpcHub.exited {
		grpcHub.exitCh <- true
	}
}

func ClientSendUnique(name string, callback func(c *GrpcClient, uit *v10.UserInfoToken)) {
	if len(grpcHub.clients) == 0 {
		noClose(name, callback)
	} else {
		go func() {
			sent := false
			for client, exists := range grpcHub.clients {
				if exists {
					client.Query(name)
					sent = true
				}
			}
			if sent == false {
				noClose(name, callback)
			}
		}()
	}
}

func ClientSend(name string, callback func(c *GrpcClient, uit *v10.UserInfoToken)) {
	if !grpcHub.exited {
		autoClose(name, callback)
	}
}

func autoClose(name string, callback func(c *GrpcClient, uit *v10.UserInfoToken)) {
	go func() {
		client := newCocoClient(callback)
		client.closeAfterQueried = true
		time.Sleep(1 * time.Second)
		client.Query(name)
	}()
}

func noClose(name string, callback func(c *GrpcClient, uit *v10.UserInfoToken)) {
	go func() {
		client := newCocoClient(callback)
		time.Sleep(1 * time.Second)
		client.Query(name)
	}()
}

func newCocoClient(callback func(c *GrpcClient, uit *v10.UserInfoToken)) (client *GrpcClient) {
	// Set up a connection to the server.
	listen, id, disabled, port := tool.LoadGRPCListen("server.grpc.vx-core")
	if disabled {
		logrus.Warnf("gRPC service '%v' is disabled.", id)
		return
	}
	if listen[0] == ':' {
		listen = "192.168.0.72" + listen
	}
	conn, err := grpc.Dial(listen, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("did not connect: %v, port is %v", err, port)
	}
	// defer conn.Close()

	c := v10.NewImCoreClient(conn)
	client = &GrpcClient{conn: conn, imCoreClient: c, callback: callback}

	// notify hub routine the new client connection object is incoming
	grpcHub.newCocoClientAdding <- client
	return
}
