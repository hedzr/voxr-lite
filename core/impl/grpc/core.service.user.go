/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
)

func (s *ImCoreService) AddFriendX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_AddFriend, v10.Op_AddFriendAck, "AddFriend", ctx, req)
}
func (s *ImCoreService) AddFriend(ctx context.Context, req *v10.AddFriendReq) (res *v10.AddFriendReply, err error) {
	r, e := s.AddFriendX(ctx, req)
	res = r.(*v10.AddFriendReply)
	err = e
	return
}
func (s *ImCoreService) UpdateFriendX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_UpdateFriend, v10.Op_UpdateFriendAck, "UpdateFriend", ctx, req)
}
func (s *ImCoreService) UpdateFriend(ctx context.Context, req *v10.UpdateFriendReq) (res *v10.UpdateFriendReply, err error) {
	r, e := s.UpdateFriendX(ctx, req)
	res = r.(*v10.UpdateFriendReply)
	err = e
	return
}
func (s *ImCoreService) DeleteFriendX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_DeleteFriend, v10.Op_DeleteFriendAck, "DeleteFriend", ctx, req)
}
func (s *ImCoreService) DeleteFriend(ctx context.Context, req *v10.DeleteFriendReq) (res *v10.DeleteFriendReply, err error) {
	r, e := s.DeleteFriendX(ctx, req)
	res = r.(*v10.DeleteFriendReply)
	err = e
	return
}
func (s *ImCoreService) GetFriendListX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_GetFriendList, v10.Op_GetFriendListAck, "GetFriendList", ctx, req)
}

func (s *ImCoreService) GetFriendList(ctx context.Context, req *v10.GetFriendListReq) (res *v10.GetFriendListReply, err error) {
	r, e := s.GetFriendListX(ctx, req)
	res = r.(*v10.GetFriendListReply)
	err = e
	return
}

//

//

func (s *ImCoreService) AddContactX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_AddContact, v10.Op_AddContactAck, "AddContact", ctx, req)
	// return s.AddContact(ctx, req.(*v10.AddContactReq))
}

func (s *ImCoreService) AddContact(ctx context.Context, req *v10.AddContactReq) (res *v10.AddContactReply, err error) {
	r, e := s.AddContactX(ctx, req)
	res = r.(*v10.AddContactReply)
	err = e
	return
}
func (s *ImCoreService) RemoveContactX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_RemoveContact, v10.Op_RemoveContactAck, "RemoveContact", ctx, req)
	// return s.RemoveContact(ctx, req.(*v10.RemoveContactReq))
}

func (s *ImCoreService) RemoveContact(ctx context.Context, req *v10.RemoveContactReq) (res *v10.RemoveContactReply, err error) {
	r, e := s.RemoveContactX(ctx, req)
	res = r.(*v10.RemoveContactReply)
	err = e
	return
}
func (s *ImCoreService) UpdateContactX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_UpdateContact, v10.Op_UpdateContactAck, "UpdateContact", ctx, req)
	// return s.UpdateContact(ctx, req.(*v10.UpdateContactReq))
}

func (s *ImCoreService) UpdateContact(ctx context.Context, req *v10.UpdateContactReq) (res *v10.UpdateContactReply, err error) {
	r, e := s.UpdateContactX(ctx, req)
	res = r.(*v10.UpdateContactReply)
	err = e
	return
}

func (s *ImCoreService) ListContactsX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_ListContacts, v10.Op_ListContactsAck, "ListContacts", ctx, req)
	// return s.ListContacts(ctx, req.(*v10.ListContactsReq))
}

func (s *ImCoreService) ListContacts(ctx context.Context, req *v10.ListContactsReq) (res *v10.ListContactsReply, err error) {
	r, e := s.ListContactsX(ctx, req)
	res = r.(*v10.ListContactsReply)
	err = e
	return
	// if req != nil && req.ProtoOp == v10.Op_ListContacts && req.Seq > 0 {
	// 	if fnRPC, ok := s.fwdrs["ListContacts"]; ok {
	//
	// 		var r proto.Message
	// 		r, err = fnRPC(ctx, req) // see also: service.BuildFwdr()
	//
	// 		if rr, ok := r.(*v10.ListContactsReply); ok {
	// 			res = rr
	// 		} else {
	// 			if err != nil {
	// 				logrus.Warnf("    [core.service] invoke backend error: %v", err)
	// 			} else {
	// 				logrus.Warn("    [core.service] invoke backend generic error, no futher details")
	// 			}
	// 			res = &v10.ListContactsReply{ProtoOp: v10.Op_ListContactsAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_INVOKE,}
	// 		}
	//
	// 	} else {
	// 		logrus.Warn("    [core.service] no backend or no backend api found")
	// 		res = &v10.ListContactsReply{ProtoOp: v10.Op_ListContactsAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_API_NOTFOUND,}
	// 	}
	// } else {
	// 	res = &v10.ListContactsReply{ProtoOp: v10.Op_ListContactsAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS,}
	// }
	// return
}

func (s *ImCoreService) GetContactX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_GetContact, v10.Op_GetContactAck, "GetContact", ctx, req)
	// return s.GetContact(ctx, req.(*v10.GetContactReq))
}

func (s *ImCoreService) GetContact(ctx context.Context, req *v10.GetContactReq) (res *v10.GetContactReply, err error) {
	r, e := s.GetContactX(ctx, req)
	res = r.(*v10.GetContactReply)
	err = e
	return
	// if req != nil && req.ProtoOp == v10.Op_GetContact && req.Seq > 0 {
	// 	if fnRPC, ok := s.fwdrs["GetContact"]; ok {
	//
	// 		var r proto.Message
	// 		r, err = fnRPC(ctx, req) // see also: service.BuildFwdr()
	//
	// 		if rr, ok := r.(*v10.GetContactReply); ok {
	// 			res = rr
	// 		} else {
	// 			if err != nil {
	// 				logrus.Warnf("    [core.service] invoke backend error: %v", err)
	// 			} else {
	// 				logrus.Warn("    [core.service] invoke backend generic error, no futher details")
	// 			}
	// 			res = &v10.GetContactReply{ProtoOp: v10.Op_GetContactAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_INVOKE,}
	// 		}
	//
	// 	} else {
	// 		logrus.Warn("    [core.service] no backend or no backend api found")
	// 		res = &v10.GetContactReply{ProtoOp: v10.Op_GetContactAck, Seq: req.Seq, ErrorCode: v10.Err_BACKEND_API_NOTFOUND,}
	// 	}
	// } else {
	// 	res = &v10.GetContactReply{ProtoOp: v10.Op_GetContactAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS,}
	// }
	// return
}

func (s *ImCoreService) GetUserContactsX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_GetUserContacts, v10.Op_GetUserContactsAck, "GetUserContacts", ctx, req)
	// return s.GetContact(ctx, req.(*v10.GetContactReq))
}

func (s *ImCoreService) GetUserContacts(ctx context.Context, req *v10.GetUserContactsReq) (res *v10.GetUserContactsReply, err error) {
	r, e := s.GetUserContactsX(ctx, req)
	res = r.(*v10.GetUserContactsReply)
	err = e
	return
}

func (s *ImCoreService) SetUserContactsX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	return s.xmas(v10.Op_SetUserContacts, v10.Op_SetUserContactsAck, "SetUserContacts", ctx, req)
	// return s.GetContact(ctx, req.(*v10.GetContactReq))
}

func (s *ImCoreService) SetUserContacts(ctx context.Context, req *v10.SetUserContactsReq) (res *v10.SetUserContactsReply, err error) {
	r, e := s.SetUserContactsX(ctx, req)
	res = r.(*v10.SetUserContactsReply)
	err = e
	return
}

func (s *ImCoreService) ContactOperate(ctx context.Context, req *v10.ContactAllReq) (res *v10.ContactAllReply, err error) {
	r, e := s.ContactOperateX(ctx, req)
	res = r.(*v10.ContactAllReply)
	err = e
	return
}

func (s *ImCoreService) ContactOperateX(ctx context.Context, req proto.Message) (res proto.Message, err error) {
	if r, ok := req.(*v10.ContactAllReq); ok {
		if r.GetAcr() != nil {
			var ret *v10.AddContactReply
			ret, err = s.AddContact(ctx, r.GetAcr())
			res = &v10.ContactAllReply{ProtoOp: v10.Op_ContactAllAck, Seq: ret.Seq, ErrorCode: ret.ErrorCode, Oneof: &v10.ContactAllReply_Acr{ret}}
			return
		}
		if r.GetRcr() != nil {
			var ret *v10.RemoveContactReply
			ret, err = s.RemoveContact(ctx, r.GetRcr())
			res = &v10.ContactAllReply{ProtoOp: v10.Op_ContactAllAck, Seq: ret.Seq, ErrorCode: ret.ErrorCode, Oneof: &v10.ContactAllReply_Rcr{ret}}
			return
		}
		if r.GetUcr() != nil {
			var ret *v10.UpdateContactReply
			ret, err = s.UpdateContact(ctx, r.GetUcr())
			res = &v10.ContactAllReply{ProtoOp: v10.Op_ContactAllAck, Seq: ret.Seq, ErrorCode: ret.ErrorCode, Oneof: &v10.ContactAllReply_Ucr{ret}}
			return
		}
		if r.GetLcr() != nil {
			var ret *v10.ListContactsReply
			ret, err = s.ListContacts(ctx, r.GetLcr())
			res = &v10.ContactAllReply{ProtoOp: v10.Op_ContactAllAck, Seq: ret.Seq, ErrorCode: ret.ErrorCode, Oneof: &v10.ContactAllReply_Lcr{ret}}
			return
		}
		if r.GetGcr() != nil {
			var ret *v10.GetContactReply
			ret, err = s.GetContact(ctx, r.GetGcr())
			res = &v10.ContactAllReply{ProtoOp: v10.Op_ContactAllAck, Seq: ret.Seq, ErrorCode: ret.ErrorCode, Oneof: &v10.ContactAllReply_Gcr{ret}}
			return
		}
		if r.GetGucr() != nil {
			var ret *v10.GetUserContactsReply
			ret, err = s.GetUserContacts(ctx, r.GetGucr())
			res = &v10.ContactAllReply{ProtoOp: v10.Op_ContactAllAck, Seq: ret.Seq, ErrorCode: ret.ErrorCode, Oneof: &v10.ContactAllReply_Gucr{ret}}
			return
		}
		if r.GetSucr() != nil {
			var ret *v10.SetUserContactsReply
			ret, err = s.SetUserContacts(ctx, r.GetSucr())
			res = &v10.ContactAllReply{ProtoOp: v10.Op_ContactAllAck, Seq: ret.Seq, ErrorCode: ret.ErrorCode, Oneof: &v10.ContactAllReply_Sucr{ret}}
			return
		}

		// return s.xmas(v10.Op_GetContact, v10.Op_GetContactAck, "GetContact", ctx, req)
	}
	// return s.GetContact(ctx, req.(*v10.GetContactReq))
	return
}
