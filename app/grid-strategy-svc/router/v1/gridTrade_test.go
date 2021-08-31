package v1

import (
	"fmt"
	"fortune-bd/app/grid-strategy-svc/model"
	"fortune-bd/app/grid-strategy-svc/util/grid"
	"testing"
	"time"


	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/k0kubun/pp"
	"github.com/zhufuyi/pkg/krand"
	"github.com/zhufuyi/pkg/logger"
)

func TestOutJSON(t *testing.T) {
	o := &strategyOut{}
	data, _ := jsoniter.MarshalIndent(o, "", "    ")
	fmt.Println(string(data))
}

var createGridJSON = []byte(`{
    "uid": "1593572486",
    "exchange": "binance",
    "symbol": "btcusdt",
    "anchorSymbol": "usdt",
    "gridIntervalType": "ASGrid",
    "gridNum": 5,
    "totalSum": 100,
    "minPrice": 11230,
    "maxPrice": 11300
}`)

func TestCreateNormalGridTradeStrategy(t *testing.T) {
	if err := model.InitStrategyCache(); err != nil {
		t.Error(err)
		return
	}

	form := &gridTradeForm{}
	err := jsoniter.Unmarshal(createGridJSON, form)
	if err != nil {
		t.Error(err)
		return
	}
	form.ApiKey = form.Exchange

	err = form.valid()
	if err != nil {
		t.Error(err)
		return
	}

	err = form.initExchangeAccount()
	if err != nil {
		t.Error(err)
		return
	}
	logger.Info("初始化交易所账号成功", logger.String("uid", form.UID))

	grids, err := form.initGrid()
	if err != nil {
		t.Error(err)
		return
	}
	logger.Info("根据参数初始化网格成功", logger.String("uid", form.UID), logger.Any("grids", grids))

	grid.PrintFormat(grids)
	return

	needBuyCoinQuantity, err := form.checkAccountBalance(grids)
	if err != nil {
		t.Error(err)
		return
	}
	logger.Info("账号余额满足网格需求", logger.String("uid", form.UID), logger.Float64("needBuyCoinQuantity", needBuyCoinQuantity))

	// 买入市价单
	marketRecord, err := form.placeMarketOrder("buy", needBuyCoinQuantity)
	if err != nil {
		t.Error(err)
		return
	}

	firstKeys, _ := splitGridKeys(grids, form.GridBasisNO)

	// 买入、卖出网格限价单
	limitRecords, err := form.placeLimitOrder(grids, firstKeys)
	if err != nil {
		t.Error(err)
		return
	}
	if len(limitRecords) == len(grids) {
		logger.Info("买入、卖出网格限价单成功", logger.String("uid", form.UID), logger.Int("orderSize", len(limitRecords)))
	} else {
		logger.Warn("买入、卖出部分网格限价单成功", logger.String("uid", form.UID), logger.Int("gridSize", len(grids)), logger.Int("orderSize", len(limitRecords)))
	}

	// 保存信息
	records := append(limitRecords, marketRecord)
	form.saveInfo(grids, records)
	logger.Info("网格策略已初始化和持久化，开始订阅和处理网格订单")

	// 订阅订单通知和处理网格订单
	gp := form.toGridProcess(grids, form.GridBasisNO, form.gridStrategyID)
	err = model.ProcessGridOrder(gp)
	if err != nil {
		logger.Error("订阅订单通知和处理网格订单失败", logger.Err(err), logger.Any("gridProcess", gp.Desensitize()))
		return
	}

	// 写入缓存
	key := model.StrategyCacheKey(form.UID, form.Exchange, form.Symbol, form.gridStrategyID)
	model.SetStrategyCache(key, gp)

	select {}
}

func TestUpdateGridStrategy(t *testing.T) {
	data := []byte(`{
		"exchange": "binance",
		"symbol":   "eosusdt",
		"gsid":     "5f3a19ebc7ef826a80046e04",
		"minPrice": 3.5,
		"maxPrice": 4.3
}`)
	form := &updateGridStrategyForm{}

	err := jsoniter.Unmarshal(data, form)
	if err != nil {
		t.Error(err)
		return
	}

	err = form.valid()
	if err != nil {
		logger.Error("参数错误", logger.Err(err), logger.Any("form", form))
		return
	}

	grids, err := form.genNewGrid()
	if err != nil {
		logger.Error("genNewGrid error", logger.Err(err), logger.Any("form", form))
		return
	}

	// 初始化交易所账号
	err = form.initExchangeAccount()
	if err != nil {
		logger.Error("初始化交易所账号失败", logger.Err(err), logger.Any("form", form))
		return
	}
	logger.Info("初始化交易所账号成功", logger.String("uid", form.gs.UID))

	err = model.UpdateRunningGridStrategy(form.gs, form.exchangeAccount, grids, form.MinPrice, form.MaxPrice, form.latestPrice, false)
	if err != nil {
		logger.Error("updateGridStrategy error", logger.Err(err), logger.String("gsid", form.Gsid))
		return
	}

}

func TestCalculateMoney(t *testing.T) {
	//input := 10100.0
	//min, max := getMinMax(input, 10)
	//fmt.Println(input, min, max, (max-min)/float64(10))
	//return

	var autoGridJSON = []byte(`{
    "exchange": "binance",
    "symbol": "ethusdt"
}`)

	form := &gridTradeForm{}
	err := jsoniter.Unmarshal(autoGridJSON, form)
	if err != nil {
		t.Error(err)
		return
	}

	err = form.valid2()
	if err != nil {
		t.Error(err)
		return
	}

	money, err := form.calculateGridNeedMoney()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(money)
}

var cancelGridJSON = []byte(`{
    "uid":"1593572486",
    "exchange":"binance",
	"symbol":"winusdt",
    "gsid":"5f362e23c7ef82201c273154",
	"isClosePosition":true
}`)

func TestCancelNormalGridTradeStrategy(t *testing.T) {
	form := &cancelGridStrategyForm{}
	err := jsoniter.Unmarshal(cancelGridJSON, form)
	if err != nil {
		t.Error(err)
		return
	}

	err = form.valid()
	if err != nil {
		t.Error(err)
		return
	}

	err = form.initExchangeAccount()
	if err != nil {
		t.Error(err)
		return
	}
	logger.Info("初始化交易所账号成功", logger.String("uid", form.UID))

	// 取消在交易所的委托
	successOrderIDs, err := form.cancelCommissionOrder()
	if err != nil {
		t.Error(err)
		return
	}
	if len(successOrderIDs) > 0 {
		logger.Info("取消成功的交易所委托订单", logger.Any("successOrderIDs", successOrderIDs))
	}

	// 更新订单状态
	failedUpdateOrderIDs := updateOrderStatus(successOrderIDs)
	if len(failedUpdateOrderIDs) > 0 {
		logger.Error("本地更新状态失败的订单", logger.String("uid", form.UID), logger.Any("failedUpdateOrderIDs", failedUpdateOrderIDs))
	} else {
		logger.Info("本地更新所有订单状态成功", logger.String("uid", form.UID), logger.Any("successOrderIDs", successOrderIDs))
	}

	// 修改策略运行状态
	err = updateStrategyOrder(form.Gsid)
	if err != nil {
		logger.Error("更新策略运行状态失败", logger.String("uid", form.UID), logger.Err(err), logger.String("gsid", form.Gsid))
	} else {
		logger.Info("更新策略运行状态成功", logger.String("uid", form.UID), logger.String("gsid", form.Gsid))
	}

	// 判断是否平仓
	if form.IsClosePosition {
		err = form.closePosition()
		if err != nil {
			logger.Info("closePosition error", logger.Err(err))
			return
		}
	}
}

func TestGetStatisticalInfo(t *testing.T) {
	uid := "1266203624401276928"
	strategyID := "xxxx"
	ek, err := getStatisticalInfo(uid, strategyID)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(ek)
}

func TestGetMinMoney(t *testing.T) {
	exchange := "huobi"
	symbols := []string{"btcusdt", "ethusdt", "eosusdt"}

	for _, symbol := range symbols {
		form := &autoGenerateForm{
			Exchange: exchange,
			Symbol:   symbol,
		}
		if err := form.valid1(); err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(form.getMinMoney())
	}
}

func TestFloatRound(t *testing.T) {
	pi := 3.14159200
	for i := 0; i < 11; i++ {
		fmt.Println(model.FloatRound(pi, i), model.FloatRoundOff(pi, i))
	}

	point := 0
	for i := 0; i < 20; i++ {
		point = i / 5
		fmt.Println(model.FloatRound(krand.Float64(6, 10, 1000), point))
		time.Sleep(time.Millisecond * 100)
	}
}

func TestFilterSymbols(t *testing.T) {
	e := exchangeLimitForm{
		Exchange:       "huobi",
		Symbols:        []string{"btcusdt", "ethusdt", "usdtbtc"},
		AnchorCurrency: "usdt",
	}
	fmt.Println(e.filterSymbols())
}

func TestGetBigGridParams(t *testing.T) {
	b := &bigGridParamForm{
		Exchange:   "huobi",
		Symbol:     "hcusdt",
		MinPrice:   "1.1",
		ProfitRate: "0.003",
	}

	err := b.valid()
	if err != nil {
		t.Error(err)
		return
	}

	bgp, err := b.calculateParams()
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(bgp)
}

func TestGetBigGridParams2(t *testing.T) {
	//minPrice := 320.0 // 最近三个月最小价格
	//maxPrice := 400.0
	//q := 1.004
	//totalSum := 300.0
	//
	//num := int(math.Log10(maxPrice/minPrice) / math.Log10(q))
	//
	//// 计算网格数，最大价格为最低价格的2倍数
	//grids := grid.GenerateGS(minPrice, q, totalSum, num, 2, 6)
	//grid.PrintFormat(grids)

	form := &calculateGridParamForm{
		Exchange:   "binance",
		Symbol:     "ethusdt",
		MinPrice:   "320",
		MaxPrice:   "360",
		ProfitRate: "0.003",
	}

	err := form.valid()
	if err != nil {
		t.Error(err)
		return
	}

	bgp, grids, err := form.calculateParams()
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(bgp)
	grid.PrintFormat(grids)
}

func TestCalculateGrid(t *testing.T) {
	symbol := "ethusdt"
	exchange := "binance"
	totalSum := 700.0
	latestPrice := 376.0

	cg := &model.CalculateGrid{
		Exchange:         exchange,
		Symbol:           symbol,
		TargetProfitRate: 0.0015,
		ParamsRange: &model.SymbolParams{
			TotalSum:    &model.ValueRange{totalSum, totalSum, 1},
			LatestPrice: &model.ValueRange{latestPrice, latestPrice, 1},
		},
	}

	limitVolume, fees := model.GetExchangeRule(exchange, symbol)
	gf := cg.Done(limitVolume, fees)
	pp.Println(gf)

	el := model.GetExchangeLimitCache(model.GetKey(cg.Exchange, cg.Symbol))
	grids, err := grid.Generate(grid.GSGrid, gf.MinPrice, gf.MaxPrice, gf.TotalSum, gf.GridNum, el.PricePrecision, el.QuantityPrecision)
	if err != nil {
		t.Error(err)
		return
	}
	grid.PrintFormat(grids)
}

func TestFloatToStr(t *testing.T) {
	x := 0.000121

	r1 := fmt.Sprintf("%v", x)
	r2 := strconv.FormatFloat(x, 'f', -1, 64)

	fmt.Println(r1, r2)
}
