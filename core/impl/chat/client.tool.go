/*
 * Copyright © 2019 Hedzr Yeh.
 */

package chat

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/vxconf/gwk"
	"github.com/hedzr/voxr-lite/core/impl/grpc"
	"github.com/hedzr/voxr-lite/core/impl/grpc/coco-client/coco"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/sirupsen/logrus"
	"io"
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
	"ux": ux,
	"ud": func(c *WsClient, msg string) (handled bool) { c.simulateCoreLoginDirectly(); handled = true; return },
	"gc": func(c *WsClient, msg string) (handled bool) { c.simulateGetContact(msg); handled = true; return },
	"lc": func(c *WsClient, msg string) (handled bool) { c.simulateListContacts(msg); handled = true; return },
	"sm": func(c *WsClient, msg string) (handled bool) { c.simulateSendMsg(msg); handled = true; return },
	"ts": ts,
}

func ux(c *WsClient, msg string) (handled bool) {
	// simulate a grpc request
	// 连接到 vx-core 服务，发起一个 Login 请求
	coco.ClientSend(msg[3:], func(cc *coco.GrpcClient, uit *v10.UserInfoToken) {
		c._writeBack(uit.String())
		// _ = c.conn.WriteMessage(1, []byte(ret))
		// cc.RequestClose()
	})
	handled = true
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

	// test for scheduler.Invoke()
	scheduler.Invoke(api.GrpcAuth, api.GrpcAuthPackageName, "UserAction", "/inx.im.user.UserAction/Login",
		&coco.DemoLoginReq, nil, func(e error, input *scheduler.Input, out proto.Message) {
			if r, ok := out.(*v10.UserInfoToken); ok {
				logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
			} else if e != nil {
				logrus.Errorf("   invoke failed, err: %v", e)
			} else {
				logrus.Warnf(">> Input: %v\nhas error??? output: %v", input, out)
			}
		})
	handled = true
	return
}

func (c *WsClient) simulateCoreLoginDirectly() {
	// simulate a grpc request
	// 直接内部调用 vx-core 服务，发起一个 Login 请求
	if ret, err := grpc.Instance.Login(context.Background(), &v10.AuthReq{
		Oneof: &v10.AuthReq_Req{Req: &coco.DemoLoginReq},
	}); err == nil {
		// b, _ := json.Marshal(res)
		c._writeBack(ret.String())
	} else {
		c._writeBack(fmt.Sprintf("ERR: %v", err))
	}
}

func prepareListContactsReq(seq uint32, cid int64) *v10.ListContactsReq {
	seq++
	return &v10.ListContactsReq{
		ProtoOp:  v10.Op_ListContacts,
		Seq:      seq,
		UidOwner: cid,
	}
}

func (c *WsClient) simulateListContacts(msg string) {
	cid, err := strconv.ParseInt(msg, 10, 64)
	if err != nil {
		logrus.Warnf("Err: %v", err)
	}

	if res, err := grpc.Instance.ListContacts(context.Background(), prepareListContactsReq(c.seq, cid)); err == nil {
		c._writeBack(res.String())
	} else {
		c._writeBack(fmt.Sprintf("ERR: %v", err))
	}
	c.seq++
}

func prepareGetContactReq(seq uint32, cid int64) *v10.GetContactReq {
	seq++
	return &v10.GetContactReq{
		ProtoOp:   v10.Op_GetContact,
		Seq:       seq,
		UidOwner:  1,
		UidFriend: cid,
	}
}

func (c *WsClient) simulateGetContact(msg string) {
	cid, err := strconv.ParseInt(msg, 10, 64)
	if err != nil {
		logrus.Warnf("Err: %v", err)
	}

	if res, err := grpc.Instance.GetContact(context.Background(), prepareGetContactReq(c.seq, cid)); err == nil {
		c._writeBack(res.String())
	} else {
		c._writeBack(fmt.Sprintf("ERR: %v", err))
	}
	c.seq++
}

func prepareSendMsgReq(seq uint32, msg string) *v10.SendMsgReq {
	seq++
	return &v10.SendMsgReq{
		ProtoOp: v10.Op_SendMsg,
		Seq:     seq,
		Body: &v10.SaveMessageRequest{
			GroupId:    0,
			FromUser:   1,
			ToUser:     2,
			MsgContent: fmt.Sprintf("自然而然 %v - %v", seq, msg),
			MsgType:    0,
		},
	}
}

func (c *WsClient) simulateSendMsg(msg string) {
	if res, err := grpc.Instance.SendMsg(context.Background(), prepareSendMsgReq(c.seq, msg)); err == nil {
		c._writeBack(res.Body.String())
	} else {
		c._writeBack(fmt.Sprintf("ERR: %v", err))
	}
	c.seq++
}

// grpc.ClientSend(msg[3:], func(ret string) {
// 	c.conn.WriteMessage(1, []byte(ret))
// })

//
// scheduler.Invoke(api.GrpcAuth, "UserAction", "Login", &loginReq, func(e error, input *s.Input, out interface{}) {
// 	if r, ok := out.(*user.UserInfoToken); ok {
// 		logrus.Debugf(">> Input: %v\n<< Output: %v", input, r)
// 	} else {
// 		logrus.Warnf(">> Input: %v\nhas error??? output: %v", input, out)
// 	}
// })

// // encode sth.
// token := &user.UserInfoToken{
// 	UserInfo: &user.UserInfo{
// 		8, "uid-1", 3, "realname", "nickname",
// 		"13801234567", "avatar://xxx", "510214xxxxxxxxxxxx",
// 		100, 1, "example@example.com", "fhsdkfhdskf",
// 		9,
// 		struct{}{}, []byte{}, 0,
// 	},
// 	Token: "fdjsfdsjfjsdkfhsd",
// }
// data, err := proto.Marshal(token)
// if err != nil {
// 	logrus.Fatalf("encode failed: ", err)
// }
//
// // decode sth.
// var target user.UserInfoToken
// err = proto.Unmarshal(data, &target)
// if err != nil {
// 	logrus.Fatalf("encode failed: ", err)
// }
// logrus.Debugf("userInfoToken = %v", target)

// func randomUserInfo() *user.UserInfo {
// 	user_login_res := user.UserInfo{Id: 1,
// 		Uid:       "op_ajkhsajk98217hjsbnp",
// 		UType:     0,
// 		UNickname: "David",
// 		UPhone:    "13323977614",
// 		UAvatar:   "https://www.ajskja.com",
// 		UIdcard:   "500129199602293301",
// 		UAge:      23,
// 		USex:      1,
// 		UEmail:    "xxxx@163.com",
// 		URealname: "佚名",
// 		UPass:     "123456"}
// 	return &user_login_res
// }

// bAD
func (c *WsClient) writeN(ok bool, type_ int, message []byte) (err error) {
	if err = c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return
	}
	if !ok {
		// The hub closed the channel.
		if err = c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			err = nil
		}
		return
	}

	var w io.WriteCloser
	w, err = c.conn.NextWriter(type_)
	if err != nil {
		return
	}
	if _, err = w.Write(message); err != nil {
		return
	}

	// Add queued chat messages to the current websocket message.
	n := len(c.textSend)
	for i := 0; i < n; i++ {
		_, err = w.Write(newline)
		if err == nil {
			_, err = w.Write(<-c.textSend)
		}
		if err != nil {
			return
		}
	}

	err = w.Close()
	return
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
