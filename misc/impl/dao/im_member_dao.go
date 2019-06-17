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

type MemberDao struct{}

func NewMemberDao() *MemberDao {
	return &MemberDao{
		// NewMemberDao(),
	}
}

func (s *MemberDao) GetById(id uint64) (ret *models.Member, err error) {
	ret = new(models.Member)
	if err = dbe.GormDb().First(&ret, id).Error; err != nil {
		logrus.Errorf("CANNOT get the record: %v", err)
		return
	}
	return
}

func (s *MemberDao) Get(tmpl *models.Member) (ret *models.Member, err error) {
	ret = new(models.Member)
	if err = dbe.GormDb().First(&ret, tmpl).Error; err != nil {
		logrus.Errorf("CANNOT get the record: %v", err)
		return
	}
	return
}

// not recommended
func (s *MemberDao) AddOrUpdate1(in *models.Member) (ret *models.Member, err error) {
	if in.Id == 0 {
		// err = exception.New(exception.InvalidParams)
		if err = dbe.GormDb().Create(in).Error; err != nil {
			logrus.Errorf("CANNOT insert the new record: %v", err)
			return
		}
		return
	}

	if err = dbe.GormDb().Where("id=?", in.Id).FirstOrCreate(in).Error; err != nil {
		logrus.Errorf("CANNOT insert/update the record: %v", err)
		return
	}

	ret = in
	return
}

// Soft Delete
func (s *MemberDao) RemoveById(id uint64) (rows int64, err error) {
	if id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(&models.Member{BaseModel: models.BaseModel{Id: id}}).Error; err != nil {
		return
	}
	return
}

func (s *MemberDao) Remove(tmpl *models.Member) (rows int64, err error) {
	if tmpl.Id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(tmpl).Error; err != nil {
		return
	}
	return
}

func (s *MemberDao) Update(in *models.Member) (rows int64, err error) {
	err = dbe.GormDb().Model(in).Omit("id").Updates(in).Error
	return
}

func (s *MemberDao) Where(query interface{}, args ...interface{}) *gorm.DB {
	return dbe.GormDb().Model(&models.Member{}).Where(query, args)
}

func (s *MemberDao) List(limit, start int64, orderBy string, query interface{}, args ...interface{}) (ret []*models.Member, err error) {
	err = dbe.GormDb().Where(query, args).Limit(limit).Offset(start).Order(orderBy).Find(&ret).Error
	return
}

func (s *MemberDao) LazyLoad(model *models.Member) (ret *v10.Member, err error) {
	ret = model.ToProto()
	if model.UserId > 0 {
		var user = &models.User{}
		if err = dbe.GormDb().First(user, model.UserId).Error; err != nil {
			logrus.Errorf("CANNOT load userinfo for member: %v", err)
			return
		}
		ret.Ui = user.ToProto()
	}
	return
}
