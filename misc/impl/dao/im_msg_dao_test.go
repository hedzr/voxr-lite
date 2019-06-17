/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"context"
	"fmt"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"github.com/hedzr/voxr-lite/misc/impl/service"
	"testing"
	"time"
)

type Product struct {
	models.BaseModel
	Code  string
	Price uint
}

type Aaa struct {
	Id        int64
	Name      string
	CreatedAt time.Time `xorm:"notnull timestampz default CURRENT_TIMESTAMP(6) created"`
	UpdatedAt time.Time `xorm:"notnull timestampz"`
}

func (r *Aaa) BeforeUpdate() {
	r.UpdatedAt = time.Now()
	fmt.Printf("updated_at -> %v\n", r.UpdatedAt)
}

func (r *Aaa) BeforeInsert() {
	r.UpdatedAt = time.Now()
	fmt.Printf("updated_at -> %v\n", r.UpdatedAt)
}

// func (*Product) TableName() string {
// 	return "hz_prod"
// }

func TestAaa(t *testing.T) {
	// db, err := prepareDb()
	// if err != nil {
	// 	t.Fatalf("open failed: %v", err)
	// }
	// defer db.Close()
	//
	// a := &models.Aaa{Name: "jdsfdl"}
	//
	// err1 := db.Engine().Sync(a)
	// if err1 != nil {
	// 	t.Fatalf("err: %v", err1)
	// }
	//
	// rows, err := db.Engine().Insert(a)
	// if err != nil {
	// 	t.Fatalf("err: %v", err)
	// }
	// if rows != 1 {
	// 	t.Fatalf("rows!=1: %v", rows)
	// }

	// db, err := gorm.Open("sqlite3", "test.db")
	// db, err := gorm.Open("mysql", "dev:123456@tcp(localhost:3306)/db_im?charset=utf8&parseTime=true")
	// if err != nil {
	// 	panic("连接数据库失败")
	// }
	// defer db.Close()

	_, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer localClose()

	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return "t_" + defaultTableName
	// }
	dbe.GormDb().SingularTable(true)

	// dbe.GormDb().AutoMigrate(&models.Organization{},
	// 	&models.Topic{},
	// 	&models.Member{},
	// 	&models.MemberSettings{},
	// 	&models.App{},
	// 	&models.TopicApp{},
	// 	&models.Filter{},
	// 	&models.TopicFilter{},
	// 	&models.Hook{},
	// 	&models.Msg{},
	// 	&models.MsgLastPt{},
	// 	&models.MsgLinkedItems{},
	// 	&models.Media{},
	// 	&models.Setting{},
	// 	&models.SettingTmpl{},
	// 	&models.SettingTmplMap{})

	// var parentId uint64 = 1
	// parent := &models.Msg{BaseModel: models.BaseModel{Id: parentId,},}
	// if err = dbe.GormDb().First(parent, parentId).Error; err != nil {
	// 	logrus.Errorf("CANNOT locate the parent record %v: %v", parentId, err)
	// 	return
	// }
	// t.Logf("p = %v", parent)

	// // 自动迁移模式
	// dbe.GormDb().AutoMigrate(&models.Organization{}, &models.MsgLastPt{}, &Product{})
	//
	// 创建
	// dbe.GormDb().Create(&Product{Code: "L1212", Price: 1000})
	//
	// // 读取
	// var product Product
	// dbe.GormDb().First(&product, 1)                   // 查询id为1的product
	// dbe.GormDb().First(&product, "code = ?", "L1212") // 查询code为l1212的product
	//
	// // 更新 - 更新product的price为2000
	// dbe.GormDb().Model(&product).Update("Price", 2000)
	//
	// // 删除 - 删除product
	// dbe.GormDb().Delete(&product)

	d := &dao.MsgDao{}
	if ret, err := d.GetById(54); err != nil {
		t.Logf("ret = %v", ret)
		t.Fatal(err)
	} else {
		t.Logf("ret = %v", ret)
	}

	if ret, err := d.ListFromParentId(54, 14); err != nil {
		t.Fatal(err)
	} else {
		for ix, r := range ret {
			t.Logf("%5d. %v", ix, r)
		}
	}
}

func TestMsgDao_AddToTest(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer localClose()

	err = db.Engine().Sync(
		&models.MsgLastPt{}, &models.Msg{},
		&models.Member{}, &models.MemberSettings{})
	if err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	doAddTo(t, db)
}

func doAddTo(t *testing.T, db *dbe.DB) {
	o := &dao.MsgDao{}

	var err error
	var toUsers []uint64
	var parentId uint64 = 3
	var topicId uint64 = 1
	var orgId uint64 = 1
	in := &models.Msg{
		Content:  "sample msg .",
		FromUser: 7,
		TopicId:  topicId,
	}
	toUsers, err = o.AddTo(in, parentId, topicId, orgId)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("msg inserted: %v", in)
	t.Logf("toUsers: %v", toUsers)
}

func TestMsgDao_AddTo(t *testing.T) {
	_, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer localClose()

	// doAddTo(t, db)

	req := &v10.GetMsgReqV12{
		ProtoOp: v10.Op_MsgsAll, Seq: 9, TopicId: 1, MsgId: 90, Newer: false, SortByAsc: false,
	}
	svc := &service.ImMsgService{}
	res, err := svc.GetMsg(context.TODO(), req)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	t.Logf("list of msgs (nextMsgId = %v)", res.NextMsgId)
	for ix, r := range res.Msgs {
		t.Logf("msg #%5d: %v", ix, *r)
	}

	// p := fmt.Println

	tm := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	// tm = time.Unix(0, -6795364578871345152)

	t.Logf("%v - %v, %v", tm, tm.Unix(), tm.UnixNano())
	t1, e := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	t.Log(t1)

	if tm1, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", "1754-08-30 22:43:41.128654848 +0000 UTC"); err != nil {
		t.Errorf("Err: %v", err)
	} else {
		fmt.Printf("tm1 = %v, %v", tm1, tm1.UnixNano())
	}

	t.Log(e)
}

func TestMsgDao_LoadApp(t *testing.T) {
	_, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer localClose()

	// x
	dx := dao.NewTopicAppDao()
	if r, err := dx.GetByIdEager(1); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("topicApp: %v", r)
	}

	if r, err := dx.ListEager("1=1"); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("results: %v", len(r))
		for ix, rt := range r {
			t.Logf("%5d. %v", ix, rt)
		}
	}

}
