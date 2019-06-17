/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

// import (
// 	"flag"
// 	"github.com/hedzr/voxr-api/api/v10"
// 	"github.com/labstack/echo"
// 	"google.golang.org/grpc/grpclog"
// 	"net/http"
//
// 	"github.com/grpc-ecosystem/grpc-gateway/runtime"
// 	"golang.org/x/net/context"
// 	"google.golang.org/grpc"
// )
//
// var (
// 	echoEndpoint = flag.String("echo_endpoint", "localhost:9090", "endpoint of YourService")
// )
//
// func BuildApiHandler(e *echo.Echo, endpoint string) (err error) {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()
//
// 	opts := []grpc.DialOption{grpc.WithInsecure()}
//
// 	conn, err := grpc.Dial(endpoint, opts...)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			if cerr := conn.Close(); cerr != nil {
// 				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
// 			}
// 			return
// 		}
// 		go func() {
// 			<-ctx.Done()
// 			if cerr := conn.Close(); cerr != nil {
// 				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
// 			}
// 		}()
// 	}()
//
// 	mux := runtime.NewServeMux()
// 	err = v10.RegisterImCoreHandler(ctx, mux, conn)
// 	if err == nil {
// 		// TODO mux.Handle();
// 	}
//
// 	return
// }
//
// func run() error {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()
//
// 	mux := runtime.NewServeMux()
// 	opts := []grpc.DialOption{grpc.WithInsecure()}
// 	err := v10.RegisterImCoreHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
// 	if err != nil {
// 		return err
// 	}
//
// 	return http.ListenAndServe(":8080", mux)
// }
