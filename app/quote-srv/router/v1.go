package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/v2/errors"
	"log"
	"net/http"
	"time"
	pb "wq-fotune-backend/api/quote"
	"wq-fotune-backend/app/quote-srv/client"
	"wq-fotune-backend/app/quote-srv/cron"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	exchange_info "wq-fotune-backend/pkg/exchange-info"
	"wq-fotune-backend/pkg/response"
)

var (
	quoteService pb.QuoteService
	upGrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func v1api(group *gin.RouterGroup) {
	quoteService = client.NewQuoteClient(env.EtcdAddr)
	group.GET("/ticks", GetTicks)
	group.GET("/ticks/realtime", SubRealTimeTickers)
}

// SubRealTimeTickers 实时获取行情数据
func SubRealTimeTickers(c *gin.Context) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		logger.Warnf("ws upgrader conn err: %v", err)
		response.NewErrorCreate(c, "网络出错", nil)
		return
	}
	go StreamHandler(conn)
}

func StreamHandler(ws1 *websocket.Conn) {
	//即便 我们不再 期望 来自websocket 更多的请求,我们仍需要 去 websocket 读取 内容，为了能获取到 close 信号
	run := func(ws *websocket.Conn, exchange string, ctxOut context.Context, cancelClose context.CancelFunc) {
		// 发送请求 给 stream server
		service, err := quoteService.StreamOkexTicks(context.Background(), &pb.GetTicksReq{Exchange: exchange})
		if err != nil {
			logger.Warnf("行情服务端连接失败 %v", err)
			errMsg := response.NewResultInternalErr("行情服务端连接失败")
			_ = ws.WriteJSON(errMsg)
			return
		}
		defer func() {
			service.Close()
		}()
		for {
			select {
			case <-ctxOut.Done():
				return
			default:
				// 1. 不断获取行情数据
				resp, err := service.Recv()
				if resp == nil {
					time.Sleep(1 * time.Second)
					continue
				}
				if err != nil {
					logger.Warnf("行情服务recv数据失败 %v", err)
					errMsg := response.NewResultInternalErr("行情服务recv数据失败")
					_ = ws.WriteJSON(errMsg)
					time.Sleep(2 * time.Second)
					continue
				}
				// 2. 转发到ws中
				//TODO 币本位行情
				//ticks := make([]map[string][]cron.Ticker, 0)
				var ticks []cron.Ticker
				if err := jsoniter.Unmarshal(resp.Ticks, &ticks); err != nil {
					logger.Warnf("StreamHandler:Unmarshal数据失败")
					errMsg := response.NewResultInternalErr("StreamHandler:Unmarshal数据失败")
					_ = ws.WriteJSON(errMsg)
					time.Sleep(5 * time.Second)
					continue
				}

				err = ws.WriteJSON(response.NewResultSuccess(ticks))
				if err != nil {
					if isExpectedClose(err) {
						logger.Warnf("expected close on socket")
					}
					logger.Warnf("ws1 writeErr: %v", err)
					cancelClose()
					return
				}
				time.Sleep(2 * time.Second)
			}

		}
	}

	go func(ws1 *websocket.Conn) {
		ctx := context.Background()
		ctxOut, cancelOut := context.WithCancel(ctx)
		ctxRun, cancelRun := context.WithCancel(ctx)
		msg := []byte(exchange_info.BINANCE)
		var err error
		go run(ws1, string(msg), ctxOut, cancelRun)
		//open := true
		time.Sleep(2 * time.Second)
		defer ws1.Close()
		for {
			select {
			case <-ctxRun.Done():
				return
			default:
				_, msg, err = ws1.ReadMessage()
				if err != nil {
					log.Println(err)
					cancelOut()
					return
				}
				if string(msg) == exchange_info.OKEX || string(msg) == exchange_info.BINANCE || string(msg) == exchange_info.HUOBI {
					cancelOut()
					time.Sleep(2 * time.Second)
					ctxOut, cancelOut = context.WithCancel(ctx)
					ctxRun, cancelRun = context.WithCancel(ctx)
					go run(ws1, string(msg), ctxOut, cancelRun)
				}
			}
		}
	}(ws1)
}

func isExpectedClose(err error) bool {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
		logger.Warnf("Unexpected websocket close: %v", err)
		return false
	}
	return true
}

// GetTicks 获取行情数据
func GetTicks(c *gin.Context) {
	resp, err := quoteService.GetTicksWithExchange(context.Background(), &pb.GetTicksReq{All: false})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	//var ticks []cron.Ticker
	var ticks = make(map[string]map[string]interface{})
	if err := jsoniter.Unmarshal(resp.Ticks, &ticks); err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, ticks)
}
