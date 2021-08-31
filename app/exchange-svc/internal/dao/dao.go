package dao

import (
	"fortune-bd/libs/cache"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
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
