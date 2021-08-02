package order

import "wq-fotune-backend/internal/grid-strategy-srv/util/huobi/pkg/response/base"

type SubscribeOrderV2Response struct {
	base.WebSocketV2ResponseBase
	Data *struct {
		EventType       string `json:"eventType"`
		Symbol          string `json:"symbol"`
		OrderId         int64  `json:"orderId"`
		ClientOrderId   string `json:"clientOrderId"`
		OrderPrice      string `json:"orderPrice"`
		OrderSize       string `json:"orderSize"`
		Type            string `json:"type"`
		OrderStatus     string `json:"orderStatus"`
		OrderCreateTime int64  `json:"orderCreateTime"`
		TradePrice      string `json:"tradePrice"`
		TradeVolume     string `json:"tradeVolume"`
		TradeId         int64  `json:"tradeId"`
		TradeTime       int64  `json:"tradeTime"`
		Aggressor       bool   `json:"aggressor"`
		RemainAmt       string `json:"remainAmt"`
		LastActTime     int64  `json:"lastActTime"`
	}
}
