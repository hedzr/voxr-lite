/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dbe

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-lite/internal/config"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

type (
	DB struct {
		config  *config.DatabaseConfig
		db      *sqlx.DB
		xengine *xorm.Engine
		gormDb  *gorm.DB
	}
)

// var db *sqlx.DB
// var engine *xorm.Engine

func New() *DB {
	return &DB{}
}

func (db *DB) Open() (err error) {
	db.config = config.InitDatabaseConfig()
	var url string
	if len(db.config.Url) > 0 {
		url = fmt.Sprintf("%s:%s@%s", db.config.Username, db.config.Password, db.config.Url)
	} else {
		url = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", db.config.Username, db.config.Password, db.config.Hosts[0], db.config.Database)
	}
	return db.OpenUrl(db.config.Name, url)
}

func (db *DB) OpenUrl(driver, url string) (err error) {
	if db.xengine == nil {
		db.xengine, err = xorm.NewEngine(driver, url)
		if err == nil {
			db.xengine.SetMaxOpenConns(50)
			db.xengine.SetMapper(core.GonicMapper{})
			db.xengine.TZLocation, _ = time.LoadLocation("UTC") // "UTC") //"Asia/Shanghai"
		}
	}
	if db.gormDb == nil {
		db.gormDb, err = gorm.Open(driver, url)
		if err != nil {
			// panic("连接数据库失败")
			db.gormDb = nil
		} else {
			// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			// 	return "t_" + defaultTableName
			// }
			db.gormDb.SingularTable(true)

			// 自动迁移模式
			db.gormDb.AutoMigrate(&models.Organization{},
				&models.Topic{},
				&models.Member{},
				&models.MemberSettings{},
				&models.App{},
				&models.TopicApp{},
				&models.Filter{},
				&models.TopicFilter{},
				&models.Hook{},
				&models.Msg{},
				&models.MsgLastPt{},
				&models.MsgLinkedItems{},
				&models.Media{},
				&models.Setting{},
				&models.SettingTmpl{},
				&models.SettingTmplMap{})
		}
	}
	err = db.OpenXUrl(driver, url)

	db.set_log()
	return
}

func (db *DB) Close() {
	if db.xengine != nil {
		if err := db.xengine.Close(); err != nil {
			logrus.Errorf("CAN'T close database connection. %v", err)
		}
		db.xengine = nil
	}
	if db.gormDb != nil {
		if err := db.gormDb.Close(); err != nil {
			logrus.Errorf("CAN'T close database connection. %v", err)
		}
		db.gormDb = nil
	}
	_ = db.CloseX()
}

func (db *DB) set_log() {
	log := vxconf.GetBoolR("server.db.debug", false)

	if db.xengine != nil {
		db.xengine.ShowSQL(log)
	}

	if db.gormDb != nil {
		db.gormDb.LogMode(log)
		db.gormDb.SetLogger(&GormLogger{})
	}

	// logWriter, err := syslog.New(syslog.LOG_DEBUG, "rest-xorm-i-core")
	// if err != nil {
	// 	logrus.Fatalf("Fail to create xorm system logger: %v", err)
	// }
	//
	// logger := xorm.NewSimpleLogger(logWriter)
	// logger.ShowSQL(true)
	// db.xengine.SetLogger(logger)
}

func (db *DB) OpenX() (err error) {
	if db.db == nil {
		db.config = config.InitDatabaseConfig()
		url := fmt.Sprintf("%s:%s@%s", db.config.Username, db.config.Password, db.config.Url)
		err = db.OpenXUrl(db.config.Name, url)
	}
	return
}

func (db *DB) OpenXUrl(driver, url string) (err error) {
	if db.db == nil {
		db.db, err = sqlx.Open(driver, url)
		if err != nil {
			logrus.Errorf("CAN'T close database: %v", err)
			db.db = nil
		}
	}
	return
}

func (db *DB) CloseX() (err error) {
	if db.db != nil {
		err = db.db.Close()
		if err != nil {
			logrus.Errorf("CAN'T close database: %v", err)
		}
		db.db = nil
	}
	return
}

func (db *DB) Engine() *xorm.Engine {
	return db.xengine
}

// func (db *DB) GormDb() *gorm.DB {
// 	return db.gormDb
// }
