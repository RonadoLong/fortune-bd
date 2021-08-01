package marketwebsocketclient

import (
	"encoding/json"
	"fmt"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/client/websocketclientbase"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/response/market"

	"github.com/zhufuyi/pkg/logger"
)

// Responsible to handle Trade data from WebSocket
type TradeWebSocketClient struct {
	websocketclientbase.WebSocketClientBase
}

// Initializer
func (p *TradeWebSocketClient) Init(host string) *TradeWebSocketClient {
	p.WebSocketClientBase.Init(host)
	return p
}

// Set callback handler
func (p *TradeWebSocketClient) SetHandler(
	connectedHandler websocketclientbase.ConnectedHandler,
	responseHandler websocketclientbase.ResponseHandler) {
	p.WebSocketClientBase.SetHandler(connectedHandler, p.handleMessage, responseHandler)
}

// Request latest 300 trade data
func (p *TradeWebSocketClient) Request(symbol string, clientId string) {
	topic := fmt.Sprintf("market.%s.trade.detail", symbol)
	req := fmt.Sprintf("{\"req\": \"%s\",\"id\": \"%s\" }", topic, clientId)

	p.Send(req)

	logger.Infof("WebSocket requested, topic=%s, clientId=%s", topic, clientId)
}

// Subscribe latest completed trade in tick by tick mode
func (p *TradeWebSocketClient) Subscribe(symbol string, clientId string) {
	topic := fmt.Sprintf("market.%s.trade.detail", symbol)
	sub := fmt.Sprintf("{\"sub\": \"%s\",\"id\": \"%s\" }", topic, clientId)

	p.Send(sub)

	logger.Infof("WebSocket subscribed, topic=%s, clientId=%s", topic, clientId)
}

// Unsubscribe trade
func (p *TradeWebSocketClient) UnSubscribe(symbol string, clientId string) {
	topic := fmt.Sprintf("market.%s.trade.detail", symbol)
	unsub := fmt.Sprintf("{\"unsub\": \"%s\",\"id\": \"%s\" }", topic, clientId)

	p.Send(unsub)

	logger.Infof("WebSocket unsubscribed, topic=%s, clientId=%s", topic, clientId)
}

func (p *TradeWebSocketClient) handleMessage(msg string) (interface{}, error) {
	result := market.SubscribeTradeResponse{}
	err := json.Unmarshal([]byte(msg), &result)
	return result, err
}
