package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/app/exchange-srv/model"
)

//func (d *Dao) GetUserStrategyList(uid string) (userStrategyList []*model.WqUserStrategy) {
//	field := "id , parent_strategy_id, user_id, group_id, strategy_id, api_key, platform, symbol, balance, state, created_at, updated_at"
//	if err := d.db.Table(TABLE_WQ_USER_STRATEGY).Select(field).Where("user_id=? and state=1", uid).Find(&userStrategyList).
//		Error; err != nil {
//		logger.Errorf("GetUserStrategyList has err %v", err)
//		return
//	}
//	return
//}

func newCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

// 如果传了bson.m 就用bson.m  查询
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
	//if err := d.db.Table(TABLE_WQ_USER_STRATEGY).Where("user_id=? and strategy_id=?", uid, strategyId).First(&strategy).
	//	Error; err != nil {
	//	logger.Errorf("GetUserStrategy has err %v", err)
	//	return nil, err
	//}
	//return strategy, nil
}

//func (d *Dao) GetUserStrategyByParentID(uid string, parentId int64) (*model.WqUserStrategy, error) {
//	strategy := &model.WqUserStrategy{}
//	if err := d.db.Table(TABLE_WQ_USER_STRATEGY).Where("user_id=? and parent_strategy_id=?", uid, parentId).First(&strategy).
//		Error; err != nil {
//		logger.Errorf("GetUserStrategyByParentID has err %v", err)
//		return nil, err
//	}
//	return strategy, nil
//}

func (d *Dao) GetUserStrategyByApiKey(uid, apiKey string) []*model.GridStrategy {
	//var strategy []*model.GridStrategy
	//if err := d.db.Table(TABLE_WQ_USER_STRATEGY).Where("user_id=? and api_key=? and state = 1", uid, apiKey).Find(&strategy).
	//	Error; err != nil {
	//	logger.Errorf("GetUserStrategyByApiKey has err %v", err)
	//	return nil, err
	//}
	//return strategy, nil

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

//func (d *Dao) SetUserStrategyApi(uid, strategyId, apiKey string) error {
//	strategy := &model.WqUserStrategy{
//		ApiKey: apiKey,
//	}
//	db := d.db.Table(TABLE_WQ_USER_STRATEGY).Where("user_id=? and strategy_id=?", uid, strategyId).Update(strategy)
//	if db.Error != nil {
//		logger.Errorf("SetUserStrategyApi has err %v", db.Error)
//		return db.Error
//	}
//	if db.RowsAffected == 0 {
//		logger.Warnf("SetUserStrategyApi no row found user_id %s strategy_id %s", uid, strategyId)
//		return errors.New("no row found")
//	}
//	return nil
//}

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
	//db := d.db.Table(TABLE_WQ_USER_STRATEGY).Where("user_id=? and platform=?", uid, platform).UpdateColumns(model.GridStrategy{ApiKey: apiKey})
	//if db.Error != nil {
	//	logger.Errorf("SetUserAllStrategyApi has err %v", db.Error)
	//	return db.Error
	//}
	//if db.RowsAffected == 0 {
	//	logger.Warnf("SetUserAllStrategyApi no row found user_id %s  platform %s", uid, platform)
	//	return errors.New("no row found")
	//}
	return nil
}

//func (d *Dao) SetUserStrategyBalance(uid, strategyId, balance string) error {
//	strategy := &model.WqUserStrategy{
//		Balance: balance,
//	}
//	db := d.db.Table(TABLE_WQ_USER_STRATEGY).Where("user_id=? and strategy_id=?", uid, strategyId).Update(strategy)
//	if db.Error != nil {
//		logger.Errorf("SetUserStrategyBalance has err %v", db.Error)
//		return db.Error
//	}
//	if db.RowsAffected == 0 {
//		logger.Warnf("SetUserStrategyBalance no row found user_id %s strategy_id %s", uid, strategyId)
//		return errors.New("no row found")
//	}
//	return nil
//}

//func (d *Dao) CreateUserStrategy(strategy *model.WqUserStrategy) error {
//	if err := d.db.Table(TABLE_WQ_USER_STRATEGY).Create(strategy).Error; err != nil {
//		logger.Errorf("CreateUserStrategy has  err %v", strategy)
//		return err
//	}
//	return nil
//}

//func (d *Dao) UpdateUserStrategy(uid string, strategy *model.WqUserStrategy) error {
//	db := d.db.Table(TABLE_WQ_USER_STRATEGY).Where("user_id = ? and strategy_id = ?", uid, strategy.StrategyID).Update(strategy)
//	if db.Error != nil {
//		logger.Errorf("UpdateUserStrategy has err %v", db.Error)
//		return db.Error
//	}
//	if db.RowsAffected == 0 {
//		return errors.New("no rows found")
//	}
//	return nil
//}
