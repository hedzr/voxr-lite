/*
 * Copyright © 2019 Hedzr Yeh.
 */

package r2

// import (
// 	"context"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/grpclog"
// 	"google.golang.org/grpc/keepalive"
// 	"os"
// 	"time"
// )
//
// func B(){
// 	//将插件注册进gRPC
// 	RegisterResolver("/xxx_server/")
//
// 	keepAlive := keepalive.ClientParameters{
// 		10 * time.Second,
// 		20 * time.Second,
// 		true,
// 	}
//
// 	etcdCluster := "etcd:///" + util.AppConfig.ETCDCluster // 指定使用etcd来做名称解析
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
//
// 	// grpc.WithInsecure: 不使用安全连接
// 	// grpc.WithBalancerName("round_robin"), 轮询机制做负载均衡
// 	// grpc.WithBlock: 握手成功才返回
// 	// grpc.WithKeepaliveParams: 连接保活，防止因为长时间闲置导致连接不可用
// 	conn, err := grpc.DialContext(ctx, etcdCluster, grpc.WithInsecure(), grpc.WithBalancerName(const_grpc_lbname),
// 		grpc.WithBlock(), grpc.WithKeepaliveParams(keepAlive))
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stderr, os.Stderr, 9))
// 	qcloudRpcClient = proto.NewQcloudServiceClient(conn)
// }
//
