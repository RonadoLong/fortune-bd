package router

import (
	"context"
	"fortune-bd/api/constant"
	pb "fortune-bd/api/quote/v1"
	"fortune-bd/api/response"
	"fortune-bd/app/quote-svc/cron"
	"fortune-bd/app/quote-svc/internal/service"
	"fortune-bd/libs/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/http"
	"time"
)

var (
	quoteService = service.NewQuoteService()
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func v1api(group *gin.RouterGroup) {
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
		for {
			select {
			case <-ctxOut.Done():
				return
			default:
				// 1. 不断获取行情数据
				resp := quoteService.GetTicks(ctxOut, &pb.GetTicksReq{Exchange: exchange}) //
				if resp == nil {
					time.Sleep(1 * time.Second)
					continue
				}
				// 2. 转发到ws中
				var ticks []cron.Ticker
				if err := jsoniter.Unmarshal(resp.Ticks, &ticks); err != nil {
					logger.Warnf("StreamHandler:Unmarshal数据失败")
					errMsg := response.NewResultInternalErr("StreamHandler:Unmarshal数据失败")
					_ = ws.WriteJSON(errMsg)
					time.Sleep(5 * time.Second)
					continue
				}

				err := ws.WriteJSON(response.NewResultSuccess(ticks))
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
	//todo 待优化代码
	go func(ws1 *websocket.Conn) {
		ctx := context.Background()
		ctxOut, cancelOut := context.WithCancel(ctx)
		ctxRun, cancelRun := context.WithCancel(ctx)
		msg := []byte(constant.BINANCE)
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
				if string(msg) == constant.OKEX || string(msg) == constant.BINANCE || string(msg) == constant.HUOBI {
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
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	var ticks = make(map[string]map[string]interface{})
	if err := jsoniter.Unmarshal(resp.Ticks, &ticks); err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, ticks)
}
