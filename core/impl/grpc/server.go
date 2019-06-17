/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/im_grpc_server"
	"google.golang.org/grpc"
)

var server *grpc.Server

// StopServer to stop the grpc server
func StopServer() {
	closeServer()
}

func closeServer() {
	if server != nil {
		server.Stop()
		server = nil
	}

	if Instance != nil {
		Instance.Shutdown()
		Instance = nil
	}
	if PrivateInstance != nil {
		PrivateInstance.Shutdown()
		PrivateInstance = nil
	}
}

// StartServer to start the grpc server
func StartServer() {
	if server != nil {
		return
	}

	// go newServer()
	server = im_grpc_server.StartServer(func(server *grpc.Server) {
		// switch id {
		// case "inx.im.apply":
		// 	pb.RegisterApplyServiceServer(server, &ApplyService{})
		// case "inx.im.core":
		cs := NewImCoreService()
		v10.RegisterImCoreServer(server, cs)
		Instance = cs

		csp := NewImCorePrivateService()
		v10.RegisterImCorePrivServer(server, csp)
		PrivateInstance = csp

		// enable grpc reflection. see also: https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md#enable-server-reflection
		// reflection.Register(server)

		// }
	})
}

// // never used
// func newServerNeverUsed() {
// 	if server != nil {
// 		return
// 	}
//
// 	var main_grpc = fmt.Sprintf("server.grpc.%v", vxconf.GetStringR("server.grpc.main", "vx-grpc"))
// 	if len(main_grpc) == 0 {
// 		return
// 	}
//
// 	grpcListen, id, disabled, port := tool.LoadGRPCListen(main_grpc)
// 	if disabled {
// 		logrus.Warnf("gRPC listen on %v but disabled.", grpcListen)
// 		logrus.Println("gRPC exiting...")
// 		return
// 	}
//
// 	// find an available port between starting port and portMax
// 	var portMax = port + 10
// 	var err error
// 	var listen net.Listener
// 	for {
// 		listen, err = net.Listen("tcp", grpcListen) // listen on tcp4 and tcp6
// 		if err != nil {
// 			logrus.Warnf("gRPC Failed to listen: %v", err)
// 			if port > portMax {
// 				logrus.Fatalf("gRPC Failed to listen: %v, port = %v", err, port)
// 				return
// 			}
// 			grpcListen, port = tool.IncGrpcListen(main_grpc)
// 		} else {
// 			break
// 		}
// 	}
//
// 	logrus.Infof("gRPC Listening at :%v....", tool.Port())
//
// 	// grpc 注册相应的 rpc 服务
// 	server = grpc.NewServer()
// 	switch id {
// 	case "inx.im.apply":
// 		v10.RegisterApplyServiceServer(server, &ApplyService{})
// 	case "inx.im.core":
// 		cs := NewImCoreService()
// 		v10.RegisterImCoreServer(server, cs)
//
// 		csp := NewImCorePrivateService()
// 		v10.RegisterImCorePrivServer(server, csp)
//
// 		// enable grpc reflection. see also: https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md#enable-server-reflection
// 		reflection.Register(server)
// 	}
// 	reflection.Register(server)
//
// 	logrus.Println("gRPC Started successfully...")
// 	// grpc Serve: it will block here
// 	if err := server.Serve(listen); err != nil {
// 		logrus.Fatalf("gRPC Failed to serve: %v", err)
// 	}
// 	logrus.Println("gRPC exiting...")
// }
