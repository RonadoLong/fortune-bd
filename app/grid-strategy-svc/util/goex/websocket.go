package goex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zhufuyi/pkg/logger"
)

type WsConfig struct {
	WsUrl                          string
	ProxyUrl                       string
	ReqHeaders                     map[string][]string //连接的时候加入的头部信息
	HeartbeatIntervalTime          time.Duration       //
	HeartbeatData                  func() []byte       //心跳数据2
	IsAutoReconnect                bool
	ProtoHandleFunc                func([]byte) error           //协议处理函数
	DecompressFunc                 func([]byte) ([]byte, error) //解压函数
	ErrorHandleFunc                func(err error)
	ConnectSuccessAfterSendMessage func() []byte //for reconnect
	IsDump                         bool
	readDeadLineTime               time.Duration
	reconnectInterval              time.Duration
}

var dialer = &websocket.Dialer{
	Proxy:             http.ProxyFromEnvironment,
	HandshakeTimeout:  30 * time.Second,
	EnableCompression: true,
}

type WsConn struct {
	c *websocket.Conn
	WsConfig
	writeBufferChan        chan []byte
	pingMessageBufferChan  chan []byte
	pongMessageBufferChan  chan []byte
	closeMessageBufferChan chan []byte
	subs                   [][]byte
	//close                  chan bool
	reConnectLock *sync.Mutex

	Ctx       context.Context
	Cancel    context.CancelFunc
	isExist   bool   // 是否已经退出了，true:是，false:否
	ListenKey string // web socket 标识
}

type WsBuilder struct {
	wsConfig *WsConfig
}

func NewWsBuilder() *WsBuilder {
	return &WsBuilder{&WsConfig{
		ReqHeaders:        make(map[string][]string, 1),
		reconnectInterval: time.Second * 10,
		IsAutoReconnect:   true,
	}}
}

func (b *WsBuilder) WsUrl(wsUrl string) *WsBuilder {
	b.wsConfig.WsUrl = wsUrl
	return b
}

func (b *WsBuilder) ProxyUrl(proxyUrl string) *WsBuilder {
	b.wsConfig.ProxyUrl = proxyUrl
	return b
}

func (b *WsBuilder) ReqHeader(key, value string) *WsBuilder {
	b.wsConfig.ReqHeaders[key] = append(b.wsConfig.ReqHeaders[key], value)
	return b
}

func (b *WsBuilder) AutoReconnect() *WsBuilder {
	b.wsConfig.IsAutoReconnect = true
	return b
}

func (b *WsBuilder) Dump() *WsBuilder {
	b.wsConfig.IsDump = true
	return b
}

func (b *WsBuilder) Heartbeat(heartbeat func() []byte, t time.Duration) *WsBuilder {
	b.wsConfig.HeartbeatIntervalTime = t
	b.wsConfig.HeartbeatData = heartbeat
	return b
}

func (b *WsBuilder) ReconnectInterval(t time.Duration) *WsBuilder {
	b.wsConfig.reconnectInterval = t
	return b
}

func (b *WsBuilder) ProtoHandleFunc(f func([]byte) error) *WsBuilder {
	b.wsConfig.ProtoHandleFunc = f
	return b
}

func (b *WsBuilder) DecompressFunc(f func([]byte) ([]byte, error)) *WsBuilder {
	b.wsConfig.DecompressFunc = f
	return b
}

func (b *WsBuilder) ErrorHandleFunc(f func(err error)) *WsBuilder {
	b.wsConfig.ErrorHandleFunc = f
	return b
}

func (b *WsBuilder) ConnectSuccessAfterSendMessage(msg func() []byte) *WsBuilder {
	b.wsConfig.ConnectSuccessAfterSendMessage = msg
	return b
}

func (b *WsBuilder) Build() *WsConn {
	wsConn := &WsConn{WsConfig: *b.wsConfig}
	return wsConn.NewWs()
}

func (ws *WsConn) NewWs() *WsConn {
	if ws.HeartbeatIntervalTime == 0 {
		ws.readDeadLineTime = time.Minute
	} else {
		ws.readDeadLineTime = ws.HeartbeatIntervalTime * 2
	}

	if err := ws.connect(); err != nil {
		logger.Error(fmt.Sprintf("[%s] %s", ws.WsUrl, err.Error()))
		return nil
	}

	//ws.close = make(chan bool, 1)
	ws.pingMessageBufferChan = make(chan []byte, 10)
	ws.pongMessageBufferChan = make(chan []byte, 10)
	ws.closeMessageBufferChan = make(chan []byte, 10)
	ws.writeBufferChan = make(chan []byte, 10)
	ws.reConnectLock = new(sync.Mutex)

	ws.Ctx, ws.Cancel = context.WithCancel(context.Background())
	ws.SetIsExist(false)

	go ws.writeRequest()
	go ws.receiveMessage()

	if ws.ConnectSuccessAfterSendMessage != nil {
		msg := ws.ConnectSuccessAfterSendMessage()
		ws.SendMessage(msg)
		logger.Infof("[ws] [%s] execute the connect success after send message=%s", ws.WsUrl, string(msg))
	}

	return ws
}

func (ws *WsConn) connect() error {
	if ws.ProxyUrl != "" {
		proxy, err := url.Parse(ws.ProxyUrl)
		if err == nil {
			logger.Infof("[ws][%s] proxy url:%s", ws.WsUrl, proxy)
			dialer.Proxy = http.ProxyURL(proxy)
		} else {
			logger.Errorf("[ws][%s]parse proxy url [%s] err %s  ", ws.WsUrl, ws.ProxyUrl, err.Error())
		}
	}

	wsConn, resp, err := dialer.Dial(ws.WsUrl, http.Header(ws.ReqHeaders))
	if err != nil {
		//logger.Errorf("[ws][%s] %s", ws.WsUrl, err.Error())
		if ws.IsDump && resp != nil {
			dumpData, _ := httputil.DumpResponse(resp, true)
			logger.Debugf("[ws][%s] %s", ws.WsUrl, string(dumpData))
		}
		return err
	}

	//wsConn.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))

	if ws.IsDump {
		dumpData, _ := httputil.DumpResponse(resp, true)
		logger.Debugf("[ws][%s] %s", ws.WsUrl, string(dumpData))
	}
	logger.Info("[ws] web socket connected success")
	ws.c = wsConn
	return nil
}

func (ws *WsConn) reconnect() {
	ws.reConnectLock.Lock()
	defer ws.reConnectLock.Unlock()

	ws.c.Close() //主动关闭一次
	var err error
	for retry := 1; retry <= 100; retry++ {
		err = ws.connect()
		if err != nil {
			logger.Warnf("[ws] [%s] websocket reconnect failed, %s", ws.WsUrl, err.Error())
		} else {
			break
		}
		time.Sleep(ws.WsConfig.reconnectInterval * time.Duration(retry))
	}

	if err != nil {
		logger.Errorf("[ws] [%s] retry connect 100 count failed, begin exiting. ", ws.WsUrl)
		ws.CloseWs()
		if ws.ErrorHandleFunc != nil {
			ws.ErrorHandleFunc(errors.New("retry reconnect fail"))
		}
	} else {
		//re subscribe
		if ws.ConnectSuccessAfterSendMessage != nil {
			msg := ws.ConnectSuccessAfterSendMessage()
			ws.SendMessage(msg)
			logger.Infof("[ws] [%s] execute the connect success after send message=%s", ws.WsUrl, string(msg))
			time.Sleep(time.Second) //wait response
		}

		for _, sub := range ws.subs {
			logger.Infof("[ws] re subscribe: %s", string(sub))
			ws.SendMessage(sub)
		}
	}
}

func (ws *WsConn) writeRequest() {
	var (
		heartTimer *time.Timer
		err        error
	)

	if ws.HeartbeatIntervalTime == 0 {
		heartTimer = time.NewTimer(time.Hour)
	} else {
		heartTimer = time.NewTimer(ws.HeartbeatIntervalTime)
	}

	for {
		select {
		//case <-ws.close:
		//	logger.Infof("[ws][%s] close websocket , exiting write message goroutine.", ws.WsUrl)
		//	return
		case <-ws.Ctx.Done():
			logger.Infof("[ws exist][%s] close websocket , exiting receive message goroutine.", ws.WsUrl)
			return
		case d := <-ws.writeBufferChan:
			err = ws.c.WriteMessage(websocket.TextMessage, d)
		case d := <-ws.pingMessageBufferChan:
			err = ws.c.WriteMessage(websocket.PingMessage, d)
		case d := <-ws.pongMessageBufferChan:
			logger.Debugf("wsClient ***%s send pong %s", ws.ListenKey, string(d))
			err = ws.c.WriteMessage(websocket.PongMessage, d)
		case d := <-ws.closeMessageBufferChan:
			err = ws.c.WriteMessage(websocket.CloseMessage, d)
		case <-heartTimer.C:
			if ws.HeartbeatIntervalTime > 0 {
				err = ws.c.WriteMessage(websocket.TextMessage, ws.HeartbeatData())
				heartTimer.Reset(ws.HeartbeatIntervalTime)
			}
		}

		if err != nil {
			logger.Errorf("[ws][%s] write message %s", ws.WsUrl, err.Error())
			return
			//time.Sleep(time.Second)
		}
	}
}

func (ws *WsConn) Subscribe(subEvent interface{}) error {
	data, err := json.Marshal(subEvent)
	if err != nil {
		logger.Errorf("[ws][%s] json encode error , %s", ws.WsUrl, err)
		return err
	}
	ws.writeBufferChan <- data
	ws.subs = append(ws.subs, data)
	return nil
}

func (ws *WsConn) SendMessage(msg []byte) {
	if ws.CheckIsExist() {
		return
	}
	ws.writeBufferChan <- msg
}

func (ws *WsConn) SendPingMessage(msg []byte) {
	if ws.CheckIsExist() {
		return
	}
	ws.pingMessageBufferChan <- msg
}

func (ws *WsConn) SendPongMessage(msg []byte) {
	if ws.CheckIsExist() {
		return
	}
	ws.pongMessageBufferChan <- msg
}

func (ws *WsConn) SendCloseMessage(msg []byte) {
	if ws.CheckIsExist() {
		return
	}
	ws.closeMessageBufferChan <- msg
}

func (ws *WsConn) SendJsonMessage(m interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	ws.writeBufferChan <- data
	return nil
}

func (ws *WsConn) receiveMessage() {
	//exit
	ws.c.SetCloseHandler(func(code int, text string) error {
		logger.Warnf("[ws][%s] websocket exiting [code=%d , text=%s]", ws.WsUrl, code, text)
		//ws.CloseWs()
		return nil
	})

	ws.c.SetPongHandler(func(pong string) error {
		//logger.Debugf("[%s] received [pong] %s", ws.WsUrl, pong)
		//ws.c.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
		return nil
	})

	ws.c.SetPingHandler(func(ping string) error {
		//logger.Debugf("[%s] received [ping] %s", ws.WsUrl, ping)
		//ws.c.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
		return nil
	})

	for {
		select {
		//case <-ws.close:
		//	logger.Infof("[ws][%s] close websocket , exiting receive message goroutine.", ws.WsUrl)
		//	return
		case <-ws.Ctx.Done():
			logger.Infof("[ws exist][%s] close websocket , exiting receive message goroutine.", ws.WsUrl)
			return
		default:
			t, msg, err := ws.c.ReadMessage()
			if err != nil {
				logger.Infof("[ws] Unexpected Closed , err=%s.", err.Error())
				if ws.CheckIsExist() {
					return
				}
				logger.Info("[ws] reconnect web socket")
				ws.reconnect()
				continue
				//if ws.IsAutoReconnect {
				//	logger.Infof("[ws][%s] Unexpected Closed , Begin Retry Connect.", ws.WsUrl)
				//	ws.reconnect()
				//	continue
				//}
				//if ws.ErrorHandleFunc != nil {
				//	ws.ErrorHandleFunc(err)
				//}
				//return
			}

			//ws.c.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
			switch t {
			case websocket.TextMessage:
				if strings.Contains(string(msg), "Invalid sign") {
					logger.Warnf("ws call err: %+v", msg)
					ws.ErrorHandleFunc(errors.New(string(msg)))
					//ws.close <- true
					continue
				}
				_ = ws.ProtoHandleFunc(msg)
			case websocket.BinaryMessage:
				if ws.DecompressFunc == nil {
					ws.ProtoHandleFunc(msg)
				} else {
					msg2, err := ws.DecompressFunc(msg)
					if err != nil {
						logger.Errorf("[ws][%s] decompress error %s", ws.WsUrl, err.Error())
					} else {
						logger.Infof("receive call back msg %s", string(msg2))
						ws.ProtoHandleFunc(msg2)
					}
				}
				//	case websocket.CloseMessage:
				//	ws.CloseWs()

			default:
				logger.Errorf("[ws][%s] error websocket message type , content is :\n %s \n", ws.WsUrl, string(msg))
				return
			}
		}
	}
}

func (ws *WsConn) CheckIsExist() bool {
	ws.reConnectLock.Lock()
	defer ws.reConnectLock.Unlock()
	return ws.isExist
}

func (ws *WsConn) SetIsExist(isClose bool) {
	ws.reConnectLock.Lock()
	defer ws.reConnectLock.Unlock()
	ws.isExist = isClose
}

func (ws *WsConn) CloseWs() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Errorf("CloseWs panic, err=%v", e)
			}
		}()

		ws.SetIsExist(true)
		ws.Cancel()
		if ws.c != nil {
			err := ws.c.Close()
			if err != nil {
				logger.Errorf("[ws][" + ws.WsUrl + "] close websocket error ," + err.Error())
			} else {
				logger.Info("close web socket success")
			}
		}
	}()
}

func (ws *WsConn) clearChannel(c chan struct{}) {
	for {
		if len(c) > 0 {
			<-c
		} else {
			break
		}
	}
}
