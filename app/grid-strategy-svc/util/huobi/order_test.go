// https://wq-grid-strategy/util/huobi

package huobi

import (
	"fmt"
	"testing"
	"github.com/k0kubun/pp"
	"github.com/zhufuyi/pkg/krand"
	"github.com/zhufuyi/pkg/logger"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/client/orderwebsocketclient"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/auth"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/order"
)

func TestPlaceLimitOrder(t *testing.T) {
	side := "buy"
	symbol := "btcusdt"
	price := "9000.0"
	amount := "0.0006"
	clientOrderID := fmt.Sprintf("g_4_b_%s", krand.String(krand.R_NUM, 12))

	orderID, err := account.PlaceLimitOrder(side, symbol, price, amount, clientOrderID)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(orderID)
}

func TestPlaceMarketOrder(t *testing.T) {
	side := "sell"
	symbol := "btcusdt"
	amount := "5.0"
	clientOrderID := fmt.Sprintf("g_4_b_%s", krand.String(krand.R_NUM, 12))

	orderID, err := account.PlaceMarketOrder(side, symbol, amount, clientOrderID)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(orderID)
}

func TestPlaceSpotOrders(t *testing.T) {
	params := PlaceParams{
		{
			AccountId:     "13867389",
			Type:          "buy-limit",
			Symbol:        "btcusdt",
			Price:         "9000.0",
			Amount:        "0.0001",
			ClientOrderId: fmt.Sprintf("g_1_b_%s", krand.String(krand.R_NUM, 12)),
		},
		{
			AccountId:     "13867389",
			Type:          "sell-market",
			Symbol:        "btcusdt",
			Amount:        "1.0",
			ClientOrderId: fmt.Sprintf("g_2_b_%s", krand.String(krand.R_NUM, 12)),
		},
	}

	results, err := account.PlaceSpotOrders(params)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(results)
}

func TestCancelOrder(t *testing.T) {
	err := account.CancelOrder("63374048844630", "")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("cancel order success")
}

func TestCancelOrders(t *testing.T) {
	orderIDs := []string{"48273545127461", "48274518257692"}
	results, err := account.CancelOrders(orderIDs)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(results)
}

func TestGetOrderInfo(t *testing.T) {
	orderInfo, err := account.GetOrder("77588038046460")
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(orderInfo)
}

func TestGetHistoryOrders(t *testing.T) {
	symbol := "btcusdt"
	states := OrderStateSubmitted
	types := "" // 空表示忽略

	results, err := account.GetHistoryOrders(symbol, states, types)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(len(results), results)
}

func TestWsSubscribeOrder(t *testing.T) {
	wsClientID := "e40001afeffd_574656_btcusdt_" + string(krand.String(7, 4))

	// Initialize a new instance
	cli := new(orderwebsocketclient.SubscribeOrderWebSocketV2Client).Init(accessKey, secretKey, host)

	// Connected biz
	connectedHandler := func(resp *auth.WebSocketV2AuthenticationResponse) {
		if resp.IsSuccess() {
			// Subscribe if authentication passed
			cli.Subscribe("btcusdt", wsClientID)
		} else {
			logger.Errorf("authentication error, code: %d, message:%s", resp.Code, resp.Message)
		}
	}

	// Response biz
	responseHandler := func(resp interface{}) {
		subResponse, ok := resp.(order.SubscribeOrderV2Response)
		if !ok {
			logger.Errorf("received unknown response: %v", resp)
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
			logger.Infof("order update, event: %s, symbol: %s, type: %s, status: %s", data.EventType, data.Symbol, data.Type, data.OrderStatus)
			pp.Println(data)
			switch data.EventType {
			case "creation":
				logger.Debugf("order created, orderId: %d, clientOrderId: %s", data.OrderId, data.ClientOrderId)
			case "cancellation":
				logger.Debugf("order cancelled, orderId: %d, clientOrderId: %s", data.OrderId, data.ClientOrderId)
			case "trade":
				logger.Debugf("order filled, orderId: %d, clientOrderId: %s, fill type: %s", data.OrderId, data.ClientOrderId, data.OrderStatus)

			default:
				logger.Warnf("unknown eventType, should never happen, orderId: %d, clientOrderId: %s, eventType: %s", data.OrderId, data.ClientOrderId, data.EventType)
			}
		}
	}

	// Set the callback handlers
	cli.SetHandler(connectedHandler, responseHandler)

	// Connect to the server and wait for the biz to handle the response
	cli.Connect(wsClientID)

	fmt.Println("waiting ......")

	wait := make(chan bool)
	<-wait
}

type ug struct {
}

func (u *ug) UpdateGridOrder(tradeTime int64, tradeAmount float64, orderID, clientOrderID, orderStatus string) error {
	fmt.Println("---------UpdateGridOrder--------", tradeTime, tradeAmount, orderID, clientOrderID, orderStatus)
	return nil
}

func TestWsSubscribeOrder2(t *testing.T) {
	symbol := "btcusdt"
	wsClientID := "e40001afeffd_574656_btcusdt_" + string(krand.String(7, 4))
	u := &ug{}

	processFun := ProcessTradeOrder(u)
	_, err := WsSubscribeOrder(wsClientID, symbol, processFun, accessKey, secretKey)
	if err != nil {
		t.Error(err)
		return
	}

	select {}
}
