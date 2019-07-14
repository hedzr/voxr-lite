/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/pool"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	Pooled struct {
		pool          *pool.PoolWithFunc
		pooledEntries map[v10.Op]*CmdInfo
	}

	CmdInfo struct {
		DoFor func(*Request) error
	}

	Request struct {
		// op    v10.Op
		CmdInfo  *CmdInfo
		Model    interface{}
		Res      proto.Message
		Callback func(error, *Request)
	}
)

var afterOkPooled *Pooled

func init() {
	afterOkPooled = &Pooled{}
	afterOkPooled.Init()
}

func (s *Pooled) Init() {
	s.initPool()
	s.initMappings()
	logrus.Debugf("pool inside service initialized.")
}

func (s *Pooled) initMappings() {
	s.pooledEntries = map[v10.Op]*CmdInfo{
		v10.Op_OrgsAll:   {},
		v10.Op_TopicsAll: {},
	}
}

func (s *Pooled) initPool() {

	var (
		err            error
		size           = vxconf.GetIntR("server.grpc.settings.pools.after.size", 32)
		expireDuration = vxconf.GetDurationR("server.grpc.settings.pools.after.expire", time.Hour)
	)

	if size == 0 {
		size = 10
	}
	expireDuration = 1 * time.Hour

	s.pool, err = pool.NewTimingPoolWithFunc(size, expireDuration, s.poolWorker)
	if err != nil {
		logrus.Errorf("CAN'T initialize the worker pool. err: %v", err)
	}

}

func PooledInvoke(op v10.Op, modelOut interface{}, pbOut proto.Message, cb func(error, *Request)) (err error) {
	err = afterOkPooled.PooledInvoke(op, modelOut, pbOut, cb)
	return
}

func (s *Pooled) PooledInvoke(op v10.Op, model interface{}, res proto.Message, cb func(error, *Request)) (err error) {
	if ci, ok := s.pooledEntries[op]; ok {
		err = s.pooledInvokeInternal(&Request{ci, model, res, cb})
	} else {
		err = fmt.Errorf("CAN'T found MainOp: %v", op)
	}
	return
}

// pooledInvokeInternal 将来自 ws client 的 core api 请求解释到 ImCoreService 调用
func (s *Pooled) pooledInvokeInternal(payload *Request) (err error) {
	err = s.pool.Invoke(payload)
	return
}

func (s *Pooled) poolWorker(payload interface{}) {
	if req, ok := payload.(*Request); !ok {
		logrus.Warnf("unexpected payload for poolWorker: invalid Result{} structure. payload = %v", payload)
		return

	} else {
		if req.CmdInfo.DoFor != nil {
			req.Callback(req.CmdInfo.DoFor(req), req)
		} else {
			logrus.Debugf("poolWorker done. ignore callback: req.CmdInfo.DoFor == nil.")

			if req.Callback != nil {
				req.Callback(nil, req)
			}
		}

		// in := req.CmdInfo.InTemplate()
		// err := proto.Unmarshal(req.InParam, in)
		// if err != nil {
		// 	logrus.Errorf("CAN'T unmarshal input data package: %v", err)
		// 	return
		// }
		//
		// // invoke ImCoreService directly...
		// var ctx = context.Background()
		// res, err := req.CmdInfo.Target(ctx, in)
		// if err != nil {
		// 	logrus.Errorf("grpc.Instance.SendMsg wrong: %v", err)
		// 	return
		// }
		//
		// var b []byte
		// b, err = proto.Marshal(res)
		// if err != nil {
		// 	logrus.Errorf("CAN'T marshal output data package: %v", err)
		// 	return
		// }
		//
		// // The ack: write back to ws client
		// req.From.PostBinaryMsg(b)
	}
}
