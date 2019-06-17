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

type ImTopicService struct {
	ImBase // for createBot
	ImMemberService
	dao *dao.TopicDao
	// daoMember *dao.MemberDao
}

func NewImTopicService() *ImTopicService {
	return &ImTopicService{dao: dao.NewTopicDao()}
}

func (s *ImTopicService) List(ctx context.Context, req *v10.ListTopicsReq) (res *v10.ListTopicsReply, err error) {
	res = &v10.ListTopicsReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_TopicsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	// // 列出给定上级的第一级下级组织机构，不会递归多级机构。
	// // orgId=0 将会列举全部顶级组织机构
	var ret []*models.Topic
	ret, err = s.dao.ListFromParentId(req.TopicId)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	for _, m := range ret {
		mx := m.ToProto()
		res.Topics = append(res.Topics, mx)
	}
	res.ErrorCode = v10.Err_OK
	return
}

func (s *ImTopicService) Add(ctx context.Context, req *v10.AddTopicReq) (res *v10.AddTopicReply, err error) {
	res = &v10.AddTopicReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_TopicsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	model := (&models.Topic{}).FromProto(req.Topic)
	err = s.dao.Add(model, model.OrgId) // .AddOrUpdate(model)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	logrus.Debugf("topic added: %v, %v", model.Id, model.Name)
	res.ErrorCode = v10.Err_OK
	res.Id = model.Id

	mq.RaiseEvent(v10.GlobalEvents_EvTopicAdded, model)
	logrus.Debugf("org ebq sent: EvTopicAdded: %v", model)

	if model.BotId == 0 {
		err = s.createBot(req.ProtoOp, req.Seq, model, res, model.Id, model.Name, func(botId uint64) {
			topic := &models.Topic{BaseModel: models.BaseModel{Id: model.Id}, BotId: uint64(botId)}
			if err = s.dao.Update(topic); err != nil {
				logrus.Errorf("Err: link bot id %v to topic %v failed. err = %v.", topic.Id, topic.BotId, err)
			} else {
				logrus.Debugf("topic add: bot created: bot=%v, topic=%v", topic.BotId, topic)
				mq.RaiseEvent(v10.GlobalEvents_EvTopicUpdated, topic)
			}
		})
	}
	logrus.Debugf("topic add end, err=%v", err)

	return
}

func (s *ImTopicService) Remove(ctx context.Context, req *v10.RemoveTopicReq) (res *v10.RemoveTopicReply, err error) {
	res = &v10.RemoveTopicReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_TopicsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	err = s.dao.RemoveById(req.TopicId)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	res.ErrorCode = v10.Err_OK
	res.RowsAffected = 1

	mq.RaiseEvent(v10.GlobalEvents_EvTopicRemoved, res)

	return
}

func (s *ImTopicService) Update(ctx context.Context, req *v10.UpdateTopicReq) (res *v10.UpdateTopicReply, err error) {
	res = &v10.UpdateTopicReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_TopicsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	model := (&models.Topic{}).FromProto(req.Topic)
	err = s.dao.Update(model)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	res.ErrorCode = v10.Err_OK
	res.RowsAffected = 1

	mq.RaiseEvent(v10.GlobalEvents_EvTopicUpdated, res)

	return
}

func (s *ImTopicService) Get(ctx context.Context, req *v10.GetTopicReq) (res *v10.GetTopicReply, err error) {
	res = &v10.GetTopicReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_TopicsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	var ret *models.Topic
	var r *v10.Topic
	if req.TopicId > 0 {
		// get by id
		ret, err = s.dao.GetById(req.TopicId)
	} else {
		// get by parentId and name
		ret, err = s.dao.Get(&models.Topic{Pid: req.ParentId, Name: req.Name})
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
	res.Topic = r
	return
}

func (s *ImTopicService) GetTopicUnreads(ctx context.Context, req *v10.GetTopicUnreadsReq) (res *v10.GetTopicUnreadsReply, err error) {
	res = &v10.GetTopicUnreadsReply{ProtoOp: v10.Op_TopicsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_TopicsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	return
}
