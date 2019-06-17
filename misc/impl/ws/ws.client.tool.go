/*
 * Copyright © 2019 Hedzr Yeh.
 */

package ws

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/vxconf/gwk"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func (c *WsClient) PostBinaryMsg(message []byte) {
	if !c.exited {
		c.ppttSend <- message
	}
}

func (c *WsClient) _postBinMsg(message []byte) {
	if !c.exited {
		c.ppttSend <- message
	}
}

func (c *WsClient) _postTxtMsg(message []byte) {
	if !c.exited {
		c.textSend <- message
	}
}

func (c *WsClient) _writeBack(msg string) {
	if !c.exited {
		c.textSend <- []byte(msg)
		// _ = from.conn.WriteMessage(1, []byte(msg))
	}
}

func (c *WsClient) _writeBackBytes(msg []byte) {
	if !c.exited {
		c.textSend <- msg
		// _ = from.conn.WriteMessage(1, []byte(msg))
	}
}

func (c *WsClient) onTxtMsg(msg string) (handled bool) {
	if strings.EqualFold(msg, "ping") {
		// _ = c.conn.WriteMessage(1, []byte("pong"))
		c._writeBack("pong")
		handled = true
		return
	}

	if strings.HasPrefix(msg, "shutdown:shutdown:shutdown") {
		logrus.Infof("SHUTDOWN from websocket text command.")
		daemon.StopSelf()
		handled = true
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			logrus.Errorln("Websocket.onTxtMsg() error:", err)
		}
	}()

	// 测试用的触发器

	prefix := msg[0:2]
	if fn, ok := txtMsgFuncMaps[prefix]; ok && msg[2:3] == ":" {
		handled = fn(c, msg[3:])
		return
	}

	// if strings.HasPrefix(msg, "ux:") {
	// 	// simulate a grpc request
	// 	// 连接到 vx-core 服务，发起一个 Login 请求
	// 	coco.ClientSend(msg[3:], func(cc *coco.GrpcClient, uit *v10.UserInfoToken) {
	// 		c._writeBack(uit.String())
	// 		// _ = c.conn.WriteMessage(1, []byte(ret))
	// 		// cc.RequestClose()
	// 	})
	// 	handled = true
	// 	return
	// }
	//
	// if strings.HasPrefix(msg, "ud:") {
	// 	c.simulateCoreLoginDirectly()
	// 	handled = true
	// 	return
	// }
	//
	// //
	// if strings.HasPrefix(msg, "gc:") {
	// 	c.simulateGetContact(msg[3:])
	// 	handled = true
	// 	return
	// }
	//
	// if strings.HasPrefix(msg, "lc:") {
	// 	c.simulateListContacts(msg[3:])
	// 	handled = true
	// 	return
	// }
	//
	// if strings.HasPrefix(msg, "sm:") {
	// 	// simulate a grpc request
	// 	// 连接到 vx-core 服务，发起一个 SendMsg 请求
	// 	c.simulateSendMsg(msg[3:])
	//
	// 	handled = true
	// 	return
	// }
	//
	// if strings.HasPrefix(msg, "ts:") { // timestamp msg for testing, and make a grpc request
	// 	// test for SvrRecordResolverAll()
	// 	var r = gwk.ThisConfig.Registrar
	// 	if r.IsOpen() {
	// 		var records = r.SvrRecordResolverAll(api.GrpcCore, "grpc")
	// 		for ix, rec := range records {
	// 			logrus.Infof("%3d. id:%v, ip:%v, port:%v, what:%v", ix, rec.ID, rec.IP, rec.Port, rec.What)
	// 		}
	// 	} else {
	// 		logrus.Warn("store is NOT open.")
	// 	}
	//
	// 	// test for scheduler.Invoke()
	// 	scheduler.Invoke(api.GrpcAuth, api.GrpcAuthPackageName, "UserAction", "/inx.im.user.UserAction/Login", &coco.DemoLoginReq, func(e error, input *scheduler.Input, out interface{}) {
	// 		if r, ok := out.(*v10.UserInfoToken); ok {
	// 			logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
	// 		} else {
	// 			logrus.Warnf(">> Input: %v\nhas error??? output: %v", input, out)
	// 		}
	// 	})
	// 	handled = true
	// 	return
	// }

	return
}

type handler func(c *WsClient, msg string) (handled bool)

var txtMsgFuncMaps = map[string]handler{
	"au": func(c *WsClient, msg string) (handled bool) { c.simulateAddUser(msg); handled = true; return },
	"ao": func(c *WsClient, msg string) (handled bool) {
		c.simulateAddOrgAndMqAndEventBus(msg)
		handled = true
		return
	},
	"ro": func(c *WsClient, msg string) (handled bool) {
		c.simulateRemoveOrgAndMqAndEventBus(msg)
		handled = true
		return
	},
	"gm": func(c *WsClient, msg string) (handled bool) { c.simulateGetMsg(msg); handled = true; return },
	// "ux": ux,
	// "ud": func(c *WsClient, msg string) (handled bool) { c.simulateCoreLoginDirectly(); handled = true; return },
	// "gc": func(c *WsClient, msg string) (handled bool) { c.simulateGetContact(msg); handled = true; return },
	// "lc": func(c *WsClient, msg string) (handled bool) { c.simulateListContacts(msg); handled = true; return },
	"sm": func(c *WsClient, msg string) (handled bool) { c.simulateSendMsg(msg); handled = true; return },
	"sx": func(c *WsClient, msg string) (handled bool) { c.simulateSendMsg_notopic(msg); handled = true; return },
	"rx": func(c *WsClient, msg string) (handled bool) { c.simulateSendMsg_reply(msg); handled = true; return },
	// "ts": ts,
}

func prepareSendMsgReq(seq *uint32, msg string) *v10.SendMsgReq {
	(*seq)++
	return &v10.SendMsgReq{
		ProtoOp: v10.Op_SendMsg,
		Seq:     *seq,
		Body: &v10.SaveMessageRequest{
			GroupId:    0,
			FromUser:   1,
			ToUser:     2,
			MsgContent: fmt.Sprintf("自然而然 %v - %v", seq, msg),
			MsgType:    0,
		},
	}
}

func prepareSendMsgReqV12(seq *uint32, msg string) *v10.SendMsgReqV12 {
	(*seq)++
	var topicId uint64 = 1
	return &v10.SendMsgReqV12{
		ProtoOp: v10.Op_MsgsAll, // v10.Op_SendMsg,
		Seq:     *seq,
		// TopicId: topicId,
		Msg: &v10.Msg{
			// 0: 文字消息/markdown; 31/32/33: voice/audio/video; 5: rich 图文消息; 7: html; 8: wiki; ...
			// 前端提供最便利的方法，使能用户markdown输入(或默认支持); IM 优先支持 markdown
			Type: 0,
			// 标准的消息内容，最大物理存储尺寸被限制在4000bytes，约合1300汉字。
			// 但实际的限制有系统配置表进行限定，通常为 256-1200 汉字之间。标准的默认设置值为 512 汉字。
			// 实际的设置应该在前端被体现出来，在输入时予以足够的提示。
			Content: fmt.Sprintf("自然而然 %v - %v", *seq, msg),
			// 如果真的要写长文章，可以考虑使用 detail 字段，此字段的持久化存储超过 2^16-1 bytes。
			// 具体的表现形式由前端决定。
			// 超长文字消息的发送许可权，由用户的套餐权利所决定。
			Detail: "",
			// 发消息人
			// 标准用户，内置 fxbot，用户自定义的 bot（通过安装一个app），第三方组织机构的默认 bot
			From: 3,
			// pin 住的消息，在显示时具有特别的背景色、字体样式等，以便能够在视觉上被显著突出。
			// pin 的原意是置顶，但在 topic 窗口中置顶一系列消息可能并不是特别有意义。
			// 如何呈现在用户面前由前端最终决定，只要用户能便利地找到pinned的消息即可。
			Pinned: false,
			// ? 待定 ? 消息在持久化时使用 id+version 结构，以便维持消息的修改历史
			// Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
			// 如果指明了收件人的话，to>0；否则消息属于整个会话
			To: []uint64{},
			// 如果有特别需要提及到的对象，则传送他们的 memberId
			// 对于 @mentioned 的成员，在这里指明他们，而在 msg.content 中包含了他们的占位符 @1,@2,@3,....
			ReactMembers: []uint64{},
			// 如果是发件人正在发起新会话的话，topicId有效；否则后端会建立新会话并返回会话id
			TopicId: topicId,
			// 使用 lft，rgt 算法和相应的存储结构是为了提升递归嵌套msgs的性能。
			// Lft int32 `protobuf:"varint,6,opt,name=lft,proto3" json:"lft,omitempty"`
			// Rgt int32 `protobuf:"varint,7,opt,name=rgt,proto3" json:"rgt,omitempty"`
			// nanosec.
			// CreatedAt int64 `protobuf:"varint,21,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
			// nanosec. if the msg was modified
			// UpdatedAt int64 `protobuf:"varint,22,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
			// nanosec. if the msh was revoked
			// DeletedAt            int64    `protobuf:"varint,23,opt,name=deletedAt,proto3" json:"deletedAt,omitempty"`
		},
		// HilightMembers: []uint64{},
		ParentMsgId: 0,
	}
}

func prepareSendMsgReqV12_notopic(seq *uint32, msg string) *v10.SendMsgReqV12 {
	(*seq)++
	var topicId uint64 = 0
	return &v10.SendMsgReqV12{
		ProtoOp: v10.Op_MsgsAll, // v10.Op_SendMsg,
		Seq:     *seq,
		// TopicId: topicId,
		Msg: &v10.Msg{
			Type:         0,
			Content:      fmt.Sprintf("DM - 并非自然而然 %v @1 - %v", *seq, msg),
			TopicId:      topicId,
			From:         3,
			To:           []uint64{18},
			ReactMembers: []uint64{},
			Pinned:       false,
		},
		ParentMsgId: 0,
	}
}

var lastTopicId, lastMsgId, lastParentMsgId uint64

func (c *WsClient) simulateSendMsg(msg string) {
	var in proto.Message
	if c.seq%2 == 0 {
		in = prepareSendMsgReqV12(&c.seq, msg)
	} else {
		in = prepareSendMsgReqV12_notopic(&c.seq, msg)
	}
	in = prepareSendMsgReqV12_notopic(&c.seq, msg)

	scheduler.Invoke(api.GrpcStorage, api.GrpcStoragePackageName, api.ImMsgActionName, "SendMsg",
		in, &v10.SendMsgReplyV12{}, func(e error, input *scheduler.Input, out proto.Message) {
			if r, ok := out.(*v10.SendMsgReplyV12); ok {
				logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
				c._writeBack(out.String())
				lastMsgId = r.MsgId
			} else if e != nil {
				logrus.Errorf("    invoke failed, err: %v", e)
				c._writeBack(e.Error())
			} else {
				logrus.Warnf(">> Input: %v\n   has error??? output: %v", input, out)
				c._writeBack(out.String())
			}
		})
	// if res, err := grpc.Instance.SendMsg(context.Background(), prepareSendMsgReq(&c.seq, msg)); err == nil {
	// 	c._writeBack(res.Body.String())
	// } else {
	// 	c._writeBack(fmt.Sprintf("ERR: %v", err))
	// }
	// c.seq++
}

func (c *WsClient) simulateSendMsg_notopic(msg string) {
	var in proto.Message
	in = prepareSendMsgReqV12_notopic(&c.seq, msg)
	scheduler.Invoke(api.GrpcStorage, api.GrpcStoragePackageName, api.ImMsgActionName, "SendMsg",
		in, &v10.SendMsgReplyV12{}, func(e error, input *scheduler.Input, out proto.Message) {
			if r, ok := out.(*v10.SendMsgReplyV12); ok {
				logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
				c._writeBack(out.String())
				lastMsgId = r.MsgId
				lastTopicId = r.TopicId
			} else if e != nil {
				logrus.Errorf("    invoke failed, err: %v", e)
				c._writeBack(e.Error())
			} else {
				logrus.Warnf(">> Input: %v\n   has error??? output: %v", input, out)
				c._writeBack(out.String())
			}
		})
}

func (c *WsClient) simulateSendMsg_reply(msg string) {
	var in *v10.SendMsgReqV12
	in = prepareSendMsgReqV12_notopic(&c.seq, msg)
	in.Msg.TopicId = lastTopicId
	in.ParentMsgId = lastParentMsgId
	if in.ParentMsgId == 0 {
		lastParentMsgId = lastMsgId
		in.ParentMsgId = lastMsgId
	}
	scheduler.Invoke(api.GrpcStorage, api.GrpcStoragePackageName, api.ImMsgActionName, "SendMsg",
		in, &v10.SendMsgReplyV12{}, func(e error, input *scheduler.Input, out proto.Message) {
			if r, ok := out.(*v10.SendMsgReplyV12); ok {
				logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
				c._writeBack(out.String())
				lastMsgId = r.MsgId
			} else if e != nil {
				logrus.Errorf("    invoke failed, err: %v", e)
				c._writeBack(e.Error())
			} else {
				logrus.Warnf(">> Input: %v\n   has error??? output: %v", input, out)
				c._writeBack(out.String())
			}
		})
}

//

func prepareGetMsgReq(seq *uint32, msg string) *v10.GetMsgReqV12 {
	(*seq)++
	msgId, _ := strconv.ParseUint(msg, 10, 64)
	return &v10.GetMsgReqV12{
		ProtoOp: v10.Op_MsgsAll, Seq: *seq,
		UserId: 3, TopicId: 1, MsgId: msgId,
		Newer: false, SortByAsc: false, AutoAck: true,
	}
}

func (c *WsClient) simulateGetMsg(msg string) {
	scheduler.Invoke(api.GrpcStorage, api.GrpcStoragePackageName, api.ImMsgActionName, "GetMsg",
		prepareGetMsgReq(&c.seq, msg), &v10.GetMsgReplyV12{}, func(e error, input *scheduler.Input, out proto.Message) {
			if r, ok := out.(*v10.GetMsgReplyV12); ok {
				logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
				c._writeBack(out.String())
			} else if e != nil {
				logrus.Errorf("    invoke failed, err: %v", e)
				c._writeBack(e.Error())
			} else {
				logrus.Warnf(">> Input: %v\n   has error??? output: %v", input, out)
				c._writeBack(out.String())
			}
		})
}

//

func prepareDemoAddUserReq(seq *uint32) *v10.UserAllReq {
	(*seq)++
	return &v10.UserAllReq{
		ProtoOp: v10.Op_UserAll,
		Seq:     *seq,
		Oneof: &v10.UserAllReq_Aur{
			Aur: &v10.AddUserReq{
				User: &v10.UserInfo{Nickname: "random", Pass: "123456"}},
		},
	}
}

func (c *WsClient) simulateAddUser(msg string) {
	// 在 vx-misc 中临时实现的 AddUser 等一系列 vx-auth 的接口
	scheduler.Invoke(api.GrpcMisc, api.GrpcMiscPackageName, api.UserActionName, "AddUser",
		prepareDemoAddUserReq(&c.seq), &v10.UserAllReply{}, func(e error, input *scheduler.Input, out proto.Message) {
			if r, ok := out.(*v10.UserAllReply); ok {
				logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
				c._writeBack(out.String())
			} else if e != nil {
				logrus.Errorf("    invoke failed, err: %v", e)
				c._writeBack(e.Error())
			} else {
				logrus.Warnf(">> Input: %v\n   has error??? output: %v", input, out)
				c._writeBack(out.String())
			}
		})

	// cid, err := strconv.ParseInt(msg, 10, 64)
	// if err != nil {
	// 	logrus.Warnf("Err: %v", err)
	// }
	//
	// if res, err := grpc.Instance.ListContacts(context.Background(), prepareListContactsReq(c.seq, cid)); err == nil {
	// 	c._writeBack(res.String())
	// } else {
	// 	c._writeBack(fmt.Sprintf("ERR: %v", err))
	// }
	// c.seq++
}

func prepareDemoAddOrgReq(seq *uint32, msg string) *v10.OrgAllReq {
	(*seq)++
	return &v10.OrgAllReq{
		ProtoOp: v10.Op_OrgsAll,
		Seq:     *seq,
		Oneof: &v10.OrgAllReq_Aor{
			Aor: &v10.AddOrgReq{
				Org: &v10.Organization{Name: msg}},
		},
	}
}

func prepareDemoAddOrgReqLite(seq *uint32, msg string) *v10.AddOrgReq {
	(*seq)++
	return &v10.AddOrgReq{
		ProtoOp: v10.Op_OrgsAll,
		Seq:     *seq,
		Org:     &v10.Organization{Name: msg},
	}
}

func (c *WsClient) simulateAddOrgAndMqAndEventBus(msg string) {
	scheduler.Invoke(api.GrpcMisc, api.GrpcMiscPackageName, api.ImOrgActionName, "Add",
		prepareDemoAddOrgReqLite(&c.seq, msg), &v10.AddOrgReply{}, func(e error, input *scheduler.Input, out proto.Message) {
			if r, ok := out.(*v10.AddOrgReply); ok {
				lastOrgId = r.Id
				logrus.Debugf(">> Input: %v\n<< Output: %v\n<<         id=%v", input, r, lastOrgId)
				c._writeBack(out.String())
			} else if e != nil {
				logrus.Errorf("    invoke failed, err: %v", e)
				c._writeBack(e.Error())
			} else {
				logrus.Warnf(">> Input: %v\n   has error??? output: %v", input, out)
				c._writeBack(out.String())
			}
		})
}

var lastOrgId uint64

func prepareDemoRemoveOrgReq(seq *uint32, msg string) *v10.OrgAllReq {
	(*seq)++
	return &v10.OrgAllReq{
		ProtoOp: v10.Op_OrgsAll,
		Seq:     *seq,
		Oneof: &v10.OrgAllReq_Ror{
			Ror: &v10.RemoveOrgReq{
				OrgId: lastOrgId},
		},
	}
}

func prepareDemoRemoveOrgReqLite(seq *uint32, msg string) *v10.RemoveOrgReq {
	(*seq)++
	return &v10.RemoveOrgReq{
		ProtoOp: v10.Op_OrgsAll,
		Seq:     *seq,
		OrgId:   lastOrgId,
	}
}

func (c *WsClient) simulateRemoveOrgAndMqAndEventBus(msg string) {
	scheduler.Invoke(api.GrpcMisc, api.GrpcMiscPackageName, api.ImOrgActionName, "Remove",
		prepareDemoRemoveOrgReqLite(&c.seq, msg), &v10.RemoveOrgReply{}, func(e error, input *scheduler.Input, out proto.Message) {
			if r, ok := out.(*v10.RemoveOrgReply); ok {
				logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
				c._writeBack(out.String())
			} else if e != nil {
				logrus.Errorf("    invoke failed, err: %v", e)
				c._writeBack(e.Error())
			} else {
				logrus.Warnf(">> Input: %v\n   has error??? output: %v", input, out)
				c._writeBack(out.String())
			}
		})
}

func ux(c *WsClient, msg string) (handled bool) {
	// // simulate a grpc request
	// // 连接到 vx-core 服务，发起一个 Login 请求
	// coco.ClientSend(msg[3:], func(cc *coco.GrpcClient, uit *v10.UserInfoToken) {
	// 	c._writeBack(uit.String())
	// 	// _ = c.conn.WriteMessage(1, []byte(ret))
	// 	// cc.RequestClose()
	// })
	// handled = true
	return
}

func ts(c *WsClient, msg string) (handled bool) {
	// test for SvrRecordResolverAll()
	var r = gwk.ThisConfig.Registrar
	if r.IsOpen() {
		var records = r.SvrRecordResolverAll(api.GrpcCore, "grpc")
		for ix, rec := range records {
			logrus.Infof("%3d. id:%v, ip:%v, port:%v, what:%v", ix, rec.ID, rec.IP, rec.Port, rec.What)
		}
	} else {
		logrus.Warn("store is NOT open.")
	}

	// // test for scheduler.Invoke()
	// scheduler.Invoke(api.GrpcAuth, api.GrpcAuthPackageName, "UserAction", "/inx.im.user.UserAction/Login", &coco.DemoLoginReq, func(e error, input *scheduler.Input, out interface{}) {
	// 	if r, ok := out.(*v10.UserInfoToken); ok {
	// 		logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
	// 	} else {
	// 		logrus.Warnf(">> Input: %v\nhas error??? output: %v", input, out)
	// 	}
	// })
	handled = true
	return
}

func (c *WsClient) simulateCoreLoginDirectly() {
	// // simulate a grpc request
	// // 直接内部调用 vx-core 服务，发起一个 Login 请求
	// if res, err := grpc.Instance.Login(context.Background(), &v10.AuthReq{
	// 	Oneof: &v10.AuthReq_Req{Req: &coco.DemoLoginReq,},
	// }); err == nil {
	// 	// b, _ := json.Marshal(res)
	// 	c._writeBack(res.String())
	// } else {
	// 	c._writeBack(fmt.Sprintf("ERR: %v", err))
	// }
}

func prepareListContactsReq(seq *uint32, cid int64) *v10.ListContactsReq {
	(*seq)++
	return &v10.ListContactsReq{
		ProtoOp:  v10.Op_ListContacts,
		Seq:      *seq,
		UidOwner: cid,
	}
}

func (c *WsClient) simulateListContacts(msg string) {
	// cid, err := strconv.ParseInt(msg, 10, 64)
	// if err != nil {
	// 	logrus.Warnf("Err: %v", err)
	// }
	//
	// if res, err := grpc.Instance.ListContacts(context.Background(), prepareListContactsReq(&c.seq, cid)); err == nil {
	// 	c._writeBack(res.String())
	// } else {
	// 	c._writeBack(fmt.Sprintf("ERR: %v", err))
	// }
	// c.seq++
}

func prepareGetContactReq(seq *uint32, cid int64) *v10.GetContactReq {
	(*seq)++
	return &v10.GetContactReq{
		ProtoOp:   v10.Op_GetContact,
		Seq:       *seq,
		UidOwner:  1,
		UidFriend: cid,
	}
}

func (c *WsClient) simulateGetContact(msg string) {
	// cid, err := strconv.ParseInt(msg, 10, 64)
	// if err != nil {
	// 	logrus.Warnf("Err: %v", err)
	// }
	//
	// if res, err := grpc.Instance.GetContact(context.Background(), prepareGetContactReq(&c.seq, cid)); err == nil {
	// 	c._writeBack(res.String())
	// } else {
	// 	c._writeBack(fmt.Sprintf("ERR: %v", err))
	// }
	// c.seq++
}

func (c *WsClient) writeB(type_ int, data []byte) (err error) {
	if err = c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		logrus.Warnf("error occurs at ws SetWriteDeadline/writeWait: %v", err)
		return
	}

	if err = c.conn.WriteMessage(type_, data); err != nil {
		logrus.Warnf("error occurs at ws WriteMessage/%d: %v", type_, err)
	}
	return
}
