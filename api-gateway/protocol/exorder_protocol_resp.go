package protocol

import (
	"time"
	pb "wq-fotune-backend/app/exchange-srv/proto"
)

type ExchangeResp struct {
	ID        int64     `json:"id"`
	Exchange  string    `json:"exchange"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExchangeApiResp struct {
	ID            int64             `json:"id"`
	UserID        string            `json:"user_id"`
	ExchangeID    int32             `json:"exchange_id"`
	ExchangeName  string            `json:"exchange_name"`
	ApiKey        string            `json:"api_key"`
	Secret        string            `json:"secret"`
	Passphrase    string            `json:"passphrase"`
	TotalUsdt     string            `json:"total_usdt"`
	TotalRmb      string            `json:"total_rmb"`
	UsdtBalance   string            `json:"usdt_balance"`
	BtcBalance    string            `json:"btc_balance"`
	BalanceDetail []*pb.ExchangePos `json:"balance_detail"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type TradeResp struct {
	ID           int64     `gorm:"column:id" json:"id"`
	UserID       string    `gorm:"column:user_id" json:"user_id"`
	TradeID      string    `gorm:"column:trade_id" json:"trade_id"`
	ApiKey       string    `gorm:"column:api_key" json:"api_key"`
	StrategyID   string    `gorm:"column:strategy_id" json:"strategy_id"`
	Symbol       string    `gorm:"column:symbol" json:"symbol"`
	OpenPrice    float64   `gorm:"column:open_price" json:"open_price"`
	ClosePrice   float64   `gorm:"column:close_price" json:"close_price"`
	AvgPrice     float64   `gorm:"column:avg_price" json:"avg_price"`
	Volume       string    `gorm:"column:volume" json:"volume"`
	Commission   string    `gorm:"column:commission" json:"commission"`
	Profit       string    `gorm:"column:profit" json:"profit"`
	Pos          string    `gorm:"column:pos" json:"pos"`
	PosDirection string    `gorm:"column:pos_direction" json:"pos_direction"`
	Direction    string    `gorm:"column:direction" json:"direction"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type ProfitResp struct {
	ID              int64     `json:"id"`
	UserID          string    `json:"user_id"`
	ApiKey          string    `json:"api_key"`
	StrategyID      string    `json:"strategy_id"`
	Symbol          string    `json:"symbol"`
	RealizeProfit   string    `json:"realize_profit"`
	UnRealizeProfit string    `json:"un_realize_profit"`
	Position        int64     `json:"position"`
	RateReturn      float64   `json:"rate_return"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserStrategyResp struct {
	UserID           string `json:"user_id"`
	StrategyID       string `json:"strategy_id"`
	ParentStrategyID int64  `json:"parent_strategy_id"`
	ApiKey           string `json:"api_key"`
	Platform         string `json:"platform"`
	Balance          string `json:"balance"`
	State            int32  `json:"state"`
}

type WqStrategy struct {
	ID           int64  `json:"id"`
	Tag          string `json:"tag"`
	Level        int32  `json:"level"`
	Name         string `json:"name"`
	Remark       string `json:"remark"`
	ExchangeName string `json:"exchange_name"`
	ExchangeID   int32  `json:"exchange_id"`
	State        int32  `json:"state"`
}

type UserStrategyEvaResp struct {
	RealizeProfit  string        `json:"realize_profit"`
	RateReturnYear string        ` json:"rate_return_year"`
	RateReturn     string        `json:"rate_return"`
	EvaDaily       []interface{} `json:"evaDaily"`
}

type ExchangeApiInfoResp struct {
	UserId       string `json:"user_id"`
	ExchangeId   int64  `json:"exchange_id"`
	ExchangeName string `json:"exchange_name"`
	ApiKey       string `json:"api_key"`
	Secret       string `json:"secret"`
	Passphrase   string `json:"passphrase"`
}
