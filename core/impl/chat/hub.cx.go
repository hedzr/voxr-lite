/*
 * Copyright © 2019 Hedzr Yeh.
 */

package chat

func (h *Hub) preInitCommands() {

	// apps/clients 通过 websocket 发送 PB 消息，Hub 负责将这些消息通过 Hub.commands 分发到具体的处理逻辑
	// h.commands = make(map[core.MainOp]grpc.CmdFunc)
	// for ix, ci := range []grpc.CmdInfo{
	// 	{core.MainOp_SendMsg, grpc.Instance.SendMsgX, func() proto.Message {
	// 		return new(core.SendMsgReq)
	// 	}},
	// } {
	// 	h.commands[ci.Op] = build(ix, &ci)
	// }

	// h.commands[core.MainOp_SendMsg] = h.cxSendMsg
	// h.commands[core.MainOp_GetMsgList] = h.cxGetMsgList
	// h.commands[core.MainOp_GetOffLineConversationSet] = h.cxGetOfflineConversationSet
	// h.commands[core.MainOp_Subscribe] = h.cxSubscribe
	// h.commands[core.MainOp_Unsubscribe] = h.cxUnsubscribe

}

// func build(ix int, info *grpc.CmdInfo) grpc.CmdFunc {
// 	return func(from grpc.WsClientSkel, body []byte) {
// 		grpc.Instance.PoolAdd(func(payload interface{}) {
// 			if req, ok:=payload.(*grpc.Request); !ok {
// 				return
// 			}else{
// 				req.From.PostBinaryMsg(req.InParam)
// 			}
// 		})
// 		// grpc.Instance.Add
// 		in := info.InTemplate()
// 		err := proto.Unmarshal(body, in)
// 		if err != nil {
// 			logrus.Errorf("CAN'T unmarshal input data package: %v", err)
// 		}
//
// 		// invoke ImCoreService directly...
// 		var ctx = context.Background()
// 		res, err := info.Target(ctx, in)
// 		if err != nil {
// 			logrus.Errorf("grpc.Instance.SendMsg wrong: %v", err)
// 		}
//
// 		var b []byte
// 		b, err = proto.Marshal(res)
// 		if err != nil {
// 			logrus.Errorf("CAN'T marshal output data package: %v", err)
// 		}
//
// 		from.PostBinaryMsg(b)
// 	}
// }
//
// // Never Used
// // SendMsg was being requested by apps/clients.
// func (h *Hub) cxSendMsg(from *WsClient, body []byte) {
// 	go func() {
//
// 		in := new(v10.SendMsgReq)
// 		err := proto.Unmarshal(body, in)
// 		if err != nil {
// 			logrus.Errorf("CAN'T unmarshal input data package: %v", err)
// 		}
//
// 		// invoke ImCoreService directly...
// 		var ctx = context.Background()
// 		res, err := grpc.Instance.SendMsg(ctx, in)
// 		if err != nil {
// 			logrus.Errorf("grpc.Instance.SendMsg wrong: %v", err)
// 		}
//
// 		var b []byte
// 		b, err = proto.Marshal(res)
// 		if err != nil {
// 			logrus.Errorf("CAN'T marshal output data package: %v", err)
// 		}
// 		from._postBinMsg(b)
//
// 	}()
// }
//
// // Never Used
// func (h *Hub) cxGetMsgList(from *WsClient, body []byte) {
// 	go func() {
//
// 		in := new(v10.GetMsgListReq)
// 		err := proto.Unmarshal(body, in)
// 		if err != nil {
// 			logrus.Errorf("CAN'T unmarshal input data package: %v", err)
// 		}
//
// 		// invoke ImCoreService directly...
// 		var ctx = context.Background()
// 		res, err := grpc.Instance.GetMsgList(ctx, in)
// 		if err != nil {
// 			logrus.Errorf("grpc.Instance.SendMsg wrong: %v", err)
// 		}
//
// 		var b []byte
// 		b, err = proto.Marshal(res)
// 		if err != nil {
// 			logrus.Errorf("CAN'T marshal output data package: %v", err)
// 		}
// 		from._postBinMsg(b)
//
// 	}()
// }
//
// func (h *Hub) cxGetOfflineConversationSet(from *WsClient, body []byte) {
// 	go func() {
//
// 		in := new(core.SendMsgReq)
// 		err := proto.Unmarshal(body, in)
// 		if err != nil {
// 			logrus.Errorf("CAN'T unmarshal input data package: %v", err)
// 		}
//
// 		// invoke ImCoreService directly...
// 		var ctx = context.Background()
// 		res, err := grpc.Instance.SendMsg(ctx, in)
// 		if err != nil {
// 			logrus.Errorf("grpc.Instance.SendMsg wrong: %v", err)
// 		}
//
// 		var b []byte
// 		b, err = proto.Marshal(res)
// 		if err != nil {
// 			logrus.Errorf("CAN'T marshal output data package: %v", err)
// 		}
// 		from._postBinMsg(b)
//
// 	}()
// }
//
// func (h *Hub) cxSubscribe(from *WsClient, body []byte) {
// 	go func() {
//
// 		in := new(core.SendMsgReq)
// 		err := proto.Unmarshal(body, in)
// 		if err != nil {
// 			logrus.Errorf("CAN'T unmarshal input data package: %v", err)
// 		}
//
// 		// invoke ImCoreService directly...
// 		var ctx = context.Background()
// 		res, err := grpc.Instance.SendMsg(ctx, in)
// 		if err != nil {
// 			logrus.Errorf("grpc.Instance.SendMsg wrong: %v", err)
// 		}
//
// 		var b []byte
// 		b, err = proto.Marshal(res)
// 		if err != nil {
// 			logrus.Errorf("CAN'T marshal output data package: %v", err)
// 		}
// 		from._postBinMsg(b)
//
// 	}()
// }
//
// func (h *Hub) cxUnsubscribe(from *WsClient, body []byte) {
// 	go func() {
//
// 		in := new(core.SendMsgReq)
// 		err := proto.Unmarshal(body, in)
// 		if err != nil {
// 			logrus.Errorf("CAN'T unmarshal input data package: %v", err)
// 		}
//
// 		// invoke ImCoreService directly...
// 		var ctx = context.Background()
// 		res, err := grpc.Instance.SendMsg(ctx, in)
// 		if err != nil {
// 			logrus.Errorf("grpc.Instance.SendMsg wrong: %v", err)
// 		}
//
// 		var b []byte
// 		b, err = proto.Marshal(res)
// 		if err != nil {
// 			logrus.Errorf("CAN'T marshal output data package: %v", err)
// 		}
// 		from._postBinMsg(b)
//
// 	}()
// }
