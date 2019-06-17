/*
 * Copyright © 2019 Hedzr Yeh.
 */

package models

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/dc"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	BxModel struct {
		Id        uint64    `gorm:"primary_key" xorm:"pk autoincr"`
		CreatedAt time.Time `gorm:"type:timestamp(6)" xorm:"created"`
		UpdatedAt time.Time `gorm:"type:timestamp(6)" xorm:"updated"`
	}

	BaseModel struct {
		Id        uint64     `gorm:"primary_key" xorm:"pk autoincr"`
		CreatedAt time.Time  `gorm:"type:timestamp(6)" xorm:"created"`
		UpdatedAt time.Time  `gorm:"type:timestamp(6)" xorm:"updated"`
		DeletedAt *time.Time `gorm:"type:timestamp(6)" xorm:"deleted" sql:"index"`
	}

	Circle struct {
		BaseModel
		Pid     uint64 `gorm:"not null;default:0" xorm:"bigint notnull default 0"`
		UserId  uint64 `gorm:"not null;default:0" xorm:"bigint notnull default 0"`
		Header  string `gorm:"null;type:varchar(128)" xorm:"varchar(128) null"`
		Title   string `gorm:"null;type:varchar(128)" xorm:"varchar(128) null"`
		Content string `gorm:"null;type:text" xorm:"text null"`
		Footer  string `gorm:"null;type:text" xorm:"text null"`
		HeadUrl string `gorm:"null;type:text" xorm:"text null"`
		Remarks string `gorm:"null;type:text" xorm:"text null"`
	}

	CircleImage struct {
		BaseModel
		CircleId  uint64 `gorm:"not null;default:0;column:cid" xorm:"bigint notnull default 0 'cid'"`
		UserId    uint64 `gorm:"not null;default:0;column:uid" xorm:"bigint notnull default 0 'uid'"`
		BaseName  string `gorm:"null;type:text" xorm:"text null"`
		Mime      string `gorm:"null;type:text" xorm:"text null"`
		Size      int64  `gorm:"not null;default:0" xorm:"bigint notnull default 0"`
		LocalPath string `gorm:"null;type:text" xorm:"text null"`
		Url       string `gorm:"null;type:text" xorm:"text null"`
	}
)

func (Circle) TableName() string {
	return "t_circle"
}

func (CircleImage) TableName() string {
	return "t_circle_image"
}

func (r *Circle) FromProto(pb *v10.Circle) *Circle {
	if err := dc.StandardCopier.Copy(r, pb); err != nil {
		logrus.Errorf("FromProto failed: %v", err)
		return nil
	}
	if pb.DeletedAt == 0 {
		r.DeletedAt = nil
	}
	return r
}

func (r *Circle) ToProto() (ret *v10.Circle) {
	if err := dc.StandardCopier.Copy(ret, r); err != nil {
		logrus.Errorf("ToProto failed: %v", err)
		return nil
	}
	return
}
