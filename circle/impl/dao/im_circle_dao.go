/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao

import (
	"github.com/hedzr/voxr-lite/circle/impl/models"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/jinzhu/gorm"
)

type CircleDao struct {
	// daoMember *MemberDao
}

func NewCircleDao() *CircleDao {
	return &CircleDao{
		// &MemberDao{},
	}
}

func (s *CircleDao) GetImageById(id uint64) (ret *models.CircleImage, err error) {
	ret = new(models.CircleImage)
	err = dbe.GormDb().First(ret, id).Error
	return
}

func (s *CircleDao) GetById(id uint64) (ret *models.Circle, err error) {
	ret = new(models.Circle)
	err = dbe.GormDb().First(&ret, id).Error
	return
}

func (s *CircleDao) Get(tmpl *models.Circle) (ret *models.Circle, err error) {
	ret = new(models.Circle)
	err = dbe.GormDb().First(&ret, tmpl).Error
	return
}

// Soft Delete
func (s *CircleDao) RemoveById(id uint64) (err error) {
	err = dbe.GormDb().Delete(&models.Circle{BaseModel: models.BaseModel{Id: id}}).Error
	return
}

func (s *CircleDao) Remove(tmpl *models.Circle) (err error) {
	if tmpl.Id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}
	err = dbe.GormDb().Delete(tmpl).Error
	return
}

func (s *CircleDao) Update(in *models.Circle) (err error) {
	err = dbe.GormDb().Model(in).Omit("id").Updates(in).Error
	return
}

func (s *CircleDao) Where(query interface{}, args ...interface{}) *gorm.DB {
	return dbe.GormDb().Model(&models.Circle{}).Where(query, args)
}

func (s *CircleDao) SearchUserFriends(userId uint64) (ret []uint64, err error) {
	err = dbe.GormDb().Exec(`SELECT c.uid FROM t_contact c
	INNER JOIN t_contact_relation cr ON cr.cid=c.id AND cr.deleted_at IS NULL
	INNER JOIN t_user uown ON cr.uid_owner=uown.id AND uown.deleted_at IS NULL
	WHERE uown.id=? AND c.deleted_at IS NULL`, userId).Scan(&ret).Error
	return
}

func (s *CircleDao) SaveImage(in *models.CircleImage) (err error) {
	err = dbe.GormDb().Model(in).Where("url=?", in.Url).FirstOrCreate(in).Error
	return
}

// func (s *CircleDao) Updates(attrs ...interface{}) (err error) {
// 	err = dbe.GormDb().Where(in).Update(attrs).Error
// 	return
// }

func (s *CircleDao) List(limit, start int64, orderBy string, query interface{}, args ...interface{}) (ret []*models.Circle, err error) {
	err = dbe.GormDb().Where(query, args).Limit(limit).Offset(start).Order(orderBy).Find(&ret).Error
	// err = dbe.DBE.Engine().Limit(limit, start).Find(&ret, condition...)
	return
}

func (s *CircleDao) ListBy(userId uint64, limit, start int64, orderBy, whereCondition string, whereArgs ...interface{}) (ret []*models.Circle, err error) {
	var (
		model = &models.Circle{}
	)

	if err = dbe.GormDb().Model(model).Limit(limit).Offset(start).Where(whereCondition, whereArgs...).Order(orderBy).Find(&ret).Error; err != nil {
		return
	}

	return
}

func (s *CircleDao) Add(in *models.Circle, parentId uint64) (err error) {
	if in.Pid == 0 && parentId != 0 {
		in.Pid = parentId
	}
	err = dbe.GormDb().Create(in).Error
	return
}
