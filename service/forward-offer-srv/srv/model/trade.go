package model

import (
	"time"
	"wq-fotune-backend/service/forward-offer-srv/global"
)

type WqTrade struct {
	ID         int64     `gorm:"column:id" json:"id"`
	UserID     string    `gorm:"column:user_id" json:"user_id"`
	LedgerId   string    `json:"ledger_id"`
	OrderId    string    `json:"order_id"`
	ApiKey     string    `gorm:"column:api_key" json:"api_key"`
	StrategyID string    `gorm:"column:strategy_id" json:"strategy_id"`
	Symbol     string    `gorm:"column:symbol" json:"symbol"`
	OpenPrice  float64   `gorm:"column:open_price" json:"open_price"`
	ClosePrice float64   `gorm:"column:close_price" json:"close_price"`
	AvgPrice   float64   `gorm:"column:avg_price" json:"avg_price"`
	Volume     string    `gorm:"column:volume" json:"volume"`
	Commission string    `gorm:"column:commission" json:"commission"`
	Profit     string    `gorm:"column:profit" json:"profit"`
	Direction  string    `gorm:"column:direction" json:"direction"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// DelayerReq 用于延迟队列
type DelayerReq struct {
	ExchangeInfo
	OrderReq
	StrategyID string `json:"strategyID" binding:"required"`
	UserID     string `json:"userID"`
	Type       int    `json:"type"` // 1 重试单  0为第一次发单
	PqTime     time.Time
}

// ExchangeInfo 交易所信息
type ExchangeInfo struct {
	APIKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
	EcPass    string `json:"ec_pass"`
}

// OrderReq 下单实体
type OrderReq struct {
	OrderID   string  `json:"orderID" binding:"required"`  // 自定义订单ID
	OrdType   string  `json:"ordType" binding:"required"`  // Limit or Market
	Symbol    string  `json:"symbol" binding:"required"`   // 品种
	Direction string  `json:"direction"`                   // Buy or Sell
	OrderQty  float64 `json:"orderQty" binding:"required"` // 订单数量
	Price     float64 `json:"price"`                       // 价格
	TryCount  int     `json:"tryCount"`                    // 重试次数
}

func CreateDelayerReq(orderReq OrderReq, userID, strategyID string) *DelayerReq {
	req := &DelayerReq{}
	req.OrderID = orderReq.OrderID
	req.OrderQty = orderReq.OrderQty
	req.Price = orderReq.Price
	req.Direction = orderReq.Direction
	req.StrategyID = strategyID
	req.Symbol = orderReq.Symbol
	req.TryCount = orderReq.TryCount
	req.UserID = userID
	req.PqTime = global.GetCurrentTime()
	return req
}

type CheckError struct {
	TypeCode int
	Msg      string
}

func CreateCheckError(errType int, msg string) *CheckError {
	return &CheckError{
		TypeCode: errType,
		Msg:      msg,
	}
}
