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

type HookDao struct{}

func NewHookDao() *HookDao {
	return &HookDao{
		// NewMemberDao(),
	}
}

func (s *HookDao) GetById(id uint64) (ret *models.Hook, err error) {
	ret = new(models.Hook)
	if err = dbe.GormDb().First(&ret, id).Error; err != nil {
		logrus.Errorf("CANNOT get the record: %v", err)
		return
	}
	return
}

func (s *HookDao) Get(tmpl *models.Hook) (ret *models.Hook, err error) {
	ret = new(models.Hook)
	if err = dbe.GormDb().First(&ret, tmpl).Error; err != nil {
		logrus.Errorf("CANNOT get the record: %v", err)
		return
	}
	return
}

// not recommended
func (s *HookDao) AddOrUpdate1(in *models.Hook) (ret *models.Hook, err error) {
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
func (s *HookDao) RemoveById(id uint64) (rows int64, err error) {
	if id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(&models.Hook{BaseModel: models.BaseModel{Id: id}}).Error; err != nil {
		return
	}
	return
}

func (s *HookDao) Remove(tmpl *models.Hook) (rows int64, err error) {
	if tmpl.Id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(tmpl).Error; err != nil {
		return
	}
	return
}

func (s *HookDao) Update(in *models.Hook) (rows int64, err error) {
	err = dbe.GormDb().Model(in).Omit("id").Updates(in).Error
	return
}

func (s *HookDao) Where(query interface{}, args ...interface{}) *gorm.DB {
	return dbe.GormDb().Model(&models.Hook{}).Where(query, args)
}

func (s *HookDao) List(limit, start int64, orderBy string, query interface{}, args ...interface{}) (ret []*models.Hook, err error) {
	err = dbe.GormDb().Where(query, args).Limit(limit).Offset(start).Order(orderBy).Find(&ret).Error
	return
}
