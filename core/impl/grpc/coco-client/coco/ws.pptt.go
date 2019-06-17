/*
 * Copyright © 2019 Hedzr Yeh.
 */

package coco

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/sirupsen/logrus"
	"time"
)

func (h *WsHub) testForSendMsg(tm time.Time) (err error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var msg []byte
	for client := range h.clients {
		to := client
		for xc := range h.clients {
			if xc.uid != client.uid {
				to = xc
				break
			}
		}

		req := h.prepareSendMsgReq(client.uid, to.uid, fmt.Sprintf("%v->%v|vvndfs,dsa 抵 %v", client.uid, to.uid, tm))
		msg, err = proto.Marshal(req)

		if err == nil {
			logrus.Debugf("  [%v] send-msg: %v", req.Seq, req.Body.MsgContent)
			err = client.sendPpttMessage(msg)
			if err != nil {
				logrus.Errorf("ERR [%v] send-msg: %v", req.Seq, err)
				h.unregister <- client
			}
		}
	}
	return
}

func (h *WsHub) testForGetOffline(tm time.Time) (err error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// cid, uid, endSort uint64
	var msg []byte
	for client := range h.clients {
		req := h.prepareGetOfflineMsgReq(0, client.uid, 0)
		msg, err = proto.Marshal(req)

		if err == nil {
			logrus.Debugf("  [%v] get-off-msg: %v", req.Seq, req.Body)
			err = client.sendPpttMessage(msg)
			if err != nil {
				logrus.Errorf("ERR [%v] get-off-msg: %v", req.Seq, err)
				h.unregister <- client
			}
		}
	}
	return
}

func (h *WsHub) prepareSendMsgReq(from, to uint64, msg string) *v10.SendMsgReq {
	h.seq++
	return &v10.SendMsgReq{
		ProtoOp: v10.Op_SendMsg,
		Seq:     h.seq,
		Body: &v10.SaveMessageRequest{
			GroupId:    0,
			FromUser:   from,
			ToUser:     to,
			MsgContent: fmt.Sprintf("自然而然 %v - %v", h.seq, msg),
			MsgType:    0,
		},
	}
}

func (h *WsHub) prepareGetMsgListReq(cid, uid, endSort uint64) *v10.GetMsgListReq {
	h.seq++
	return &v10.GetMsgListReq{
		ProtoOp: v10.Op_GetMsgList,
		Seq:     h.seq,
		Body: &v10.GetMessageRequest{
			ConversationSection: 0,
			ConversationId:      cid,
			ReceiveUser:         uid,
			// EndSort:             endSort,
			MaxCount: 20,
		},
	}
}

func (h *WsHub) prepareGetOfflineMsgReq(cid, uid, endSort uint64) *v10.GetOffLineConversationSetReq {
	h.seq++
	return &v10.GetOffLineConversationSetReq{
		ProtoOp: v10.Op_GetOffLineConversationSet,
		Seq:     h.seq,
		Body: &v10.GetOffLineConversationSetRequest{
			UserId: uid,
			// ConversationSection: 0,
			// ConversationId:      cid,
			// ReceiveUser:         uid,
			// EndSort:             endSort,
			// MaxCount:            20,
		},
	}
}

func (h *WsHub) pullNewMessages(po *PullMessages) (err error) {
	if po != nil && po.notify != nil && po.notify.ProtoOp == v10.Op_NotifyAck && po.notify.GetTalking() != nil {
		talking := po.notify.GetTalking()
		req := h.prepareGetMsgListReq(talking.SubscribeId, po.uid, talking.SortNum)

		h.mutex.RLock()
		defer h.mutex.RUnlock()

		for client := range h.clients {
			if client.uid == po.uid && client.did == po.did {
				var msg []byte
				msg, err = proto.Marshal(req)

				if err == nil {
					logrus.Debugf("  [%v] pull-new-msg: %v", req.Seq, req.Body)
					err = client.sendPpttMessage(msg)
					if err != nil {
						logrus.Errorf("ERR [%v] pull-new-msg: %v", req.Seq, err)
						h.unregister <- client
					}
				}
			}
		}
	}
	return
}

func (h *WsHub) pullOfflineMessages(po *PullMessages) (err error) {
	if po != nil && po.notify != nil && po.notify.ProtoOp == v10.Op_NotifyAck && po.notify.GetTalking() != nil {
		talking := po.notify.GetTalking()
		req := h.prepareGetOfflineMsgReq(talking.SubscribeId, po.uid, talking.SortNum)

		h.mutex.RLock()
		defer h.mutex.RUnlock()

		for client := range h.clients {
			if client.uid == po.uid && client.did == po.did {
				var msg []byte
				msg, err = proto.Marshal(req)

				if err == nil {
					logrus.Debugf("  [%v] pull-offline-msg: %v", req.Seq, req.Body)
					err = client.sendPpttMessage(msg)
					if err != nil {
						logrus.Errorf("ERR [%v] pull-offline-msg: %v", req.Seq, err)
						h.unregister <- client
					}
				}
			}
		}
	}
	return
}

func (c *WsClient) ppttProcess(bin []byte) (err error) {
	if bin[0] == 8 {
		var op = v10.Op(int(bin[1]))
		var ret proto.Message
		var fmtstr string

		switch op {
		case v10.Op_SendMsgAck:
			ret = &v10.SendMsgReply{}
			fmtstr = "    [ws][SendMsgAck    ]: %v"
		case v10.Op_NotifyAck:
			ret = &v10.NotifyMessage{}
			fmtstr = "    [ws][NotifyAck     ]: %v"

		case v10.Op_SubscribeAck:
			ret = &v10.SubscribeReply{}
			fmtstr = "    [ws][SubscribeAck  ]: %v"
		case v10.Op_UnsubscribeAck:
			ret = &v10.UnsubscribeReply{}
			fmtstr = "    [ws][UnsubscribeAck]: %v"

		case v10.Op_GetMsgListAck:
			ret = &v10.GetMsgListReply{}
			fmtstr = "    [ws][GetMsgListAck ]: %v"
		case v10.Op_GetOffLineConversationSetAck:
			ret = &v10.GetOffLineConversationSetReply{}
			fmtstr = "    [ws][GetOffLineConversationSetAck]: %v"

		case v10.Op_SetReminderAck:
			// ret = &v10.SubscribeReply{}
			fmtstr = "    [ws][SetReminderAck]: %v"
		case v10.Op_RemoveReminderAck:
			// ret = &v10.RemoveR{}
			fmtstr = "    [ws][RemoveReminderAck]: %v"
		}

		if ret != nil {
			err = proto.Unmarshal(bin, ret)
			if err == nil {
				logrus.Debugf(fmtstr, ret)
			}
		}

		switch op {
		case v10.Op_NotifyAck:
			wshub.pullMessages <- &PullMessages{c.uid, c.did, ret.(*v10.NotifyMessage)}
		}

		// case core.MainOp_NotifyAck:
		// 	var ret = &core.NotifyMessage{}
		// 	err = proto.Unmarshal(bin, ret)
		// 	if err == nil {
		// 		logrus.Debugf("    [ws][SendMsgAck    ]: %v", ret)
		// 		logrus.Debugf("    [ws][NotifyAck     ]: %v", ret)
		// 		logrus.Debugf("    [ws][HandshakeAck  ]: %v", ret)
		// 		logrus.Debugf("    [ws][HeartbeatAck  ]: %v", ret)
		// 		logrus.Debugf("    [ws][DisconnectAck ]: %v", ret)
		// 		logrus.Debugf("    [ws][AuthAck       ]: %v", ret)
		// 		logrus.Debugf("    [ws][GetMsgListAck ]: %v", ret)
		// 		logrus.Debugf("    [ws][GetOffLineConversationSetAck]: %v", ret)
		// 		logrus.Debugf("    [ws][SubscribeAck  ]: %v", ret)
		// 		logrus.Debugf("    [ws][UnsubscribeAck]: %v", ret)
		// 		logrus.Debugf("    [ws][SetReminderAck]: %v", ret)
		// 		logrus.Debugf("    [ws][RemoveReminderAck]: %v", ret)
		// 		logrus.Debugf("    [ws][]: %v", ret)
		// 		logrus.Debugf("    [ws][]: %v", ret)
		// 		logrus.Debugf("    [ws][]: %v", ret)
		// 	}
		// }
	}
	return
}
