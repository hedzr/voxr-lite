/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dbe

import "github.com/sirupsen/logrus"

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "sql", "file": v[1], "time": v[2]}).Print(v[3:]...)
	}
	if v[0] == "log" {
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "log", "file": v[1], "time": v[2]}).Print(v[2:]...)
	}
}
