package session

import (
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/internal/forward-offer-srv/global"
	"wq-fotune-backend/internal/forward-offer-srv/srv/cache_service"
	"wq-fotune-backend/internal/forward-offer-srv/srv/model"
)

// ProcessCallBackOrder 处理订单回调
func (c *Client) ProcessCallBackOrder(order *goex.FutureOrder) {
	var orderID = order.ClientOid
	if orderID == "" {
		logger.Warnf("【 接收到订单回报 】订单不是在平台交易: %v", global.StructToJsonStr(order))
		return
	}
	logger.Warnf(" 【 接收到订单回报 】%v", global.StructToJsonStr(order))
	if order.Status == goex.ORDER_FINISH || order.Status == goex.ORDER_PART_FINISH {
		c.getTadesList(order)
		//todo 统计交易
	}
}

func (c *Client) getTadesList(order *goex.FutureOrder) {
	uid, sid := cache_service.GetOrderIDBindUserVal(order.ClientOid)
	logger.Infof("实时结算用户的成交: %s", uid)
	id := cache_service.GetCacheAccountLeID(sid, order.ContractName)
	hisotry, err := c.ApiClient.GetAccountTradeHisotry(order.ContractName, id)
	if err != nil {
		return
	}
	var beforeId string
	for _, hTrade := range hisotry {
		// 判断是否为这个订单的交易流水
		if hTrade.Details != nil && hTrade.Details.OrderID == order.OrderID2 {
			if beforeId == "" {
				beforeId = hTrade.LedgerID
			}
			trade := model.WqTrade{
				UserID:     uid,
				LedgerId:   hTrade.LedgerID,
				OrderId:    order.OrderID2,
				ApiKey:     c.ApiClient.GetApiKey(),
				StrategyID: sid,
				Symbol:     order.ContractName,
				OpenPrice:  order.Price,
				AvgPrice:   order.AvgPrice,
				Commission: global.Float64ToString(order.Fee),
				Profit:     hTrade.Amount,
				Direction:  parseState(order.OType),
				CreatedAt:  global.GetCurrentTime(),
				UpdatedAt:  global.GetCurrentTime(),
			}
			// 判断是开仓还是平仓
			if global.StringToFloat64(hTrade.Balance) < 0 {
				trade.ClosePrice = order.AvgPrice
			}
			// 成交记录
			trade.Volume = hTrade.Balance
			logger.Infof("实时结算用户的成交: %s", global.StructToJsonStr(trade))
		}
	}
	if beforeId != "" {
		cache_service.GetCacheAccountLeID(sid, order.ContractName)
	}
}

func parseState(s int) string {
	switch s {
	case 1:
		return global.BuyType // "开多"
	case 2:
		return global.SellType // 开空
	case 3:
		return global.SellType // 平多
	case 4:
		return global.BuyType // 平空
	default:
		return "未知"
	}
}
