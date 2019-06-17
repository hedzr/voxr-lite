/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-lite/internal/dbe"
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
