/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"context"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"github.com/sirupsen/logrus"
)

type UserContactServer struct{}

//

//

func (s *UserContactServer) AddContact(ctx context.Context, req *v10.AddContactReq) (res *v10.Result, err error) {
	var r *v10.AddContactReply
	var rc *v10.Contact
	res = &v10.Result{Ok: false}
	rc, err = dao.AddContactFromUserInternal(req.UidOwner, req.UidFriend, req.GroupName, req.IsGroup,
		req.Relationship, req.RemarkName, req.Tags, req.Remarks)
	if err != nil {
		if e, ok := err.(*exception.XmError); ok {
			res.ErrCode = int32(e.Code)
			res.Msg = e.Msg
		}
		return
	}

	r = &v10.AddContactReply{ProtoOp: v10.Op_AddContactAck, Seq: req.Seq, Contact: rc}
	if res, err = M(r); err != nil {
		logrus.Debugf("ERR: %v", err)
	}
	return
}

func (s *UserContactServer) GetContact(ctx context.Context, req *v10.GetContactReq) (res *v10.Result, err error) {
	var r *v10.GetContactReply
	res = &v10.Result{Ok: false}
	r, err = dao.GetContact(req)
	if err != nil {
		if e, ok := err.(*exception.XmError); ok {
			res.ErrCode = int32(e.Code)
			res.Msg = e.Msg
		}
		return
	}

	if res, err = M(r); err != nil {
		logrus.Debugf("ERR: %v", err)
	}
	return
}

func (s *UserContactServer) RemoveContact(ctx context.Context, req *v10.RemoveContactReq) (res *v10.Result, err error) {
	var r *v10.RemoveContactReply
	res = &v10.Result{Ok: false}

	r, err = dao.RemoveContact(req)
	if err != nil {
		if e, ok := err.(*exception.XmError); ok {
			res.ErrCode = int32(e.Code)
			res.Msg = e.Msg
		}
		return
	}

	if res, err = M(r); err != nil {
		logrus.Debugf("ERR: %v", err)
	}
	return
}

func (s *UserContactServer) ListContacts(ctx context.Context, req *v10.ListContactsReq) (res *v10.Result, err error) {
	var r *v10.ContactGroups
	res = &v10.Result{Ok: false}

	r, err = dao.ListContacts(req.UidOwner)
	if err != nil {
		if e, ok := err.(*exception.XmError); ok {
			res.ErrCode = int32(e.Code)
			res.Msg = e.Msg
		}
		return
	}

	rr := &v10.ListContactsReply{ProtoOp: v10.Op_ListContactsAck, Seq: req.Seq, ErrorCode: v10.Err_OK, Groups: r}
	if res, err = M(rr); err != nil {
		logrus.Debugf("ERR: %v", err)
	}
	return
}

func (s *UserContactServer) UpdateContact(ctx context.Context, req *v10.UpdateContactReq) (res *v10.Result, err error) {
	var r *v10.UpdateContactReply
	res = &v10.Result{Ok: false}

	r, err = dao.UpdateContact(req)
	if err != nil {
		if e, ok := err.(*exception.XmError); ok {
			res.ErrCode = int32(e.Code)
			res.Msg = e.Msg
		}
		return
	}

	if res, err = M(r); err != nil {
		logrus.Debugf("ERR: %v", err)
	}
	return
}
