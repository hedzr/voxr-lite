/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-api/util"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/sirupsen/logrus"
)

type ImBase struct{}

func (s *ImBase) createBot(reqProtoOp v10.Op, reqSeq uint32, modelOut interface{}, pbOut proto.Message, modelId uint64, modelName string, done func(botId uint64)) (err error) {
	err = PooledInvoke(reqProtoOp, modelOut, pbOut, func(e error, request *Request) {
		// sth.
		// 	// TOD O add bot for new organization
		// 	// TOD O link fxbot into the members of the new organization

		// logrus.Debugf("org add: creating bot: req=%v", req)
		logrus.Debugf("org add: creating bot: cmdInfo=%v", request.CmdInfo)

		gen := util.DefaultGenerator
		bot := &v10.UserInfo{
			Pass:      gen.NewPassword(),
			LoginName: fmt.Sprintf("org.%v.%v", modelName, gen.NewRandomString(8)),
			Nickname:  fmt.Sprintf("%v.bot", modelName),
			Mobile:    gen.NewRandomMobile(),
			Email:     gen.NewRandomEmail(),
			Uid:       gen.NewRandomString(24),
			Type:      models.UserTypeBotMask + models.UserTypeOrgBotMask + models.UserTypeSpecialMask, // 添加 org 专用 bot
			Lang:      models.DefaultLanguage,
			Tz:        models.UTCTimezone,
			Status:    models.StatusStandardValid,
		}

		in := &v10.UserAllReq{ProtoOp: v10.Op_UserAll, Seq: reqSeq, Oneof: &v10.UserAllReq_Aur{Aur: &v10.AddUserReq{User: bot}}}
		scheduler.Invoke(api.GrpcUser, api.GrpcAuthPackageName, api.UserActionName, "AddUser",
			in, &v10.UserAllReply{}, func(e error, input *scheduler.Input, out proto.Message) {
				if uar, ok := out.(*v10.UserAllReply); ok {
					done(uar.GetAur().Id)
				} else if e != nil {
					logrus.Errorf("   invoke failed, err: %v", e)
				}
			})
	})
	return
}
