package handler

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CreateConnection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	var _ = fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s",
		host, user, dbName, password,
	)

	var sql = fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host,
	)
	return gorm.Open(dbName, sql)
}
