package binance

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sort"
	"testing"
	"time"
	"wq-fotune-backend/app/grid-strategy-srv/util/goex"

	"strings"

	"github.com/k0kubun/pp"
)

var ba = NewWithConfig(
	&goex.APIConfig{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse("socks5://127.0.0.1:10808")
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
	stream, err := ba.GetUserDataStream()
	if err != nil {
		t.Log(stream, err)
	}
}

func TestBinance_WS(t *testing.T) {
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
		log.Printf("%+v \n", order)
	})
	select {}
}

func TestBinance_GetUnfinishOrders(t *testing.T) {
	orders, err := ba.GetUnfinishOrders(goex.ETH_BTC)
	t.Log(orders, err)
}

func TestBinance_GetKlineRecords(t *testing.T) {
	//before := time.Now().Add(-time.Hour).Unix() * 1000
	//kline, _ := ba.GetKlineRecords(goex.ETH_BTC, goex.KLINE_PERIOD_5MIN, 100, int(before))
	before := time.Now().Add(-time.Hour*24*30*3).Unix() * 1000
	kline, _ := ba.GetKlineRecords(goex.ETH_USDT, goex.KLINE_PERIOD_1MONTH, 1, int(before))
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

// -------------------------------------------------------------

var (
	account   *BinanceAccount
	accessKey = "p3nbQkUhKsD2vD6Nt6tsv5OQ8OK8IJVGrjDD6ZDx28Iganzha3gVIN6UOPTIWXR2"
	secretKey = "HvgxFc2dtMYYKmgWjm90E7mEWrzbJfNyZ4yhPwbW0n0VBom4l9iJHlB96HPLcWq3"
	proxyAddr = "socks5://127.0.0.1:10808"
)

func init() {
	var err error
	account, err = InitAccount(accessKey, secretKey, proxyAddr)
	if err != nil {
		panic(err)
	}
}

func TestBinanceAccount_GetAccountBalance(t *testing.T) {
	resp, err := account.GetAccountBalance()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(resp)
}

func TestBinanceAccount_GetCurrencyBalance(t *testing.T) {
	resp, err := account.GetCurrencyBalance("usdt")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(resp)
}

func TestBinanceAccount_PlaceLimitOrder(t *testing.T) {
	side := "buy"
	symbol := "btcusdt"
	price := "10101.0"
	amount := "0.001"
	clientOrderID := BrokerID + "glob_1_8275c0c86db3_1565"

	orderID, err := account.PlaceLimitOrder(side, symbol, price, amount, clientOrderID)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(orderID) // 2773009701
}

func TestBinanceAccount_PlaceMarketOrder(t *testing.T) {
	side := "sell"
	symbol := "winusdt"
	amount := "89743"
	clientOrderID := fmt.Sprintf("mo_%d", time.Now().Unix())

	orderID, err := account.PlaceMarketOrder(side, symbol, amount, clientOrderID)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(orderID) // 2774355049
}

func TestBinanceAccount_CancelOrder(t *testing.T) {
	err := account.CancelOrder("2773009701", "btcusdt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestBinanceAccount_GetOrderInfo(t *testing.T) {
	orderInfo, err := account.GetOrderInfo("2773009701", "btcusdt")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(orderInfo)
}

func TestBinanceAccount_GetHistoryOrdersInfo(t *testing.T) {
	symbol := "btcusdt"
	states := "canceled"
	types := ""

	results, err := account.GetHistoryOrdersInfo(symbol, states, types)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(len(results.([]*OrderInfo)), results)
}

func TestGetSymbolLimit(t *testing.T) {
	symbols := []string{"btcusdt", "ethusdt"}

	for _, symbol := range symbols {
		sl, err := GetSymbolLimit(symbol, proxyAddr)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(*sl)
	}
}

func TestGetLatestPrice(t *testing.T) {
	price, err := GetLatestPrice("ethusdt", proxyAddr)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(price)
}

func TestGetKlines(t *testing.T) {
	klines, err := GetKlines("atomusdt", goex.KLINE_PERIOD_1MONTH, 3, proxyAddr)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(klines)
}

func TestGetLowerPrice(t *testing.T) {
	lowerPrice, highPrice, err := Get3CyclePriceRange("ethusdt", goex.KLINE_PERIOD_1MONTH, proxyAddr)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(lowerPrice, highPrice)
}

// 测试获取交易所品种参数
func TestGetSymbols(t *testing.T) {
	anchorCurrency := "usdt"
	symbols := []string{"btcusdt", "ethusdt", "eosusdt"}
	eis, err := GetSelectSymbols(symbols, proxyAddr)
	if err != nil {
		t.Error(err)
		return
	}
	pp.Println(eis)

	symbolStr := "\""
	acLen := len(anchorCurrency)
	for _, v := range eis {
		v.Symbol = strings.ToLower(v.Symbol)
		if len(v.Symbol) > acLen {
			if v.Symbol[len(v.Symbol)-acLen:] == anchorCurrency {
				symbolStr += v.Symbol + "\"" + ", \""
			}
		}
	}

	fmt.Println(symbolStr)
}
func TestGetAnchorCurrencySymbols(t *testing.T) {
	anchorCurrency := "btc"
	eis, err := GetAnchorCurrencySymbols(anchorCurrency, proxyAddr)
	if err != nil {
		t.Error(err)
		return
	}

	symbols := []string{}
	for _, v := range eis {
		v.Symbol = strings.ToLower(v.Symbol)
		symbols = append(symbols, v.Symbol)

	}
	sort.Strings(symbols)

	symbolStr := "\""
	for _, symbol := range symbols {
		symbolStr += symbol + "\"" + ", \""
	}

	fmt.Println(symbolStr)
}

// 添加品种
func TestGetGridParams(t *testing.T) {
	symbols := []string{
		"btcusdt", "ethusdt", "bnbusdt", "bccusdt", "neousdt", "ltcusdt", "qtumusdt", "adausdt", "xrpusdt", "eosusdt", "tusdusdt", "iotausdt", "xlmusdt", "ontusdt", "trxusdt", "etcusdt", "icxusdt", "venusdt", "nulsusdt", "vetusdt", "paxusdt", "bchabcusdt", "bchsvusdt", "usdcusdt", "linkusdt", "wavesusdt", "bttusdt", "usdsusdt", "ongusdt", "hotusdt", "zilusdt", "zrxusdt", "fetusdt", "batusdt", "xmrusdt", "zecusdt", "iostusdt", "celrusdt", "dashusdt", "nanousdt", "omgusdt", "thetausdt", "enjusdt", "mithusdt", "maticusdt", "atomusdt", "tfuelusdt", "oneusdt", "ftmusdt", "algousdt", "usdsbusdt", "gtousdt", "erdusdt", "dogeusdt", "duskusdt", "ankrusdt", "winusdt", "cosusdt", "npxsusdt", "cocosusdt", "mtlusdt", "tomousdt", "perlusdt", "dentusdt", "mftusdt", "keyusdt", "stormusdt", "dockusdt", "wanusdt", "funusdt", "cvcusdt", "chzusdt", "bandusdt", "busdusdt", "beamusdt", "xtzusdt", "renusdt", "rvnusdt", "hcusdt", "hbarusdt", "nknusdt", "stxusdt", "kavausdt", "arpausdt", "iotxusdt", "rlcusdt", "mcousdt", "ctxcusdt", "bchusdt", "troyusdt", "viteusdt", "fttusdt", "eurusdt", "ognusdt", "drepusdt", "bullusdt", "bearusdt", "ethbullusdt", "ethbearusdt", "tctusdt", "wrxusdt", "btsusdt", "lskusdt", "bntusdt", "ltousdt", "eosbullusdt", "eosbearusdt", "xrpbullusdt", "xrpbearusdt", "stratusdt", "aionusdt", "mblusdt", "cotiusdt", "bnbbullusdt", "bnbbearusdt", "stptusdt", "wtcusdt", "datausdt", "xzcusdt", "solusdt", "ctsiusdt", "hiveusdt", "chrusdt", "btcupusdt", "btcdownusdt", "gxsusdt", "ardrusdt", "lendusdt", "mdtusdt", "stmxusdt", "kncusdt", "repusdt", "lrcusdt", "pntusdt", "compusdt", "bkrwusdt", "scusdt", "zenusdt", "snxusdt", "ethupusdt", "ethdownusdt", "adaupusdt", "adadownusdt", "linkupusdt", "linkdownusdt", "vthousdt", "dgbusdt", "gbpusdt", "sxpusdt", "mkrusdt", "daiusdt", "dcrusdt", "storjusdt", "bnbupusdt", "bnbdownusdt", "xtzupusdt", "xtzdownusdt", "manausdt", "audusdt", "yfiusdt", "balusdt", "blzusdt", "irisusdt", "kmdusdt", "jstusdt", "srmusdt",
	}

	USDT := "USDT"
	sort.Strings(symbols)

	currencies := []string{}
	underlineCurrencyPairs := []string{}
	currencyPairs1 := []string{}
	currencyPairs2 := []string{}
	for _, v := range symbols {
		v := strings.ToUpper(v)
		currencyPair := v
		currency := strings.Replace(v, USDT, "", -1)
		underlineCurrencyPair := currency + "_" + USDT

		currencies = append(currencies, fmt.Sprintf("%s = Currency{\"%s\", \"\"}", currency, currency))
		underlineCurrencyPairs = append(underlineCurrencyPairs, fmt.Sprintf("%s = CurrencyPair{%s, %s}", underlineCurrencyPair, currency, USDT))

		currencyPairs1 = append(currencyPairs1, fmt.Sprintf("\"%s\": %s,", currencyPair, underlineCurrencyPair))
		currencyPairs2 = append(currencyPairs2, fmt.Sprintf("\"%s\": %s,", underlineCurrencyPair, underlineCurrencyPair))
	}

	//for _, v := range currencies {
	//	fmt.Println(v)
	//}
	////fmt.Println("-------------------------------------------")
	//for _, v := range underlineCurrencyPairs {
	//	fmt.Println(v)
	//}
	//fmt.Println("-------------------------------------------")
	//for _, v := range currencyPairs1 {
	//	fmt.Println(v)
	//}
	//fmt.Println("-------------------------------------------")
	for _, v := range currencyPairs2 {
		fmt.Println(v)
	}
	fmt.Println("-------------------------------------------")
}
