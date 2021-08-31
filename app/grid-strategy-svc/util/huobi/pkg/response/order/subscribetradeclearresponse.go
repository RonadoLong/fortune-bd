package order

import "fortune-bd/app/grid-strategy-svc/util/huobi/pkg/response/base"

type SubscribeTradeClearResponse struct {
	base.WebSocketV2ResponseBase
	Data *struct {
		Symbol        string `json:"symbol"`
		OrderId       int64  `json:"orderId"`
		TradePrice    string `json:"tradePrice"`
		TradeVolume   string `json:"tradeVolume"`
		OrderSide     string `json:"orderSide"`
		OrderType     string `json:"orderType"`
		Aggressor     bool   `json:"aggressor"`
		TradeId       int64  `json:"tradeId"`
		TradeTime     int64  `json:"tradeTime"`
		TransactFee   string `json:"transactFee"`
		FeeDeduct     string `json:"feeDeduct"`
		FeeDeductType string `json:"feeDeductType"`
	}
}
