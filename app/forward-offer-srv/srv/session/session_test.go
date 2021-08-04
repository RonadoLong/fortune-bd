package session

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
	"wq-fotune-backend/app/forward-offer-srv/global"
	"wq-fotune-backend/app/forward-offer-srv/srv/model"
)

var (
	Passphrase = "test123"
	apikey     = "62dd3fc3-03d6-4e7e-a4fc-ce238df6aa40"
	secretkey  = "ABD4DC058DE1AEAAD5A7BC920BEB8E46" //6
	Symbol     = "ETH-USD-SWAP"
)

const TradeRequestQueue = "trade:%s:request"

func TestLogin(t *testing.T) {
	c := initClient(apikey, secretkey, Passphrase)
	c.SubscriptExchangeEvent()
}

type Data struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func TestClient_CreateReliableOrder(t *testing.T) {
	global.InitRedisClient("127.0.0.1:6379", "")
	val := &model.OrderReq{}
	val.OrderID = "OK" + global.GetUUID()
	val.Symbol = "ETH-USDT-SWAP"
	val.OrderQty = 0.1
	val.Price = float64(233.20)
	val.TryCount = 10
	val.Direction = "buy"

	marshalToString, _ := jsoniter.MarshalToString(val)
	d := &Data{
		Type:  "autoAdd",
		Value: marshalToString,
	}
	DataString, _ := jsoniter.MarshalToString(d)

	var tradeReq = model.ExchangeReq{}
	tradeReq.Exchange = "okex"
	tradeReq.Data = DataString
	tradeReq.UserID = "1258438045758132224"
	tradeReq.StrategyID = "yyyy"
	tradeReq.ExchangeInfo = &model.ExchangeInfo{
		APIKey:    apikey,
		SecretKey: secretkey,
		EcPass:    Passphrase,
	}
	key := fmt.Sprintf(TradeRequestQueue, "okex")
	bytes, _ := jsoniter.Marshal(tradeReq)
	global.PushReqOrderMessage(key, bytes)
}
