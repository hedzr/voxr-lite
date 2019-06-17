/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
)

type (
	// CmdInfo 为 ws 端提供调用转换。
	// ws 服务端收到报文后，按照 Op 指令的不同直接调用 ImCoreService.PooledInvoke() 从而实现到GRPC调用的转换
	CmdInfo struct {
		Op          v10.Op
		OpAck       v10.Op
		Target      func(ctx context.Context, req proto.Message) (res proto.Message, err error)
		InTemplate  func() proto.Message
		OutTemplate func(op v10.Op, seq uint32, errCode v10.Err) proto.Message // by xmas()
		// Out    func(out proto.Message) (b []byte)
	}

	CmdFunc func(from WsClientSkel, body []byte)

	WsClientSkel interface {
		PostBinaryMsg([]byte) // 向WS客户端写回响应报文
	}

	Request struct {
		CmdInfo *CmdInfo
		From    WsClientSkel
		InParam []byte
		// Result  chan []byte
	}
)

//
// 定义 ws 请求映射表。将来自于 ws client 的 API 请求转换为后端适用的 proto.Message 并完成调用，然后返回调用结果。
//
// 注意，这里提供的一切接口一律以 core.proto 为蓝本。
// core.proto 提供公开接口，并处理收到的请求报文，转换为后端报文后执行 grpc 调用。
//
func (s *ImCoreService) aList() []*CmdInfo {
	return []*CmdInfo{

		{v10.Op_SendMsg, v10.Op_SendMsg, s.SendMsgX, func() proto.Message { return new(v10.SendMsgReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.SendMsgReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_UpdateMsg, v10.Op_UpdateMsgAck, s.UpdateMsgX, func() proto.Message { return new(v10.SendMsgReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.SendMsgReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_GetMsgList, v10.Op_GetMsgListAck, s.GetMsgListX, func() proto.Message { return new(v10.GetMsgListReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.GetMsgListReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_GetMsgHistory, v10.Op_GetMsgHistoryAck, s.GetMsgHistoryX, func() proto.Message { return new(v10.GetMsgHistoryReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.GetMsgHistoryReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_AckMsg, v10.Op_AckMsgAck, s.AckMsgX, func() proto.Message { return new(v10.AckMsgReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.AckMsgReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_GetOffLineConversationSet, v10.Op_GetOffLineConversationSetAck, s.GetOffLineConversationSetX, func() proto.Message { return new(v10.GetOffLineConversationSetReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.GetOffLineConversationSetReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_GetNotReadConversationSet, v10.Op_GetNotReadConversationSetAck, s.GetNotReadConversationSetX, func() proto.Message { return new(v10.GetNotReadConversationSetReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.GetNotReadConversationSetReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		// Contact v11

		{v10.Op_ContactAll, v10.Op_ContactAllAck, s.ContactOperateX, func() proto.Message { return new(v10.ContactAllReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.ContactAllReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		// 以下 contact 转发将考虑废弃
		{v10.Op_GetContact, v10.Op_GetContactAck, s.GetContactX, func() proto.Message { return new(v10.GetContactReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.GetContactReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_AddContact, v10.Op_AddContactAck, s.AddContactX, func() proto.Message { return new(v10.AddContactReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.AddContactReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_RemoveContact, v10.Op_RemoveContactAck, s.RemoveContactX, func() proto.Message { return new(v10.RemoveContactReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.RemoveContactReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_UpdateContact, v10.Op_UpdateContactAck, s.UpdateContactX, func() proto.Message { return new(v10.UpdateContactReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.UpdateContactReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_ListContacts, v10.Op_ListContactsAck, s.ListContactsX, func() proto.Message { return new(v10.ListContactsReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.ListContactsReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_GetUserContacts, v10.Op_GetUserContactsAck, s.GetUserContactsX, func() proto.Message { return new(v10.GetUserContactsReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.GetUserContactsReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_SetUserContacts, v10.Op_SetUserContactsAck, s.SetUserContactsX, func() proto.Message { return new(v10.SetUserContactsReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.SetUserContactsReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		//
		// Auth v11
		//

		{v10.Op_LoginV11, v10.Op_LoginV11Ack, s.LoginV11X, func() proto.Message { return new(v10.AuthReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.AuthReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_RefreshTokenV11, v10.Op_RefreshTokenV11Ack, s.RefreshTokenV11X, func() proto.Message { return new(v10.AuthReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.AuthReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_RegisterV11, v10.Op_RegisterV11Ack, s.RegisterV11X, func() proto.Message { return new(v10.AuthReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.AuthReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		// Circle

		{v10.Op_CirclesAll, v10.Op_CirclesAllAck, s.CircleOperateX, func() proto.Message { return new(v10.CircleAllReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.CircleAllReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		// User

		{v10.Op_UserAll, v10.Op_UserAllAck, s.UserOperateX, func() proto.Message { return new(v10.UserAllReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.UserAllReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		// Search

		{v10.Op_SearchGlobal, v10.Op_SearchGlobalAck, s.SearchGlobalX, func() proto.Message { return new(v10.SearchGlobalReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.SearchGlobalReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		// QrCode

		{v10.Op_GenQrCode, v10.Op_GenQrCodeAck, s.GenerateQrCodeX, func() proto.Message { return new(v10.GenerateQrCodeReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.GenerateQrCodeReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		// Sender,Notifier

		{v10.Op_SendSMS, v10.Op_SendSMSAck, s.SendSMSX, func() proto.Message { return new(v10.Empty) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.SearchGlobalReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_SendMail, v10.Op_SendMailAck, s.SendMailX, func() proto.Message { return new(v10.Empty) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.SearchGlobalReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		// Verifier

		{v10.Op_VerifyIdCard, v10.Op_VerifyIdCardAck, s.VerifyIdCardX, func() proto.Message { return new(v10.Empty) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.SearchGlobalReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_VerifyMobileNumber, v10.Op_VerifyMobileNumberAck, s.VerifyMobileNumberX, func() proto.Message { return new(v10.Empty) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.SearchGlobalReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		//
		// Friend
		//

		{v10.Op_AddFriend, v10.Op_AddFriendAck, s.AddFriendX, func() proto.Message { return new(v10.AddFriendReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.AddFriendReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_UpdateFriend, v10.Op_UpdateFriendAck, s.UpdateFriendX, func() proto.Message { return new(v10.UpdateFriendReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.UpdateFriendReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_DeleteFriend, v10.Op_DeleteFriendAck, s.DeleteFriendX, func() proto.Message { return new(v10.DeleteFriendReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.DeleteFriendReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_GetFriendList, v10.Op_GetFriendListAck, s.GetFriendListX, func() proto.Message { return new(v10.GetFriendListReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.GetFriendListReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		// Sub/Unsub

		{v10.Op_Subscribe, v10.Op_SubscribeAck, s.SubscribeX, func() proto.Message { return new(v10.SubscribeReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.SubscribeReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},
		{v10.Op_Unsubscribe, v10.Op_UnsubscribeAck, s.UnsubscribeX, func() proto.Message { return new(v10.UnsubscribeReq) }, func(op v10.Op, seq uint32, errCode v10.Err) proto.Message {
			return &v10.UnsubscribeReply{ProtoOp: op, Seq: seq, ErrorCode: errCode}
		}},

		//
		//
		// VX-CORE Basic APIs
		//
		//

		{v10.Op_Nothing, v10.Op_Nothing, s.NothingX, func() proto.Message { return new(v10.Empty) }, nil},
		{v10.Op_Ping, v10.Op_Ping, s.PingX, func() proto.Message { return new(v10.Empty) }, nil},
		{v10.Op_Close, v10.Op_Close, s.CloseX, func() proto.Message { return new(v10.Empty) }, nil},
		{v10.Op_Handshake, v10.Op_HandshakeAck, s.HandshakeX, func() proto.Message { return new(v10.HandshakeReq) }, nil},

		//
		// VX-CORE Reserved APIs
		//

		{v10.Op_PushMsg, v10.Op_PushMsg, s.PushMsgX, func() proto.Message { return new(v10.PushMsgReq) }, nil},
		{v10.Op_Broadcast, v10.Op_Broadcast, s.BroadcastX, func() proto.Message { return new(v10.BroadcastReq) }, nil},
		{v10.Op_BroadcastRoom, v10.Op_BroadcastRoom, s.BroadcastRoomX, func() proto.Message { return new(v10.BroadcastRoomReq) }, nil},
		{v10.Op_Rooms, v10.Op_Rooms, s.RoomsX, func() proto.Message { return new(v10.RoomsReq) }, nil},

		// {v10.Op_Ping, v10.Op_Ping, s.PingX, func() proto.Message { return new(v10.Empty) }, nil,},
	}
}
