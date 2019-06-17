/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dbe

import "github.com/jinzhu/gorm"

var DBE *DB

func OpenDbConnection() (err error) {
	if DBE == nil {
		DBE = New()
		err = DBE.Open()
	}
	return
}

func CloseDbConnection() {
	if DBE != nil {
		DBE.Close()
		DBE = nil
	}
}

func GormDb() *gorm.DB {
	return DBE.gormDb
}
