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
	"github.com/hedzr/voxr-lite/misc/impl/filters"
	"github.com/hedzr/voxr-lite/misc/impl/mq"
	"github.com/sirupsen/logrus"
)

type ImMsgService struct {
	dao *dao.MsgDao
	// daoMember *dao.MemberDao
}

// var topicDao = dao.NewTopicDao()

func NewImMsgService() *ImMsgService {
	return &ImMsgService{dao.NewMsgDao()}
}

func (s *ImMsgService) applyPreFilters(model *models.Msg) (ret *models.Msg, err error) {
	// pre-filters
	ret, err = filters.CallPre(v10.GlobalEvents_EvMsgIncoming, model)
	return
}

func (s *ImMsgService) applyPostFilters(model *models.Msg) (ret *models.Msg, err error) {
	// post-filters
	ret, err = filters.CallPost(v10.GlobalEvents_EvMsgPulling, model)
	return
}

//
//
//

func (s *ImMsgService) SendMsg(ctx context.Context, req *v10.SendMsgReqV12) (res *v10.SendMsgReplyV12, err error) {
	res = &v10.SendMsgReplyV12{ProtoOp: v10.Op_MsgsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS, TopicId: req.Msg.TopicId, UserId: req.Msg.From, DeviceId: req.DeviceId}
	if req == nil || req.ProtoOp != v10.Op_MsgsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	model := (&models.Msg{}).FromProto(req.Msg)

	// if len(req.HilightMembers) > 0 {
	// 	model.ToUsers = model.FromUsersArray(req.HilightMembers)
	// 	model.TopicId = req.TopicId
	// }

	if model, err = s.applyPreFilters(model); err != nil {
		return
	}

	res.ToUsers, err = s.dao.AddTo(model, req.DeviceId, req.ParentMsgId, req.OrgId)
	if err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	res.MsgId = model.Id
	res.OrgId = req.OrgId
	res.TopicId = model.TopicId // model.TopicId might be updated by creating a new room/directMessage
	res.ErrorCode = v10.Err_OK
	logrus.Debugf("msg added: id=%v, %v", model.Id, model)

	mq.RaiseEvent(v10.GlobalEvents_EvMsgIncoming, model)
	logrus.Debugf("msg ebq sent: EvMsgIncoming: %v", model)

	// vx-core:
	// 	   notify all members under topic

	return
}

func (s *ImMsgService) UpdateMsg(ctx context.Context, req *v10.UpdateMsgReqV12) (res *v10.UpdateMsgReplyV12, err error) {
	res = &v10.UpdateMsgReplyV12{ProtoOp: v10.Op_MsgsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS}
	if req == nil || req.ProtoOp != v10.Op_MsgsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	res.ErrorCode = v10.Err_OK
	return
}

// get the newest msg list in a topic, the page size is default 10
func (s *ImMsgService) GetMsg(ctx context.Context, req *v10.GetMsgReqV12) (res *v10.GetMsgReplyV12, err error) {
	res = &v10.GetMsgReplyV12{ProtoOp: v10.Op_MsgsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS, TopicId: req.TopicId, Newer: req.Newer}
	if req == nil || req.ProtoOp != v10.Op_MsgsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	var (
		// model    *models.Msg
		list     []*models.Msg
		cond     string
		limit    int = 10
		orderBy  string
		ret      []*v10.Msg
		from, to uint64
		// maxId   uint64
	)

	// // count=1，提取一条消息; 否则提取 count 条消息；
	// // count=0，提取自动多条消息（根据系统负载情况返回10-20条消息）；
	// // 不应试图一次性提取大量消息数据。为了保证整体服务效率，核心服务最多只会返回20条消息，即使指定了更大的 count 也不会产生效果。
	// int32 count = 11;
	// // newer = false, 向前提取; newer = true, 向后提取。

	if req.Count > 10 && req.Count <= 20 {
		limit = int(req.Count)
	}

	if req.Newer {
		cond = "t_msg.id >= ? and t_msg.topic_id = ?"
	} else {
		cond = "t_msg.id <= ? and t_msg.topic_id = ?"
	}

	if req.SortByAsc {
		orderBy = "t_msg.lft"
	} else {
		orderBy = "t_msg.lft desc"
	}

	if list, err = s.dao.ListBy(req.AutoAck, req.UserId, req.TopicId, limit, orderBy, cond, req.MsgId, req.TopicId); err != nil {
		err = exception.NewWith(exception.DaoError, err)
		return
	}

	if len(list) > 0 {
		ret = make([]*v10.Msg, len(list))
		for i, m := range list {
			if m, err = s.applyPostFilters(m); err == nil {
				ret[i] = m.ToProto()
				// tm := time.Unix(0, 0)
				// logrus.Debugf("%5d. pb.deleteAt=%v (%v), model.deleteAt=%v", i, ret[i].DeletedAt, tm, m.DeletedAt)
			}
		}
		res.Msgs = ret

		from = ret[len(ret)-1].Id
		to = ret[0].Id
		if from > to {
			from += to
			to = from - to
			from -= to
		}

		// newer = false, 向前(时间更早)提取; newer = true, 向后(时间更晚)提取。
		if req.Newer {
			res.NextMsgId = to + 1
		} else {
			res.NextMsgId = from - 1
		}
	}

	res.ErrorCode = v10.Err_OK
	res.Count = int32(len(res.Msgs))
	// res.TopicId = req.TopicId
	// res.Newer = req.Newer

	// err = s.applyPostFilters(model)
	return
}

func (s *ImMsgService) AckMsg(ctx context.Context, req *v10.AckMsgReqV12) (res *v10.AckMsgReplyV12, err error) {
	res = &v10.AckMsgReplyV12{ProtoOp: v10.Op_MsgsAllAck, Seq: req.Seq, ErrorCode: v10.Err_INVALID_PARAMS, TopicId: req.TopicId}
	if req == nil || req.ProtoOp != v10.Op_MsgsAll {
		err = exception.New(exception.InvalidParams)
		return
	}

	err = s.dao.MsgAcked(req.MemberId, req.TopicId, int(req.Count), int(req.Count), req.MsgIds...)
	return
}
