package biz

import (
	"fortune-bd/app/exchange-svc/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (e *ExOrderRepo) GetUserStrategyOfRun() []*model.GridStrategy {
	return e.dao.GetUserStrategyOfRun(nil)
}

func (e *ExOrderRepo) GetUserStrategyByUID(userID string) []*model.GridStrategy {
	sql := bson.M{"uid": userID}
	return e.dao.GetUserStrategyOfRun(sql)
}

func (e *ExOrderRepo) UpdateProfit(strategyID string, profit *model.WqProfit) error {
	return e.dao.UpdateProfit("", strategyID, profit)
}
