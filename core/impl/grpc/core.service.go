/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/pool"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-lite/core/impl/service"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	CoreProtocolVersion    = "1.1"
	CoreProtocolVersionInt = 0x010100
	GatewayVersionInt      = 0x010000

	userActionName    = "UserAction"
	friendActionName  = "UserAction"
	messageActionName = "MessageGreeter"
	groupActionName   = "GroupGreeter"
	loginMethod       = "Login"
)

type (
	// ImCoreService 是 vx-core 的主要 GRPC 服务
	ImCoreService struct {
		fwdrs         map[string]service.FwdrFunc // 一组转发器，restful收到对vx-core的请求时，通过此转发器调用后端接口并返回结果
		pooledEntries map[v10.Op]*CmdInfo         // 一组转发器，ws收到对vx-core的请求时通过此pool执行grpc接口调用并返回结果
		pool          *pool.PoolWithFunc
		// biList map[string]*service.BuildInf
	}
)

var (
	Instance *ImCoreService
)

func NewImCoreService() *ImCoreService {
	s := &ImCoreService{}
	s.init()
	return s
}

func (s *ImCoreService) Shutdown() {
	if s.pool != nil {
		if err := s.pool.Release(); err != nil {
			logrus.Warnf("release pool failed: %v", err)
		}
		s.pool = nil
	}
}

func (s *ImCoreService) init() *ImCoreService {

	// 转发 restful 请求到 backends
	s.fwdrs = make(map[string]service.FwdrFunc)
	service.Init(func(bi *service.BuildInf) {
		s.fwdrs[bi.FuncName] = service.BuildFwdr(bi)
		// s.biList[bi.FuncName] = bi
	})

	// 通过 v10.Op 执行 ImCoreService 的接口
	// 被用于收到来自于 websocket 的请求之后完成请求的执行
	s.pooledEntries = make(map[v10.Op]*CmdInfo)
	for _, ci := range s.aList() {
		s.pooledEntries[ci.Op] = ci
	}

	var (
		err            error
		size           = vxconf.GetIntR("server.vx-core.pool.size", 100)
		expireDuration = vxconf.GetDurationR("server.vx-core.pool.expire", time.Hour)
	)
	s.pool, err = pool.NewTimingPoolWithFunc(size, expireDuration, s.poolWorker)
	if err != nil {
		logrus.Errorf("CAN'T initialize the worker pool. err: %v", err)
	}
	return s
}

func (s *ImCoreService) PooledInvoke(op v10.Op, from WsClientSkel, body []byte) (err error) {
	if ci, ok := s.pooledEntries[op]; ok {
		err = s.pooledInvokeInternal(&Request{ci, from, body})
	} else {
		err = fmt.Errorf("CAN'T found MainOp: %v", op)
	}
	return
}

// pooledInvokeInternal 将来自 ws client 的 core api 请求解释到 ImCoreService 调用
func (s *ImCoreService) pooledInvokeInternal(payload *Request) (err error) {
	err = s.pool.Invoke(payload)
	return
}

func (s *ImCoreService) poolWorker(payload interface{}) {
	if req, ok := payload.(*Request); !ok {
		logrus.Warnf("unexpected payload for poolWorker: invalid Result{} structure. payload = %v", payload)
		return

	} else {
		in := req.CmdInfo.InTemplate()
		err := proto.Unmarshal(req.InParam, in)
		if err != nil {
			logrus.Errorf("CAN'T unmarshal input data package: %v", err)
			return
		}

		// invoke ImCoreService directly...
		var ctx = context.Background()
		res, err := req.CmdInfo.Target(ctx, in)
		if err != nil {
			logrus.Errorf("grpc.Instance.SendMsg wrong: %v", err)
			return
		}

		var b []byte
		b, err = proto.Marshal(res)
		if err != nil {
			logrus.Errorf("CAN'T marshal output data package: %v", err)
			return
		}

		// The ack: write back to ws client
		req.From.PostBinaryMsg(b)
	}
}

//
//
//
//
//
