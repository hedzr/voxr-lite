/*
 * Copyright © 2019 Hedzr Yeh.
 */

package chat

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-lite/core/impl/grpc"
	"github.com/hedzr/voxr-lite/core/impl/service"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

func (h *Hub) OnMessageSaveDone(msg proto.Message) {
	//
}

// after
func (h *Hub) onNotifyUsersNewMsgIncoming(users *service.NotifiedUsers) {
	for _, uid := range users.Users {
		push := &v10.TalkingPush{
			SubscribeId:         users.ConversationId,
			ConversationSection: users.ConversationSection,
			SortNum:             users.SortNum,
			Ts:                  users.TS, MsgIds: []uint64{users.MsgId},
		}
		nm := &v10.NotifyMessage{ProtoOp: v10.Op_NotifyAck, Oneof: &v10.NotifyMessage_Talking{Talking: push}}
		h.ppttPushing <- &newMsgIncoming{uid, nm}
	}
}

// storage 不再需要明确调用 SaveOffLineConversation() , 这没有意义。
var noNeedToSaveOffline = true

func (h *Hub) ppttDoPushing(msg *newMsgIncoming) {
	// find all clients in all zone by uid
	// and push msg.nm to him.
	go func() {
		if h.findClientsByUidAndPush(msg.uid, msg.nm) {

			// offline user
			if noNeedToSaveOffline {

				in := &v10.SaveOffLineConversationRequest{
					UserId:              msg.uid,
					ConversationSection: msg.nm.GetTalking().ConversationSection,
					ConversationId:      msg.nm.GetTalking().SubscribeId,
					SortInConversation:  msg.nm.GetTalking().SortNum,
				}
				scheduler.Invoke(api.GrpcStorage, api.GrpcStoragePackageName, api.ConversationActionName,
					"SaveOffLineConversation", in, nil,
					func(e error, input *scheduler.Input, out proto.Message) {
						if e == nil {
							if r, ok := out.(*v10.Result); ok && r.Ok && len(r.Data) > 0 {
								logrus.Debugf("SaveOffLineConversation return: %v", r)

								res := new(v10.SaveOffLineConversationResponse)
								if res != nil {
									if err := ptypes.UnmarshalAny(r.Data[0], res); err != nil {
										logrus.Warnf("cannot decode to %v: %v", reflect.TypeOf(res), r)
									} else {
										logrus.Infof("offline msg saved. uid(%v) has %v msgs. sortnum: %v", msg.uid, res.MsgCount, res.SortInConversation)
									}
								}
							} else {
								logrus.Errorf("SaveOffLineConversation return failed: %v, %v", r.ErrCode, r.Msg)
							}
						} else {
							logrus.Errorf("SaveOffLineConversation return: %v", e)
							// res = nil // ErrorResult(http.StatusUnauthorized, "Please provide valid credentials")
							// err = echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
						}
						// ch <- true
					})

			} else {
				logrus.Debugf("SaveOffLineConversation ignore and...")
			}
			return

		}
		logrus.Debugf("find all clients in all zone by uid %v, and push msg.nm to him: nm=%v", msg.uid, msg.nm)

	}()
}

// find the zones entries by user id, and find out its ws-client, and push new-msg-incoming to those clients.
// if not found, return offline is true.
func (h *Hub) findClientsByUidAndPush(uid uint64, nm *v10.NotifyMessage) (offline bool) {
	var zid = vxconf.GetStringR("server.id", "") // see also: id.GenerateInstanceId()
	mapZones := redis_op.FindLivingUids(uid)

	offline = true
	for zid1, keys := range mapZones {
		if zid1 == zid {
			for _, key := range keys {
				ss := strings.Split(key, ":")
				did := ss[6]
				if c, ok := h.clientsMap[did]; ok {
					// one client found, push msg.nm to him...
					c.needPushing <- nm
					offline = false
					break
				}
			}
		} else {
			// T O D O: push from other zones
			req := &v10.ExchangeNotifyingUsersReq{Mesg: nm, UserId: []uint64{uid}}
			_, err := grpc.PrivateInstance.ExchangeNotifyingUsers(context.Background(), req)
			if err != nil {
				logrus.Warnf("ExchangeNotifyingUsers failed: %v", err)
			}
		}
	}

	return
}
