/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"context"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"github.com/hedzr/voxr-lite/misc/impl/mq"
	"github.com/sirupsen/logrus"
)

type ImOrgService struct {
	ImBase // for createBot
	ImMemberService
	dao *dao.OrgDao
	// daoMember *dao.MemberDao
}

// var imOrgService = &ImOrgService{dao.NewOrgDao(),}

func NewImOrgService() *ImOrgService {
	return &ImOrgService{dao: dao.NewOrgDao()}
}

func (s *ImOrgService) List(ctx context.Context, req *v10.ListOrgsReq) (res *v10.ListOrgsReply, err error) {
	res = &v10.ListOrgsReply{ProtoOp: v10.Op_OrgsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil {
		err = exception.New(exception.InvalidParams)
		return
	}

	// 列出给定上级的第一级下级组织机构，不会递归多级机构。
	// orgId=0 将会列举全部顶级组织机构
	var ret []*models.Organization
	ret, err = s.dao.ListFromParentId(req.OrgId)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	for _, m := range ret {
		mx := m.ToProto()
		res.Orgs = append(res.Orgs, mx)
	}
	res.ErrorCode = v10.Err_OK
	return
}

func (s *ImOrgService) Add(ctx context.Context, req *v10.AddOrgReq) (res *v10.AddOrgReply, err error) {
	res = &v10.AddOrgReply{ProtoOp: v10.Op_OrgsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil {
		err = exception.New(exception.InvalidParams)
		return
	}

	model := (&models.Organization{}).FromProto(req.Org)
	model, err = s.dao.AddOrUpdate(model, model.Pid)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	logrus.Debugf("org added: %v, %v", model.Id, model.Name)
	res.ErrorCode = v10.Err_OK
	res.Id = model.Id

	mq.RaiseEvent(v10.GlobalEvents_EvOrgAdded, model)
	logrus.Debugf("org ebq sent: EvOrgAdded: %v", model)

	if model.BotId == 0 {
		err = s.createBot(req.ProtoOp, req.Seq, model, res, model.Id, model.Name, func(botId uint64) {
			org := &models.Organization{BaseModel: models.BaseModel{Id: model.Id}, BotId: botId}
			if err = s.dao.Update(org); err != nil {
				logrus.Errorf("Err: link bot id %v to org %v failed. err = %v.", org.Id, org.BotId, err)
			} else {
				logrus.Debugf("org add: bot created: bot=%v, org=%v", org.BotId, org)
				mq.RaiseEvent(v10.GlobalEvents_EvTopicUpdated, org)
			}
		})
	}
	logrus.Debugf("org add end, err=%v", err)

	return
}

func (s *ImOrgService) Remove(ctx context.Context, req *v10.RemoveOrgReq) (res *v10.RemoveOrgReply, err error) {
	res = &v10.RemoveOrgReply{ProtoOp: v10.Op_OrgsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil {
		err = exception.New(exception.InvalidParams)
		return
	}

	var rows int64
	rows, err = s.dao.RemoveById(req.OrgId)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	res.ErrorCode = v10.Err_OK
	res.RowsAffected = rows

	mq.RaiseEvent(v10.GlobalEvents_EvOrgRemoved, res)

	return
}

func (s *ImOrgService) Update(ctx context.Context, req *v10.UpdateOrgReq) (res *v10.UpdateOrgReply, err error) {
	res = &v10.UpdateOrgReply{ProtoOp: v10.Op_OrgsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil {
		err = exception.New(exception.InvalidParams)
		return
	}

	model := (&models.Organization{}).FromProto(req.Org)
	err = s.dao.Update(model)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	res.ErrorCode = v10.Err_OK
	res.RowsAffected = 1

	mq.RaiseEvent(v10.GlobalEvents_EvOrgUpdated, res)

	return
}

func (s *ImOrgService) Get(ctx context.Context, req *v10.GetOrgReq) (res *v10.GetOrgReply, err error) {
	res = &v10.GetOrgReply{ProtoOp: v10.Op_OrgsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil {
		err = exception.New(exception.InvalidParams)
		return
	}

	var ret *models.Organization
	var r *v10.Organization
	if req.OrgId > 0 {
		// get by id
		ret, err = s.dao.GetById(req.OrgId)
	} else {
		// get by parentId and name
		ret, err = s.dao.Get(&models.Organization{Pid: req.ParentId, Name: req.Name})
	}

	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	if req.LoadAll {
		r, err = s.dao.LazyLoad(ret)
		if err != nil {
			err = exception.NewWith(exception.DaoError, err)
			return
		}
	}

	res.ErrorCode = v10.Err_OK
	res.Org = r
	return
}
