package session

import (
	"context"
	"fmt"
	"strings"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/pkg/utils"
	"wq-fotune-backend/internal/exchange-srv/client"
	fotune_srv_exchange "wq-fotune-backend/internal/exchange-srv/proto"
	"wq-fotune-backend/internal/forward-offer-srv/config"
	"wq-fotune-backend/internal/forward-offer-srv/global"
	"wq-fotune-backend/internal/forward-offer-srv/srv/cache_service"
	"wq-fotune-backend/internal/forward-offer-srv/srv/model"
)

const (
	logKey = "【 可靠交易订单 】"
)

// ProcessReliableOrderEvent 处理可靠交易
func (c *Client) ProcessReliableOrderEvent(delayerReq *model.DelayerReq) {
	if delayerReq == nil {
		return
	}
	logger.Infof("%s开始处理任务: %v", logKey, delayerReq.OrderID)
	defer logger.Infof("%s结束处理任务: %v", logKey, delayerReq.OrderID)
	c.checkOrderIsDeal(delayerReq)
}

// cancelOrderFromExchange 撤销订单
func (c *Client) cancelOrderFromExchange(instrumentId, orderID string) bool {
	logger.Warnf("%s撤销订单: %v 合约: %v", logKey, orderID, instrumentId)
	var err error
	err = utils.ReTryFunc(10, func() (bool, error) {
		_, err := c.ApiClient.CancelOrderByID(instrumentId, orderID)
		// 您当前没有未成交的订单（用户撤销未成交订单时,报错32004, 不需要重试
		if err != nil && !strings.Contains(err.Error(), "32004") {
			logger.Warnf("CancelOrderByID warn: %s", err.Error())
			return false, nil
		}
		return false, nil
	})
	return err == nil
}

// getAllUnFinishOrders 获取未成交的交易
// checkOrderIsDeal 检查订单是否正常
func (c *Client) checkOrderIsDeal(delayerReq *model.DelayerReq) {
	// 根据订单查询交易所成交记录
	c.cancelOrderFromExchange(delayerReq.Symbol, delayerReq.OrderID)
	order := c.ApiClient.FindOrderByOrderID(delayerReq.Symbol, delayerReq.OrderID)
	if order.Status == goex.ORDER_FINISH {
		c.calTrade(order.OrderID2, delayerReq)
		consoleFinish(delayerReq)
		return
	}
	if order.Status == goex.ORDER_CANCEL {
		var unDealCount = order.Amount - order.DealAmount
		logger.Warnf("%s检查到订单并没有完全成交: %v 剩余数量: %v", logKey, delayerReq.OrderID, unDealCount)
		c.calTrade(order.OrderID2, delayerReq)
		if unDealCount > 0 {
			c.RetrySendOrder(delayerReq, unDealCount)
		}
	}
	if order.Status == goex.ORDER_REJECT {
		logger.Warnf("%s检查到订单被拒绝: %v 数据: %v", logKey, delayerReq.OrderID, global.StructToJsonStr(order))
	}
	if order.Status == goex.ORDER_FAIL {
		logger.Warnf("%s检查到订单被拒绝: %v 数据: %v", logKey, delayerReq.OrderID, global.StructToJsonStr(order))
	}
}

func (c *Client) calTrade(orderId string, delayerReq *model.DelayerReq) {
	trades := c.ApiClient.FindTradesByOrderID(delayerReq.Symbol, orderId)
	var exchangeClient = client.NewExOrderClient(config.Config.EtcdAddr)
	for _, trade := range trades {
		t := &fotune_srv_exchange.TradeReq{
			TradeId:    string(trade.Tid),
			UserId:     delayerReq.UserID,
			ApiKey:     delayerReq.APIKey,
			OrderId:    orderId,
			StrategyId: delayerReq.StrategyID,
			Direction:  delayerReq.Direction,
			Volume:     global.Float64ToString(trade.Amount),
			Commission: trade.Fee,
			Symbol:     delayerReq.Symbol,
			Price:      global.Float64ToString(trade.Price),
			Exchange:   "okex",
		}
		//发送成交记录到统计服务
		_, err := exchangeClient.Evaluation(context.Background(), t)
		if err != nil {
			logger.Warnf("成交实时计算失败: %s", global.StructToJsonStr(t))
			continue
		}
	}

}

// RetrySendOrder 重试发单
func (c *Client) RetrySendOrder(deq *model.DelayerReq, count float64) {
	if deq.TryCount == 0 {
		msg := fmt.Sprintf("%s检查到订单已最大次数重试: %v", logKey, deq.OrderID)
		logger.Warnf(msg)
		return
	}
	req := deq.OrderReq
	pOrder := req.OrderID
	req.OrderID = global.GetUUID()
	req.TryCount -= 1
	req.OrderQty = count
	orderBook := c.ApiClient.GetOrderBook(deq.Symbol)
	if orderBook == nil {
		msg := fmt.Sprintf("%s获取最新的对手价格失败: %v InstrumentName: %v", logKey, deq.OrderID, deq.Symbol)
		logger.Warn(msg)
		return
	}
	var ask = orderBook.AskList[1]
	var bid = orderBook.BidList[1]
	// 刷新订单价 = 对手价 + 滑价
	if strings.ToLower(req.Direction) == global.BuyType {
		// 卖价
		req.Price = ask.Price
	} else {
		// 买家
		req.Price = bid.Price
	}
	msg := fmt.Sprintf("%s正在重试发单: %v 父类ID: %v 价格为: %v", logKey, req.OrderID, pOrder, req.Price)
	logger.Warn(msg)
	cache_service.CacheOrderIDBindUser(req.OrderID, deq.UserID, deq.StrategyID)
	checkError := c.ProcessCreateOrder(req, 1)
	if checkError != nil {
		msg := fmt.Sprintf("%s提交订单出错, %v, 订单ID: %v", logKey, checkError.Msg, deq.OrderID)
		logger.Warn(msg)
	}
	c.pushReliableOrderToWaitQueue(req, deq)
}

func consoleFinish(delayerReq *model.DelayerReq) {
	logger.Warnf("%s检查到订单完全成交, strategyID: %s  orderID: %s", logKey, delayerReq.StrategyID, delayerReq.OrderID)
}
