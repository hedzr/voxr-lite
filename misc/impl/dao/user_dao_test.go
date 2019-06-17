/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"testing"
	"time"
)

type Ux struct {
	Id              int64
	LoginName       string    `xorm:"varchar(64) notnull unique 'login_name'"`
	Mobile          string    `xorm:"varchar(64) null unique 'mobile'"`
	Email           string    `xorm:"varchar(64) null unique 'email'"`
	Nickname        string    `xorm:"varchar(64) null 'nickname'"`
	FullName        string    `xorm:"varchar(64) null 'full_name'"`
	Avatar          string    `xorm:"text"`
	Sex             int16     `xorm:"smallint null default -1"`
	Birthday        string    `xorm:"varchar(32) null"`
	Country         string    `xorm:"varchar(64) null"`
	Province        string    `xorm:"varchar(64) null"`
	City            string    `xorm:"varchar(64) null"`
	Lang            string    `xorm:"varchar(64) null default 'zh-cn'"`
	Tz              string    `xorm:"varchar(64) null default '+00:00'"`
	GivenName       string    `xorm:"varchar(64) null"`
	GivenMobile     string    `xorm:"varchar(64) null"`
	GivenSn         string    `xorm:"varchar(64) null"`
	Type            int       `xorm:"int null default 0"`
	CreatedAt       time.Time `xorm:"created"` // `xorm:"timestamp(6)"` // `xorm:"created"`
	UpdatedAt       time.Time `xorm:"updated"` // `xorm:"timestamp(6)"` // `xorm:"updated"`
	DeletedAt       time.Time `xorm:"deleted"` // `xorm:"datetime(6)"` // `xorm:"deleted"`
	Password        string    `xorm:"varchar(64) null"`
	Salt            string    `xorm:"varchar(64) null"`
	Ip              string    `xorm:"varchar(64) null"`
	Channel         string    `xorm:"varchar(64) null"`
	UniqueId        string    `xorm:"varchar(64) null"`
	PlatId          uint64    `xorm:"bigint null"`
	UserId          uint64    `xorm:"bigint null"`
	Openid          string    `xorm:"varchar(128) null"`
	Unionid         string    `xorm:"varchar(128) null"`
	Status          int       `xorm:"int not null default 1"`
	Blocked         bool      `xorm:"tinyint(1) not null default 0"`
	Forbidden       bool      `xorm:"tinyint(1) not null default 0"`
	Privilege       string    `xorm:"varchar(64) null"`
	JsonProfile     string    `xorm:"text null"`
	InvokeStatus    bool      `xorm:"tinyint(1) not null default 0"`
	Followers       string    `xorm:"text null"`
	Token           string    `xorm:"varchar(512) null"`
	Secret          string    `xorm:"varchar(512) null"`
	ExpiresIn       int       `xorm:"int not null"`
	ExpiresTime     time.Time `xorm:"datetime(6) null utc"`
	AccessStatus    bool      `xorm:"bit(1) not null default 0"`
	WebRouters      string    `xorm:"text null"`
	ApiRouters      string    `xorm:"text null"`
	DefaultWebRoute string    `xorm:"text null"`
}

func TestUx(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	if err = db.Engine().Sync(new(Ux)); err != nil {
		t.Fatal(err)
	}

	var rows int64
	rows, err = db.Engine().Unscoped().Where("login_name=?", "xxx").Delete(new(Ux))
	if err != nil {
		t.Fatal(err)
	}
	ux := &Ux{LoginName: "xxx"}
	rows, err = db.Engine().Insert(ux)
	t.Logf("rows = %d, err = %v, id = %v", rows, err, ux.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserLogin(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	// TODO dao.UserLoginByPhoneCode()
}

func TestUserRegister(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	// drop user, and register her

	err = db.UserDrop("vvv")
	if err != nil {
		t.Fatal(err)
	}

	ui := &v10.UserInfo{
		Type: 1, Realname: "vvv", Nickname: "vvv", Pass: "123456", Phone: "13000000000",
	}
	ui, err = dao.UserRegister(ui)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("new user id = %d", ui.Id)
	}
}

func TestUserUpdate(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	ui := &v10.UserInfo{
		Type: 2, Realname: "system robot", Phone: "13000000000",
		Sex:  models.SexUnspecified,
		Pass: models.SystemBotUnique,
	}
	ok, err, obj := dao.UpdateUserInfo(ui)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	} else {
		t.Logf("%v, user = %v", ok, obj)
	}
}

func TestGetUserInfoByUid(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	ui, err := dao.GetUserInfoByUid(models.SystemBotUnique)
	if err == nil {
		t.Log(ui)
	} else {
		t.Fatal(err)
	}
}
