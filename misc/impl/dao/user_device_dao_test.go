/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"testing"
)

func TestUserDeviceUpdate(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	di := &v10.Device{
		SystemOs: "IOS", SystemVersion: "IOS12", Model: "iPhoneXr",
		Nickname: "Tom的手x机", Unique: "sasjliuygt6789uiojkhgvfc^&rt5y67",
	}
	if obj, err := dao.SaveOrUpdateDevice(1, di); err != nil {
		t.Errorf("update failed: %v", err)
	} else {
		t.Logf("device updated: %v", obj)
	}

	if ret, err := dao.GetUserDdevices(1); err != nil {
		t.Errorf("GetUserDdevices failed: %v", err)
	} else {
		for ix, r := range ret {
			t.Logf("%5d: %v", ix, r)
		}
	}
}

func TestUserDeviceInsert(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	di := &v10.Device{
		SystemOs: "IOS", SystemVersion: "IOS12", Model: "iPhoneXr",
		Nickname: "Tom的手x机", Unique: "dashjdasdhkadsj^&rt5y67",
	}
	if obj, err := dao.SaveOrUpdateDevice(1, di); err != nil {
		t.Errorf("update failed: %v", err)
	} else {
		t.Logf("device updated: %v", obj)
	}

	if ret, err := dao.GetUserDdevices(1); err != nil {
		t.Errorf("GetUserDdevices failed: %v", err)
	} else {
		for ix, r := range ret {
			t.Logf("%5d: %v", ix, r)
		}
	}
}
