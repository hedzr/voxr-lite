/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-lite/core/impl/service"
	"github.com/sirupsen/logrus"
)

func (s *ImCoreService) SendMsgX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.SendMsg(ctx, req.(*v10.SendMsgReq))
}

// SendMsg send the message of a client/user's device.
// rpc SendMsg (SendMsgReq) returns (SendMsgReply)
func (s *ImCoreService) SendMsg(ctx context.Context, req *v10.SendMsgReq) (res *v10.SendMsgReply, err error) {
	// if fnRPC, ok := s.fwdrs["SendMsg"]; ok {
	// 	var ret proto.Message
	// 	ret, err = fnRPC(ctx, req)
	//
	// 	if r, ok := ret.(*base.Result); ok && r.Ok && r.Count == 1 && len(r.Data) == 1 {
	// 		res = new(core.SendMsgReply)
	// 		if err = ptypes.UnmarshalAny(r.Data[0], res); err == nil {
	// 			s.afterMsgSaved(res)
	// 			return
	// 		} else {
	// 			logrus.Warnf("cannot decode to %v: %v", reflect.TypeOf(res), r)
	// 		}
	// 	}
	// 	// logrus.Debugf("fn = %v", fnRPC)
	// }
	if req != nil && req.ProtoOp == v10.Op_SendMsg && req.Seq > 0 {
		if fnRPC, ok := s.fwdrs["SaveMessage"]; ok {

			var r proto.Message
			r, err = fnRPC(ctx, req.Body) // see also: service.BuildFwdr()

			if rr, ok := r.(*v10.SaveMessageResponse); ok {
				res = &v10.SendMsgReply{
					ProtoOp: v10.Op_SendMsgAck, Seq: req.Seq, ErrorCode: v10.Err_OK,
					SubscribeId: rr.ConversationId, //
					MsgId:       rr.MsgId,          //
					Body:        rr,
				}
				s.afterMsgSaved(rr)

			} else {
				if err != nil {
					logrus.Warnf("    [core.service] invoke backend error: %v", err)
				} else {
					logrus.Warn("    [core.service] invoke backend generic error, no futher details")
				}
				res = &v10.SendMsgReply{ProtoOp: v10.Op_SendMsgAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_INVOKE}
			}

		} else {
			logrus.Warn("    [core.service] no backend or no backend api found")
			res = &v10.SendMsgReply{ProtoOp: v10.Op_SendMsgAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_API_NOTFOUND}
		}
	} else {
		res = &v10.SendMsgReply{ProtoOp: v10.Op_SendMsgAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	}
	return
}

func (s *ImCoreService) afterMsgSaved(res *v10.SaveMessageResponse) {
	// push notifications
	service.UsersNeedNotified <- service.MakeNotifiedUsers(res)
}

//
//
//

func (s *ImCoreService) GetMsgListX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.GetMsgList(ctx, req.(*v10.GetMsgListReq))
}

func (s *ImCoreService) GetMsgList(ctx context.Context, req *v10.GetMsgListReq) (res *v10.GetMsgListReply, err error) {
	// panic("implement me")
	if req != nil && req.ProtoOp == v10.Op_GetMsgList && req.Seq > 0 {
		if fnRPC, ok := s.fwdrs["GetMessage"]; ok {

			var r proto.Message
			r, err = fnRPC(ctx, req.Body) // see also: service.BuildFwdr()

			if rr, ok := r.(*v10.GetMessageResponse); ok {
				res = &v10.GetMsgListReply{
					ProtoOp: v10.Op_GetMsgListAck, Seq: req.Seq, ErrorCode: v10.Err_OK,
					SubscribeId: req.Body.ConversationId, //
					Body:        rr,                      //
				}

			} else {
				if err != nil {
					logrus.Warnf("    [core.service] invoke backend error: %v", err)
				} else {
					logrus.Warn("    [core.service] invoke backend generic error, no futher details")
				}
				res = &v10.GetMsgListReply{ProtoOp: v10.Op_GetMsgListAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_INVOKE}
			}

		} else {
			logrus.Warn("    [core.service] no backend or no backend api found")
			res = &v10.GetMsgListReply{ProtoOp: v10.Op_GetMsgListAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_API_NOTFOUND}
		}
	} else {
		res = &v10.GetMsgListReply{ProtoOp: v10.Op_GetMsgListAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	}
	return
}

func (s *ImCoreService) GetMsgHistoryX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.GetMsgHistory(ctx, req.(*v10.GetMsgHistoryReq))
}

func (s *ImCoreService) GetMsgHistory(ctx context.Context, req *v10.GetMsgHistoryReq) (res *v10.GetMsgHistoryReply, err error) {
	// panic("implement me")
	if req != nil && req.ProtoOp == v10.Op_GetMsgHistory && req.Seq > 0 {
		if fnRPC, ok := s.fwdrs["GetMessageHistory"]; ok {

			var r proto.Message
			r, err = fnRPC(ctx, req.Body) // see also: service.BuildFwdr()

			if rr, ok := r.(*v10.GetMessageHistoryResponse); ok {
				res = &v10.GetMsgHistoryReply{
					ProtoOp: v10.Op_GetMsgHistoryAck, Seq: req.Seq, ErrorCode: v10.Err_OK,
					Body: rr, //
				}

			} else {
				if err != nil {
					logrus.Warnf("    [core.service] invoke backend error: %v", err)
				} else {
					logrus.Warn("    [core.service] invoke backend generic error, no futher details")
				}
				res = &v10.GetMsgHistoryReply{ProtoOp: v10.Op_GetMsgHistoryAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_INVOKE}
			}

		} else {
			logrus.Warn("    [core.service] no backend or no backend api found")
			res = &v10.GetMsgHistoryReply{ProtoOp: v10.Op_GetMsgHistoryAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_API_NOTFOUND}
		}
	} else {
		res = &v10.GetMsgHistoryReply{ProtoOp: v10.Op_GetMsgHistoryAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	}
	return
}

func (s *ImCoreService) AckMsgX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.AckMsg(ctx, req.(*v10.AckMsgReq))
}

func (s *ImCoreService) AckMsg(ctx context.Context, req *v10.AckMsgReq) (res *v10.AckMsgReply, err error) {
	// panic("implement me")
	if req != nil && req.ProtoOp == v10.Op_AckMsg && req.Seq > 0 {
		if fnRPC, ok := s.fwdrs["AckMessage"]; ok {

			var r proto.Message
			r, err = fnRPC(ctx, req.Body) // see also: service.BuildFwdr()

			if rr, ok := r.(*v10.AckMessageResponse); ok {
				res = &v10.AckMsgReply{
					ProtoOp: v10.Op_AckMsgAck, Seq: req.Seq, ErrorCode: v10.Err_OK,
					Body: rr, //
				}

			} else {
				if err != nil {
					logrus.Warnf("    [core.service] invoke backend error: %v", err)
				} else {
					logrus.Warn("    [core.service] invoke backend generic error, no futher details")
				}
				res = &v10.AckMsgReply{ProtoOp: v10.Op_AckMsgAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_INVOKE}
			}

		} else {
			logrus.Warn("    [core.service] no backend or no backend api found")
			res = &v10.AckMsgReply{ProtoOp: v10.Op_AckMsgAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_API_NOTFOUND}
		}
	} else {
		res = &v10.AckMsgReply{ProtoOp: v10.Op_AckMsgAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	}
	return
}

func (s *ImCoreService) GetOffLineConversationSetX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.GetOffLineConversationSet(ctx, req.(*v10.GetOffLineConversationSetReq))
}

func (s *ImCoreService) GetOffLineConversationSet(ctx context.Context, req *v10.GetOffLineConversationSetReq) (res *v10.GetOffLineConversationSetReply, err error) {
	// panic("implement me")
	if req != nil && req.ProtoOp == v10.Op_GetOffLineConversationSet && req.Seq > 0 {
		if fnRPC, ok := s.fwdrs["GetOffLineConversationSet"]; ok {

			var r proto.Message
			r, err = fnRPC(ctx, req.Body) // see also: service.BuildFwdr()

			if rr, ok := r.(*v10.GetOffLineConversationSetResponse); ok {
				res = &v10.GetOffLineConversationSetReply{
					ProtoOp: v10.Op_GetOffLineConversationSetAck, Seq: req.Seq, ErrorCode: v10.Err_OK,
					SubscribeId: 0,  //
					Body:        rr, //
				}

			} else {
				if err != nil {
					logrus.Warnf("    [core.service] invoke backend error: %v", err)
				} else {
					logrus.Warn("    [core.service] invoke backend generic error, no futher details")
				}
				res = &v10.GetOffLineConversationSetReply{ProtoOp: v10.Op_GetOffLineConversationSetAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_INVOKE}
			}

		} else {
			logrus.Warn("    [core.service] no backend or no backend api found")
			res = &v10.GetOffLineConversationSetReply{ProtoOp: v10.Op_GetOffLineConversationSetAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_API_NOTFOUND}
		}
	} else {
		res = &v10.GetOffLineConversationSetReply{ProtoOp: v10.Op_GetOffLineConversationSetAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	}
	return
}

func (s *ImCoreService) GetNotReadConversationSetX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.GetNotReadConversationSet(ctx, req.(*v10.GetNotReadConversationSetReq))
}

func (s *ImCoreService) GetNotReadConversationSet(ctx context.Context, req *v10.GetNotReadConversationSetReq) (res *v10.GetNotReadConversationSetReply, err error) {
	// panic("implement me")
	if req != nil && req.ProtoOp == v10.Op_GetNotReadConversationSet && req.Seq > 0 {
		if fnRPC, ok := s.fwdrs["GetNotReadConversationSet"]; ok {

			var r proto.Message
			r, err = fnRPC(ctx, req.Body) // see also: service.BuildFwdr()

			if rr, ok := r.(*v10.GetNotReadConversationSetResponse); ok {
				res = &v10.GetNotReadConversationSetReply{
					ProtoOp: v10.Op_GetNotReadConversationSetAck, Seq: req.Seq, ErrorCode: v10.Err_OK,
					Body: rr, //
				}

			} else {
				if err != nil {
					logrus.Warnf("    [core.service] invoke backend error: %v", err)
				} else {
					logrus.Warn("    [core.service] invoke backend generic error, no futher details")
				}
				res = &v10.GetNotReadConversationSetReply{ProtoOp: v10.Op_GetNotReadConversationSetAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_INVOKE}
			}

		} else {
			logrus.Warn("    [core.service] no backend or no backend api found")
			res = &v10.GetNotReadConversationSetReply{ProtoOp: v10.Op_GetNotReadConversationSetAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_API_NOTFOUND}
		}
	} else {
		res = &v10.GetNotReadConversationSetReply{ProtoOp: v10.Op_GetNotReadConversationSetAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	}
	return
}

//
//
//

func (s *ImCoreService) SubscribeX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_Subscribe, v10.Op_SubscribeAck, "Subscribe", ctx, req)
}

// rpc Subscribe (SubscribeReq) returns (SubscribeReply);
// rpc Unsubscribe (UnsubscribeReq) returns (UnsubscribeReply);
func (s *ImCoreService) Subscribe(ctx context.Context, req *v10.SubscribeReq) (res *v10.SubscribeReply, err error) {
	r, e := s.SubscribeX(ctx, req)
	res = r.(*v10.SubscribeReply)
	err = e
	return
}

func (s *ImCoreService) UnsubscribeX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_Unsubscribe, v10.Op_UnsubscribeAck, "Unsubscribe", ctx, req)
}

func (s *ImCoreService) Unsubscribe(ctx context.Context, req *v10.UnsubscribeReq) (res *v10.UnsubscribeReply, err error) {
	r, e := s.UpdateMsgX(ctx, req)
	res = r.(*v10.UnsubscribeReply)
	err = e
	return
}

func (s *ImCoreService) UpdateMsgX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_UpdateMsg, v10.Op_UpdateMsgAck, "UpdateMsg", ctx, req)
}

func (s *ImCoreService) UpdateMsg(ctx context.Context, req *v10.UpdateMsgReq) (res *v10.UpdateMsgReply, err error) {
	r, e := s.UpdateMsgX(ctx, req)
	res = r.(*v10.UpdateMsgReply)
	err = e
	return
}
