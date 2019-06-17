/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-lite/circle/impl/dao"
	"github.com/hedzr/voxr-lite/circle/impl/models"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"testing"
)

func TestAaa(t *testing.T) {
	_, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer localClose()

	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return "t_" + defaultTableName
	// }
	dbe.GormDb().SingularTable(true)

	if err = dbe.GormDb().AutoMigrate(&models.Circle{}, &models.CircleImage{}).Error; err != nil {
		t.Fatal(err)
	}

	dx := dao.NewCircleDao()
	in := &models.Circle{Title: "hello", Header: "header", Footer: "footer", Content: "xxx world", Remarks: "remarks", UserId: 3, HeadUrl: "cover"}
	if err := dx.Add(in, 1); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("ret = %v", in)
	}

	if err := dx.SaveImage(&models.CircleImage{CircleId: 2, UserId: 3, BaseName: "vx.png", Mime: "image/png", Size: 3, LocalPath: "vx.png", Url: "vx.png"}); err != nil {
		t.Fatal(err)
	} else {
		//
	}

	if obj, err := dx.GetById(1); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("obj = %v", obj)
	}

	if obj, err := dx.Get(&models.Circle{BaseModel: models.BaseModel{Id: 1}}); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("obj = %v", obj)
	}

	in.Header = "hoho"
	if err := dx.Update(in); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("in.updated = %v", in)
	}

	if ret, err := dx.List(20, 0, "", ""); err != nil {
		t.Fatal(err)
	} else {
		for ix, r := range ret {
			t.Logf("%5d. %v", ix, r)
		}
	}

	if err := dx.Remove(in); err != nil {
		t.Fatal(err)
	} else {
		t.Log("removed")
	}
}
