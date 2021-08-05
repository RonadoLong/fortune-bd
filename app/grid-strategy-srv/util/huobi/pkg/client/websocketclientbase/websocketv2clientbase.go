package websocketclientbase

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal/gzip"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal/model"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal/requestbuilder"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/auth"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/base"

	"github.com/gorilla/websocket"
	"github.com/zhufuyi/pkg/logger"
)

const (
	websocketV2Path = "/ws/v2"
)

// It will be invoked after websocket v2 authentication response received
type AuthenticationV2ResponseHandler func(resp *auth.WebSocketV2AuthenticationResponse)

// The base class that responsible to get data from websocket authentication v2
type WebSocketV2ClientBase struct {
	host string
	conn *websocket.Conn

	authenticationResponseHandler AuthenticationV2ResponseHandler
	messageHandler                MessageHandler
	responseHandler               ResponseHandler

	stopReadChannel chan int
	//stopTickerChannel chan int
	//ticker            *time.Ticker
	//lastReceivedTime time.Time
	sendMutex *sync.Mutex

	requestBuilder *requestbuilder.WebSocketV2RequestBuilder

	WsClientID string // websocket标识id
}

// Initializer
func (p *WebSocketV2ClientBase) Init(accessKey string, secretKey string, host string) *WebSocketV2ClientBase {
	p.host = host
	p.stopReadChannel = make(chan int, 1)
	//p.stopTickerChannel = make(chan int, 1)
	p.requestBuilder = new(requestbuilder.WebSocketV2RequestBuilder).Init(accessKey, secretKey, host, websocketV2Path)
	p.sendMutex = &sync.Mutex{}
	return p
}

// Set callback biz
func (p *WebSocketV2ClientBase) SetHandler(authHandler AuthenticationV2ResponseHandler, msgHandler MessageHandler, repHandler ResponseHandler) {
	p.authenticationResponseHandler = authHandler
	p.messageHandler = msgHandler
	p.responseHandler = repHandler
}

// Connect to websocket server
// if autoConnect is true, then the connection can be re-connect if no data received after the pre-defined timeout
func (p *WebSocketV2ClientBase) Connect(wsClientID string) error {
	p.WsClientID = wsClientID
	if err := p.connectWebSocket(); err != nil {
		return err
	}
	p.startReadLoop()

	return nil
	//if autoConnect {
	//	p.startTicker()
	//}
}

// Send data to websocket server
func (p *WebSocketV2ClientBase) Send(data string) error {
	if p.conn == nil {
		return errors.New("WebSocket sent error: no connection available")
	}

	p.sendMutex.Lock()
	err := p.conn.WriteMessage(websocket.TextMessage, []byte(data))
	p.sendMutex.Unlock()

	if err != nil {
		logger.Error("WebSocket sent error", logger.Err(err), logger.String("data", data), logger.String("wsClientID", p.WsClientID))
		return err
	}

	return nil
}

// Close the connection to server
func (p *WebSocketV2ClientBase) Close() {
	p.disconnectWebSocket()
}

// connect to server
func (p *WebSocketV2ClientBase) connectWebSocket() error {
	var err error
	url := fmt.Sprintf("wss://%s%s", p.host, websocketV2Path)
	p.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Warn("WebSocket connected error", logger.Err(err), logger.String("wsClientID", p.WsClientID))
		return err
	}
	logger.Info("WebSocket connected success", logger.String("wsClientID", p.WsClientID))

	auth, err := p.requestBuilder.Build()
	if err != nil {
		logger.Warn("Signature generated error", logger.Err(err), logger.String("wsClientID", p.WsClientID))
		return err
	}

	err = p.Send(auth)
	if err != nil {
		logger.Warn("web socket auth failed", logger.Err(err), logger.String("wsClientID", p.WsClientID))
		return err
	}

	return nil
}

// disconnect with server
func (p *WebSocketV2ClientBase) disconnectWebSocket() {
	if p.conn == nil {
		return
	}

	// start a new goroutine to send a signal
	go p.stopReadLoop()

	err := p.conn.Close()
	if err != nil {
		logger.Error("WebSocket disconnect error", logger.Err(err), logger.String("wsClientID", p.WsClientID))
		return
	}

	logger.Info("WebSocket disconnected success", logger.String("wsClientID", p.WsClientID))
}

/*
// initialize a ticker and start a goroutine tickerLoop()
func (p *WebSocketV2ClientBase) startTicker() {
	p.ticker = time.NewTicker(TimerIntervalSecond * time.Second)
	p.lastReceivedTime = time.Now()

	go p.tickerLoop()
}

// stop ticker and stop the goroutine
func (p *WebSocketV2ClientBase) stopTicker() {
	p.ticker.Stop()
	p.stopTickerChannel <- 1
}

// defines a for loop that will run based on ticker's frequency
// It checks the last data that received from server, if it is longer than the threshold,
// it will force disconnect server and connect again.
func (p *WebSocketV2ClientBase) tickerLoop() {
	for {
		select {
		// start a goroutine readLoop()
		case <-p.stopTickerChannel:
			logger.Info("tickerLoop stopped")
			return

		// Receive tick from tickChannel
		case <-p.ticker.C:
			elapsedSecond := time.Now().Sub(p.lastReceivedTime).Seconds()
			//logger.Infof("WebSocket received data %f sec ago", elapsedSecond)

			if elapsedSecond > ReconnectWaitSecond {
				logger.Info("WebSocket reconnect ......")
				p.disconnectWebSocket()
				if err := p.connectWebSocket(); err != nil {
					logger.Error("connectWebSocket error", logger.Err(err))
					continue
				}
				p.startReadLoop()
			}
		}
	}
}
*/
// start a goroutine readLoop()
func (p *WebSocketV2ClientBase) startReadLoop() {
	go p.readLoop()
}

// stop the goroutine readLoop()
func (p *WebSocketV2ClientBase) stopReadLoop() {
	p.stopReadChannel <- 1
}

// defines a for loop to read data from server
// it will stop once it receives the signal from stopReadChannel
func (p *WebSocketV2ClientBase) readLoop() {
	count := 0
	isExit := false

	defer func() {
		if err := recover(); err != nil {
			logger.Error("receive market panic", logger.Any("error", err), logger.String("wsClientID", p.WsClientID))
		}

		if !isExit {
			time.Sleep(2 * time.Second) // 等待2秒
			if err := p.connectWebSocket(); err != nil {
				logger.Warn("reconnect websocket failed", logger.Err(err), logger.String("wsClientID", p.WsClientID), logger.String("wsAddr", fmt.Sprintf("wss://%s%s", p.host, websocketV2Path)))
			} else {
				logger.Info("reconnect websocket success", logger.String("wsClientID", p.WsClientID))
			}
			time.Sleep(3 * time.Second) // 等待3秒后重试
			p.readLoop()                // 迭代
		}
	}()

	if p.conn == nil {
		logger.Warn("no websocket connection available", logger.String("wsClientID", p.WsClientID))
		return
	}

	for {
		select {
		// Receive data from stopChannel
		case <-p.stopReadChannel:
			logger.Info("readLoop stopped", logger.String("wsClientID", p.WsClientID))
			isExit = true
			return

		default:
			msgType, buf, err := p.conn.ReadMessage()
			if err != nil {
				logger.Warn("websocket read error", logger.Err(err), logger.String("wsClientID", p.WsClientID))
				return
			}

			message, err := getWsMsg(msgType, buf)
			if err != nil {
				logger.Error("getWsMsg error", logger.Err(err), logger.String("wsClientID", p.WsClientID))
			}

			if p.checkIsPingMsg(&count, message) {
				continue
			}

			wsV2Resp := base.ParseWSV2Resp(message)
			if wsV2Resp != nil {
				switch wsV2Resp.Action {
				case "req":
					authResp := auth.ParseWSV2AuthResp(message)
					if authResp != nil && p.authenticationResponseHandler != nil {
						p.authenticationResponseHandler(authResp)
					}

				case "sub", "push":
					result, err := p.messageHandler(string(message))
					if err != nil {
						logger.Error("Handle message error", logger.Err(err), logger.String("wsClientID", p.WsClientID))
						continue
					}
					if p.responseHandler != nil {
						p.responseHandler(result)
					}
				}
			}
		}
	}
}

func getWsMsg(msgType int, buf []byte) ([]byte, error) {
	switch msgType {
	case websocket.BinaryMessage:
		message, err := gzip.GZipDecompress(buf)
		if err != nil {
			return nil, fmt.Errorf("GZipDecompress error, %s", err.Error())
		}
		return []byte(message), nil

	case websocket.TextMessage:
		return buf, nil

	default:
		return nil, fmt.Errorf("unknown websocket message type %d", msgType)
	}
}

func (p *WebSocketV2ClientBase) checkIsPingMsg(count *int, message []byte) bool {
	pingV2Msg := model.ParsePingV2Message(message)
	if pingV2Msg.IsPing() {
		pongMsg := fmt.Sprintf("{\"action\": \"pong\", \"data\": { \"ts\": %d } }", pingV2Msg.Data.Timestamp)
		err := p.Send(pongMsg)
		if err != nil {
			logger.Error("websocket send error", logger.Err(err), logger.String("data", pongMsg), logger.String("wsClientID", p.WsClientID))
		} else {
			*count++
			if *count%9 == 0 { // 大概3分钟打印一次
				*count = 0
				logger.Info("ping-pong success", logger.String("wsClientID", p.WsClientID))
			}
		}
		return false
	}

	return false
}
