package binance

import (
	"fortune-bd/libs/goex"
	"log"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var ba = NewWithConfig(
	&goex.APIConfig{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse("socks5://127.0.0.1:7891")
				},
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
			},
			Timeout: 10 * time.Second,
		},
		Endpoint:     GLOBAL_API_BASE_URL,
		ApiKey:       "p3nbQkUhKsD2vD6Nt6tsv5OQ8OK8IJVGrjDD6ZDx28Iganzha3gVIN6UOPTIWXR2",
		ApiSecretKey: "HvgxFc2dtMYYKmgWjm90E7mEWrzbJfNyZ4yhPwbW0n0VBom4l9iJHlB96HPLcWq3",
	})

func TestBinance_GetTicker(t *testing.T) {
	ticker, err := ba.GetTicker(goex.NewCurrencyPair2("USDT_USD"))
	t.Log(ticker, err)
}

func TestBinance_LimitBuy(t *testing.T) {
	order, err := ba.LimitBuy("0.005", "10300", goex.BTC_USDT)
	t.Log(order, err)
}

func TestBinance_MarketBuy(t *testing.T) {
	order, err := ba.MarketBuy("0.005", "", goex.BTC_USDT)
	t.Log(order, err)
}

func TestBinance_LimitSell(t *testing.T) {
	order, err := ba.LimitSell("0.00499", "10200", goex.BTC_USDT)
	t.Log(order, err)
}

func TestBinance_MarketSell(t *testing.T) {
	order, err := ba.MarketSell("0.00499", "", goex.BTC_USDT)
	t.Log(order, err)
}

func TestBinance_CancelOrder(t *testing.T) {
	t.Log(ba.CancelOrder("1156274704", goex.BTC_USDT))
}

func TestBinance_CancelAllOrders(t *testing.T) {
	ok, err := ba.CancelAllOrders(goex.BTC_USDT)
	if err != nil {
		t.Log(ok, err)
	}
}

func TestBinance_GetOneOrder(t *testing.T) {
	t.Log(ba.GetOneOrder("1156274704", goex.BTC_USDT))
}

func TestBinance_GetDepth(t *testing.T) {
	//return
	dep, err := ba.GetDepth(5, goex.ETH_BTC)
	t.Log(err)
	if err == nil {
		t.Log(dep.AskList)
		t.Log(dep.BidList)
	}
}

func TestBinance_GetAccount(t *testing.T) {
	account, err := ba.GetAccount()
	if err != nil {
		t.Log(err)
		return
	}
	for currency, subAccount := range account.SubAccounts {
		if subAccount.Balance > 0 || subAccount.Amount > 0 {
			t.Logf("currency: %+v, val: %+v", currency, subAccount)
		}
	}
}

func TestBinance_GetUserDataStream(t *testing.T) {
	lKey := "boz16htWU3gjF7hkoRUG3EZponFxRpbQ3XDDUywmKO80krZAzoY13tzAuKIQ"
	go func() {
		for {
			stream, err := ba.PutUserDataStream(lKey)
			if err != nil {
				t.Log(stream, err)
			}
			time.Sleep(time.Minute * 1)
		}
	}()
	_ = bnWs.SubscribeExecutionReport(lKey, func(order *goex.Order) {
		log.Println(order)
	})
	select {}
}

func TestBinance_GetUnfinishOrders(t *testing.T) {
	orders, err := ba.GetUnfinishOrders(goex.ETH_BTC)
	t.Log(orders, err)
}

func TestBinance_GetKlineRecords(t *testing.T) {
	before := time.Now().Add(-time.Hour).Unix() * 1000
	kline, _ := ba.GetKlineRecords(goex.ETH_BTC, goex.KLINE_PERIOD_5MIN, 100, int(before))
	for _, k := range kline {
		tt := time.Unix(k.Timestamp, 0)
		t.Log(tt, k.Open, k.Close, k.High, k.Low, k.Vol)
	}
}

func TestBinance_GetTrades(t *testing.T) {
	t.Log(ba.GetTrades(goex.BTC_USDT, 0))
}

func TestBinance_GetTradeSymbols(t *testing.T) {
	t.Log(ba.GetTradeSymbol(goex.BTC_USDT))
}

func TestBinance_SetTimeOffset(t *testing.T) {
	//t.Log(ba.setTimeOffset())
	t.Log(ba.timeOffset)
	t.Log(ba.GetExchangeInfo())
}

func TestBinance_GetOrderHistorys(t *testing.T) {
	orders, err := ba.GetOrderHistorys(goex.BTC_USDT, 1, 20)
	if err != nil {
		return
	}
	t.Log(len(orders), err)

}