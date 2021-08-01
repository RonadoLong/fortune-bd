package mongoClient

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
	"wq-fotune-backend/libs/logger"
)

var ticker = time.NewTicker(time.Minute)

func InitMongo(addr string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() //养成良好的习惯，在调用WithTimeout之后defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		logger.Errorf("mongodb %s 连接失败 %v", addr, err)
		return nil, err
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	if err = client.Ping(ctx2, readpref.Primary()); err != nil {
		logger.Errorf("mongodb  %s 连接失败 %v", addr, err)
		return nil, err
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		for {
			select {
			case <-ticker.C:
				err2 := client.Ping(ctx, readpref.Primary())
				if err2 != nil {
					logger.Err(err2)
				}
			}
		}
	}()
	return client, nil
}
