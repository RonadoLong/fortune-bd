package cache

import (
	"context"
	"fortune-bd/libs/env"
	"fortune-bd/libs/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"

)

var mongodb *mongo.Client

func Mongo() *mongo.Client {
	if mongodb == nil {
		mongodb = InitMongo(env.MongoAddr)
	}
	return mongodb
}

func InitMongo(addr string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() //养成良好的习惯，在调用WithTimeout之后defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		logger.Errorf("mongodb %s 连接失败 %v", addr, err)
		return nil
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	if err = client.Ping(ctx2, readpref.Primary()); err != nil {
		logger.Errorf("mongodb  %s 连接失败 %v", addr, err)
		return nil
	}
	return client
}
