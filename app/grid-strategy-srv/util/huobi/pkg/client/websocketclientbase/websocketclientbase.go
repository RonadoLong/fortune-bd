package websocketclientbase

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal/gzip"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal/model"

	"github.com/gorilla/websocket"
	"github.com/zhufuyi/pkg/logger"
)

const (
	TimerIntervalSecond = 5
	ReconnectWaitSecond = 60

	path = "/ws"
)

// It will be invoked after websocket connected
type ConnectedHandler func()

// It will be invoked after valid message received
type MessageHandler func(message string) (interface{}, error)

// It will be invoked after response is parsed
type ResponseHandler func(response interface{})

// The base class that responsible to get data from websocket
type WebSocketClientBase struct {
	host              string
	conn              *websocket.Conn
	connectedHandler  ConnectedHandler
	messageHandler    MessageHandler
	responseHandler   ResponseHandler
	stopReadChannel   chan int
	stopTickerChannel chan int
	ticker            *time.Ticker
	lastReceivedTime  time.Time
	sendMutex         *sync.Mutex
}

// Initializer
func (p *WebSocketClientBase) Init(host string) *WebSocketClientBase {
	p.host = host
	p.stopReadChannel = make(chan int, 1)
	p.stopTickerChannel = make(chan int, 1)
	p.sendMutex = &sync.Mutex{}

	return p
}

// Set callback biz
func (p *WebSocketClientBase) SetHandler(connHandler ConnectedHandler, msgHandler MessageHandler, repHandler ResponseHandler) {
	p.connectedHandler = connHandler
	p.messageHandler = msgHandler
	p.responseHandler = repHandler
}

// Connect to websocket server
// if autoConnect is true, then the connection can be re-connect if no data received after the pre-defined timeout
func (p *WebSocketClientBase) Connect(autoConnect bool) {
	p.connectWebSocket()

	if autoConnect {
		p.startTicker()
	}
}

// Send data to websocket server
func (p *WebSocketClientBase) Send(data string) {
	if p.conn == nil {
		logger.Errorf("WebSocket sent error: no connection available")
		return
	}

	p.sendMutex.Lock()
	err := p.conn.WriteMessage(websocket.TextMessage, []byte(data))
	p.sendMutex.Unlock()

	if err != nil {
		logger.Errorf("WebSocket sent error: data=%s, error=%s", data, err)
	}
}

// Close the connection to server
func (p *WebSocketClientBase) Close() {
	p.stopTicker()
	p.disconnectWebSocket()
}

// connect to server
func (p *WebSocketClientBase) connectWebSocket() {
	var err error
	url := fmt.Sprintf("wss://%s%s", p.host, path)
	p.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Error("WebSocket connected error", logger.Err(err))
		return
	}
	logger.Info("WebSocket connected success")

	p.startReadLoop()

	if p.connectedHandler != nil {
		p.connectedHandler()
	}
}

// disconnect with server
func (p *WebSocketClientBase) disconnectWebSocket() {
	if p.conn == nil {
		return
	}

	p.stopReadLoop()

	err := p.conn.Close()
	if err != nil {
		logger.Error("WebSocket disconnect error: %s", logger.Err(err))
		return
	}

	logger.Info("WebSocket disconnected success")
}

// initialize a ticker and start a goroutine tickerLoop()
func (p *WebSocketClientBase) startTicker() {
	p.ticker = time.NewTicker(TimerIntervalSecond * time.Second)
	p.lastReceivedTime = time.Now()

	go p.tickerLoop()
}

// stop ticker and stop the goroutine
func (p *WebSocketClientBase) stopTicker() {
	p.ticker.Stop()
	p.stopTickerChannel <- 1
}

// defines a for loop that will run based on ticker's frequency
// It checks the last data that received from server, if it is longer than the threshold,
// it will force disconnect server and connect again.
func (p *WebSocketClientBase) tickerLoop() {
	for {
		select {
		// Receive data from stopChannel
		case <-p.stopTickerChannel:
			logger.Info("tickerLoop stopped")
			return

		// Receive tick from tickChannel
		case <-p.ticker.C:
			elapsedSecond := time.Now().Sub(p.lastReceivedTime).Seconds()
			logger.Infof("WebSocket received data %f sec ago", elapsedSecond)

			if elapsedSecond > ReconnectWaitSecond {
				logger.Info("WebSocket reconnect ......")
				p.disconnectWebSocket()
				p.connectWebSocket()
			}
		}
	}
}

// start a goroutine readLoop()
func (p *WebSocketClientBase) startReadLoop() {
	go p.readLoop()
}

// stop the goroutine readLoop()
func (p *WebSocketClientBase) stopReadLoop() {
	p.stopReadChannel <- 1
}

// defines a for loop to read data from server
// it will stop once it receives the signal from stopReadChannel
func (p *WebSocketClientBase) readLoop() {

	for {
		select {
		// Receive data from stopChannel
		case <-p.stopReadChannel:
			logger.Info("readLoop stopped success")
			return

		default:
			if p.conn == nil {
				logger.Error("Read error: no connection available")
				time.Sleep(TimerIntervalSecond * time.Second)
				continue
			}

			msgType, buf, err := p.conn.ReadMessage()
			if err != nil {
				logger.Error("Read error: %s", logger.Err(err))
				time.Sleep(TimerIntervalSecond * time.Second)
				continue
			}

			p.lastReceivedTime = time.Now()

			// decompress gzip data if it is binary message
			if msgType == websocket.BinaryMessage {
				message, err := gzip.GZipDecompress(buf)
				if err != nil {
					logger.Error("UnGZip data error: %s", logger.Err(err))
					continue
				}

				// Try to pass as PingMessage
				pingMsg := model.ParsePingMessage(message)

				// If it is Ping then respond Pong
				if pingMsg != nil && pingMsg.Ping != 0 {
					logger.Infof("Received Ping: %d", pingMsg.Ping)
					pongMsg := fmt.Sprintf("{\"pong\": %d}", pingMsg.Ping)
					p.Send(pongMsg)
					logger.Infof("Replied Pong: %d", pingMsg.Ping)
				} else if strings.Contains(message, "tick") || strings.Contains(message, "data") {
					// If it contains expected string, then invoke message biz and response biz
					result, err := p.messageHandler(message)
					if err != nil {
						logger.Error("Handle message error", logger.Err(err))
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
