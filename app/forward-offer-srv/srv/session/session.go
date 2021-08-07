package session

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wq-fotune-backend/app/forward-offer-srv/global"
	"wq-fotune-backend/app/forward-offer-srv/srv/cache_service"
	"wq-fotune-backend/app/forward-offer-srv/srv/model"
	"wq-fotune-backend/libs/exchangeclient"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/goex"

	jsoniter "github.com/json-iterator/go"
)

const (
	maxReqErr  = "too many requests"
	maxReqCode = "30014"
)

type Client struct {
	OnlineTime time.Time
	ApiClient  *exchangeclient.OKClient
	cancelFunc context.CancelFunc
	cxt        context.Context
}

func initClient(apiKey, secret, Passphrase string) *Client {
	var (
		err      error
		okClient *exchangeclient.OKClient
	)
	var tryCount = 3
	for {
		if tryCount == 0 {
			return nil
		}
		okClient = exchangeclient.InitOKEX(apiKey, secret, Passphrase)
		if err = okClient.BuildWS(); err != nil {
			if strings.Contains(err.Error(), maxReqCode) || strings.Contains(err.Error(), maxReqErr) {
				<-time.After(time.Millisecond * 50)
				logger.Warnf("登录交易所API受限制重试，%s", apiKey)
				tryCount--
				continue
			}
		}
		break
	}

	if err != nil {
		logger.Errorf("登录交易所失败: %s, content: %s", apiKey, err.Error())
		return nil
	}
	logger.Infof("登录交易所成功: %s", apiKey)
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		OnlineTime: time.Now().UTC(),
		ApiClient:  okClient,
		cxt:        ctx,
		cancelFunc: cancel,
	}
}

func InitClientNormal(info *model.ExchangeInfo) *Client {
	okClient := exchangeclient.InitOKEX(info.APIKey, info.SecretKey, info.EcPass)
	return &Client{ApiClient: okClient}
}

// SubscriptExchangeEvent 订阅交易事件
func (c *Client) SubscriptExchangeEvent() {
	// 避免还未登录
	time.Sleep(time.Second * 2)

	defer func() {
		if err := recover(); err != nil {
			logger.Warnf("SubscriptExchangeEvent panic: %v", err)
		} else {
			logger.Warnf("SubscriptExchangeEvent ws error exit func: %v", c.ApiClient.GetApiKey())
		}
		delClient(c.ApiClient.GetApiKey())
	}()
	cancel, cancelFunc := context.WithCancel(context.Background())
	c.ApiClient.SetOrderCallBack(func(order *goex.FutureOrder, s string) {
		c.ProcessCallBackOrder(order)
	})
	c.ApiClient.SetErrorCallBack(func(err error) {
		logger.Warnf("SetErrorCallBack: %v", err)
		cancelFunc()
	})
	c.ApiClient.SubscribeOrdersEvent("")
	<-cancel.Done()
}

// BroadCastOrder 转发订单
func (c *Client) BroadCastOrder(req *model.ExchangeReq) {
	if req.TradeReq.Value == "" {
		logger.Error("订单数据不能为空")
		return
	}
	var orderType = req.TradeReq.Type
	switch orderType {
	case global.CreateType:
		_ = c.CreateNormalOrder(req)
	case global.AutoAddType:
		_ = c.CreateReliableOrder(req)
	case global.CancelType:
		_ = c.CancelOrder(req)
	default:
		logger.Errorf("未知的订单类型: %v", req.Data)
	}
}

// CreateNormalOrder send normal order to exchange
func (c *Client) CreateNormalOrder(req interface{}) error {
	eReq := req.(*model.ExchangeReq)
	var orderReq model.OrderReq
	_ = jsoniter.UnmarshalFromString(eReq.TradeReq.Value, &orderReq)
	if orderReq.OrderQty <= float64(0) {
		logger.Errorf("提交订单数量不能为0")
		return nil
	}
	// 缓存订单id绑定用户id
	cache_service.CacheOrderIDBindUser(orderReq.OrderID, eReq.UserID, eReq.StrategyID)
	checkError := c.ProcessCreateOrder(orderReq, 1)
	if checkError != nil {
		logger.Errorf("提交订单错误: %s", checkError.Msg)
	}
	return nil
}

// CreateReliableOrder send reliable order to exchange
func (c *Client) CreateReliableOrder(req interface{}) error {
	eReq := req.(*model.ExchangeReq)
	var orderReq model.OrderReq
	_ = jsoniter.UnmarshalFromString(eReq.TradeReq.Value, &orderReq)
	if orderReq.OrderQty <= float64(0) {
		logger.Errorf("提交订单数量不能为0")
		return nil
	}
	// 发送订单
	cache_service.CacheOrderIDBindUser(orderReq.OrderID, eReq.UserID, eReq.StrategyID)
	checkError := c.ProcessCreateOrder(orderReq, 1)
	if checkError != nil {
		msg := fmt.Sprintf("提交订单出错, %v, 订单ID: %v", checkError.Msg, orderReq.OrderID)
		logger.Warnf("ProcessCreateOrder: %s", msg)
		if checkError.TypeCode != global.CollectOrderLimiterErrorCode {
			return nil
		}
	}
	// 10秒后检查订单
	delayerReq := model.CreateDelayerReq(orderReq, eReq.UserID, eReq.StrategyID)
	delayerReq.APIKey = c.ApiClient.GetApiKey()
	delayerReq.SecretKey = c.ApiClient.GetApiSecretKey()
	delayerReq.EcPass = c.ApiClient.GetApiPassphrase()
	cache_service.PushDelayerInfoToQueue(delayerReq)
	return nil
}

// CancelOrder send a cancel order to exchange
func (c *Client) CancelOrder(req interface{}) error {
	eReq := req.(*model.ExchangeReq)
	var orderReq model.OrderReq
	_ = jsoniter.UnmarshalFromString(eReq.TradeReq.Value, &orderReq)
	return nil
}
