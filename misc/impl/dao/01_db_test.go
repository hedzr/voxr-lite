/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-api/models"
	"testing"
)

func TestDbOpen(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Errorf("open failed: %v", err)
	}
	defer localClose()

	// Get User Testing

	user := new(models.User)
	name := "Alice"
	ok, err := db.Engine().Alias("o").Where("o.login_name = ?", name).Get(user)
	if ok {
		t.Logf("[ok=%v] %v: %v", ok, name, *user)
	} else {
		t.Errorf("[ok=%v] %v ERR: %v", ok, name, err)
	}

	var roles []*models.Role
	roles, err = db.LoadRoles(user)
	if err != nil {
		t.Errorf("open failed: %v", err)
	}
	for ix, role := range roles {
		t.Logf("%5d. %v", ix, role)
	}

	// Password Testing

	// ret := user.encode("123456")
	// t.Logf("ret = %v | %d", ret, len(ret))
	//
	// matched := user.PasswordMatched("123456", ret)
	// t.Logf("matched: %v", matched)

	// Reset all user's passwords
	resetAllUsersPassword(t, db)

	// Enable IM roles

	// err = db.Engine().Where("1=1").Iterate(new(models.User), func(i int, bean interface{}) error {
	// 	user := bean.(*models.User)
	// 	err = db.EnsureUserRole(user.LoginName, "user,im.c.user")
	// 	if len(user.UniqueId) == 0 {
	// 		user.UniqueId = util.UUid()
	// 		if affected, err := db.Engine().Id(user.Id).Update(user); affected != 1 || err != nil {
	// 			t.Error(err)
	// 		}
	// 	}
	// 	return err
	// })
	// if err != nil {
	// 	t.Error(err)
	// }

	// M-to-M: 1. list and build

	// err = db.Engine().Where("1=1").Iterate(new(models.User), func(i int, bean interface{}) error {
	// 	user := bean.(*models.User)
	// 	t.Logf("    user = %v", user.LoginName)
	// 	user.Password = "123456"
	// 	user.Salt = "1234"
	// 	user.Channel = "imx"
	// 	if affected, err := db.Engine().Id(user.Id).Update(user); affected != 1 || err != nil {
	// 		t.Error(err)
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	t.Error(err)
	// }

	err = db.EnsureUserRole("Dave", "user,im.c.user")
	if err != nil {
		t.Error(err)
	}
	err = db.EnsureUserRole("Alice", "user,im.c.user")
	if err != nil {
		t.Error(err)
	}
	err = db.EnsureUserRole("Bob", "user,im.c.user")
	if err != nil {
		t.Error(err)
	}
	err = db.EnsureUserRole("Tom", "user,im.c.user")
	if err != nil {
		t.Error(err)
	}

	// M-to-M: 2. lists

	var urs = make([]models.UserRole, 0)
	err = db.Engine().Table(&models.UserRole{}).Alias("m").
		Join("LEFT", []string{(&models.User{}).TableName(), "u"}, "m.user_id=u.id").
		Join("LEFT", []string{(&models.Role{}).TableName(), "r"}, "m.role_id=r.id").
		Find(&urs)
	for ix, ur := range urs {
		t.Logf("%5d. %v", ix, ur)
	}

	// dump all

	// db.Engine().DumpAllToFile("all.sql")
}
