package cache

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
)
var db *gorm.DB

func Mysql() *gorm.DB {
	if db == nil {
		db = InitMysql(env.DbDSN)
	}
	return db
}

func InitMysql(host string) *gorm.DB {
	logger.Infof("mysql host: %s", host)
	var err error
	DB, err := gorm.Open("mysql", host)
	if nil != err {
		log.Panic("opens database failed: " + err.Error())
	}
	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(5)
	DB.DB().SetMaxOpenConns(10)
	return DB
}

func Close() {
	if db != nil {
		_ = db.Close()
	}
	if rdb != nil {
		_ = rdb.Close()
	}
}
