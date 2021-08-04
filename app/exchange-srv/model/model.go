package model

import (
	"time"
)

type WqExchange struct {
	ID        int64     `gorm:"column:id" json:"id"`
	Exchange  string    `gorm:"column:exchange" json:"exchange"`
	Status    string    `gorm:"column:status" json:"status"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// 用户的apikey表
type WqExchangeApi struct {
	ID           int64     `gorm:"column:id" json:"id"`
	UserID       string    `gorm:"column:user_id" json:"user_id"`
	ExchangeID   int64     `gorm:"column:exchange_id" json:"exchange_id"`
	ExchangeName string    `gorm:"column:exchange_name" json:"exchange_name"`
	ApiKey       string    `gorm:"column:api_key" json:"api_key"`
	Secret       string    `gorm:"column:secret" json:"secret"`
	Passphrase   string    `gorm:"column:passphrase" json:"passphrase"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

//合约的交易记录
type WqTrade struct {
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

func NewWqTrade(
	userID string,
	tradeID string,
	apiKey string,
	strategyID string,
	symbol string,
	volume string,
	commission string,
	direction string) *WqTrade {
	return &WqTrade{
		UserID:       userID,
		TradeID:      tradeID,
		ApiKey:       apiKey,
		StrategyID:   strategyID,
		Symbol:       symbol,
		OpenPrice:    0,
		ClosePrice:   0,
		AvgPrice:     0,
		Volume:       volume,
		Commission:   commission,
		Profit:       "",
		Pos:          "",
		PosDirection: "",
		Direction:    direction,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

//盈亏表
type WqProfit struct {
	ID              int64     `gorm:"column:id" json:"id"`
	UserID          string    `gorm:"column:user_id" json:"user_id"`
	ApiKey          string    `gorm:"column:api_key" json:"api_key"`
	StrategyID      string    `gorm:"column:strategy_id" json:"strategy_id"`
	Symbol          string    `gorm:"column:symbol" json:"symbol"`
	RealizeProfit   string    `gorm:"column:realize_profit" json:"realize_profit"`
	UnRealizeProfit string    `gorm:"column:un_realize_profit" json:"un_realize_profit"`
	Position        int64     `gorm:"column:position" json:"position"`
	RateReturn      float64   `gorm:"column:rate_return" json:"rate_return"`
	RateReturnYear  float64   `gorm:"column:rate_return_year" json:"rate_return_year"`
	Commission      string    `gorm:"column:commission" json:"commission"`
	Unit            string    `gorm:"column:unit" json:"unit"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type WqProfitRank struct {
	Symbol         string  `gorm:"column:symbol" json:"symbol"`
	RateReturnYear float64 `gorm:"column:ANY_VALUE(rate_return_year)" json:"rate_return_year"`
}

type WqProfitDaily struct {
	WqProfit
	Date time.Time `gorm:"column:date" json:"date"`
}

//type WqUserStrategy struct {
//	ID               int64     `gorm:"column:id" json:"id"`
//	UserID           string    `gorm:"column:user_id" json:"user_id"`
//	GroupID          string    `gorm:"column:group_id" json:"group_id"`
//	StrategyID       string    `gorm:"column:strategy_id" json:"strategy_id"`
//	ParentStrategyID int64     `gorm:"column:parent_strategy_id" json:"parent_strategy_id"`
//	ApiKey           string    `gorm:"column:api_key" json:"api_key"`
//	Platform         string    `gorm:"column:platform" json:"platform"`
//	Balance          string    `gorm:"column:balance" json:"balance"`
//	State            int32     `gorm:"column:state" json:"state"`
//	Symbol           string    `gorm:"column:symbol" json:"symbol"`
//	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
//	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
//}

// GridStrategy 网格策略
type GridStrategy struct {
	ID   string `json:"id" bson:"_id"`
	UID  string `json:"uid" bson:"uid"`   // 用户id
	Type int    `json:"type" bson:"type"` // 0:网格交易，1:杠杆网格，2:借贷网格，3:反向网格，4:反向杠杆网格，5:无线网格

	// 交易币网格信息
	ApiKey   string `json:"apiKey" bson:"apiKey"`     // 交易所apikey
	Exchange string `json:"exchange" bson:"exchange"` // 交易所
	Symbol   string `json:"symbol" bson:"symbol"`     // 交易的品种

	GridIntervalType string  `json:"gridIntervalType" bson:"gridIntervalType"` // 网格间隔类型，AQGrid:等差, PQGrid:等比
	GridDesc         string  `json:"gridDesc" bson:"gridDesc"`                 // 网格说明，例如xxx震荡上涨
	MinPrice         float64 `json:"minPrice" bson:"minPrice"`                 // 网格最低价格
	MaxPrice         float64 `json:"maxPrice" bson:"maxPrice"`                 // 网格最高价格
	GridNum          int     `json:"gridNum" bson:"gridNum"`                   // 网格数量
	EachGridMoney    float64 `json:"eachGridMoney" bson:"eachGridMoney"`       // 每格买卖金额
	EachGridProfit   string  `json:"eachGridProfit" bson:"eachGridProfit"`     // 每格利润

	StopProfitPrice float64 `json:"stopProfitPrice" bson:"stopProfitPrice"` // 止盈价格
	StopLossPrice   float64 `json:"stopLossPrice" bson:"stopLossPrice"`     // 止损价格
	BasisPrice      float64 `json:"basisPrice" bson:"basisPrice"`           // 网格开单价格，当价格小于等于基准价格时买入，否则卖出，用来统计需要的资金和已买入的币，当基准价格为0时，使用当前币价格

	// 锚定币信息
	AnchorSymbol string    `json:"anchorSymbol" bson:"anchorSymbol"` // 锚定币的品种名称，例如USDT
	TotalSum     float64   `json:"totalSum" bson:"totalSum"`         // 网格投资总额
	IsRun        bool      `json:"isRun" bson:"isRun"`               // 网格策略是否运行中
	CreatedAt    time.Time `json:"created_at" bson:"createdAt"`
}

type WqStrategy struct {
	ID           int64     `gorm:"column:id" json:"id"`
	Name         string    `gorm:"column:name" json:"name"`
	Remark       string    `gorm:"column:remark" json:"remark"`
	ExchangeName string    `gorm:"column:exchange_name" json:"exchange_name"`
	ExchangeID   int64     `gorm:"column:exchange_id" json:"exchange_id"`
	Tag          string    `gorm:"column:tag" json:"tag"`
	Level        int32     `gorm:"column:level" json:"level"`
	State        int32     `gorm:"column:state" json:"state"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type WqSymbol struct {
	ID       int64  `gorm:"column:id" json:"id"`
	Symbol   string `gorm:"column:symbol" json:"symbol"`
	Exchange string `gorm:"column:exchange" json:"exchange"`
	State    int32  `gorm:"column:state" json:"state"`
	Unit     string `gorm:"column:unit" json:"unit"`
}

type WqSymbolRecommend struct {
	ID             int64  `gorm:"column:id" json:"id"`
	Symbol         string `gorm:"column:symbol" json:"symbol"`
	RateReturnYear string `gorm:"column:rate_return_year" json:"rate_return_year"`
	State          int32  `gorm:"column:state" json:"state"`
	Url            string `gorm:"column:url" json:"url"`
}

type RateRank struct {
	ID             int    `json:"id"`
	UserId         string `json:"user_id"`
	Avatar         string `json:"avatar"`
	Name           string `json:"name"`
	RateReturn     string `json:"rate_return"`
	RateReturnYear string `json:"rate_return_year"`
}

// GridTradeRecord 网格交易记录
type GridTradeRecord struct {
	ID                  string    `json:"id" bson:"_id"`
	GSID                string    `json:"gsid" bson:"gsid"`                               // 网格策略id
	GID                 int       `json:"gid" bson:"gid"`                                 // 网格编号
	OrderID             string    `json:"orderID" bson:"orderID"`                         // 委托订单id，用户查询、取消订单
	ClientOrderID       string    `json:"clientOrderID" bson:"clientOrderID"`             // 用户自定义订单id
	OrderType           string    `json:"orderType" bson:"orderType"`                     // 订单类型，limit:限价单，market:市价单
	Side                string    `json:"side" bson:"side"`                               // 买入卖出，buy:买入，sell:卖出
	Price               float64   `json:"price" bson:"price"`                             // 成交价格
	Quantity            float64   `json:"quantity" bson:"quantity"`                       // 买卖数量
	Volume              float64   `json:"volume" bson:"volume"`                           // 成交额
	Unit                string    `json:"unit" bson:"unit"`                               // 成交额单位
	Fees                float64   `json:"fees" bson:"fees"`                               // 买入卖出手续费
	OrderState          string    `json:"orderState" bson:"orderState"`                   // 当前订单状态 submitted:委托中, canceled:取消, filled:已成交
	StateTime           time.Time `json:"stateTime" bson:"stateTime"`                     // 订单状态变化时间
	GridPeerSellOrderID string    `json:"gridPeerSellOrderID" bson:"gridPeerSellOrderID"` // 成功卖出后的记录id，此字段只有买入订单记录才有值
	GridPeerBuyOrderID  string    `json:"gridPeerBuyOrderID" bson:"gridPeerBuyOrderID"`   // 成功买入后的记录id，此字段只有卖出订单记录才有值
	Exchange            string    `json:"exchange" bson:"exchange"`                       // 交易所
	Symbol              string    `json:"symbol" bson:"symbol"`                           // 品种名称
}

type WqTradeRecord struct {
	ID            int64     `gorm:"column:id" json:"id"`
	UserID        string    `gorm:"column:user_id" json:"user_id"`
	ApiKey        string    `gorm:"column:api_key" json:"api_key"`
	StrategyID    string    `gorm:"column:strategy_id" json:"strategy_id"`
	OrderID       string    `gorm:"column:order_id" json:"order_id"`
	Symbol        string    `gorm:"column:symbol" json:"symbol"`
	RealizeProfit string    `gorm:"column:realize_profit" json:"realize_profit"`
	BuyPrice      string    `gorm:"column:buy_price" json:"buy_price"`
	SellPrice     string    `gorm:"column:sell_price" json:"sell_price"`
	Commission    string    `gorm:"column:commission" json:"commission"`
	Volume        string    `gorm:"column:volume" json:"volume"`
	Unit          string    `gorm:"column:unit" json:"unit"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}
