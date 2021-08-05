package protocol

// TradeSignal 交易信号
type TradeSignal struct {
	UID       string      `json:"uid"`       // 用户id
	FileID    string      `json:"fileID"`    // 策略文件id
	SharedID  string      `json:"sharedID"`  // 共享id
	TradeType int         `json:"tradeType"` // 交易信号类型，0:普通下单，2:挂单
	Orders    interface{} `json:"orders"`// 订单内容
}

// OrdinaryOrder 普通订单
type OrdinaryOrder struct {
	Side               string `json:"side"`                  // 交易动作
	OrdType            string `json:"ordType"`               // 订单类型
	OrderQty           int    `json:"orderQty"`              // 订单数量
	OrderID            string `json:"orderID"`               // 订单id
	DelayerTime        int    `json:"delayerTime"`           // 延时
	TryCount           int    `json:"tryCount"`              // 重试次数
	Exchange           string `json:"exchange"`              // 交易所
	Symbol             string `json:"symbol"`                // 品种
	Price              int    `json:"price"`                 // 价格
	SlipPrice          int    `json:"slipPrice"`             // 滑价
	TradingCoefficient string `json:"tradingQtyCoefficient"` // 交易量系数，默认为1
}
