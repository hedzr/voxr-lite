/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

// service.BuildEchoHandlerFunc() 使用对照表建立 restful 路由。
// 客户端的 restful 请求通过这些路由，执行对应的后端接口，并将结果返回给客户端。
// 部分地实现了 grpc-restful 的通用转换。
//
// 向后端接口的调用，使用 scheduler.Invoke() 异步地进行。
// 异步调用被 scheduler 内部的负载均衡器按照配置指定的算法进行调度。
func BuildEchoHandlerFunc(bi *BuildInf) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		// 	fmt.Println("c")
		// 	if err := recover(); err != nil {
		// 		fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		// 	}
		// 	fmt.Println("d")
		// }()

		if in, err := bi.preparingInputParam(c); err == nil {
			ch := make(chan bool)
			scheduler.Invoke(bi.svc, bi.pkg, bi.pbsvc, bi.FuncName, in, bi.Result(), func(e error, input *scheduler.Input, out proto.Message) {
				if e == nil {

					// trigger the post-process function:
					if r, ok := out.(*v10.Result); ok {
						res := bi.realResultTemplate()
						if res != nil && len(r.Data) > 0 {
							if err := ptypes.UnmarshalAny(r.Data[0], res); err != nil {
								logrus.Warnf("CANNOT decode `Any` to %v: %v", reflect.TypeOf(res), r)
								res = r
							} else {
								if bi.onEverythingOk != nil {
									bi.onEverythingOk(res)
								}
							}
						} else {
							logrus.Warnf("i did not write this code: bi.realResultTemplate() return nil currently.")
							res = r
						}
					} else {
						if bi.onEverythingOk != nil {
							bi.onEverythingOk(r)
						}
					}

					err = c.JSON(api.HttpOk, out)

				} else {
					// grpc invoke failed
					logrus.Errorf("grpc invoker return failed: %v", e)
					err = echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
				}
				ch <- true
			})

			<-ch
		} else {
			logrus.Errorf("bi.preparingInputParam failed. %v", err)
		}
		return
	}
}

// service.BuildFwdr() 为 grpc.Init() 所调用。
//
// ImCoreService 籍此建立转发器对照表，并将对 vx-core 的调用转换为到后端接口的调用，然后返回相应的结果。
//
// 向后端接口的调用，使用 scheduler.Invoke() 异步地进行。
// 异步调用被 scheduler 内部的负载均衡器按照配置指定的算法进行调度。
func BuildFwdr(bi *BuildInf) FwdrFunc {
	return func(ctx context.Context, in proto.Message) (res proto.Message, err error) {
		ch := make(chan bool)

		// logrus.Debugf("user.login invoking: %v", req.UserInfo.UNickname, req.UserInfo.UPass)
		scheduler.Invoke(bi.svc, bi.pkg, bi.pbsvc, bi.FuncName, in, bi.Result(), func(e error, input *scheduler.Input, out proto.Message) {
			if e == nil {
				logrus.Debugf("%v.%v.%v %v return: %v", bi.svc, bi.pkg, bi.pbsvc, bi.FuncName, out)
				if r, ok := out.(*v10.Result); ok && r.Ok && len(r.Data) > 0 {

					// trigger the post-process function:
					res = bi.realResultTemplate()
					if res != nil && len(r.Data) > 0 {
						if err := ptypes.UnmarshalAny(r.Data[0], res); err != nil {
							logrus.Warnf("CANNOT decode to %v: %v", reflect.TypeOf(res), r)
							res = r
						} else {
							if bi.onEverythingOk != nil {
								bi.onEverythingOk(res)
							}
						}
					} else {
						logrus.Warnf("i did not write this code: bi.realResultTemplate() return nil currently.")
						res = r
					}

				} else {
					if bi.onEverythingOk != nil {
						bi.onEverythingOk(res)
					}
				}
			} else {
				logrus.Errorf("%v.%v.%v %v return error: %v", bi.svc, bi.pkg, bi.pbsvc, bi.FuncName, e)
				res = nil // ErrorResult(http.StatusUnauthorized, "Please provide valid credentials")
				err = echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
			}
			ch <- true
		})

		<-ch
		return
	}
}

// func BuildDirectInvoker(bi *BuildInf) FwdrFunc {
// 	return func(ctx context.Context, in proto.Message) (res proto.Message, err error) {
//
// 	}
// }
