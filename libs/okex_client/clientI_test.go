package api

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/internal/forward-offer-srv/global"
	"wq-fotune-backend/internal/forward-offer-srv/srv/model"
)

var (
	//apiKey  = "a527c56a-c990-4228-a455-f4e2b945d809"
	//apiSe   = "99F073DA29350F99F989F931835824A0"
	//pass    = "abc123"
	//TestSym = "TBTC-USD-200626"

	Passphrase = "scWk20511"
	apikey     = "67a3ea44-32d7-49bd-80f6-65e240b5580a"
	secretkey  = "84160D4F65890BDA14EB63D4530A229A" //6
	Symbol     = "ETH-USDT-SWAP"

	//apikey = "f71f4c5c-3988-47e2-b00e-459239c8149f"
	//secretkey = "54D0400EDD6A85EA78D0A55F835C5654"

)

func TestInitClient(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)

	orderBook := client.GetOrderBook(Symbol)
	var ask = orderBook.AskList[1]
	var bid = orderBook.BidList[1]
	log.Println(ask, bid)
}

func TestOKClient_FindTradesByOrderID(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	ret := client.FindTradesByOrderID(Symbol, "OK6670108577312899072")
	for _, trade := range ret {
		log.Printf("%+v", trade)
	}
}

func TestOKClient_GetAccount(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	account := client.GetAccount()
	if account == nil {
		log.Println("account call bank")
		return
	}
	for currency, subAccount := range account.FutureSubAccounts {
		log.Printf("currency: %v account: %v", currency, global.StructToJsonStr(subAccount))
	}
}

func TestOKClient_GetPosition(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	position := client.GetPosition(Symbol)
	for _, futurePosition := range position {
		log.Printf("%+v", futurePosition)
	}
}

func TestOKClient_PostMatchOrder(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	orderBook := client.GetOrderBook(Symbol)
	var ask = orderBook.AskList[1] // 卖价
	var bid = orderBook.BidList[1] // 买价
	log.Println(ask.Price, bid.Price)
	req := model.OrderReq{
		OrderID:  global.GetUUID(),
		OrdType:  global.BuyType,
		Symbol:   Symbol,
		OrderQty: 1,
		TryCount: 0,
	}
	ret, err := client.PostOrder(req)
	log.Println(ret, err)
}

func TestOKClient_PostOrder(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	orderBook := client.GetOrderBook(Symbol)
	var ask = orderBook.AskList[1] // 卖价
	var bid = orderBook.BidList[1] // 买价
	log.Println(ask.Price, bid.Price)
	req := model.OrderReq{
		OrderID:   global.GetUUID(),
		Direction: global.BuyType,
		Symbol:    Symbol,
		OrderQty:  1,
		TryCount:  0,
	}
	ret, err := client.PostOrder(req)
	log.Println(ret, err)
}

func TestOKClient_CancelOrderByID(t *testing.T) {
	//
	client := InitClient(apikey, secretkey, Passphrase)

	ret, err := client.CancelOrderByID(Symbol, "OK6670108577312899072")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(ret)
}

// long 为平多， short为平空单
func TestOKClient_CallAll(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	err := client.CallAll(Symbol, "short")
	if err != nil {
		log.Println(err)
	}
}

func TestOKClient_GetAllOrders(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	os, err := client.GetAllUnFinishOrders(Symbol)
	if err != nil {
		log.Println(err)
	}
	log.Println(os)
	for _, o := range os {
		fmt.Printf("%+v \n", o)
	}
}

func TestOKClient_GetLastFinishOrders(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	os, err := client.GetLastFinishOrders(Symbol)
	if err != nil {
		log.Println(err)
	}
	for _, o := range os {
		fmt.Printf("%+v \n", o)
	}
}

func TestOKClient_GetAccountTradeHistory(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	tradeHisotry, err := client.GetAccountTradeHisotry(Symbol, "520037974051155968")
	if err != nil {
		log.Println(err)
	}
	for _, o := range tradeHisotry {
		fmt.Printf("%+v \n", global.StructToJsonStr(o))
	}
}

func TestOKClient_GetAccountWithSymbol(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	ac := client.GetAccountWithSymbol(Symbol)
	log.Println(global.StructToJsonStr(ac))
}

func TestTick(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)

	tickers, e := client.APIClient.OKExSpot.GetAllTicker()
	if e != nil {
		t.Logf("%v", e)
		return
	}
	for _, v := range tickers {
		if strings.Contains(v.Symbol, "USDT") {
			t.Log(v)
		}
	}
}

func TestOKClient_SubscribeOrdersEvent(t *testing.T) {
	client := InitClient(apikey, secretkey, Passphrase)
	err := client.BuildWS()
	if err != nil {
		log.Println(err)
		return
	}
	cancel, cancelFunc := context.WithCancel(context.Background())
	client.SetOrderCallBack(func(order *goex.FutureOrder, s string) {
		log.Printf("SetOrderCallBack: %+v", order)
	})
	client.WS.SetCallbacks(func(ticker *goex.FutureTicker) {
		log.Printf("%+v", ticker.Ticker)
	}, func(depth *goex.Depth) {
		log.Println(depth)
	}, func(trade *goex.Trade, contract string) {
		log.Println(contract, trade)
	}, nil)
	client.SetErrorCallBack(func(err error) {
		log.Println("SetErrorCallBack：", err)
		cancelFunc()
	})
	time.Sleep(time.Second * 2)
	client.SubscribeOrdersEvent(Symbol)
	go func() {
		for {
			<-time.After(time.Second * 2)
			//log.Println(runtime.NumGoroutine())
		}
	}()
	<-cancel.Done()
}

func TestOKClient_GetAccountSpot(t *testing.T) {
	Clt := InitClient(apikey, secretkey, Passphrase)
	account, err := Clt.GetAccountSpot()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(account.Asset, account.NetAsset)
	for _, v := range account.SubAccounts {
		log.Printf("%+v", v)
	}
}
