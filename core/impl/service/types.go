/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/labstack/echo"
)

const (
	LoginMethod        = "Login"
	RefreshTokenMethod = "RefreshToken"
)

type (
	BuildInf struct {
		Entry               string
		svc, pkg, pbsvc     string
		FuncName            string                                        // grpc 调用函数名
		preparingInputParam func(ctx echo.Context) (proto.Message, error) // 从RESTful请求中抽出参数形成grpc调用入参
		realResultTemplate  func() (out proto.Message)                    // 生成一个空白的结果类型，被用于 pbtypes.UnmarshalAny
		wrappedResult       bool                                          // true 如果接口返回类型为 v10.Result.
		onEverythingOk      func(ret proto.Message)                       // restful 如需后处理返回消息的话，提供 onEverythinOk; grpc 调用时，等待返回消息抵达并作为返回值，故无需 onEverythingOk
	}

	FwdrFunc func(ctx context.Context, in proto.Message) (res proto.Message, err error)

	NotifiedUsers struct {
		ConversationId      uint64
		ConversationSection uint32
		MsgId               uint64
		Users               []uint64
		TS                  int64
		SortNum             uint64
	}

	Hooker interface {
		OnMessageSaveDone(ret proto.Message)
	}
)

var (
	// allow 500 messages buffered in UsersNeedNotified queue.
	UsersNeedNotified = make(chan *NotifiedUsers, 500)
)

func (s *BuildInf) Result() proto.Message {
	if s.wrappedResult {
		return nil
	}
	return s.realResultTemplate()
}

func MakeNotifiedUsersFromNotifyMessage(uid uint64, nm *v10.NotifyMessage) (ret *NotifiedUsers) {
	if nm.ProtoOp == v10.Op_NotifyAck {
		tp := nm.GetTalking()
		if tp != nil {
			ret = MakeNotifiedUsersFromTalkingPush(uid, tp)
		}
	}
	return
}

func MakeNotifiedUsersFromTalkingPush(uid uint64, tp *v10.TalkingPush) (ret *NotifiedUsers) {
	ret = new(NotifiedUsers)

	ret.ConversationId = tp.SubscribeId
	ret.ConversationSection = tp.ConversationSection
	if len(tp.MsgIds) > 0 {
		ret.MsgId = tp.MsgIds[0]
	}
	ret.SortNum = tp.SortNum

	ret.Users = make([]uint64, 1)
	ret.Users[0] = uint64(uid)

	ret.TS = tp.Ts
	return
}

// MakeNotifiedUsers prepares a NotifiedUsers structure from storage.SaveMessageResponse{}
func MakeNotifiedUsers(rr *v10.SaveMessageResponse) (ret *NotifiedUsers) {
	ret = new(NotifiedUsers)

	ret.ConversationId = rr.ConversationId
	ret.ConversationSection = rr.ConversationSection
	ret.MsgId = uint64(rr.MsgId)
	ret.SortNum = rr.SortNum

	ret.Users = make([]uint64, len(rr.ReceiveUserList))
	for i, uid := range rr.ReceiveUserList {
		ret.Users[i] = uint64(uid)
	}

	ret.TS = rr.MsgTime // api.Int64ToTimestamp(rr.MsgTime)
	return
}
