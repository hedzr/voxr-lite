/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type OrgDao struct{ daoMember *MemberDao }

func NewOrgDao() *OrgDao {
	return &OrgDao{
		NewMemberDao(),
	}
}

func (s *OrgDao) GetById(id uint64) (ret *models.Organization, err error) {
	ret = new(models.Organization)
	if err = dbe.GormDb().First(&ret, id).Error; err != nil {
		return
	}
	// tmpl := &models.Organization{BaseModel: models.BaseModel{Id: id},}
	// var yes bool
	// yes, err = dbe.DBE.Engine().Get(tmpl)
	// if yes && err == nil {
	// 	ret = tmpl
	// } else {
	// 	logrus.Errorf("CANNOT get the record: %v", err)
	// }
	return
}

func (s *OrgDao) Get(tmpl *models.Organization) (ret *models.Organization, err error) {
	ret = new(models.Organization)
	if err = dbe.GormDb().First(&ret, tmpl).Error; err != nil {
		return
	}
	// var yes bool
	// yes, err = dbe.DBE.Engine().Get(tmpl)
	// if yes && err == nil {
	// 	ret = tmpl
	// } else {
	// 	logrus.Errorf("CANNOT get the record: %v", err)
	// }
	return
}

func (s *OrgDao) AddOrUpdate(in *models.Organization, parentId uint64) (ret *models.Organization, err error) {
	if parentId > 0 {
		cnt := 0
		tmpl := &models.Organization{BaseModel: models.BaseModel{Id: parentId}}
		if err = dbe.GormDb().Model(tmpl).Where(tmpl).Count(&cnt).Error; cnt != 1 || err != nil {
			logrus.Errorf("CANNOT locate the parent record %v: cnt=%d, %v", parentId, cnt, err)
		}
	}

	in.Pid = parentId

	if err = dbe.GormDb().Where("name=? and pid=?", in.Name, parentId).FirstOrCreate(in).Error; err != nil {
		logrus.Errorf("CANNOT insert the new record: %v", err)
		return
	}

	ret = in
	return
}

// Soft Delete
func (s *OrgDao) RemoveById(id uint64) (rows int64, err error) {
	if id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(&models.Organization{BaseModel: models.BaseModel{Id: id}}).Error; err != nil {
		return
	}
	return
}

func (s *OrgDao) Remove(tmpl *models.Organization) (rows int64, err error) {
	if tmpl.Id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(tmpl).Error; err != nil {
		return
	}
	return
}

func (s *OrgDao) UpdateAttrs(attrs ...interface{}) (err error) {
	err = dbe.GormDb().Model(&models.Organization{}).Update(attrs...).Error
	return
}

func (s *OrgDao) Update(in *models.Organization) (err error) {
	err = dbe.GormDb().Model(in).Omit("id").Updates(in).Error
	return
}

func (s *OrgDao) Where(query interface{}, args ...interface{}) *gorm.DB {
	return dbe.GormDb().Model(&models.Organization{}).Where(query, args)
}

func (s *OrgDao) List(limit, start int64, orderBy string, query interface{}, args ...interface{}) (ret []*models.Organization, err error) {
	err = dbe.GormDb().Where(query, args).Limit(limit).Offset(start).Order(orderBy).Find(&ret).Error
	return
}

// ListFromParentId 列出给定上级组织的全部下级组织
func (s *OrgDao) ListFromParentId(parentId uint64) (ret []*models.Organization, err error) {
	if err = dbe.GormDb().Where("pid=?", parentId).Find(&ret).Error; err != nil {
		return
	}
	return
}

// ListFromParent 列出给定上级组织的全部下级组织
func (s *OrgDao) ListFromParent(parent *models.Organization) (ret []*models.Organization, err error) {
	if err = dbe.GormDb().Where("pid=?", parent.Id).Find(&ret).Error; err != nil {
		return
	}
	return
}

func (s *OrgDao) AddTo1(in *models.Organization, parentId uint64) (err error) {
	// var yes bool
	//
	// tmpl := &models.Organization{BaseModel: models.BaseModel{Id: parentId,},}
	// yes, err = dbe.DBE.Engine().Exist(tmpl)
	// if err != nil {
	// 	logrus.Errorf("CANNOT locate the parent record %v: %v", parentId, err)
	// 	return
	// }
	//
	// tmpl = &models.Organization{BaseModel: models.BaseModel{Id: in.Id,},}
	// yes, err = dbe.DBE.Engine().Exist(tmpl)
	// if err == nil {
	// 	var rows int64
	// 	if yes {
	// 		in.Pid = parentId
	// 		rows, err = dbe.DBE.Engine().Update(in, tmpl)
	// 		if rows > 1 {
	// 			logrus.Errorf("multiple (%v) rows updated, it might be wrong.", rows)
	// 		} else if err != nil {
	// 			logrus.Errorf("CANNOT update the record: %v", err)
	// 		}
	// 	} else {
	// 		in.Pid = parentId
	// 		rows, err = dbe.DBE.Engine().Insert(in)
	// 		if rows == 0 || err != nil {
	// 			logrus.Errorf("CANNOT insert the new record: %v", err)
	// 		}
	// 		// s.addBotAndLink(in)
	// 	}
	// } else {
	// 	logrus.Errorf("CANNOT check exists for the record: %v", err)
	// }
	return
}

// load the lazy objects
func (s *OrgDao) LazyLoad(model *models.Organization) (ret *v10.Organization, err error) {
	var obj = model.ToProto()

	if obj.Pid > 0 {
		tmpl := &models.Organization{}
		if err = dbe.GormDb().First(tmpl, "pid=?", obj.Pid).Error; err != nil {
			logrus.Errorf("CANNOT locate the parent record %v: %v", obj.Pid, err)
			return
		}
		obj.ParentOrg = tmpl.ToProto()
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
	if err = dbe.GormDb().Find(&members, "org_id=?", obj.Id).Error; err != nil {
		logrus.Errorf("CANNOT load members of org %v: %v", obj.Id, err)
		return
	}
	for _, m := range members {
		if mx, err := s.daoMember.LazyLoad(m); err != nil {
			logrus.Errorf("CANNOT load member's userinfo of topic %v: %v", m.Id, err)
		} else {
			ret.Members = append(ret.Members, mx)
		}
	}

	// TODO load OrgClassified, OrgContact

	return
}

//
