/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func SaveOrUpdateDevice(uid int64, device *v10.Device) (dur *v10.DeviceUserRelation, err error) {
	var (
		dl      []*models.UserDevice
		d, tmpl *models.UserDevice
		rows    int64
		s       string
	)

	tmpl = &models.UserDevice{Uid: uid, Unique: device.Unique}
	if err = dbe.DBE.Engine().Find(&dl, tmpl); err != nil {
		logrus.Errorf("SaveOrUpdateDevice error: %v", err)
		return
	}

	if len(dl) > 0 {
		d = &models.UserDevice{}
		dx := dl[0]
		u := false
		if !strings.EqualFold(device.SystemOs, dx.SystemOs) {
			d.SystemOs = device.SystemOs
			u = true
		}
		if !strings.EqualFold(device.SystemVersion, dx.SystemVersion) {
			d.SystemVersion = device.SystemVersion
			u = true
		}
		if !strings.EqualFold(device.Model, dx.Model) {
			d.Model = device.Model
			u = true
		}
		if !strings.EqualFold(device.Nickname, dx.Nickname) {
			d.Nickname = device.Nickname
			u = true
		}
		if u {
			rows, err = dbe.DBE.Engine().Update(d, tmpl)
		}
		s = "updated"
	} else {
		d = &models.UserDevice{
			Uid: uid, SystemOs: device.SystemOs, SystemVersion: device.SystemVersion,
			Model: device.Model, Nickname: device.Nickname,
			Unique:    device.Unique,
			CreatedAt: time.Now(),
		}
		rows, err = dbe.DBE.Engine().Insert(d)
		s = "inserted"
	}

	if rows > 1 || err != nil {
		logrus.Errorf("SaveOrUpdateDevice error: %v | rows = %v", err, rows)
	} else {
		dur = &v10.DeviceUserRelation{
			Id: d.Id, Uid: d.Uid,
			Device: &v10.Device{
				SystemOs: d.SystemOs, SystemVersion: d.SystemVersion,
				Model: d.Model, Nickname: d.Nickname,
				Unique: d.Unique,
			},
			CreateTime: d.CreatedAt.Unix(),
		}
		logrus.Debugf("SaveOrUpdateDevice %s: %d | ret = %v", s, rows, dur)
	}
	return
}

func GetUserDdevices(uid int64) (ret []*v10.DeviceUserRelation, err error) {
	var (
		dl []*models.UserDevice
		d  *models.UserDevice
	)

	d = &models.UserDevice{Uid: uid}
	if err = dbe.DBE.Engine().Find(&dl, d); err != nil {
		logrus.Errorf("GetUserDdevices error: %v", err)
		return
	}

	if len(dl) > 0 {
		for _, d := range dl {
			dur := &v10.DeviceUserRelation{
				Id: d.Id, Uid: d.Uid,
				Device: &v10.Device{
					SystemOs: d.SystemOs, SystemVersion: d.SystemVersion,
					Model: d.Model, Nickname: d.Nickname,
					Unique: d.Unique,
				},
				CreateTime: dl[0].CreatedAt.Unix(),
			}
			ret = append(ret, dur)
		}
		logrus.Debugf("GetUserDdevices ret = %v", ret)
	}
	return
}

func IsExistDevice(uid int64, deviceId string) (yes bool, dur *v10.DeviceUserRelation) {
	var (
		dl []*models.UserDevice
		d  *models.UserDevice
	)

	d = &models.UserDevice{Uid: uid, Unique: deviceId}
	if err := dbe.DBE.Engine().Find(&dl, d); err != nil {
		logrus.Errorf("GetUserDdevices error: %v", err)
		return
	}

	if len(dl) > 0 {
		yes = true
		dur = &v10.DeviceUserRelation{
			Id: d.Id, Uid: d.Uid,
			Device: &v10.Device{
				SystemOs: d.SystemOs, SystemVersion: d.SystemVersion,
				Model: d.Model, Nickname: d.Nickname,
				Unique: d.Unique,
			},
			CreateTime: dl[0].CreatedAt.Unix(),
		}
	}
	return
}
