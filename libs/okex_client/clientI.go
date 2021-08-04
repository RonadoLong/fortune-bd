package api

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/http"
	"strings"
	"time"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/pkg/goex/okex"
	"wq-fotune-backend/pkg/utils"
	"wq-fotune-backend/app/forward-offer-srv/global"

	"wq-fotune-backend/app/forward-offer-srv/srv/model"
)

var (
	instrumentMap = map[string]goex.CurrencyPair{
		"BTC-USD-SWAP":  goex.BTC_USD,
		"ETH-USD-SWAP":  goex.ETH_USD,
		"EOS-USD-SWAP":  goex.EOS_USD,
		"BTC-USDT-SWAP": goex.BTC_USDT,
		"ETH-USDT-SWAP": goex.ETH_USDT,
		"EOS-USDT-SWAP": goex.EOS_USDT,
	}
	client = &http.Client{
		Timeout:   time.Second * 5,
		Transport: &http.Transport{
			//Proxy: func(req *http.Request) (*url.URL, error) {
			//	return &url.URL{
			//		Scheme: "socks5",
			//		Host:   "192.168.101.220:1080"}, nil
			//},
		},
	}
)

type ClientI interface {
	Authenticate(key, secret string) error
	GetPublicGetInstruments(instrumentName string, king string) interface{}
	PutOrderToExchange(req model.OrderReq) (interface{}, error)
}

type OKClient struct {
	APIClient *okex.OKEx
	WS        *okex.OKExV3FuturesWs
}

type ErrResp struct {
	Code int     `json:"code"`
	Desc ErrDesc `json:"desc"`
}

type ErrDesc struct {
	ErrorMessage string `json:"error_message"`
	Result       string `json:"result"`
	ErrorCode    string `json:"error_code"`
	OrderID      string `json:"order_id"`
}

func InitClient(apiKey, apiSecretKey, apiPassphrase string) *OKClient {
	okexClt := okex.NewOKEx(&goex.APIConfig{
		HttpClient:    client,
		ApiKey:        apiKey,
		ApiSecretKey:  apiSecretKey,
		ApiPassphrase: apiPassphrase,
		Endpoint:      "https://www.okex.com",
	})
	return &OKClient{
		APIClient: okexClt,
	}
}

func (c *OKClient) GetApiKey() string {
	return c.APIClient.Config.ApiKey
}
func (c *OKClient) GetApiSecretKey() string {
	return c.APIClient.Config.ApiSecretKey
}
func (c *OKClient) GetApiPassphrase() string {
	return c.APIClient.Config.ApiPassphrase
}

func (c *OKClient) BuildWS() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("BuildWS panic %+v", e)
		}
	}()

	base := okex.NewOKEx(c.APIClient.Config)
	ws := okex.NewOKExV3FuturesWs(base)
	err = ws.Login()
	if err != nil {
		return err
	}
	c.WS = ws
	return err
}

func (c *OKClient) SetOrderCallBack(setCallBack func(*goex.FutureOrder, string)) {
	c.WS.OrderCallback(setCallBack)
}

func (c *OKClient) SetErrorCallBack(setCallBack func(err error)) {
	c.WS.ErrorCallbacks(setCallBack)
}

func (c *OKClient) SubscribeOrdersEvent(symbol string) {
	for key, _ := range instrumentMap {
		_ = c.WS.SubscribeOrder(key, goex.SWAP_CONTRACT)
	}
}

func (c *OKClient) CancelOrderByID(instrumentId, oID string) (bool, error) {
	pair, ok := instrumentMap[instrumentId]
	if !ok {
		return false, errors.New(global.StringJoinString("暂不支持该合约交易: ", instrumentId))
	}
	return c.APIClient.OKExSwap.FutureCancelOrder(pair, instrumentId, oID)
}

func (c *OKClient) FindOrderByOrderID(instrumentId, oID string) *goex.FutureOrder {
	pair, ok := instrumentMap[instrumentId]
	if !ok {
		return nil
	}
	order, err := c.APIClient.OKExSwap.GetFutureOrder(oID, pair, instrumentId)
	if err != nil {
		logger.Errorf("FindOrderByOrderID error: %s", err.Error())
		return nil
	}
	return order
}

func (c *OKClient) FindTradesByOrderID(symbol string, id string) []goex.Trade {
	var cs = goex.CurrencyPair{}
	os := fmt.Sprintf("%s:%s", id, symbol)
	trades, err := c.APIClient.OKExSwap.GetTrades(os, cs, 0)
	if err != nil {
		logger.Errorf("FindTradesByOrderID err: %s", err.Error())
	}
	return trades
}

func (c *OKClient) GetAccount() *goex.FutureAccount {
	account, err := c.APIClient.OKExSwap.GetFutureUserinfo()
	if err != nil {
		logger.Warnf("GetAccount: %s", err.Error())
		return nil
	}
	return account
}

func (c *OKClient) GetAccountWithSymbol(symbol string) *okex.AccountInfo {
	var cu goex.CurrencyPair
	if strings.Contains(symbol, "ETH-USD") {
		cu = goex.ETC_USDT
	}
	if strings.Contains(symbol, "BTC-USD") {
		cu = goex.BTC_USDT
	}
	log.Println(global.StructToJsonStr(cu))
	var err error
	var account *okex.AccountInfo
	err = utils.ReTryFunc(3, func() (bool, error) {
		account, err = c.APIClient.OKExSwap.GetFutureAccountInfo(cu)
		if err != nil {
			logger.Warnf("GetAccountWithSymbol: %s", err.Error())
		}
		return false, err
	})
	return account
}

// 获取持仓
func (c *OKClient) GetPosition(symbol string) []goex.FuturePosition {
	var cs = goex.CurrencyPair{}
	pos, err := c.APIClient.OKExSwap.GetFuturePosition(cs, symbol)
	if err != nil {
		logger.Errorf("GetPosition error: %s", err.Error())
		return nil
	}
	return pos
}

func (c *OKClient) GetOrderBook(symbol string) *goex.Depth {
	var cs = goex.CurrencyPair{}
	depth, err := c.APIClient.OKExSwap.GetFutureDepth(cs, symbol, 10)
	if err != nil {
		logger.Errorf("GetOrderBook error: %s", err.Error())
		return nil
	}
	return depth
}

func (c *OKClient) PostOrder(req model.OrderReq) (string, *ErrResp) {
	var start = time.Now()
	var errResp ErrResp
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("PostOrder: %v", err)
		}
		sub := time.Now().Sub(start)
		logger.Infof("提交订单成功回调: %s 所花时间：%v", req.OrderID, sub)
	}()
	_, ok := instrumentMap[req.Symbol]
	if !ok {
		return "", &ErrResp{
			Code: 500,
			Desc: ErrDesc{
				ErrorMessage: global.StringJoinString("暂不支持该合约交易: ", req.Symbol),
			},
		}
	}
	var params = make(map[string]interface{})
	params["clientOid"] = req.OrderID
	if req.Price == 0.0 {
		params["matchPrice"] = "1" // 是否以对手价下单。0:不是; 1:是
	} else {
		params["price"] = global.Float64ToString(req.Price)
		params["matchPrice"] = "0" // 是否以对手价下单。0:不是; 1:是
	}
	params["size"] = global.Float64ToString(req.OrderQty)
	params["instrumentId"] = req.Symbol
	if strings.ToLower(req.Direction) == global.SellType {
		params["type"] = "2" // 开空
	} else if strings.ToLower(req.Direction) == global.BuyType {
		params["type"] = "1" // 开多
	} else {
		return "", &ErrResp{
			Code: 500,
			Desc: ErrDesc{
				ErrorMessage: global.StringJoinString("请填写准确的订单方向: ", req.Symbol),
			},
		}
	}
	var cs = goex.CurrencyPair{}
	order, err := c.APIClient.OKExSwap.PostFutureOrder(cs, req.Symbol, params)
	if err != nil {
		err := jsoniter.UnmarshalFromString(err.Error(), &errResp)
		if err != nil {
			errResp = ErrResp{
				Code: 1500,
				Desc: ErrDesc{
					ErrorMessage: err.Error(),
				},
			}
		}
		return "", &errResp
	}
	return order, nil
}

func (c *OKClient) PostMatchOrder(req model.OrderReq) (string, error) {
	_, ok := instrumentMap[req.Symbol]
	if !ok {
		return "", errors.New(global.StringJoinString("暂不支持该合约交易: ", req.Symbol))
	}
	var params = make(map[string]interface{})
	params["clientOid"] = "okex" + req.OrderID
	params["matchPrice"] = "1" // 是否以对手价下单。0:不是; 1:是
	params["size"] = global.Float64ToString(req.OrderQty)
	params["instrumentId"] = req.Symbol
	if strings.ToLower(req.Direction) == global.SellType {
		params["type"] = "2" // 开空
	} else if strings.ToLower(req.Direction) == global.BuyType {
		params["type"] = "1" // 开空
	}
	var cs = goex.CurrencyPair{}
	return c.APIClient.OKExSwap.PostFutureOrder(cs, req.Symbol, params)
}

func (c *OKClient) CallAll(symbol, de string) error {
	cancelOrder, err := c.APIClient.OKExSwap.ClosePosition(symbol, de)
	if cancelOrder {
		logger.Info("cancelOrder success")
	}
	return err
}

func (c *OKClient) GetAllUnFinishOrders(symbol string) ([]goex.FutureOrder, error) {
	var cp = goex.CurrencyPair{}
	var err error
	var orders []goex.FutureOrder
	err = utils.ReTryFunc(10, func() (bool, error) {
		orders, err = c.APIClient.OKExSwap.GetUnfinishFutureOrders(cp, symbol)
		if err != nil {
			time.Sleep(time.Second)
			return false, err
		}
		return false, nil
	})
	return orders, err
}

func (c *OKClient) GetLastFinishOrders(symbol string) ([]goex.FutureOrder, error) {
	var err error
	var orders []goex.FutureOrder
	err = utils.ReTryFunc(10, func() (bool, error) {
		orders, err = c.APIClient.OKExSwap.GetFinishFutureOrders(symbol, 1)
		if err != nil {
			time.Sleep(time.Second)
			return false, err
		}
		return false, nil
	})
	return orders, err
}

//transfer	String	转入/转出
//match	String	交易产生的资金变动
//settlement	String	清算/分摊
//liquidation	String	强平/减仓
//funding	String	资金费
//margin	String	修改保证金的资金变动
// GetAccountTradeHisotry 查询账户流水
func (c *OKClient) GetAccountTradeHisotry(symbol, beforeId string) ([]goex.AccountTrade, error) {
	var err error
	var trades []goex.AccountTrade
	err = utils.ReTryFunc(10, func() (bool, error) {
		trades, err = c.APIClient.OKExSwap.GetAccountTrades(symbol, beforeId, 100)
		if err != nil {
			time.Sleep(time.Second * 2)
			return false, err
		}
		return false, nil
	})
	return trades, err
}

func (c *OKClient) GetAccountSpot() (*goex.Account, error) {
	return c.APIClient.OKExSpot.GetAccount()
}

func (c *OKClient) GetAccountSwap() (*goex.Account, error) {
	resp := &goex.Account{}
	userinfo, err := c.APIClient.OKExSwap.GetFutureUserinfo()
	if err != nil {
		return nil, err
	}
	resp.Exchange = "okex"
	for currency, account := range userinfo.FutureSubAccounts {
		resp.SubAccounts[currency] = goex.SubAccount{
			Currency:     account.Currency,
			Amount:       account.TotalAvailBalance,
			ForzenAmount: account.MarginFrozen,
			LoanAmount:   0,
			Balance:      account.AccountRights,
		}
	}
	return resp, nil
}

func (c *OKClient) CheckIfApiValid() error {
	_, err := c.APIClient.GetAccount()
	return err
}
