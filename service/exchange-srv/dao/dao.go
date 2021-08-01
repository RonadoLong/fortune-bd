package dao

import (
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/mongoClient"
	"wq-fotune-backend/pkg/dbclient"
)

type Dao struct {
	db    *gorm.DB
	mongo *mongo.Client
}

func New() *Dao {
	mgoClient, err := mongoClient.InitMongo(env.MongoAddr)
	if err != nil {
		os.Exit(-1)
	}
	return &Dao{
		db:    dbclient.NewDB(env.DBDSN),
		mongo: mgoClient,
	}
}

var RowNotFoundErr = gorm.ErrRecordNotFound
