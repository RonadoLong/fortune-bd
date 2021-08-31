package orderwebsocketclient

import (
	"encoding/json"
	"fmt"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/client/websocketclientbase"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/response/order"

	"github.com/zhufuyi/pkg/logger"
)

// Responsible to handle order subscription from WebSocket
// This need authentication version 2
type SubscribeOrderWebSocketV2Client struct {
	websocketclientbase.WebSocketV2ClientBase
}

// Initializer
func (p *SubscribeOrderWebSocketV2Client) Init(accessKey string, secretKey string, host string) *SubscribeOrderWebSocketV2Client {
	p.WebSocketV2ClientBase.Init(accessKey, secretKey, host)
	return p
}

// Set callback biz
func (p *SubscribeOrderWebSocketV2Client) SetHandler(
	authHandler websocketclientbase.AuthenticationV2ResponseHandler,
	responseHandler websocketclientbase.ResponseHandler) {
	p.WebSocketV2ClientBase.SetHandler(authHandler, p.handleMessage, responseHandler)
}

// Subscribe all order updates of the current account
func (p *SubscribeOrderWebSocketV2Client) Subscribe(symbol string, clientId string) {
	channel := fmt.Sprintf("orders#%s", symbol)
	sub := fmt.Sprintf("{\"action\":\"sub\", \"ch\":\"%s\", \"cid\": \"%s\"}", channel, clientId)

	p.Send(sub)

	logger.Infof("WebSocket subscribed, channel=%s, clientId=%s", channel, clientId)
}

// Unsubscribe order updates
func (p *SubscribeOrderWebSocketV2Client) UnSubscribe(symbol string, clientId string) {
	channel := fmt.Sprintf("orders#%s", symbol)
	unsub := fmt.Sprintf("{\"action\":\"unsub\", \"ch\":\"%s\", \"cid\": \"%s\"}", channel, clientId)

	p.Send(unsub)

	logger.Infof("WebSocket unsubscribed, channel=%s, clientId=%s", channel, clientId)
}

func (p *SubscribeOrderWebSocketV2Client) handleMessage(msg string) (interface{}, error) {
	result := order.SubscribeOrderV2Response{}
	err := json.Unmarshal([]byte(msg), &result)
	return result, err
}
