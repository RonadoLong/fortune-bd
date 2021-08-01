package dao

import (
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/service/exchange-srv/model"
)

func (d *Dao) GetAllSymbolWithState(state int32, exchange, unit string) (symbols []*model.WqSymbol) {
	if err := d.db.Table(TABLE_WQ_SYMBOL).Where("state = ? and exchange = ? and unit = ?", state, exchange, unit).
		Order("id desc").Find(&symbols).Error; err != nil {
		logger.Infof("GetAllSymbol 没有找到数据")
		return
	}
	return
}

func (d *Dao) GetSymbolRecommend(state int32) (symbols []*model.WqSymbolRecommend) {
	if err := d.db.Table(TABLE_WQ_SYMBOL_RECOMMEND).Where("state = ?", state).Find(&symbols).Error; err != nil {
		logger.Infof("GetSymbolRecommend 没有找到数据")
		return
	}
	return
}
