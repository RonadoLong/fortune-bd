package global

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
	"wq-fotune-backend/libs/logger"
)

//var DB *gorm.DB
var ticker = time.NewTicker(time.Minute)

func CreateConnectionByHost(host string) *gorm.DB {
	var err error
	DB, err := gorm.Open("mysql", host)
	if nil != err {
		log.Panic("opens database failed: " + err.Error())
	}
	//if config.Config.RuntimeMode == "dev"{
	//	DB.LogMode(true)
	//} else {
	//	DB.LogMode(false)
	//}
	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(20)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := DB.DB().Ping()
				logger.Err(err)
			}
		}
	}()
	return DB
}

func DisconnectDB(DB *gorm.DB) {
	if err := DB.Close(); nil != err {
		logger.Infof("Disconnect from client failed: " + err.Error())
	}
}
