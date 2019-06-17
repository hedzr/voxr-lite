/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"context"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-lite/internal/exception"
)

type ImMemberService struct {
}

func (s *ImMemberService) AddMember(ctx context.Context, req *v10.AddMemberReq) (res *v10.AddMemberReply, err error) {
	res = &v10.AddMemberReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || (req.ProtoOp != v10.Op_TopicsAll && req.ProtoOp != v10.Op_OrgsAll) {
		err = exception.New(exception.InvalidParams)
		return
	}

	return
}

func (s *ImMemberService) RemoveMember(ctx context.Context, req *v10.RemoveMemberReq) (res *v10.RemoveMemberReply, err error) {
	res = &v10.RemoveMemberReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || (req.ProtoOp != v10.Op_TopicsAll && req.ProtoOp != v10.Op_OrgsAll) {
		err = exception.New(exception.InvalidParams)
		return
	}

	return
}

func (s *ImMemberService) InviteMember(ctx context.Context, req *v10.InviteMemberReq) (res *v10.InviteMemberReply, err error) {
	res = &v10.InviteMemberReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || (req.ProtoOp != v10.Op_TopicsAll && req.ProtoOp != v10.Op_OrgsAll) {
		err = exception.New(exception.InvalidParams)
		return
	}

	return
}

func (s *ImMemberService) UpdateMember(ctx context.Context, req *v10.UpdateMemberReq) (res *v10.UpdateMemberReply, err error) {
	res = &v10.UpdateMemberReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || (req.ProtoOp != v10.Op_TopicsAll && req.ProtoOp != v10.Op_OrgsAll) {
		err = exception.New(exception.InvalidParams)
		return
	}

	return
}

func (s *ImMemberService) ListMembers(ctx context.Context, req *v10.ListMembersReq) (res *v10.ListMembersReply, err error) {
	res = &v10.ListMembersReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || (req.ProtoOp != v10.Op_TopicsAll && req.ProtoOp != v10.Op_OrgsAll) {
		err = exception.New(exception.InvalidParams)
		return
	}

	return
}
