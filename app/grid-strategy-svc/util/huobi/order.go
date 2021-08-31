package huobi

import (
	"errors"
	"fmt"
	"fortune-bd/app/grid-strategy-svc/util/goex"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/client"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/client/orderwebsocketclient"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/client/websocketclientbase"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/getrequest"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/postrequest"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/response/auth"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/response/order"

	"github.com/zhufuyi/pkg/krand"
	"github.com/zhufuyi/pkg/logger"
)

const (
	// OrderStateSubmitted 已委托
	OrderStateSubmitted = "submitted"
	// OrderStateSubmitted 已取消
	OrderStateCanceled = "canceled"
	// OrderStateFilled 已成交
	OrderStateFilled = "filled"
)

// 是否模拟执行
var isSimulation = false

// PlaceLimitOrder 买入、卖出限价单，限制频率为50次/2s
func (a *Account) PlaceLimitOrder(side string, symbol string, price string, amount string, clientOrderID string) (string, error) {
	// test 模拟成功订单
	if isSimulation {
		return string(krand.String(krand.R_NUM, 14)), nil
	}

	if a.AccountID == "" {
		return "", errors.New("accountID is empty, need init first")
	}

	cli := new(client.OrderClient).Init(a.AccessKey, a.SecretKey, host)
	resp, err := cli.PlaceOrder(&postrequest.PlaceOrderRequest{
		AccountId:     a.AccountID,
		Type:          side + "-limit", // buy-limit:买入限价单, sell-limit:卖出限价单
		Source:        "spot-api",
		Symbol:        symbol,
		Price:         price,  // 最小为5
		Amount:        amount, // 最小为0.0001
		ClientOrderId: clientOrderID,
	})
	if err != nil {
		return "", err
	}

	if resp.Status != "ok" {
		return "", fmt.Errorf("%s, %s, %s", resp.Status, resp.ErrorCode, resp.ErrorMessage)
	}

	return resp.Data, nil
}

// PlaceMarketOrder 买入、卖出市价单，限制频率为50次/2s
func (a *Account) PlaceMarketOrder(side string, symbol string, amount string, clientOrderID string) (string, error) {
	// test 模拟成功订单
	if isSimulation {
		return string(krand.String(krand.R_NUM, 14)), nil
	}

	if a.AccountID == "" {
		return "", errors.New("accountID is empty, need init first")
	}

	cli := new(client.OrderClient).Init(a.AccessKey, a.SecretKey, host)
	resp, err := cli.PlaceOrder(&postrequest.PlaceOrderRequest{
		AccountId:     a.AccountID,
		Type:          side + "-market", // buy-market:买入市价单, sell-market:卖出市价单
		Source:        "spot-api",
		Symbol:        symbol,
		Amount:        amount, // 市价单为交易额，不是数量，例如btcusdt的最小为5
		ClientOrderId: clientOrderID,
	})
	if err != nil {
		return "", err
	}

	if resp.Status != "ok" {
		return "", fmt.Errorf("%s, %s, %s", resp.Status, resp.ErrorCode, resp.ErrorMessage)
	}

	return resp.Data, nil
}

// PlaceParams 下单参数
type PlaceParams []postrequest.PlaceOrderRequest

// BatchOrders 批量下单(限价或市价单)，最大10单，限制频率为5次/2s
func (a *Account) PlaceSpotOrders(p PlaceParams) ([]order.PlaceOrderResult, error) {
	ids := []order.PlaceOrderResult{}

	if len(p) > 10 {
		return ids, errors.New("cannot place more than 10 orders")
	}

	cli := new(client.OrderClient).Init(a.AccessKey, a.SecretKey, host)
	resp, err := cli.PlaceOrders(p)
	if err != nil {
		return ids, err
	}

	if resp.Status != "ok" {
		return ids, fmt.Errorf("%s, %s, %s", resp.Status, resp.ErrorCode, resp.ErrorMessage)
	}

	return resp.Data, nil
}

// -------------------------------------------------------------------------------------------------

// CancelOrder 取消委托订单
func (a *Account) CancelOrder(orderID string, symbol string) error {
	cli := new(client.OrderClient).Init(a.AccessKey, a.SecretKey, host)
	resp, err := cli.CancelOrderById(orderID)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("%s, %s, %s", resp.Status, resp.ErrorCode, resp.ErrorMessage)
	}

	return nil
}

// CancelResult 取消订单结果
type CancelResult struct {
	Success []string   `json:"success"` // 取消成功订单id
	Failed  []struct { // 取消失败的订单id
		OrderId       string `json:"order-id"`
		ClientOrderId string `json:"client-order-id"`
		ErrorCode     string `json:"err-code"`
		ErrorMessage  string `json:"err-msg"`
	} `json:"failed"`
}

// CancelOrders 批量取消订单
func (a *Account) CancelOrders(orderIDs []string) (*CancelResult, error) {
	result := &CancelResult{}

	cli := new(client.OrderClient).Init(a.AccessKey, a.SecretKey, host)
	resp, err := cli.CancelOrdersByIds(&postrequest.CancelOrdersByIdsRequest{
		OrderIds: orderIDs,
	})
	if err != nil {
		return result, err
	}

	if resp.Status != "ok" {
		return result, fmt.Errorf("%s, %s, %s", resp.Status, resp.ErrorCode, resp.ErrorMessage)
	}

	return &CancelResult{
		Success: resp.Data.Success,
		Failed:  resp.Data.Failed,
	}, nil
}

// -------------------------------------------------------------------------------------------------

// OrderInfo 订单信息
type OrderInfo struct {
	ID               int64  `json:"id"`                // 订单id
	ClientOrderId    string `json:"client-order-id"`   // 用户自定义id
	AccountId        int    `json:"account-id"`        // 账号id
	Symbol           string `json:"symbol"`            // 品种
	Price            string `json:"price"`             // 价格
	Amount           string `json:"amount"`            // 交易数量
	CreatedAt        int64  `json:"created-at"`        // 创建订单时间
	Type             string `json:"type"`              // 订单类型，buy-limit:买入限价单, sell-limit:卖出限价单, buy-market:买入市价单, sell-market:卖出市价单
	FilledAmount     string `json:"field-amount"`      // 已完成交易数量
	FilledCashAmount string `json:"field-cash-amount"` // 已完成交易金额
	FilledFees       string `json:"field-fees"`        // 手续费
	Source           string `json:"source"`            // 账号类型
	State            string `json:"state"`             // 状态

	FinishedAt int64  `json:"finishedAt"`
	StopPrice  string `json:"stopPrice"` //止盈止损订单触发价格
	Operator   string `json:"operator"`  // 止盈止损订单触发价运算符	gte,lte
}

//GetOrder 获取订单信息
func (a *Account) GetOrder(orderID string) (*OrderInfo, error) {
	cli := new(client.OrderClient).Init(a.AccessKey, a.SecretKey, host)
	resp, err := cli.GetOrderById(orderID)
	if err != nil {
		return &OrderInfo{}, err
	}

	if resp.Status != "ok" {
		return &OrderInfo{}, fmt.Errorf("%s, %s, %s", resp.Status, resp.ErrorCode, resp.ErrorMessage)
	}

	return &OrderInfo{
		ID:               resp.Data.Id,
		ClientOrderId:    resp.Data.ClientOrderId,
		AccountId:        resp.Data.AccountId,
		Symbol:           resp.Data.Symbol,
		Price:            resp.Data.Price,
		Amount:           resp.Data.Amount,
		CreatedAt:        resp.Data.CreatedAt,
		Type:             resp.Data.Type,
		FilledAmount:     resp.Data.FilledAmount,
		FilledCashAmount: resp.Data.FilledCashAmount,
		FilledFees:       resp.Data.FilledFees,
		Source:           resp.Data.Source,
		State:            resp.Data.State,
		FinishedAt:       resp.Data.FinishedAt,
		StopPrice:        resp.Data.StopPrice,
		Operator:         resp.Data.Operator,
	}, nil
}

// GetHistoryOrdersInfo 获取历史订单列表
func (a *Account) GetOrderInfo(orderID string, symbol string) (interface{}, error) {
	return a.GetOrder(orderID)
}

// HistoryOrders 历史订单信息
type HistoryOrders struct {
	ID               int64  // 订单id
	ClientOrderID    string `json:"client-order-id"`  // 用户自定义订单id
	AccountID        int    `json:"account-id"`       // 账号id
	UserID           int    `json:"user-id"`          // 用户id
	Symbol           string `json:"symbol"`           // 品种
	Price            string `json:"price"`            // 价格
	Amount           string `json:"amount"`           // 交易数量
	CreatedAt        int64  `json:"created-at"`       // 创建订单时间时间戳，精确到毫秒
	CanceledAt       int64  `json:"canceled-at"`      // 取消订单时间戳
	FinishedAt       int64  `json:"finished-at"`      // 成交时间戳
	Type             string `json:"type"`             // 订单类型，buy-limit:买入限价单, sell-limit:卖出限价单, buy-market:买入市价单, sell-market:卖出市价单
	FilledAmount     string `json:"filledAmount"`     // 已成交数量
	FilledCashAmount string `json:"filledCashAmount"` // 已成交额
	FilledFees       string `json:"filledFees"`       // 已成交手续费
	Source           string `json:"source"`           // 账号类型
	State            string `json:"state"`            // 状态
	Exchange         string `json:"exchange"`         // 交易所
	Batch            string `json:"batch"`
	StopPrice        string `json:"stop-price"`
	Operator         string `json:"operator"`
}

// GetHistoryOrders 获取历史订单信息
func (a *Account) GetHistoryOrders(symbol string, states string, types string) ([]*HistoryOrders, error) {
	hos := []*HistoryOrders{}

	request := new(getrequest.GetRequest).Init()
	request.AddParam("symbol", symbol)
	request.AddParam("states", states) // 多个状态可以用逗号分隔，submitted 已提交, canceled 已撤销,  created 已创建, filled 完全成交, partial-filled 部分成交, partial-canceled 部分成交撤销
	if types != "" {                   // 可选参数
		request.AddParam("types", types) // buy-limit 买入限价单, sell-limit 卖出限价单, buy-market 买入市价单, sell-market 卖出市价单
	}

	cli := new(client.OrderClient).Init(a.AccessKey, a.SecretKey, host)
	resp, err := cli.GetHistoryOrders(request)
	if err != nil {
		return hos, err
	}

	for _, v := range resp.Data {
		hos = append(hos, &HistoryOrders{
			ID:               v.Id,
			ClientOrderID:    v.ClientOrderId,
			AccountID:        v.AccountId,
			UserID:           v.UserId,
			Symbol:           v.Symbol,
			Price:            v.Price,
			Amount:           v.Amount,
			CreatedAt:        v.CreatedAt,
			CanceledAt:       v.CanceledAt,
			FinishedAt:       v.FinishedAt,
			Type:             v.Type,
			FilledAmount:     v.FilledAmount,
			FilledCashAmount: v.FilledCashAmount,
			FilledFees:       v.FilledFees,
			Source:           v.Source,
			State:            v.State,
			Exchange:         v.Exchange,
			Batch:            v.Batch,
			StopPrice:        v.StopPrice,
			Operator:         v.Operator,
		})
	}

	return hos, nil
}

// GetHistoryOrdersInfo 获取历史订单列表
func (a *Account) GetHistoryOrdersInfo(symbol string, states string, types string) (interface{}, error) {
	return a.GetHistoryOrders(symbol, states, types)
}

// -------------------------------------------------------------------------------------------------

// Processer 处理订单接口
//type Processer interface {
//	UpdateGridOrder(tradeTime int64, tradeAmount float64, orderID, clientOrderID, orderStatus string) error
//}

// ProcessTradeOrder 处理订单
func ProcessTradeOrder(p goex.Processer) func(resp interface{}) {
	// Response biz
	return func(resp interface{}) {
		defer func() {
			if e := recover(); e != nil {
				logger.Error("ProcessTradeOrder panic", logger.Any("err", e))
			}
		}()

		subResponse, ok := resp.(order.SubscribeOrderV2Response)
		if !ok {
			logger.Warnf("received unknown response: %v", resp)
			return
		}

		switch subResponse.Action {
		case "sub":
			if !subResponse.IsSuccess() {
				logger.Errorf("subscription topic %s error, code: %d, message: %s", subResponse.Ch, subResponse.Code, subResponse.Message)
				return
			}
			logger.Infof("subscription topic %s successfully", subResponse.Ch)

		case "push":
			if subResponse.Data == nil {
				logger.Warnf("subscription response is empty")
				return
			}

			data := subResponse.Data
			switch data.EventType {
			case "creation":
				logger.Infof("[huobi] create order success, orderId=%d, clientOrderId=%s, symbol=%s, type=%s, status=%s", data.OrderId, data.ClientOrderId, data.Symbol, data.Type, data.OrderStatus)
			case "cancellation":
				logger.Infof("[huobi] cancel order success, orderId=%d, clientOrderId=%s", data.OrderId, data.ClientOrderId)
			case "trade":
				logger.Infof("[huobi] trade order success, orderId=%d, clientOrderId=%s, filledType=%s", data.OrderId, data.ClientOrderId, data.OrderStatus)

				// 处理订单
				err := p.UpdateGridOrder(
					data.TradeTime,
					str2Float64(data.TradePrice)*str2Float64(data.TradeVolume),
					fmt.Sprintf("%d", data.OrderId),
					data.ClientOrderId,
					data.OrderStatus,
				)
				if err != nil {
					logger.Error("UpdateGridOrder error",
						logger.Err(err),
						logger.String("param", fmt.Sprintf("%v,%v,%v,%v,%v",
							data.TradeTime,
							str2Float64(data.TradePrice)*str2Float64(data.TradeVolume),
							data.OrderId,
							data.ClientOrderId,
							data.OrderStatus,
						)))
				}

			default:
				logger.Warnf("unknown eventType, should never happen, orderId: %d, clientOrderId: %s, eventType: %s", data.OrderId, data.ClientOrderId, data.EventType)
			}
		}
	}
}

// WsSubscribeOrder web socket 订阅订单通知
func WsSubscribeOrder(wsClientID string, symbol string, responseHandler websocketclientbase.ResponseHandler, accessKey string, secretKey string) (*orderwebsocketclient.SubscribeOrderWebSocketV2Client, error) {
	var err error

	// Initialize a new instance
	cli := new(orderwebsocketclient.SubscribeOrderWebSocketV2Client).Init(accessKey, secretKey, host)

	// Connected biz
	connectedHandler := func(resp *auth.WebSocketV2AuthenticationResponse) {
		if !resp.IsSuccess() {
			err = fmt.Errorf("websocket authentication error, code: %d, message:%s", resp.Code, resp.Message)
			cli = nil
			return
		}
		// Subscribe if authentication passed
		cli.Subscribe(symbol, wsClientID)
	}

	// Set the callback handlers
	cli.SetHandler(connectedHandler, responseHandler)

	// Connect to the server and wait for the biz to handle the response
	cli.Connect(wsClientID)

	return cli, err
}
