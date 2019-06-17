/*
 * Copyright © 2019 Hedzr Yeh.
 */

package scheduler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// Usage:
// scheduler.Invoke("vx-user", "user", "UserAction", "Login", &user.LoginReq{...}, nil, func(e error, input *scheduler.Input, out interface{}) {
//    ...
// })
//
// 注意：out 传入 nil 意味着对方接口是返回 *v10.Result 结构的；如果返回类型不是，则应传入正确的类型和实例以便容纳返回值
func Invoke(serviceName, packageName, pbServiceName, methodName string, in proto.Message, out proto.Message, callback func(error, *Input, proto.Message)) {
	m := methodName
	if !strings.HasPrefix(methodName, "/") {
		m = fmt.Sprintf("/inx.im.%s.%s/%s", packageName, pbServiceName, methodName)
	}

	// InvokeDirectViaHub(serviceName, pbServiceName, m, in, callback)

	if isPreferredToRealService {
		if _, ok := grpcHub.byIds[serviceName]; !ok {
			if s, ok := synonym[serviceName]; ok {
				serviceName = s
			}
		}
	} else {
		if s, ok := synonym[serviceName]; ok {
			serviceName = s
		}
	}

	if c, ok := grpcHub.byIds[serviceName]; ok { // this line might need to be locked by RWLock???
		if len(c.Peers) == 0 {
			logrus.Errorf("NO Peers FOUND for %s", serviceName)
			if callback != nil {
				callback(fmt.Errorf("NO Peers FOUND for %s", serviceName), nil, nil)
			}
			RequestRefreshClient(serviceName)
			return
		}

		// using the nolock invoker for each client now.
		c.Send(&Input{
			ServiceName:   serviceName,
			PBServiceName: pbServiceName,
			MethodName:    m,
			In:            in,
			Out:           out,
			Callback:      callback,
		})
	} else {
		logrus.Errorf("Unrecognized serviceName '%v'", serviceName)
		if callback != nil {
			callback(fmt.Errorf("Unrecognized serviceName '%v'", serviceName), nil, nil)
		}
	}

	return
}

// InvokeDirectViaHub invokes grpc api via `grpcHub` main loop and a child go-routine
// 除非是 vx-core, 并且目标grpc服务在同一进程空间内，否则不要使用此功能(InvokeDirectViaHub)
func InvokeDirectViaHub(serviceName, pbServiceName, methodName string, in proto.Message, out proto.Message, callback func(error, *Input, proto.Message)) {
	grpcHub.invoking <- &Input{
		ServiceName: serviceName, PBServiceName: pbServiceName, MethodName: methodName,
		In: in, Out: out, Callback: callback}
	return
}

//
//
//
//
//

func grpcInvokeHandler(c echo.Context) (err error) {
	// 暂未使用
	return
}

func invoke_nolock_(input *Input, opts ...grpc.CallOption) {
	peer, err := input.client.balancer.Pick(input.client.Peers)

	if err != nil {
		logrus.Errorf("CAN'T Pick peer via balancer: %v", err)
		return
	}

	var timeout = vxconf.GetDurationR("server.grpc.settings.query-timeout", 10*time.Second)
	logrus.Debugf("    invoke_nolock_(%v): picked peer: %v", timeout, peer.Record)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var out = input.Out
	if out == nil {
		out = new(v10.Result)
	}
	err = peer.Conn.Invoke(ctx, input.MethodName, input.In, out, opts...)
	if err != nil {
		if input.Callback != nil {
			input.Callback(err, input, nil)
		} else {
			logrus.Errorf("    invoke_nolock_(%v): error after invoked: %v", timeout, err)
		}
	} else {
		logrus.Debugf("    invoke_nolock_(%v): invoked: out=%v", timeout, out)
		if input.Callback != nil {
			input.Callback(nil, input, out)
		}
	}
}

// never used
func invoke__(input *Input, opts ...grpc.CallOption) {
	if s, ok := synonym[input.ServiceName]; ok {
		input.ServiceName = s
	}

	if c, ok := grpcHub.byIds[input.ServiceName]; ok {
		if len(c.Peers) == 0 {
			return
		}

		input.client = c

		go func() {
			ix := rand.Intn(len(c.Peers))
			peer := c.Peers[ix]

			logrus.Debugf("    invoke__: pick peer: %v", peer.Record)

			ctx, cancel := context.WithTimeout(context.Background(), vxconf.GetDurationR("server.grpc.query-timeout", 10*time.Second))
			defer cancel()

			// var out = make([]byte, 2048) //TODO waiting for Unify ResultReply
			// var out = new(user.UserInfoToken)
			var out = new(v10.Result)
			err := peer.Conn.Invoke(ctx, input.MethodName, input.In, out, opts...)
			if err != nil {
				if input.Callback != nil {
					input.Callback(err, input, nil)
				}
			} else {
				if input.Callback != nil {
					input.Callback(nil, input, out)
				}
			}
		}()
	}
}

// func invoke_(input *Input, opt grpc.CallOption) {
// 	if c, ok := grpcHub.byIds[input.ServiceName]; ok {
// 		go func() {
// 			ix := rand.Intn(len(c.Peers))
// 			peer := c.Peers[ix]
//
// 			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*len(grpcHub.clients))*time.Second)
// 			defer cancel()
//
// 			switch input.PBServiceName {
// 			case api.GrpcCore:
// 			case "UserAction":
// 				invoke_user_UserAction_(peer, ctx, input, opt)
// 			case "FriendAction":
// 				invoke_user_FriendAction_(peer, ctx, input, opt)
//
// 			}
// 		}()
// 	}
// }
//
// func invoke_user_FriendAction_(peer *GrpcPeer, ctx context.Context, input *Input, opt grpc.CallOption) {
// 	c := user.NewFriendActionClient(peer.Conn)
// 	ret, err := _invoke(c, input.MethodName, ctx, input.In, opt)
// 	if err != nil {
// 		logrus.Warnf("querying '%s' failed: %v", input.MethodName, err)
// 	} else {
// 		input.Callback(input, ret.Interface())
// 	}
//
// 	//switch input.MethodName {
// 	//case "AddFriend":
// 	//	ret, err := c.AddFriend(ctx, input.In.(*user.Relation))
// 	//	if err != nil {
// 	//		logrus.Warnf("querying '%s' failed: ", input.MethodName, err)
// 	//	} else {
// 	//		input.Callback(input, ret)
// 	//	}
// 	//case "UpdateFriend":
// 	//	ret, err := c.UpdateFriend(ctx, input.In.(*user.Relation))
// 	//	if err != nil {
// 	//		logrus.Warnf("querying '%s' failed: ", input.MethodName, err)
// 	//	} else {
// 	//		input.Callback(input, ret)
// 	//	}
// 	//case "DeleteFriend":
// 	//	ret, err := c.DeleteFriend(ctx, input.In.(*user.Relation))
// 	//	if err != nil {
// 	//		logrus.Warnf("querying '%s' failed: ", input.MethodName, err)
// 	//	} else {
// 	//		input.Callback(input, ret)
// 	//	}
// 	//case "GetFriendList":
// 	//	ret, err := c.GetFriendList(ctx, input.In.(*user.Relation))
// 	//	if err != nil {
// 	//		logrus.Warnf("querying '%s' failed: ", input.MethodName, err)
// 	//	} else {
// 	//		input.Callback(input, ret)
// 	//	}
// 	//}
// }
//
// func invoke_user_UserAction_(peer *GrpcPeer, ctx context.Context, input *Input, opt grpc.CallOption) {
// 	c := user.NewUserActionClient(peer.Conn)
// 	ret, err := _invoke(c, input.MethodName, ctx, input.In, opt)
// 	if err != nil {
// 		logrus.Warnf("querying '%s' failed: %v", input.MethodName, err)
// 	} else {
// 		input.Callback(input, ret.Interface())
// 	}
//
// 	//switch input.MethodName {
// 	//case "Login":
// 	//	token, err := c.Login(ctx, input.In.(*user.LoginReq))
// 	//	if err != nil {
// 	//		logrus.Warnf("querying '%s' failed: ", input.MethodName, err)
// 	//	} else {
// 	//		input.Callback(input, token)
// 	//	}
// 	//}
// }

// Invoke - firstResult, err := Invoke(AnyStructInterface, MethodName, Params...)
func _invoke(any interface{}, name string, args ...interface{}) (reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	methodType := method.Type()
	numIn := methodType.NumIn()
	if numIn > len(args) {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have minimum %d params. Have %d", name, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic() {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have %d params. Have %d", name, numIn, len(args))
	}
	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var inType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argValue.String())
		}
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			in[i] = argValue.Convert(inType)
		} else {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argType)
		}
	}
	return method.Call(in)[0], nil
}
