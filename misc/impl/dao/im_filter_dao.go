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

type FilterDao struct{}

func NewFilterDao() *FilterDao {
	return &FilterDao{
		// NewMemberDao(),
	}
}

func (s *FilterDao) GetById(id uint64) (ret *models.Filter, err error) {
	ret = new(models.Filter)
	if err = dbe.GormDb().First(&ret, id).Error; err != nil {
		logrus.Errorf("CANNOT get the record: %v", err)
		return
	}
	return
}

func (s *FilterDao) Get(tmpl *models.Filter) (ret *models.Filter, err error) {
	ret = new(models.Filter)
	if err = dbe.GormDb().First(&ret, tmpl).Error; err != nil {
		logrus.Errorf("CANNOT get the record: %v", err)
		return
	}
	return
}

// not recommended
func (s *FilterDao) AddOrUpdate(in *models.Filter) (ret *models.Filter, err error) {
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
func (s *FilterDao) RemoveById(id uint64) (rows int64, err error) {
	if id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(&models.Filter{BaseModel: models.BaseModel{Id: id}}).Error; err != nil {
		return
	}
	return
}

func (s *FilterDao) Remove(tmpl *models.Filter) (rows int64, err error) {
	if tmpl.Id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(tmpl).Error; err != nil {
		return
	}
	return
}

func (s *FilterDao) Update(in *models.Filter) (rows int64, err error) {
	err = dbe.GormDb().Model(in).Omit("id").Updates(in).Error
	return
}

func (s *FilterDao) Where(query interface{}, args ...interface{}) *gorm.DB {
	return dbe.GormDb().Model(&models.Filter{}).Where(query, args)
}

func (s *FilterDao) List(limit, start int64, orderBy string, query interface{}, args ...interface{}) (ret []*models.Filter, err error) {
	err = dbe.GormDb().Where(query, args).Limit(limit).Offset(start).Order(orderBy).Find(&ret).Error
	return
}

func (s *FilterDao) ListFast(query interface{}, args ...interface{}) (ret []*models.Filter, err error) {
	err = dbe.GormDb().Where(query, args).Find(&ret).Error
	return
}
