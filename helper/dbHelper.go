package helper

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"os"
)

var DB *gorm.DB

func GetDbByHost(user string, password string, host string,)  (*gorm.DB, error){
	var mysql = fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host,
	)

	log.Print(mysql)
	DB, err := gorm.Open("mysql", mysql)
	if err != nil {
		return nil, err
	}
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	return DB, nil
}

func CreateConnection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	return GetDbByHost(user, password, host)
}
//
//func ConnectDB(mysql string, logLevel string) {
//	var err error
//	DB, err = gorm.Open("mysql", mysql)
//
//	//if err = DB.AutoMigrate(helper.Models...).Error; nil != err {
//	//	logger.Fatal("auto migrate tables failed: " + err.Error())
//	//}
//	if nil != err {
//		logger.Fatalf("opens client failed: " + err.Error())
//	}
//	if logLevel == "dev" {
//		DB.LogMode(true)
//	}else {
//		DB.LogMode(false)
//	}
//	DB.SingularTable(true)
//	DB.DB().SetMaxIdleConns(10)
//	DB.DB().SetMaxOpenConns(100)
//}
//
//func DisconnectDB() {
//	if err := DB.Close(); nil != err {
//		logger.Errorf("Disconnect from client failed: " + err.Error())
//	}
//}
