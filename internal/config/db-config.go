/*
 * Copyright © 2019 Hedzr Yeh.
 */

package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
)

type DatabaseConfig struct {
	Name              string   `yaml:"-"`
	Username          string   `yaml:"username"`
	Password          string   `yaml:"password"`
	Hosts             []string `yaml:"hosts"`
	Database          string   `yaml:"database"`
	Desc              string   `yaml:"desc"`
	Url               string   `yaml:"url"`
	ConnectionTimeout int64    `yaml:"connectionTimeout"`
	MaxOpenConns      int      `yaml:"maxOpenConns"`
	MaxIdleConns      int      `yaml:"maxIdleConns"`
}

func InitDatabaseConfig() *DatabaseConfig {
	// return &InitConfig().Db

	prefix := "server.pub.deps.db"
	backend := vxconf.GetStringRP(prefix, "backend", "mysql")
	env := vxconf.GetStringRP(prefix, "env", vxconf.GetStringR("runmode", "devel"))
	keypath := fmt.Sprintf("%v.backends.%s.%s", prefix, backend, env)

	thisConfig := &DatabaseConfig{}
	vxconf.LoadSectionTo(keypath, thisConfig)
	thisConfig.Name = backend
	logrus.Debugf("DB config Got #1: %v\n\n", thisConfig)

	return thisConfig
}

// var MyDb *sql.DB

// // for testing
// func LocalOpen(url string) {
// 	// myDb, _ := sql.Open("mysql", url)
// 	// MyDb = myDb
//
// 	// db, err := gorm.Open("mysql", "dev:123456@tcp(localhost:3306)/db_im?charset=utf8&parseTime=true")
// 	// if err != nil {
// 	// 	panic("连接数据库失败")
// 	// }
// 	// GormDB = db
// }
//
// // for testing
// func LocalClose() {
// 	if MyDb != nil {
// 		MyDb.Close()
// 		MyDb = nil
// 	}
//
// 	// if GormDB != nil {
// 	// 	GormDB.Close()
// 	// 	GormDB = nil
// 	// }
// }

// var GormDB *gorm.DB

//
func InitDBConn() {
	// dbconfig := InitDatabaseConfig()
	//
	// // if GormDB == nil {
	// // 	db, err := gorm.Open("mysql", "dev:123456@tcp(localhost:3306)/db_im?charset=utf8&parseTime=true")
	// // 	if err != nil {
	// // 		panic("连接数据库失败")
	// // 	}
	// // 	GormDB = db
	// // }
	//
	// if MyDb == nil {
	// 	username := dbconfig.Username
	// 	url := dbconfig.Url
	// 	pass := dbconfig.Password
	// 	connTimeOut := dbconfig.ConnectionTimeout
	// 	maxOpenConns := dbconfig.MaxOpenConns
	// 	maxIdleCoons := dbconfig.MaxIdleConns
	// 	myDb, _ := sql.Open("mysql", username+":"+pass+"@"+url)
	// 	myDb.SetMaxIdleConns(maxIdleCoons)
	// 	myDb.SetMaxOpenConns(maxOpenConns)
	// 	myDb.SetConnMaxLifetime(time.Duration(connTimeOut))
	// 	MyDb = myDb
	// }
}
