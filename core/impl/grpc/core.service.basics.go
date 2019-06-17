/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-lite/core"
)

//
// ImCore services implementations
//

func (s *ImCoreService) GetX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.GetOffLineConversationSet(ctx, req.(*v10.GetOffLineConversationSetReq))
}

func (s *ImCoreService) Get(ctx context.Context, req *v10.GetOffLineConversationSetReq) (res *v10.GetOffLineConversationSetReply, err error) {
	panic("implement me")
}

func (s *ImCoreService) NothingX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.Nothing(ctx, req.(*v10.Empty))
}

func (s *ImCoreService) Nothing(ctx context.Context, req *v10.Empty) (res *v10.Result, err error) {
	res = BaseResult(true, 0, "", 0, nil)
	return
}

func (s *ImCoreService) PingX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.Ping(ctx, req.(*v10.Empty))
}

func (s *ImCoreService) Ping(ctx context.Context, req *v10.Empty) (empty *v10.Empty, err error) {
	empty = new(v10.Empty)
	return
}

func (s *ImCoreService) CloseX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.Close(ctx, req.(*v10.Empty))
}

func (s *ImCoreService) Close(ctx context.Context, req *v10.Empty) (empty *v10.Empty, err error) {
	empty = new(v10.Empty)
	return
}

func (s *ImCoreService) HandshakeX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.Handshake(ctx, req.(*v10.HandshakeReq))
}

func (s *ImCoreService) Handshake(ctx context.Context, req *v10.HandshakeReq) (res *v10.ServiceDiag, err error) {
	// ws clients by device-id
	// adding requestingProtocolVersion

	var pv int32 = CoreProtocolVersionInt
	if req.RequestingProtocolVersion > 0 {
		pv = req.RequestingProtocolVersion
	}

	res = &v10.ServiceDiag{
		ProtoOp:                   v10.Op_HandshakeAck,
		GatewayVer:                GatewayVersionInt,
		CoreVer:                   core.VersionInt,
		ProtocolVer:               CoreProtocolVersionInt,
		Encryption:                "none",
		Compress:                  "none",
		NegotiatedProtocolVersion: pv,
	}
	return
}

func (s *ImCoreService) PushMsgX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.PushMsg(ctx, req.(*v10.PushMsgReq))
}

// //PushMsg push by key or mid
// rpc PushMsg (PushMsgReq) returns (PushMsgReply);
// // Broadcast send to every entity
// rpc Broadcast (BroadcastReq) returns (BroadcastReply);
// // BroadcastRoom broadcast to one room
// rpc BroadcastRoom (BroadcastRoomReq) returns (BroadcastRoomReply);
// // Rooms get all rooms
// rpc Rooms (RoomsReq) returns (RoomsReply);
func (s *ImCoreService) PushMsg(ctx context.Context, req *v10.PushMsgReq) (res *v10.PushMsgReply, err error) {
	return
}
func (s *ImCoreService) BroadcastX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.Broadcast(ctx, req.(*v10.BroadcastReq))
}
func (s *ImCoreService) Broadcast(ctx context.Context, req *v10.BroadcastReq) (res *v10.BroadcastReply, err error) {
	return
}
func (s *ImCoreService) BroadcastRoomX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.BroadcastRoom(ctx, req.(*v10.BroadcastRoomReq))
}
func (s *ImCoreService) BroadcastRoom(ctx context.Context, req *v10.BroadcastRoomReq) (res *v10.BroadcastRoomReply, err error) {
	return
}
func (s *ImCoreService) RoomsX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.Rooms(ctx, req.(*v10.RoomsReq))
}
func (s *ImCoreService) Rooms(ctx context.Context, req *v10.RoomsReq) (res *v10.RoomsReply, err error) {
	return
}
