package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/service/exchange-srv/model"
)

func (d *Dao) TradeCount(uid, strategyID string) (count int32, err error) {
	if err = d.db.Table(TABLE_WQ_TRADE).Where("user_id=? and strategy_id=?", uid, strategyID).Count(&count).
		Error; err != nil {
		logger.Errorf("TradeCount has err %v", err)
		return
	}
	return
}

func (d *Dao) GetTradeList(uid, strategyID string, pageNum, pageSize int32) (tradeList []*model.WqTrade) {
	if err := d.db.Table(TABLE_WQ_TRADE).Where("user_id=? and strategy_id=?", uid, strategyID).
		Order("id desc").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&tradeList).Error; err != nil {
		logger.Errorf("GetTradeList has err %v", err)
		return
	}
	return
}

func (d *Dao) GetTradeByLastID(uid, strategyID string) (*model.WqTrade, error) {
	trade := &model.WqTrade{}
	if err := d.db.Table(TABLE_WQ_TRADE).Where("user_id=? and strategy_id=?", uid, strategyID).
		Order("id desc").First(trade).Error; err != nil {
		logger.Warnf("GetTradeByLastID has err  %v", err)
		return nil, err
	}
	return trade, nil
}

func (d *Dao) CreateTrade(trade *model.WqTrade) error {
	if err := d.db.Table(TABLE_WQ_TRADE).Create(&trade).Error; err != nil {
		logger.Warnf("CreateTrade has err %v uid %s trade_id %s", err, trade.UserID, trade.TradeID)
		return err
	}
	return nil
}

func (d *Dao) GetProfitRealTime(uid, strategyID string) (wqProfit []*model.WqProfit) {
	if err := d.db.Table(TABLE_WQ_PROFIT).Where("strategy_id=?", strategyID).Find(&wqProfit).
		Error; err != nil {
		logger.Errorf("GetProfitRealTime has err %v", err)
		return
	}
	return
}

func (d *Dao) GetProfitByID(uid, strategyID string) (*model.WqProfit, error) {
	wqProfit := &model.WqProfit{}
	if err := d.db.Table(TABLE_WQ_PROFIT).Where("strategy_id=?", strategyID).First(wqProfit).
		Error; err != nil {
		logger.Warnf("GetProfitByID has err %v", err)
		return nil, err
	}
	return wqProfit, nil
}

func (d *Dao) CreateProfit(profit *model.WqProfit) error {
	if err := d.db.Table(TABLE_WQ_PROFIT).Create(profit).Error; err != nil {
		logger.Errorf("CreateProfit has err %v", err)
		return err
	}
	return nil
}

func (d *Dao) CreateProfitDaily(profit *model.WqProfitDaily) error {
	if err := d.db.Table(TABLE_WQ_PROFIT_DAILY).Create(profit).Error; err != nil {
		logger.Errorf("CreateProfitDaily has err %v", err)
		return err
	}
	return nil
}

func (d *Dao) GetLastProfitByStrategyId(strategyId string) (*model.WqProfitDaily, error) {
	profitDaily := &model.WqProfitDaily{}
	if err := d.db.Table(TABLE_WQ_PROFIT_DAILY).Where("strategy_id=?", strategyId).Order("date desc").
		First(profitDaily).Error; err != nil {
		return nil, err
	}
	return profitDaily, nil
}

func (d *Dao) UpdateProfit(uid, strategyID string, profit *model.WqProfit) error {
	if err := d.db.Table(TABLE_WQ_PROFIT).Where("strategy_id=?", strategyID).Update(profit).Error; err != nil {
		logger.Errorf("UpdateProfit has err %v", err)
		return err
	}
	return nil
}

func (d *Dao) UpdateProfitDaily(profit *model.WqProfitDaily) error {
	if err := d.db.Table(TABLE_WQ_PROFIT_DAILY).Where("id=?", profit.ID).Update(profit).Error; err != nil {
		logger.Errorf("UpdateProfitDaily has err %v", err)
		return err
	}
	return nil
}

func (d *Dao) GetProfitDailyList(uid, strategyID string, limit int) (wqProfitDailyList []*model.WqProfitDaily) {
	if err := d.db.Table(TABLE_WQ_PROFIT_DAILY).Where("strategy_id=?", strategyID).Order("date asc").Limit(limit).Find(&wqProfitDailyList).
		Error; err != nil {
		logger.Errorf("GetProfitDailyList has err %v", err)
		return
	}
	return
}

func (d *Dao) GetProfitListBySql(sql string, orderBy string) (wqProfitList []*model.WqProfit) {
	if err := d.db.Table(TABLE_WQ_PROFIT).Where(sql).Order(orderBy).Find(&wqProfitList).
		Error; err != nil {
		logger.Errorf("GetProfitListBySql has err %v", err)
		return
	}
	return
}

func (d *Dao) GetOrderRecord(orderID string) (*model.GridTradeRecord, error) {
	trade := &model.GridTradeRecord{}
	ctx, cancelFunc := newCtx()
	defer cancelFunc()
	err := d.mongo.Database(DATABASE).Collection(TABLE_MONGO_GRIDTRADERECORD).FindOne(ctx, bson.M{"orderID": orderID}).Decode(trade)
	if err != nil {
		logger.Warnf("GetOrderRecord has err %v", err)
		return nil, err
	}
	return trade, nil
}

func (d *Dao) CreateWqTradeRecord(trade *model.WqTradeRecord) error {
	if err := d.db.Table(TABLE_WQ_TRADE_RECORD).Create(trade).Error; err != nil {
		logger.Errorf("CreateWqTradeRecord has err %v", err)
		return err
	}
	return nil
}
