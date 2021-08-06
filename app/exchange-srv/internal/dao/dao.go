package dao

import (
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
	"wq-fotune-backend/libs/cache"
)

type Dao struct {
	db    *gorm.DB
	mongo *mongo.Client
}

func New() *Dao {
	return &Dao{
		db:   cache.Mysql(),
		mongo: cache.Mongo(),
	}
}

var RowNotFoundErr = gorm.ErrRecordNotFound
