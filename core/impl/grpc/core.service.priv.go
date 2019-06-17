/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-lite/core/impl/service"
	"github.com/sirupsen/logrus"
)

type (
	ImCorePrivateService struct {
		// fwdrs map[string]service.FwdrFunc
	}
)

var PrivateInstance *ImCorePrivateService

func NewImCorePrivateService() *ImCorePrivateService {
	s := &ImCorePrivateService{}
	s.init()
	return s
}

func (s *ImCorePrivateService) Shutdown() {

}

func (s *ImCorePrivateService) init() *ImCorePrivateService {
	return s
}

func (s *ImCorePrivateService) ExchangeNotifyingUsers(ctx context.Context, req *v10.ExchangeNotifyingUsersReq) (res *v10.ExchangeNotifyingUsersReply, err error) {
	if len(req.UserId) > 0 {
		for _, uid := range req.UserId {
			logrus.Debugf("ExchangeNotifyingUsers loop for %v", uid)
			service.UsersNeedNotified <- service.MakeNotifiedUsersFromNotifyMessage(uid, req.Mesg)
		}
	}
	return
}
