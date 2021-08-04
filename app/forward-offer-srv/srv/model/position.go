package model

import "time"

// RedisCachePosition 持仓
type RedisCachePosition struct {
	UserID     string    `json:"userID"`
	StrategyID string    `json:"strategyID"` //策略id
	Symbol     string    `json:"symbol"`     //品种
	Direction  string    `json:"direction"`  //买 - 多 / 卖 -空 /
	Volume     float64   `json:"volume"`     //持仓数量
	Price      float64   `json:"price"`      //价格
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type WqAllPosition struct {
	ID         int64     `json:"id"`
	ApiKey     string    `json:"api_key"`
	UserID     string    `json:"user_id"`
	StrategyID string    `json:"strategy_id"` //策略id
	Exchange   string    `json:"exchange"`    //交易所名称
	Symbol     string    `json:"symbol"`      //品种
	Direction  string    `json:"direction"`   //买 - 多 / 卖 -空 /
	Volume     string    `json:"volume"`      //持仓数量
	Price      string    `json:"price"`       //价格
	Status     int       `json:"status"`
	CreateAt   time.Time `json:"create_at"`
	UpdateAt   time.Time `json:"update_at"`
}
