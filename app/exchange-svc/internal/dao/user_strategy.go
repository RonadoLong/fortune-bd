package dao

import (
	"context"
	"fortune-bd/app/exchange-svc/internal/model"
	"fortune-bd/libs/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

)



func newCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

// GetUserStrategyOfRun 如果传了bson.m 就用bson.m  查询
func (d *Dao) GetUserStrategyOfRun(sql bson.M) []*model.GridStrategy {
	var userStrategyList []*model.GridStrategy
	ctx, cancelFunc := newCtx()
	defer cancelFunc()
	bsonSearch := bson.M{"isRun": true}
	if sql != nil {
		bsonSearch = sql
	}
	find, err := d.mongo.Database(DATABASE).Collection(TABLE_MONGO_GRIDSTRATEGY).Find(ctx, bsonSearch)
	if err != nil {
		logger.Errorf("GetUserStrategyOfRun has err %v", err)
		return userStrategyList
	}
	defer find.Close(ctx)
	for find.Next(ctx) {
		var strategy *model.GridStrategy
		if err = find.Decode(&strategy); err != nil {
			logger.Warnf("GetUserStrategyOfRun has err %v", err)
			return userStrategyList
		}
		userStrategyList = append(userStrategyList, strategy)
	}
	if err := find.Err(); err != nil {
		logger.Warnf("GetUserStrategyOfRun has err %v", err)
		return userStrategyList
	}
	return userStrategyList
}

func (d *Dao) GetUserStrategy(uid, strategyId string) (*model.GridStrategy, error) {
	strategy := &model.GridStrategy{}
	ctx, cancelFunc := newCtx()
	defer cancelFunc()
	objID, _ := primitive.ObjectIDFromHex(strategyId)
	err := d.mongo.Database(DATABASE).Collection(TABLE_MONGO_GRIDSTRATEGY).FindOne(ctx, bson.M{"_id": objID}).Decode(strategy)
	if err != nil {
		logger.Warnf("GetUserStrategy has err %v", err)
		return nil, err
	}
	return strategy, nil

}


func (d *Dao) GetUserStrategyByApiKey(uid, apiKey string) []*model.GridStrategy {

	var userStrategyList []*model.GridStrategy
	ctx, cancelFunc := newCtx()
	defer cancelFunc()
	find, err := d.mongo.Database(DATABASE).Collection(TABLE_MONGO_GRIDSTRATEGY).Find(ctx, bson.M{"uid": uid, "isRun": true, "apiKey": apiKey})
	if err != nil {
		logger.Warnf("GetUserStrategyByApiKey has err %v", err)
		return userStrategyList
	}
	defer find.Close(ctx)
	for find.Next(ctx) {
		var strategy *model.GridStrategy
		if err = find.Decode(&strategy); err != nil {
			logger.Warnf("GetUserStrategyByApiKey has err %v", err)
			return userStrategyList
		}
		userStrategyList = append(userStrategyList, strategy)
	}
	if err := find.Err(); err != nil {
		logger.Warnf("GetUserStrategyByApiKey has err %v", err)
		return userStrategyList
	}
	return userStrategyList
}

func (d *Dao) SetUserAllStrategyApi(uid, apiKey, platform string) error {
	ctx, cancelFunc := newCtx()
	defer cancelFunc()
	many, err := d.mongo.Database(DATABASE).Collection(TABLE_MONGO_GRIDSTRATEGY).
		UpdateMany(ctx, bson.M{"uid": uid, "exchange": platform}, bson.M{"$set": bson.M{"apiKey": apiKey}})
	if err != nil {
		return err
	}
	if many.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
