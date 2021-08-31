package model

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"fortune-bd/app/grid-strategy-svc/util/gocrypto"
	"fortune-bd/app/grid-strategy-svc/util/goex"
	"fortune-bd/app/grid-strategy-svc/util/goex/binance"
	"fortune-bd/app/grid-strategy-svc/util/grid"
	"fortune-bd/app/grid-strategy-svc/util/huobi"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/client/orderwebsocketclient"
	"fortune-bd/libs/env"
	"strconv"
	"strings"
	"sync"
	"time"


	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/gohttp"
	"github.com/zhufuyi/pkg/krand"
	"github.com/zhufuyi/pkg/logger"
)

//  已运行的网格策略信息，用户.交易所作为key，对象GridProcess作为值
var strategyRunInfo = new(sync.Map)

// StrategyCacheKey 运行策略key
func StrategyCacheKey(uid string, exchange string, symbol string, gsid string) string {
	return uid + "." + exchange + "." + symbol + "." + gsid
}

// GetStrategyCache 获取策略运行信息
func GetStrategyCache(key string) (*GridProcess, bool) {
	if value, ok := strategyRunInfo.Load(key); ok {
		return value.(*GridProcess), true
	}
	return &GridProcess{}, false
}

// SetStrategyCache 设置策略运行信息
func SetStrategyCache(key string, value *GridProcess) {
	strategyRunInfo.Store(key, value)
}

// DeleteStrategyCache 删除策略运行信息
func DeleteStrategyCache(key string) {
	strategyRunInfo.Delete(key)
}

// GetStrategyCaches 获取所有运行的策略
func GetStrategyCaches() []*GridProcess {
	gps := []*GridProcess{}
	strategyRunInfo.Range(func(key, value interface{}) bool {
		gps = append(gps, value.(*GridProcess))
		return true
	})
	return gps
}

// InitStrategyCache 把正在运行的网格策略信息加载到缓存
func InitStrategyCache() error {
	if strategyRunInfo == nil {
		strategyRunInfo = new(sync.Map)
	}

	query := bson.M{"isRun": true}
	gss, err := FindGridStrategies(query, bson.M{}, 0, 5000)
	if err != nil {
		return err
	}

	var key string
	count := 0
	for _, v := range gss {
		accessKey, secretKey, account, err := InitExchangeAccount(v.UID, v.Exchange, v.ApiKey)
		if err != nil {
			logger.Error("InitExchangeAccount error", logger.Err(err), logger.String("uid", v.UID), logger.String("exchange", v.Exchange))
			continue
		}

		query = bson.M{"gsid": v.ID}
		gpo, err := FindGridPendingOrder(query, bson.M{})
		if err != nil {
			logger.Error("FindGridPendingOrder error", logger.Err(err), logger.Any("query", query))
			continue
		}

		ctx, cancel := context.WithCancel(context.Background())

		gp := &GridProcess{
			Gsid:            v.ID.Hex(),
			UID:             v.UID,
			Exchange:        v.Exchange,
			Symbol:          v.Symbol,
			Type:            v.Type,
			AccessKey:       accessKey,
			SecretKey:       secretKey,
			AnchorSymbol:    v.AnchorSymbol,
			Grids:           gpo.Grids,
			BasisGridNO:     gpo.BasisGridNO,
			ExchangeAccount: account,

			Ctx:    ctx,
			Cancel: cancel,
		}
		err = ProcessGridOrder(gp)
		if err != nil {
			logger.Error("ProcessGridOrder error", logger.Err(err), logger.Any("gridProcess", gp.Desensitize()))
			continue
		}
		go gp.ListenOrder() // 定时检查委托订单

		key = StrategyCacheKey(v.UID, v.Exchange, v.Symbol, gp.Gsid)
		SetStrategyCache(key, gp)
		count++

		time.Sleep(40 * time.Millisecond) // 限制初始化速度
	}

	logger.Infof("InitStrategyCache finish, success=%d, total=%d", count, len(gss))

	return nil
}

// -------------------------------------------------------------------------------------------------

// GridProcess 网格处理
type GridProcess struct {
	Gsid            string         // 网格策略id
	UID             string         // 用户id
	Exchange        string         // 交易所
	Symbol          string         // 品种
	Type            int            // 策略类型
	AnchorSymbol    string         // 锚定币
	AccessKey       string         // 访问key
	SecretKey       string         // 访问密钥
	Grids           []*grid.Grid   // 网格
	BasisGridNO     int            // 网格线基准编号
	ExchangeAccount goex.Accounter // 交易所账号接口

	HuobiWsCli   *orderwebsocketclient.SubscribeOrderWebSocketV2Client // 火币的web socket
	BinanceWsCli *binance.BinanceWs                                    // 火币的web socket

	Ctx    context.Context
	Cancel context.CancelFunc
}

// TradeRecord 网格交易记录
type TradeRecord struct {
	GID           int     // 网格编号
	OrderID       string  // 委托订单id，用户查询、取消订单
	ClientOrderID string  // 用户自定义订单id
	OrderType     string  // 订单类型，limit:限价单，market:市价单
	Side          string  // 买入卖出，buy:买入，sell:卖出
	Price         float64 // 成交价格
	Quantity      float64 // 买卖数量
	Volume        float64 // 成交额
	Unit          string  // 成交额单位
	OrderState    string  // 当前订单状态 submitted:委托中, canceled:取消, filled:已成交
}

func (g *GridProcess) Desensitize() GridProcess {
	gp := *g
	gp.SecretKey = ""
	gp.ExchangeAccount = nil
	return gp
}

// 网格买单成交后处理
func (g *GridProcess) processBuy(gridNO int) (*TradeRecord, string, error) {
	// 判断将要添加委托的卖单是否超过网格界限
	if gridNO < 1 {
		logger.Warn("order is illegal")
		return nil, "", nil
	}
	if g.ExchangeAccount == nil {
		return nil, "", errors.New("not found account")
	}
	grids := g.Grids
	grid := grids[gridNO-1]

	var sellOrderID string
	// 判断将要添加委托卖单是否已经存在
	if grid.Side == "sell" && grid.OrderID != "" {
		logger.Infof("sell order(%s) already exists", grid.OrderID)
		sellOrderID = grid.OrderID
		g.Grids[gridNO].Side = ""
		g.Grids[gridNO].OrderID = ""
		g.BasisGridNO = gridNO - 1
		return nil, sellOrderID, nil
	}

	// 委托卖单
	side := "sell"
	//price := fmt.Sprintf("%v", grid.Price)
	price := strconv.FormatFloat(grid.Price, 'f', -1, 64)
	amount := fmt.Sprintf("%v", grid.SellQuantity)
	//clientOrderID := fmt.Sprintf("glos_%d_%s_%s", gridNO-1, g.Gsid, krand.String(krand.R_NUM|krand.R_LOWER, 4)) // lob_%d表示grid limit order buy 网格编号
	clientOrderID := NewGridClientOrderID(g.Exchange, side, g.Gsid, gridNO-1)
	orderID, err := g.ExchangeAccount.PlaceLimitOrder(side, g.Symbol, price, amount, clientOrderID)
	if err != nil {
		logger.Error("placeLimitOrder error", logger.Err(err), logger.Any("grid", grid))
		return nil, sellOrderID, err
	}

	sellOrderID = orderID
	grid.Side = side
	grid.OrderID = orderID
	g.Grids[gridNO-1] = grid
	g.BasisGridNO = gridNO - 1

	// 清空当前值
	g.Grids[gridNO].Side = ""
	g.Grids[gridNO].OrderID = ""

	recordOrder := &TradeRecord{
		GID:           grid.GID,
		OrderID:       orderID,
		ClientOrderID: clientOrderID,
		OrderType:     "limit",
		Side:          side,
		Price:         grid.Price,
		Quantity:      grid.SellQuantity,
		Volume:        grid.Price * grid.SellQuantity,
		Unit:          g.AnchorSymbol,
		OrderState:    huobi.OrderStateSubmitted,
	}

	return recordOrder, sellOrderID, nil
}

// 网格卖单成交后处理
func (g *GridProcess) processSell(gridNO int) (*TradeRecord, error) {
	// 判断将要添加的委托买单是否超过网格界限
	if gridNO >= len(g.Grids)-1 {
		logger.Warn("order is illegal")
		return nil, nil
	}
	if g.ExchangeAccount == nil {
		return nil, errors.New("not found account")
	}

	grids := g.Grids
	grid := grids[gridNO+1]

	// 判断将要添加委托卖单是否已经存在
	if grid.Side == "buy" && grid.OrderID != "" {
		logger.Infof("sell order(%s) already exists", grid.OrderID)
		g.Grids[gridNO].Side = ""
		g.Grids[gridNO].OrderID = ""
		g.BasisGridNO = gridNO + 1
		return nil, nil
	}

	// 委托买单
	side := "buy"
	//price := fmt.Sprintf("%v", grid.Price)
	price := strconv.FormatFloat(grid.Price, 'f', -1, 64)
	amount := fmt.Sprintf("%v", grid.BuyQuantity)
	//clientOrderID := fmt.Sprintf("glob_%d_%s_%s", gridNO+1, g.Gsid, krand.String(krand.R_NUM|krand.R_LOWER, 4)) // lob_%d表示grid limit order buy 网格编号
	clientOrderID := NewGridClientOrderID(g.Exchange, side, g.Gsid, gridNO+1)

	orderID, err := g.ExchangeAccount.PlaceLimitOrder(side, g.Symbol, price, amount, clientOrderID)
	if err != nil {
		logger.Error("placeLimitOrder error", logger.Err(err), logger.Any("grid", grid))
		return nil, err
	}

	grid.Side = side
	grid.OrderID = orderID
	grids[gridNO+1] = grid
	g.BasisGridNO = gridNO + 1 // 基准序号往前一个网格

	g.Grids[gridNO].Side = ""
	g.Grids[gridNO].OrderID = ""

	recordOrder := &TradeRecord{
		GID:           grid.GID,
		OrderID:       orderID,
		ClientOrderID: clientOrderID,
		OrderType:     "limit",
		Side:          side,
		Price:         grid.Price,
		Quantity:      grid.BuyQuantity,
		Volume:        grid.Price * grid.BuyQuantity,
		Unit:          g.AnchorSymbol,
		OrderState:    huobi.OrderStateSubmitted,
	}

	return recordOrder, nil
}

// 更新网格存储信息
func (g *GridProcess) processSaveInfo(gridNO int, orderID string, orderStatus string, tradeTime int64, tradeAmount float64, newRecord *TradeRecord, sellOrderID string) error {
	buyPrice := 0.0

	fees := 0.0
	//区分不同交易所
	switch g.Exchange {
	case ExchangeHuobi:
		fees = huobi.FilledFees
	case ExchangeBinance:
		fees = binance.FilledFees
	}

	// 更新订单状态和时间
	query := bson.M{"gsid": bson.ObjectIdHex(g.Gsid), "orderID": orderID}
	update := bson.M{
		"$set": bson.M{
			"orderState":  orderStatus,
			"stateTime":   time.Unix(tradeTime/1000, tradeTime%1000),
			"sellOrderID": sellOrderID, // 如果成交的是委托卖单，为空
			"fees":        tradeAmount * fees,
		},
	}
	gtr, err := FindAndModifyGridTradeRecord(query, update)
	if err != nil {
		logger.Error("FindAndModifyGridTradeRecord error", logger.Err(err), logger.Any("query", query), logger.Any("update", update))
		return err
	}
	logger.Info("update gridTradeRecord success", logger.String("orderID", orderID), logger.Any("update", update), logger.Float64("fees", fees))

	// 如果成交的是卖单，需要通知统计
	if gtr.Side == "sell" {
		buyOrderID := ""
		// 判断是否启动网格时委托的卖单
		if gtr.IsStartUpOrder {
			// 查找市价单价格
			query = bson.M{"gsid": gtr.GSID, "side": "buy", "orderType": "market"}
			mgtrs, err := FindGridTradeRecords(query, bson.M{}, 0, 1)
			if err != nil {
				logger.Error("FindGridTradeRecord error", logger.Err(err), logger.Any("query", query))
			} else {
				if len(mgtrs) > 0 {
					buyPrice = mgtrs[0].Price
					buyOrderID = mgtrs[0].OrderID
				}
			}
		} else {
			// 下一格的买入价格
			if gridNO < len(g.Grids)-1 {
				buyPrice = g.Grids[gridNO+1].Price
			}
		}

		g.notifyStatistics(gtr, buyPrice)

		query = bson.M{"gsid": bson.ObjectIdHex(g.Gsid), "orderID": orderID}
		update = bson.M{
			"$set": bson.M{
				"buyPrice":   buyPrice,
				"buyOrderID": buyOrderID,
			},
		}
		UpdateGridTradeRecord(query, update)
	}

	// 更新网格参数
	query = bson.M{"gsid": gtr.GSID}
	update = bson.M{
		"$set": bson.M{
			"basisGridNO":                                  g.BasisGridNO,
			fmt.Sprintf("grids.%d.side", gridNO):           g.Grids[gridNO].Side,
			fmt.Sprintf("grids.%d.orderId", gridNO):        g.Grids[gridNO].OrderID,
			fmt.Sprintf("grids.%d.side", g.BasisGridNO):    g.Grids[g.BasisGridNO].Side,
			fmt.Sprintf("grids.%d.orderId", g.BasisGridNO): g.Grids[g.BasisGridNO].OrderID,
		},
	}
	err = UpdateGridPendingOrder(query, update)
	if err != nil {
		logger.Error("UpdateGridPendingOrder error", logger.Err(err), logger.Any("query", query), logger.Any("update", update))
		return err
	}
	logger.Info("update gridPendingOrder success", logger.String("gsid", gtr.GSID.Hex()))

	// 添加新订单记录
	if newRecord != nil {
		gridTradeRecord := &GridTradeRecord{
			GSID: gtr.GSID,
			GID:  newRecord.GID,

			OrderID:       newRecord.OrderID,
			ClientOrderID: newRecord.ClientOrderID,
			OrderType:     newRecord.OrderType,
			Side:          newRecord.Side,
			Price:         newRecord.Price,
			Quantity:      newRecord.Quantity,
			Volume:        newRecord.Volume,
			Unit:          newRecord.Unit,

			OrderState: newRecord.OrderState,
			StateTime:  time.Now(),

			Exchange: gtr.Exchange,
			Symbol:   gtr.Symbol,
		}
		err = gridTradeRecord.Insert()
		if err != nil {
			logger.Error("gridTradeRecord.Insert error", logger.Err(err), logger.Any("query", query), logger.Any("update", update))
			return err
		}

		logger.Info("add new GridTradeRecord success", logger.String("orderID", gridTradeRecord.OrderID), logger.String("clientOrderID", newRecord.ClientOrderID))
	}

	return nil
}

// UpdateGridOrder 更新网格订单
func (g *GridProcess) UpdateGridOrder(tradeTime int64, tradeAmount float64, orderID, clientOrderID, orderStatus string) error {
	if !strings.Contains(clientOrderID, cutGsid(g.Gsid)) {
		return nil
	}

	logger.Info("start to UpdateGridOrder",
		logger.Int64("tradeTime", tradeTime),
		logger.Float64("tradeAmount", tradeAmount),
		logger.String("orderID", orderID),
		logger.String("clientOrderID", clientOrderID),
		logger.String("orderStatus", orderStatus),
	)

	// 判断是否完全成交
	if orderStatus != huobi.OrderStateFilled {
		return fmt.Errorf("not full-filled order(orderID=%s, orderStatus=%s), ignore processing", orderID, orderStatus)
	}

	var newRecord *TradeRecord
	var err error
	var sellOrderID string
	prefix, gridNO := parseClientOrderID(g.Exchange, clientOrderID)
	if gridNO >= len(g.Grids) {
		return fmt.Errorf("clientOrderID=%s not match %s gsid=%s, gridNum=%d", clientOrderID, g.Symbol, g.Gsid, len(g.Grids)-1)
	}

	switch prefix {
	case "glob": // 处理成功的买单
		if orderID != g.Grids[gridNO].OrderID {
			updateOrderState(orderID, huobi.OrderStateFilled)
			logger.Errorf("[glos] grid order not math, filled orderID=%s, grid[%d] order=%s", orderID, gridNO, g.Grids[gridNO].OrderID)
			return errors.New("[glos] not match order id")
		}

		newRecord, sellOrderID, err = g.processBuy(gridNO)
		if err != nil {
			logger.Error("processSell error", logger.Err(err), logger.Any("GridProcess", g))
			return err
		}

	case "glos": // 处理成功的卖单
		if orderID != g.Grids[gridNO].OrderID {
			logger.Errorf("[glos] grid order not math, filled orderID=%s, grid[%d] order=%s", orderID, gridNO, g.Grids[gridNO].OrderID)
			return errors.New("[glos] not match order id")
		}

		newRecord, err = g.processSell(gridNO)
		if err != nil {
			logger.Error("processSell error", logger.Err(err), logger.Any("GridProcess", g))
			return err
		}

	default:
		logger.Warnf("this is not grid order(orderID=%s, clientOrderID=%s)", orderID, clientOrderID)
		return nil
	}

	return g.processSaveInfo(gridNO, orderID, orderStatus, tradeTime, tradeAmount, newRecord, sellOrderID)
}

type notifyStatisticsResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	//Date interface{} `json:"date"`
}

func (g *GridProcess) notifyStatistics(gtr *GridTradeRecord, buyPrice float64) error {
	el := GetExchangeLimitCache(GetKey(g.Exchange, g.Symbol))
	//priceUnit := gtr.Unit
	//if g.Type == GridTypeReverse {
	//	priceUnit = goex.GetTradeCurrency(gtr.Symbol)
	//}

	params := map[string]interface{}{
		"user_id":     g.UID,
		"api_key":     g.AccessKey,
		"strategy_id": gtr.GSID.Hex(),
		"order_id":    gtr.OrderID,
		"exchange":    gtr.Exchange,
		"symbol":      gtr.Symbol,
		"unit":        gtr.Unit,
		//"to_unit":     priceUnit,
		//"buy_price":  fmt.Sprintf("%v", FloatRound(buyPrice, el.PricePrecision)),
		"buy_price": strconv.FormatFloat(FloatRound(buyPrice, el.PricePrecision), 'f', -1, 64),
		//"sell_price": fmt.Sprintf("%v", FloatRound(gtr.Price, el.PricePrecision)), // 最新价格
		"sell_price": strconv.FormatFloat(FloatRound(gtr.Price, el.PricePrecision), 'f', -1, 64), // 最新价格
		"volume":     fmt.Sprintf("%v", FloatRound(gtr.Quantity, el.QuantityPrecision)),
		"trade_at":   strconv.FormatInt(gtr.StateTime.UTC().Unix(), 10),
	}

	resp := &notifyStatisticsResp{}
	url := env.NotifyStatisticsURL
	err := gohttp.PostJSON(resp, url, params)
	if err != nil {
		logger.Warn("notifyStatistics error", logger.Err(err))
		return err
	}
	if resp.Code != 0 {
		err = fmt.Errorf("code=%d, msg=%s", resp.Code, resp.Msg)
		logger.Error("notifyStatistics error", logger.Err(err), logger.String("url", url), logger.Any("params", params))
		return err
	}

	logger.Info("notifyStatistics success", logger.String("orderID", gtr.OrderID))

	return nil
}

// ProcessGridOrder ws订阅订单通知和处理网格订单
func ProcessGridOrder(gp *GridProcess) error {
	wsClientID := getWsClientUid(gp.Gsid, gp.UID, gp.Symbol)
	switch gp.Exchange {
	case ExchangeHuobi:
		processFun := huobi.ProcessTradeOrder(gp)
		wsCli, err := huobi.WsSubscribeOrder(wsClientID, gp.Symbol, processFun, gp.AccessKey, gp.SecretKey)
		if err != nil {
			return err
		}
		gp.HuobiWsCli = wsCli

	case ExchangeBinance:
		processFun := binance.ProcessTradeOrder(gp)
		binanceAccount, _ := binance.InitAccount(gp.AccessKey, gp.SecretKey, env.ProxyAddr)
		wsCli, err := binance.WsSubscribeOrder(gp.Ctx, binanceAccount, gp.Symbol, processFun)
		if err != nil {
			time.Sleep(time.Second)
			wsCli, err = binance.WsSubscribeOrder(gp.Ctx, binanceAccount, gp.Symbol, processFun)
			if err != nil {
				return err
			}
		}
		gp.BinanceWsCli = wsCli

	default:
		return fmt.Errorf("unknown exchange %s, can not proccess grid order", gp.Exchange)
	}

	return nil
}

// Close 关闭websocket和清除缓存
func (g *GridProcess) Close() {
	if g.Cancel != nil {
		g.Cancel()
	}

	// 关闭websocket
	switch g.Exchange {
	case ExchangeHuobi:
		if g.HuobiWsCli != nil {
			//g.HuobiWsCli.Close()
		}
	case ExchangeBinance:
		if g.BinanceWsCli != nil {
			g.BinanceWsCli.Close()
		}
	}

	// 从缓存中删除
	key := StrategyCacheKey(g.UID, g.Exchange, g.Symbol, g.Gsid)
	DeleteStrategyCache(key)
}

// 解析ClientOrderID
func parseClientOrderID(exchange string, clientOrderID string) (string, int) {
	// 区分不同交易所
	switch exchange {
	case ExchangeHuobi: // 保持原样

	case ExchangeBinance: // 去掉broken id
		clientOrderID = strings.Replace(clientOrderID, binance.BrokerID, "", -1)
	}

	ss := strings.Split(clientOrderID, "_")
	if len(ss) < 2 {
		return "", 0
	}

	gridNO, err := strconv.Atoi(ss[1])
	if err != nil {
		return "", 0
	}

	return ss[0], gridNO
}

func getWsClientUid(gsid, uid, symbol string) string {
	if len(uid) > 5 {
		uid = uid[len(uid)-6:]
	}
	return fmt.Sprintf("%s_%s_%s_%s", cutGsid(gsid), uid, symbol, krand.String(krand.R_All, 4))
}

// -------------------------------------------------------------------------------------------------
// 初始化交易所账号
type exchangeKey struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		UserID       string `json:"user_id"`
		ExchangeID   int    `json:"exchange_id"`
		ExchangeName string `json:"exchange_name"`
		APIKey       string `json:"api_key"`
		Secret       string `json:"secret"`
		Passphrase   string `json:"passphrase"`
	} `json:"data"`
}

// GetExchangeKey 通过用户和交易所名获取访问权限密钥
func GetExchangeKey(uid string, apiKey string) (string, string, error) {
	if strings.Contains(env.MongoAddr, "192.168.101.88") {
		if apiKey == ExchangeBinance || strings.Contains(apiKey, "aMzW4Trd") { // 币安测试账号
			//return "p3nbQkUhKsD2vD6Nt6tsv5OQ8OK8IJVGrjDD6ZDx28Iganzha3gVIN6UOPTIWXR2", "HvgxFc2dtMYYKmgWjm90E7mEWrzbJfNyZ4yhPwbW0n0VBom4l9iJHlB96HPLcWq3", nil
			return "aMzW4TrdSCfvjPswXnBIcAVkyyAJwaGByX86B2822hbgCsNRs9d8VESB2SBOzPed", "nf8hNnZSplDijLYtKCAc6trIuhHdukaJgkuDCmSTRhTCtR1WAPtp970F5JGrHJFd", nil
		}
		if apiKey == ExchangeHuobi || strings.Contains(apiKey, "dbuqg6hk") { // 火币测试账号
			return "dbuqg6hkte-23df26bc-aa64857f-7930e", "bf7c125a-1ed0243f-716143f5-bdf6c", nil
		}
	}

	ek := &exchangeKey{}
	url := env.ExchangeAccessURL + fmt.Sprintf("/%s/%s", uid, apiKey)
	err := gohttp.GetJSON(ek, url, nil)
	if err != nil {
		return "", "", err
	}
	if ek.Code != 0 {
		return "", "", errors.New(ek.Msg)
	}
	if ek.Data.Secret == "" {
		return "", "", errors.New("secretKey is empty")
	}

	//解密
	data, _ := hex.DecodeString(ek.Data.Secret)
	secretByte, _ := gocrypto.AesDecrypt(data)

	return ek.Data.APIKey, string(secretByte), nil
}

// InitExchangeAccount 初始化交易所账号
func InitExchangeAccount(uid string, exchange string, apiKey string) (string, string, goex.Accounter, error) {
	accessKey, secretKey, err := GetExchangeKey(uid, apiKey)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to obtain exchange access key, uid=%s, apiKey=%s", uid, apiKey)
	}

	// 区分不同交易所
	switch exchange {
	case ExchangeHuobi:
		account, err := huobi.InitAccount(accessKey, secretKey)
		return accessKey, secretKey, account, err
	case ExchangeBinance:
		account, err := binance.InitAccount(accessKey, secretKey, env.ProxyAddr)
		return accessKey, secretKey, account, err
	}

	return "", "", nil, fmt.Errorf("failed to initialize exchange %s account", exchange)
}

// GenerateClientOrderID 根据交易所生成对应用户订单id
func GenerateClientOrderID(exchange string, prefix string, gsid string) string {
	clientOrderID := fmt.Sprintf("%s_%s_%s", prefix, cutGsid(gsid), krand.String(krand.R_All, 4))
	// 区分不同交易所
	switch exchange {
	case ExchangeHuobi:
		return clientOrderID
	case ExchangeBinance:
		return binance.BrokerID + clientOrderID
	}

	return ""
}

// NewGridClientOrderID 生成网格特定的用户订单id
func NewGridClientOrderID(exchange string, side string, gsid string, gridNO int) string {
	var prefix string
	if side == "buy" {
		prefix = PrefixIDGlob
	} else if side == "sell" {
		prefix = PrefixIDGlos
	}

	clientOrderID := fmt.Sprintf("%s_%d_%s_%s", prefix, gridNO, cutGsid(gsid), krand.String(krand.R_All, 4))

	// 区分不同交易所
	switch exchange {
	case ExchangeHuobi:
		return clientOrderID
	case ExchangeBinance:
		return binance.BrokerID + clientOrderID
	}

	return ""
}

func cutGsid(gsid string) string {
	if len(gsid) > 12 {
		gsid = gsid[len(gsid)-12:]
	}
	return gsid
}
