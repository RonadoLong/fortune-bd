package model

// ExchangeReq 策略订单请求数据结构
type ExchangeReq struct {
	Exchange     string        `json:"exchange"`                        // ctp or bitmex
	UserID       string        `json:"userID" binding:"required"`       // 用户唯一标识符
	StrategyID   string        `json:"strategyID" binding:"required"`   // 策略ID
	ExchangeInfo *ExchangeInfo `json:"exchangeInfo" binding:"required"` // 交易所信息
	Data         string        `json:"data"`                            // 请求的数据, 发单或者撤销订单的数据
	TradeReq     *TraderReqData
}

// TraderReqData 对应data数据
type TraderReqData struct {
	Type  string `json:"type"`  // creat cancel autoAdd
	Value string `json:"value"` // 取消 OrderCancelReq  下单 OrderReq
}
