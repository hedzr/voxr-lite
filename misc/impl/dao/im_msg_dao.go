/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao

import (
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type MsgDao struct{ daoMember *MemberDao }

func NewMsgDao() *MsgDao {
	return &MsgDao{
		NewMemberDao(),
	}
}

func (s *MsgDao) GetById(id uint64) (ret *models.Msg, err error) {
	ret = new(models.Msg)
	err = dbe.GormDb().First(&ret, id).Error
	return
}

func (s *MsgDao) Get(tmpl *models.Msg) (ret *models.Msg, err error) {
	ret = new(models.Msg)
	err = dbe.GormDb().First(&ret, tmpl).Error
	return
}

// don't use this func because it's not completed on the internal affairs, such as updateing MsgLastPt, notifying others, pre/post filter, ...
func (s *MsgDao) AddOrUpdate1(in *models.Msg) (ret *models.Msg, err error) {
	// var yes bool
	// var rows int64
	//
	// if in.Id == 0 {
	// 	err = exception.New(exception.InvalidParams)
	// 	return
	// }
	//
	// tmpl := &models.Msg{BaseModel: models.BaseModel{Id: in.Id},}
	// yes, err = dbe.DBE.Engine().Exist(tmpl)
	// if err != nil {
	// 	logrus.Errorf("CANNOT check exists for the record: %v", err)
	// 	return
	// }
	//
	// if !yes {
	// 	rows, err = dbe.DBE.Engine().Insert(in)
	// 	if rows == 0 || err != nil {
	// 		logrus.Errorf("CANNOT insert the new record: %v", err)
	// 	} else {
	// 		ret = in
	// 	}
	// 	return
	// }
	//
	// tmpl = &models.Msg{BaseModel: models.BaseModel{Id: in.Id},}
	// rows, err = dbe.DBE.Engine().Update(in, tmpl)
	// if rows > 1 {
	// 	logrus.Errorf("multiple (%v) rows updated, it might be wrong.", rows)
	// } else if err != nil {
	// 	logrus.Errorf("CANNOT update the record: %v", err)
	// } else {
	// 	// ret = tmpl
	// 	yes, err = dbe.DBE.Engine().Get(tmpl)
	// 	if err != nil {
	// 		logrus.Errorf("CANNOT get the record: %v", err)
	// 	} else {
	// 		ret = tmpl
	// 	}
	// }
	return
}

// Soft Delete
func (s *MsgDao) RemoveById(id uint64) (err error) {
	err = dbe.GormDb().Delete(&models.Msg{BaseModel: models.BaseModel{Id: id}}).Error
	return
}

func (s *MsgDao) Remove(tmpl *models.Msg) (err error) {
	if tmpl.Id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	err = dbe.GormDb().Delete(tmpl).Error
	return
}

func (s *MsgDao) Update(in *models.Msg) (err error) {
	err = dbe.GormDb().Omit("id").Updates(in).Error
	return
}

func (s *MsgDao) Where(query interface{}, args ...interface{}) *gorm.DB {
	return dbe.GormDb().Model(&models.Msg{}).Where(query, args)
}

func (s *MsgDao) List(limit int, start int, query interface{}, args ...interface{}) (ret []*models.Msg, err error) {
	err = dbe.GormDb().Where(query, args).Limit(limit).Offset(start).Find(&ret).Error
	return
}

// ListFromParentId 列出给定上级msg的全部下级msg, 包含该上级
// parentId=0 时列出给定 topic 中的全部msgs
func (s *MsgDao) ListFromParentId(parentId uint64, topicId uint64) (ret []*models.Msg, err error) {
	model := &models.Msg{}
	scope := dbe.GormDb().Model(model).
		Joins("LEFT JOIN t_msg p ON t_msg.topic_id=p.topic_id AND t_msg.lft BETWEEN p.lft AND p.rgt").
		Where("p.topic_id=?", topicId)
	if parentId > 0 {
		scope = scope.Where("p.id=?", parentId)
	}
	err = scope.Group("t_msg.id").Order("lft").Find(&ret).Error
	return
}

// ListFromParent 列出给定上级msg的全部下级msg, 包含该上级
func (s *MsgDao) ListFromParent(parent *models.Msg, topicId uint64) (ret []*models.Msg, err error) {
	return s.ListFromParentId(parent.Id, topicId)
}

func (s *MsgDao) ListBy(autoAck bool, userId, topicId uint64, limit int, orderBy, whereCondition string, whereArgs ...interface{}) (ret []*models.Msg, err error) {
	var (
		model    = &models.Msg{}
		from, to uint64
		distance int
	)

	if err = dbe.GormDb().Model(model).Limit(limit).Where(whereCondition, whereArgs...).Order(orderBy).Find(&ret).Error; err != nil {
		return
	}

	if autoAck && len(ret) > 0 {
		from = ret[0].Id
		to = ret[len(ret)-1].Id
		if from > to {
			from += to
			to = from - to
			from -= to
		}
		distance = int(to - from + 1)
		err = s.MsgAcked(userId, topicId, len(ret), distance, from)
	}
	return
}

func TxRollback(tx *gorm.DB, err error) {
	if err != nil {
		tx.Rollback()
	}
}

func (s *MsgDao) MsgAcked(memberId, topicId uint64, count, distance int, msgIds ...uint64) (err error) {
	logrus.Debugf("MsgAcked: memberId=%v, topicId=%v, count=%v, distance=%v, msgIds = %v", memberId, topicId, count, distance, msgIds)

	var yes bool

	tx := dbe.GormDb().Begin()
	defer TxRollback(tx, err)

	var mlp = &models.MsgLastPt{}
	if err = tx.Where("tid=? and mid=?", topicId, memberId).First(mlp).Error; err != nil {
		logrus.Errorf("CANNOT get: %v, %v", yes, err)
		return
	}

	if distance <= 1 {
		// process msgIds array

		yes = false
		for _, x := range msgIds {
			if mlp.AddRange(x, x) {
				yes = true
			}
		}

	} else {
		// process msgIds[0] .. msgIds[0]+distance-1

		yes = mlp.AddRange(msgIds[0], msgIds[0]+uint64(distance)-1)
	}

	if yes {
		mlp.RebuildRangeString()
	}

	if err = tx.Model(mlp).Where("tid=? and mid=?", topicId, memberId).Update(mlp).Error; err != nil {
		logrus.Errorf("CANNOT get mlp: %v", err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		logrus.Errorf("CANNOT commit tx: %v", err)
		return
	}

	return
}

func (s *MsgDao) updateMsgLastPointsTx(tx *gorm.DB, in *models.Msg, toMembers ...uint64) (err error) {
	// insert into t_msg_last_pt (id, uid, topic_id, mid, tm) VALUES (default, uid, topicId, model.id, CURRENT_TIMESTAMP(6))
	//
	// id, uid, tid, midRead, mid, tm

	var (
		rows        int64
		mlps        []*models.MsgLastPt
		theUsersMap = make(map[uint64]bool)
	)

	if err = tx.Where("tid=? and mid in (?)", in.TopicId, toMembers).Select("mid").Find(&mlps).Error; err != nil {
		logrus.Errorf("CANNOT get MsgLastPt.mid list by topicId=%v: %v", in.TopicId, err)
		return
	}

	for _, to := range toMembers {
		theUsersMap[to] = true
	}
	for _, to := range mlps {
		delete(theUsersMap, to.MemberId)
	}

	// update the existed records.
	if err = tx.Table("t_msg_last_pt").Where("tid=?", in.TopicId).Update("mid_newest", in.Id).Error; err != nil {
		logrus.Errorf("CANNOT update MsgLastPt: %v, %v", rows, err)
		return
	}

	// and, insert the missed records
	for to := range theUsersMap {
		logrus.Printf("mlps - %v", to)
		if err = tx.Create(&models.MsgLastPt{MsgIdNewest: in.Id, TopicId: in.TopicId, MemberId: to}).Error; err != nil {
			logrus.Errorf("CANNOT insert MsgLastPt: %v, %v", rows, err)
			return
		}
	}
	return
}

func (s *MsgDao) findNewLeft(in *models.Msg, deviceId, parentId, topicId uint64) (newLeft int32, err error) {
	if in.TopicId == 0 {
		in.TopicId = topicId
	}

	if parentId > 0 {
		parent := &models.Msg{BaseModel: models.BaseModel{Id: parentId}}
		if err = dbe.GormDb().First(parent, parentId).Error; err != nil {
			logrus.Errorf("CANNOT locate the parent record %v: %v", parentId, err)
			return
		}
		newLeft = parent.Right
	} else if in.TopicId > 0 {
		type Rgt struct{ Rgt int32 }
		var rgt Rgt
		var tmp = &models.Msg{}
		if err = dbe.GormDb().Model(tmp).Select("max(rgt) as rgt").Where("topic_id=?", in.TopicId).Scan(&rgt).Error; err != nil {
			logrus.Errorf("CANNOT locate the max record in topic %v: %v", in.TopicId, err)
			return
		}
		newLeft = rgt.Rgt + 1
	} else {
		newLeft = 0
	}
	if newLeft == 0 {
		newLeft = 1
	}
	return
}

func (s *MsgDao) addToTopic(in *models.Msg, deviceId, parentId, topicId uint64) (toUsers []uint64, err error) {
	var (
		newLeft int32
		tmp     = &models.Msg{}
	)

	// topic mode

	if newLeft, err = s.findNewLeft(in, deviceId, parentId, topicId); err != nil {
		return
	}
	in.Left = newLeft
	in.Right = newLeft + 1

	tx := dbe.GormDb().Begin()
	defer TxRollback(tx, err)

	// sql := fmt.Sprintf("UPDATE %s SET rgt=rgt+2 WHERE rgt>=? AND topic_id=?", in.TableName())
	if err = tx.Model(tmp).Where("rgt>=? AND topic_id=?", newLeft, in.TopicId).Update("rgt", gorm.Expr("rgt+2")).Error; err != nil {
		logrus.Errorf("CANNOT update rgt: %v", err)
		return
	}
	if err = tx.Model(tmp).Where("lft>? AND topic_id=?", newLeft, in.TopicId).Update("lft", gorm.Expr("lft+2")).Error; err != nil {
		logrus.Errorf("CANNOT update lft: %v", err)
		return
	}

	if err = tx.Create(in).Error; err != nil {
		logrus.Errorf("CANNOT insert Msg: err=%v", err)
		return
	}

	// get all users in the topic
	toUsers = in.ToUsersArray()
	if len(toUsers) == 0 {
		// 	// TODO cache the list to avoid db querying
		toUsers, err = (NewTopicDao()).GetMemberPkIds(in.TopicId)
		if err != nil {
			logrus.Errorf("CANNOT GetMemberPkIds(topicId=%v): %v", in.TopicId, err)
			return
		}
	}
	// logrus.Printf("toMembers = %v", toUsers)

	// update the msgIdNewest for topic members
	if err = s.updateMsgLastPointsTx(tx, in, append(toUsers, in.FromUser)...); err != nil {
		logrus.Errorf("CANNOT updateMsgLastPoints: %v", err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		logrus.Errorf("CANNOT commit tx: %v", err)
		return
	}

	return
}

// for api.SendMsg:
// AddTo() append/insert an `in` msg into `t_msg` table.
//
// if in.TopicId == 0, AddTo() will create a Room/DirectMessage as an new Topic.
// Room is a special Topic, just like conversation group for multiple users/members.
// DirectMessage is a peer-to-peer Topic, it's special for two peers conversation.
//
func (s *MsgDao) AddTo(in *models.Msg, deviceId, parentId, orgId uint64) (toUsers []uint64, err error) {

	if in.TopicId == 0 {
		if len(in.ToUsers) == 0 {
			err = exception.New(exception.InvalidParams)
			return
		}

		// ------------ room/directMessage mode

		var topic *models.Topic
		dao := &TopicDao{}
		toUsers = in.ToUsersArray()
		if in.TopicId, topic, err = dao.CreateRoomOrDirectMessage(orgId, in.FromUser, toUsers...); err != nil {
			toUsers = nil
			logrus.Errorf("addto failed: %v, %v", err, topic)
			return
		}

		logrus.Debugf("Room/DirectMessage created as #%v, #%v: %v", in.TopicId, topic.Name, topic)
	}

	toUsers, err = s.addToTopic(in, deviceId, parentId, in.TopicId)
	return
}
