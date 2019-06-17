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

type TopicAppDao struct{}

func NewTopicAppDao() *TopicAppDao {
	return &TopicAppDao{
		// NewMemberDao(),
	}
}

func (s *TopicAppDao) GetByIdEager(id uint64) (ret *models.TopicApp, err error) {
	ret = new(models.TopicApp)
	if err = dbe.GormDb().Preload("App").First(&ret, id).Error; err != nil {
		logrus.Errorf("CANNOT get the record: %v", err)
		return
	}
	return
}

func (s *TopicAppDao) GetById(id uint64) (ret *models.TopicApp, err error) {
	ret = new(models.TopicApp)
	if err = dbe.GormDb().First(&ret, id).Error; err != nil {
		logrus.Errorf("CANNOT get the record: %v", err)
		return
	}
	return
}

func (s *TopicAppDao) Get(tmpl *models.TopicApp) (ret *models.TopicApp, err error) {
	ret = new(models.TopicApp)
	if err = dbe.GormDb().First(&ret, tmpl).Error; err != nil {
		logrus.Errorf("CANNOT get the record: %v", err)
		return
	}
	return
}

func (s *TopicAppDao) AddOrUpdate(in *models.TopicApp) (ret *models.TopicApp, err error) {
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
func (s *TopicAppDao) RemoveById(id uint64) (rows int64, err error) {
	if id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(&models.TopicApp{BaseModel: models.BaseModel{Id: id}}).Error; err != nil {
		return
	}
	return
}

func (s *TopicAppDao) Remove(tmpl *models.TopicApp) (rows int64, err error) {
	if tmpl.Id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	if err = dbe.GormDb().Delete(tmpl).Error; err != nil {
		return
	}
	return
}

func (s *TopicAppDao) Update(in *models.TopicApp) (rows int64, err error) {
	err = dbe.GormDb().Model(in).Omit("id").Updates(in).Error
	return
}

func (s *TopicAppDao) Where(query interface{}, args ...interface{}) *gorm.DB {
	return dbe.GormDb().Model(&models.TopicApp{}).Where(query, args)
}

func (s *TopicAppDao) List(limit, start int64, orderBy string, query interface{}, args ...interface{}) (ret []*models.TopicApp, err error) {
	err = dbe.GormDb().Where(query, args).Limit(limit).Offset(start).Order(orderBy).Find(&ret).Error
	return
}

func (s *TopicAppDao) ListEager(query interface{}, args ...interface{}) (ret []*models.TopicApp, err error) {
	err = dbe.GormDb().Where(query, args).Preload("App").Find(&ret).Error
	return
}
