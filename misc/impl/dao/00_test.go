/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-api/util"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"testing"
)

func prepareDb1() (db *dbe.DB, err error) {
	db = dbe.New()
	err = db.OpenUrl("mysql", "dev:123456@tcp(localhost:3306)/db_im?charset=utf8&parseTime=true")
	return
}

func prepareOldLocalDb() {
	// config.LocalOpen("dev:123456@tcp(localhost:3306)/db_im?charset=utf8&parseTime=true")
	prepareDb()
}

func localClose() {
	dbe.CloseDbConnection()
}

func prepareDb() (db *dbe.DB, err error) {
	// db = dbe.New()
	// err = db.OpenUrl("mysql", "dev:123456@tcp(localhost:3306)/db_im?charset=utf8&parseTime=true")
	err = dbe.OpenDbConnection()
	if err != nil {
		err = dbe.DBE.OpenUrl("mysql", "dev:123456@tcp(localhost:3306)/db_im?charset=utf8&parseTime=true")
	}
	if err == nil {
		db = dbe.DBE
	}
	return
}

func resetAllUsersPassword(t *testing.T, db *dbe.DB) {
	err := db.Engine().Where("channel = ? and (salt is null or salt = ?)", "im", "1234").Iterate(new(models.User), func(i int, bean interface{}) error {
		user := bean.(*models.User)
		t.Logf("    user = %v", user.LoginName)
		user.Password = "123456"
		user.Salt = "1234"
		user.Channel = "imx"
		if affected, err := db.Engine().Id(user.Id).Update(user); affected != 1 || err != nil {
			t.Fatal(err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func ensureUserContacts(t *testing.T, db *dbe.DB, user *models.User) {
	ensureUserContactGroups(t, db, user)
}

func ensureUserContactGroups(t *testing.T, db *dbe.DB, user *models.User) {
	var cc []*models.ContactGroup
	tmpl := &models.ContactGroup{UidOwner: user.Id}
	if err := db.Engine().Find(&cc, tmpl); err != nil {
		t.Fatalf("Err: %v", err)
	} else {
		if len(cc) == 0 {
			t.Logf("contact group not found")
			c := &models.ContactGroup{UidOwner: user.Id, Name: models.CGUnsorted}
			if rows, err := db.Engine().Insert(c); err != nil {
				t.Fatalf("Err: %v", err)
			} else {
				t.Logf("%d row(s) inserted: %v", rows, c)
			}
		} else {
			for ix, c := range cc {
				t.Logf("contact group found: %3d, %v", ix, c)
			}
		}
	}
}

func ensureUserRoles(t *testing.T, db *dbe.DB, user *models.User) {
	err := db.EnsureUserRoles(user.LoginName, []string{models.RoleUser, models.RoleImUser})
	if len(user.UniqueId) == 0 {
		user.UniqueId = util.UUid()
		if affected, err := db.Engine().Id(user.Id).Update(user); affected != 1 || err != nil {
			t.Fatal(err)
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}

func ensureAllUserRoles(t *testing.T, db *dbe.DB, user *models.User) {
	err := db.Engine().Where("1=1").Iterate(new(models.User), func(i int, bean interface{}) error {
		user := bean.(*models.User)
		err := db.EnsureUserRoles(user.LoginName, []string{models.RoleUser, models.RoleImUser})
		if len(user.UniqueId) == 0 {
			user.UniqueId = util.UUid()
			if affected, err := db.Engine().Id(user.Id).Update(user); affected != 1 || err != nil {
				t.Fatal(err)
			}
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
}

func ref1(t *testing.T, db *dbe.DB) {
	ui := &v10.UserInfo{
		Type: 1, Realname: "vvv2", Nickname: "vvv", Phone: "13000000000",
	}
	ok, err, obj := dao.UpdateUserInfo(ui)
	if !ok {
		t.Fatalf("update failed. %v", err)
	} else {
		t.Logf("user = %v | %v", ok, obj)
	}
}
