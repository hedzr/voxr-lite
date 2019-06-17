/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao

import (
	"fmt"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
)

type TopicDao struct{ daoMember *MemberDao }

func NewTopicDao() *TopicDao {
	return &TopicDao{
		NewMemberDao(),
	}
}

func (s *TopicDao) GetMemberUserIdsNotIn(topicId uint64, notInUsers []uint64) (ret []uint64, err error) {
	var res []*models.Member
	tmpl := &models.Member{TopicId: topicId}
	if err = dbe.GormDb().Model(tmpl).Select("user_id").Where("topic_id=?", topicId).Where("user_id not in (?)", notInUsers).Find(&res).Error; err != nil {
		return
	}

	ret = make([]uint64, len(res))
	for i, m := range res {
		ret[i] = m.UserId
	}
	return
}

func (s *TopicDao) GetMemberUserIds(topicId uint64) (ret []uint64, err error) {
	var res []*models.Member
	tmpl := &models.Member{TopicId: topicId}
	if err = dbe.GormDb().Model(tmpl).Select("user_id").Where("topic_id=?", topicId).Find(&res).Error; err != nil {
		return
	}

	ret = make([]uint64, len(res))
	for i, m := range res {
		ret[i] = m.UserId
	}
	return
}

func (s *TopicDao) GetMemberPkIds(topicId uint64) (ret []uint64, err error) {
	var res []*models.Member
	tmpl := &models.Member{TopicId: topicId}
	if err = dbe.GormDb().Model(tmpl).Select("id").Where("topic_id=?", topicId).Find(&res).Error; err != nil {
		return
	}

	ret = make([]uint64, len(res))
	for i, m := range res {
		ret[i] = m.Id
	}
	return
}

func (s *TopicDao) GetMembers(topicId uint64) (ret []models.Member, err error) {
	tmpl := &models.Member{TopicId: topicId}
	if err = dbe.GormDb().Model(tmpl).Where("topic_id=?", topicId).Find(&ret).Error; err != nil {
		return
	}
	// err = dbe.DBE.Engine().Table(tmpl).Find(&ret, tmpl)
	return
}

func (s *TopicDao) GetById(id uint64) (ret *models.Topic, err error) {
	ret = new(models.Topic)
	if err = dbe.GormDb().First(&ret, id).Error; err != nil {
		return
	}
	return
}

func (s *TopicDao) Get(tmpl *models.Topic) (ret *models.Topic, err error) {
	ret = new(models.Topic)
	if err = dbe.GormDb().First(&ret, tmpl).Error; err != nil {
		return
	}
	return
}

// Soft Delete
func (s *TopicDao) RemoveById(id uint64) (err error) {
	if id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(&models.Topic{BaseModel: models.BaseModel{Id: id}}).Error; err != nil {
		return
	}
	return
}

func (s *TopicDao) Remove(tmpl *models.Topic) (rows int64, err error) {
	if tmpl.Id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(tmpl).Error; err != nil {
		return
	}
	return
}

func (s *TopicDao) Update(in *models.Topic) (err error) {
	err = dbe.GormDb().Model(in).Omit("id").Updates(in).Error
	return
}

func (s *TopicDao) Where(query interface{}, args ...interface{}) *gorm.DB {
	return dbe.GormDb().Model(&models.Topic{}).Where(query, args)
}

func (s *TopicDao) List(limit, start int64, orderBy string, query interface{}, args ...interface{}) (ret []*models.Topic, err error) {
	err = dbe.GormDb().Where(query, args).Limit(limit).Offset(start).Order(orderBy).Find(&ret).Error
	return
}

// ListFromParentId 列出给定上级话题的全部下级话题
func (s *TopicDao) ListFromParentId(parentId uint64) (ret []*models.Topic, err error) {
	if err = dbe.GormDb().Where("pid=?", parentId).Find(&ret).Error; err != nil {
		logrus.Errorf("CANNOT list the record: %v", err)
		return
	}
	return
}

// ListFromParent 列出给定上级话题的全部下级话题
func (s *TopicDao) ListFromParent(parent *models.Topic) (ret []*models.Topic, err error) {
	if err = dbe.GormDb().Where("pid=?", parent.Id).Find(&ret).Error; err != nil {
		logrus.Errorf("CANNOT list the record: %v", err)
		return
	}
	return
}

// insert or update.
func (s *TopicDao) Add(in *models.Topic, parentId uint64) (err error) {
	if in.OrgId == 0 {
		in.OrgId = models.OrgIdPUB
	}

	if parentId > 0 {
		cnt := 0
		tmpl := &models.Topic{BaseModel: models.BaseModel{Id: parentId}, OrgId: in.OrgId}
		if err = dbe.GormDb().Model(tmpl).Where(tmpl).Count(&cnt).Error; cnt != 1 || err != nil {
			logrus.Errorf("CANNOT locate the parent record %v: cnt=%d, %v", parentId, cnt, err)
		}
	}

	in.Pid = parentId

	if err = dbe.GormDb().Where("name=? and pid=?", in.Name, parentId).FirstOrCreate(in).Error; err != nil {
		logrus.Errorf("CANNOT insert the new record: %v", err)
		return
	}

	return
}

// load the lazy objects
func (s *TopicDao) LazyLoad(model *models.Topic) (ret *v10.Topic, err error) {
	var obj = model.ToProto()

	if obj.Pid > 0 {
		tmpl := &models.Topic{}
		if err = dbe.GormDb().First(tmpl, "pid=?", obj.Pid).Error; err != nil {
			logrus.Errorf("CANNOT locate the parent record %v: %v", obj.Pid, err)
			return
		}
		obj.ParentTopic = tmpl.ToProto()
	}

	if obj.BotId > 0 {
		tmpl := &models.User{}
		if err = dbe.GormDb().First(tmpl, "id=?", obj.BotId).Error; err != nil {
			logrus.Errorf("CANNOT locate the bot record %v: %v", obj.BotId, err)
			return
		}
		obj.BotUser = tmpl.ToProto()
	}

	ret = obj

	// members
	var members []*models.Member
	if err = dbe.GormDb().Find(&members, "org_id=? AND topic_id=?", obj.OrgId, obj.Id).Error; err != nil {
		logrus.Errorf("CANNOT load members of org %v: %v", obj.OrgId, err)
		return
	}
	for _, m := range members {
		if mx, err := s.daoMember.LazyLoad(m); err != nil {
			logrus.Errorf("CANNOT load member's userinfo of topic %v: %v", m.Id, err)
		} else {
			ret.Members = append(ret.Members, mx)
		}
	}

	return
}

// Room/DirectMessage 没有专设bot，所以不要创建一个botUser。
// 但 FxBot 仍然自动属于一个 Room/DirectMessage，因此系统级的消息分发依然有效，而 Room 管理员无法通过自动化工具进行消息分发。
func (s *TopicDao) CreateRoomOrDirectMessage(orgId, moderator uint64, others ...uint64) (topicId uint64, out *models.Topic, err error) {
	var (
		in   *models.Topic
		pid  uint64
		name string
	)

	if len(others) == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}

	if orgId == 0 {
		orgId = models.OrgIdPUB
		pid = models.TopicIdPUBGeneral
	} else {
		pid = 0
	}

	in = &models.Topic{
		Name:  name,
		Pid:   pid,
		OrgId: orgId,
		Mode:  v10.Topic_DirectMessage,
	}

	a := append(others, moderator)
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	var sa []string
	for _, x := range a {
		sa = append(sa, strconv.FormatUint(x, 10))
	}

	if len(others) > 1 {
		in.Mode = v10.Topic_Room
		in.Name = fmt.Sprintf("room-%v", strings.Join(sa, "."))
	} else if len(others) == 1 {
		in.Mode = v10.Topic_DirectMessage
		in.Name = fmt.Sprintf("room-%v", strings.Join(sa, "."))
	}

	if err = s.Add(in, in.Pid); err != nil {
		logrus.Errorf("CreateRoomOrDirectMessage failed: %v", err)
		return
	}

	tx := dbe.GormDb().Begin()
	defer TxRollback(tx, err)

	if _, err = s.addUserAsMember(tx, in.OrgId, in.Id, moderator, true); err != nil {
		return
	}
	for _, r := range others {
		if _, err = s.addUserAsMember(tx, in.OrgId, in.Id, r, false); err != nil {
			return
		}
	}

	if err = tx.Commit().Error; err != nil {
		return
	}

	out = in
	topicId = out.Id
	return
}

func (s *TopicDao) addUserAsMember(tx *gorm.DB, orgId, topicId, userId uint64, moderator bool) (out *models.Member, err error) {
	out = &models.Member{
		OrgId:     orgId,
		TopicId:   topicId,
		UserId:    userId,
		Moderator: moderator,
	}

	if err = tx.FirstOrCreate(out, "org_id=? AND topic_id=? AND user_id=?", orgId, topicId, userId).Error; err != nil {
		return
	}

	type Name struct{ Nickname string }
	var name Name
	if err = tx.Model(&models.User{}).Select("nickname").Where("id=?", userId).Scan(&name).Error; err != nil {
		return
	}
	out.Name = name.Nickname
	if err = tx.Model(out).Update("name", out.Name).Error; err != nil {
		return
	}

	return
}
