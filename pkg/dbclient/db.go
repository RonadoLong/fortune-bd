package dbclient

import (
	"github.com/jinzhu/gorm"
	"sync"
	"wq-fotune-backend/libs/dbclient"
)

var DbClient *db
var once sync.Once

type db struct {
	DB *gorm.DB
}

func InitDbMysql(host string) {
	once.Do(func() {
		DbClient = &db{
			DB: dbclient.CreateConnectionByHost(host),
		}
	})
}

func NewDB(host string) *gorm.DB {
	return dbclient.CreateConnectionByHost(host)
}

func (d *db) CloseMysql() {
	dbclient.DisconnectDB(d.DB)
}
