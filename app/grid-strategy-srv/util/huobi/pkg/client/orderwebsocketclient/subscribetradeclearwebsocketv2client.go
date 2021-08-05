package orderwebsocketclient

import (
	"encoding/json"
	"fmt"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/client/websocketclientbase"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/order"

	"github.com/zhufuyi/pkg/logger"
)

// Responsible to handle trade clear from WebSocket
// This need authentication version 2
type SubscribeTradeClearWebSocketV2Client struct {
	websocketclientbase.WebSocketV2ClientBase
}

// Initializer
func (p *SubscribeTradeClearWebSocketV2Client) Init(accessKey string, secretKey string, host string) *SubscribeTradeClearWebSocketV2Client {
	p.WebSocketV2ClientBase.Init(accessKey, secretKey, host)
	return p
}

// Set callback biz
func (p *SubscribeTradeClearWebSocketV2Client) SetHandler(
	authHandler websocketclientbase.AuthenticationV2ResponseHandler,
	responseHandler websocketclientbase.ResponseHandler) {
	p.WebSocketV2ClientBase.SetHandler(authHandler, p.handleMessage, responseHandler)
}

// Subscribe trade details including transaction fee and transaction fee deduction etc.
// It only updates when transaction occurs.
func (p *SubscribeTradeClearWebSocketV2Client) Subscribe(symbol string, clientId string) {
	channel := fmt.Sprintf("trade.clearing#%s", symbol)
	sub := fmt.Sprintf("{\"action\":\"sub\", \"ch\":\"%s\", \"cid\": \"%s\"}", channel, clientId)

	p.Send(sub)

	logger.Infof("WebSocket subscribed, channel=%s, clientId=%s", channel, clientId)
}

// Unsubscribe trade update
func (p *SubscribeTradeClearWebSocketV2Client) UnSubscribe(symbol string, clientId string) {
	channel := fmt.Sprintf("trade.clearing#%s", symbol)
	unsub := fmt.Sprintf("{\"action\":\"unsub\", \"ch\":\"%s\", \"cid\": \"%s\"}", channel, clientId)

	p.Send(unsub)

	logger.Infof("WebSocket unsubscribed, channel=%s, clientId=%s", channel, clientId)
}

func (p *SubscribeTradeClearWebSocketV2Client) handleMessage(msg string) (interface{}, error) {
	result := order.SubscribeTradeClearResponse{}
	err := json.Unmarshal([]byte(msg), &result)
	return result, err
}
