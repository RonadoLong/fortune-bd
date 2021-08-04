package v1

import (
	"fmt"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/app/grid-strategy-srv/model"
	"wq-fotune-backend/app/grid-strategy-srv/util/goex"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/logger"
	"github.com/zhufuyi/pkg/render"


)

// GetReverseMinMoney 获取反向网格最小投资金额
func GetReverseMinMoney(c *gin.Context) {
	form := &reverseGridForm{
		Exchange: c.Query("exchange"),
		Symbol:   c.Query("symbol"),
	}

	err := form.valid()
	if err != nil {
		render.Err400Msg(c, err.Error())
		return
	}

	minMoney, unit := form.getMinMoney()

	render.OK(c, gin.H{
		"minMoney": minMoney,
		"unit":     unit,
	})
}

// CalculateReverseGridParams 根据投入资金自动生成反向网格参数
func CalculateReverseGridParams(c *gin.Context) {
	totalSum := str2Float64(c.Query("totalSum"))
	form := &autoGenerateForm{
		Exchange:  c.Query("exchange"),
		Symbol:    c.Query("symbol"),
		TotalSum:  totalSum,
		IsReverse: true,
	}

	err := form.valid2()
	if err != nil {
		logger.Error("参数无效", logger.Err(err))
		render.Err400Msg(c, err.Error())
		return
	}

	minPrice, maxPrice := 0.0, 0.0
	gridNum := 0
	averageProfitRate := ""

	gf, err := form.getBestGrid()
	if err != nil {
		logger.Warn("not found bestGrid, use default params", logger.Err(err), logger.Any("form", form))
		minPrice, maxPrice = getMinMax(form.latestPrice, env.GridNum)
		gridNum = model.GetDefaultGridNum(form.TotalSum)
	} else {
		minPrice = gf.MinPrice
		maxPrice = gf.MaxPrice
		gridNum = gf.GridNum

		if form.TotalSum/float64(gridNum)-form.esl.VolumeLimit == 0.0 {
			gridNum--
			if gridNum < 5 {
				gridNum = 5
			}
			minPrice += gf.AverageIntervalPrice / 2
			maxPrice += gf.AverageIntervalPrice / 2
		}
		averageProfitRate = fmt.Sprintf("%v", FloatRound(gf.AverageProfitRate*100)) + "%"
	}

	pricePrecision := form.esl.PricePrecision
	if pricePrecision == 0 {
		pricePrecision = 6
	}

	tradeCurrency, anchorCurrency := goex.SplitSymbol(form.Symbol)

	out := gin.H{
		"minPrice":          FloatRound(minPrice, pricePrecision),
		"maxPrice":          FloatRound(maxPrice, pricePrecision),
		"priceUnit":         anchorCurrency,
		"gridNum":           gridNum,
		"totalSum":          totalSum,
		"totalSumUnit":      tradeCurrency,
		"averageProfitRate": averageProfitRate,
	}

	render.OK(c, out)
}

// StartupReverseGridStrategy 启动反向网格策略
func StartupReverseGridStrategy(c *gin.Context) {
	form := &gridTradeForm{}
	err := render.BindJSON(c, form)
	if err != nil {
		logger.Warn("json解析失败", logger.Err(err))
		render.Err400Msg(c, "json解析失败")
		return
	}

	err = form.valid()
	if err != nil {
		logger.Error("参数错误", logger.Err(err), logger.Any("form", form))
		render.Err400Msg(c, err.Error())
		return
	}

	// 初始化交易所账号
	err = form.initExchangeAccount()
	if err != nil {
		logger.Error("初始化交易所账号失败", logger.Err(err), logger.Any("form", form))
		render.Err500(c, "初始化交易所账号失败")
		return
	}
	logger.Info("初始化交易所账号成功", logger.String("uid", form.UID))

	// 初始化网格
	grids, err := form.initReverseGrid()
	if err != nil {
		logger.Error("根据参数初始化网格失败", logger.Err(err), logger.Any("form", form.desensitize()))
		render.Err500(c, err.Error())
		return
	}
	logger.Info("根据参数初始化网格成功", logger.String("uid", form.UID), logger.Any("grids", grids))

	// 检测账号下持仓量和资金余额是否能够满足网格
	needSellCoinQuantity, err := form.calculateTradeCurrencySize(grids)
	if err != nil {
		logger.Error("账号余额不足", logger.Err(err), logger.Any("form", form.desensitize()))
		render.Err500(c, "账号余额不足")
		return
	}
	logger.Info("账号余额满足网格策略", logger.String("uid", form.UID), logger.Float64("needSellCoinQuantity", needSellCoinQuantity))

	// 卖出市价单，如果返回的needBuyCoinQuantity为0，则忽略
	marketRecord, err := form.placeMarketOrder("sell", needSellCoinQuantity)
	if err != nil {
		logger.Error("卖出市价单失败", logger.Err(err), logger.Float64("needSellCoinQuantity", needSellCoinQuantity), logger.Any("form", form.desensitize()))
		render.Err500(c, "启动网格失败，原因是卖出市价单失败")
		return
	}

	// 拆分挂单
	firstKey, delayKey := splitGridKeys(grids, form.GridBasisNO)
	// 买入、卖出网格限价单
	limitRecords, err := form.placeLimitOrder(grids, firstKey)
	if err != nil {
		logger.Error("买入、卖出网格限价单失败", logger.Err(err), logger.Any("grids", grids), logger.Any("form", form.desensitize()))
		render.Err500(c, "买入、卖出网格限价单失败")
		return
	}
	//if len(limitRecords) == len(grids)-1 {
	//	logger.Info("买入、卖出网格所有限价单成功", logger.String("uid", form.UID), logger.Int("orderSize", len(limitRecords)))
	//} else {
	//	logger.Warn("有部分买入、卖出网格限价单失败", logger.String("uid", form.UID), logger.Int("gridSize", len(grids)), logger.Int("orderSize", len(limitRecords)))
	//}

	// 保存信息
	records := append(limitRecords, marketRecord)
	gsid := form.saveInfo(grids, records)
	logger.Infof("网格策略(%s)已初始化和持久化成功，开始订阅和处理网格订单", form.gridStrategyID)
	go strategyStartUpNotify(form.UID, form.gridStrategyID)

	go func() {

		gp := form.toGridProcess(grids, form.GridBasisNO, gsid)
		// 写入缓存
		key := model.StrategyCacheKey(form.UID, form.Exchange, form.Symbol, gsid)
		model.SetStrategyCache(key, gp)

		// 订阅订单通知和处理网格订单
		err = model.ProcessGridOrder(gp)
		if err != nil {
			logger.Error("订阅订单通知和处理网格订单失败", logger.Err(err), logger.Any("gridProcess", gp.Desensitize()))
			return
		}

		form.placeLimitOrderAndSave(gp, grids, delayKey) // 延时下单
		gp.ListenOrder()                                 // 定时检查委托订单
	}()

	render.OK(c)
}

// StopReverseGridStrategy 停止反向网格交易策略
func StopReverseGridStrategy(c *gin.Context) {
	form := &cancelGridStrategyForm{}
	err := render.BindJSON(c, form)
	if err != nil {
		render.Err400Msg(c, "json解析失败")
		return
	}

	err = form.valid()
	if err != nil {
		logger.Error("参数错误", logger.Err(err), logger.Any("form", form))
		render.Err400Msg(c, err.Error())
		return
	}

	// 初始化交易所账号
	err = form.initExchangeAccount()
	if err != nil {
		logger.Error("初始化交易所账号失败", logger.Err(err), logger.Any("form", form))
		render.Err500(c, "停止策略失败，原因是初始化交易所账号失败")
		return
	}
	logger.Info("初始化交易所账号成功", logger.String("uid", form.UID))

	// 取消在交易所的委托
	successOrderIDs, err := form.cancelCommissionOrder()
	if err != nil {
		logger.Error("取消交易所的委托订单失败", logger.Err(err), logger.Any("form", form))
		render.Err500(c, "取消交易所的委托订单失败")
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
		logger.Info("本地更新所有订单状态成功", logger.String("uid", form.UID))
	}

	// 修改策略运行状态
	err = updateStrategyOrder(form.Gsid)
	if err != nil {
		logger.Error("更新策略运行状态失败", logger.String("uid", form.UID), logger.Err(err), logger.String("gsid", form.Gsid))
	} else {
		logger.Info("更新策略运行状态成功", logger.String("uid", form.UID), logger.String("gsid", form.Gsid))
	}

	// 删除策略运行的缓存信息
	key := model.StrategyCacheKey(form.UID, form.Exchange, form.Symbol, form.Gsid)
	sc, ok := model.GetStrategyCache(key)
	if ok {
		sc.Close()
	} else {
		logger.Warn("not found strategy in StrategyCache", logger.String("key", key))
	}
	model.DeleteStrategyCache(key)

	out := "停止策略成功"
	// 判断是否平仓
	if form.IsClosePosition {
		closePosition := "平仓成功"
		err = form.closeReverseGridPosition()
		if err != nil {
			logger.Info("closePosition error", logger.Err(err))
			closePosition = "平仓失败"
		}
		out += ", " + closePosition
	}

	render.OK(c, out)
}
