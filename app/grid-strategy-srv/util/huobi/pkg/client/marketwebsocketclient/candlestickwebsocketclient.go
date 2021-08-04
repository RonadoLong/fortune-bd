package marketwebsocketclient

import (
	"encoding/json"
	"fmt"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/client/websocketclientbase"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/market"

	"github.com/zhufuyi/pkg/logger"
)

// Responsible to handle candlestick data from WebSocket
type CandlestickWebSocketClient struct {
	websocketclientbase.WebSocketClientBase
}

// Initializer
func (p *CandlestickWebSocketClient) Init(host string) *CandlestickWebSocketClient {
	p.WebSocketClientBase.Init(host)
	return p
}

// Set callback service
func (p *CandlestickWebSocketClient) SetHandler(
	connectedHandler websocketclientbase.ConnectedHandler,
	responseHandler websocketclientbase.ResponseHandler) {
	p.WebSocketClientBase.SetHandler(connectedHandler, p.handleMessage, responseHandler)
}

// Request the full candlestick data according to specified criteria
func (p *CandlestickWebSocketClient) Request(symbol string, period string, from int64, to int64, clientId string) {
	topic := fmt.Sprintf("market.%s.kline.%s", symbol, period)
	req := fmt.Sprintf("{\"req\": \"%s\", \"from\":%d, \"to\":%d, \"id\": \"%s\" }", topic, from, to, clientId)

	p.Send(req)

	logger.Infof("WebSocket requested, topic=%s, clientId=%s", topic, clientId)
}

// Subscribe candlestick data
func (p *CandlestickWebSocketClient) Subscribe(symbol string, period string, clientId string) {
	topic := fmt.Sprintf("market.%s.kline.%s", symbol, period)
	sub := fmt.Sprintf("{\"sub\": \"%s\", \"id\": \"%s\"}", topic, clientId)

	p.Send(sub)

	logger.Infof("WebSocket subscribed, topic=%s, clientId=%s", topic, clientId)
}

// Unsubscribe candlestick data
func (p *CandlestickWebSocketClient) UnSubscribe(symbol string, period string, clientId string) {
	topic := fmt.Sprintf("market.%s.kline.%s", symbol, period)
	unsub := fmt.Sprintf("{\"unsub\": \"%s\", \"id\": \"%s\" }", topic, clientId)

	p.Send(unsub)

	logger.Infof("WebSocket unsubscribed, topic=%s, clientId=%s", topic, clientId)
}

func (p *CandlestickWebSocketClient) handleMessage(msg string) (interface{}, error) {
	result := market.SubscribeCandlestickResponse{}
	err := json.Unmarshal([]byte(msg), &result)
	return result, err
}
