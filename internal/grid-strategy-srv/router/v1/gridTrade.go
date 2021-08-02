package v1

import (
	"fmt"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/internal/grid-strategy-srv/model"
	"wq-fotune-backend/internal/grid-strategy-srv/util/goex"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/render"
)

// StartupGridStrategy 启动网格交易策略
func StartupGridStrategy(c *gin.Context) {
	form := &gridTradeForm{}
	err := render.BindJSON(c, form)
	if err != nil {
		//render.Err400Msg(c, err.Error())
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
	grids, err := form.initGrid()
	if err != nil {
		logger.Error("根据参数初始化网格失败", logger.Err(err), logger.Any("form", form.desensitize()))
		render.Err500(c, err.Error())
		return
	}
	logger.Info("根据参数初始化网格成功", logger.String("uid", form.UID), logger.Any("grids", grids))

	// 检测账号下持仓量和资金余额是否能够满足网格所需金额
	needBuyCoinQuantity, err := form.checkAccountBalance(grids)
	if err != nil {
		logger.Error("账号余额不足", logger.Err(err), logger.Any("form", form.desensitize()))
		render.Err500(c, "账号余额不足")
		return
	}
	logger.Info("账号余额满足网格策略", logger.String("uid", form.UID), logger.Float64("needBuyCoinQuantity", needBuyCoinQuantity))

	// 买入市价单，如果返回的needBuyCoinQuantity为0，则忽略
	marketRecord, err := form.placeMarketOrder("buy", needBuyCoinQuantity)
	if err != nil {
		logger.Error("买入市价单失败", logger.Err(err), logger.Float64("needBuyCoinQuantity", needBuyCoinQuantity), logger.Any("form", form.desensitize()))
		render.Err500(c, "买入市价单失败")
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

// StopGridStrategy 停止网格交易策略
func StopGridStrategy(c *gin.Context) {
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
		//render.Err500(c, err.Error())
		render.Err500(c, "初始化交易所账号失败")
		return
	}
	logger.Info("初始化交易所账号成功", logger.String("uid", form.UID))

	// 取消在交易所的委托
	successOrderIDs, err := form.cancelCommissionOrder()
	if err != nil {
		logger.Error("取消交易所的委托订单失败", logger.Err(err), logger.Any("form", form))
		//render.Err500(c, err.Error())
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
		err = form.closePosition()
		if err != nil {
			logger.Info("closePosition error", logger.Err(err))
			closePosition = "平仓失败"
		}
		out += ", " + closePosition
	}

	render.OK(c, out)
}

// UpdateGridStrategy 更新网格策略
func UpdateGridStrategy(c *gin.Context) {
	form := &updateGridStrategyForm{}

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

	grids, err := form.genNewGrid()
	if err != nil {
		logger.Error("genNewGrid error", logger.Err(err), logger.Any("form", form))
		render.Err500(c, err.Error())
		return
	}

	// 初始化交易所账号
	err = form.initExchangeAccount()
	if err != nil {
		logger.Error("初始化交易所账号失败", logger.Err(err), logger.Any("form", form))
		render.Err500(c, "初始化交易所账号失败")
		return
	}
	logger.Info("初始化交易所账号成功", logger.String("uid", form.gs.UID))

	//err = updateGridStrategy(form.gs, form.exchangeAccount, grids, form.MinPrice, form.MaxPrice, form.latestPrice)
	err = model.UpdateRunningGridStrategy(form.gs, form.exchangeAccount, grids, form.MinPrice, form.MaxPrice, form.latestPrice, form.IsClosePosition)
	if err != nil {
		logger.Error("updateGridStrategy error", logger.Err(err), logger.String("gsid", form.Gsid))
		render.Err500(c, err.Error())
		return
	}

	render.OK(c)
}

// CalculateMoney 计算整个网格交易所需总资金
func CalculateMoney(c *gin.Context) {
	form := &gridTradeForm{}
	err := render.BindJSON(c, form)
	if err != nil {
		render.Err400Msg(c, "json解析失败")
		return
	}

	err = form.valid2()
	if err != nil {
		logger.Error("参数错误", logger.Err(err), logger.Any("form", form))
		render.Err400Msg(c, err.Error())
		return
	}

	needMoney, err := form.calculateGridNeedMoney()
	if err != nil {
		logger.Error("calculateGridNeedMoney error", logger.Err(err), logger.Any("form", form.desensitize()))
		render.Err400Msg(c, "计算网格金额失败")
		return
	}
	el := model.GetExchangeLimitCache(model.GetKey(form.Exchange, form.Symbol))
	out := gin.H{
		"minPrice": FloatRound(form.MinPrice, el.PricePrecision),
		"maxPrice": FloatRound(form.MaxPrice, el.PricePrecision),
		"gridNum":  form.GridNum,
		"totalSum": needMoney,
	}

	render.OK(c, out)
}

// CalculateGridParams 通过参数计算网格参数
func CalculateGridParams(c *gin.Context) {
	form := &calculateGridParamForm{
		Exchange:   c.Query("exchange"),
		Symbol:     c.Query("symbol"),
		MinPrice:   c.Query("minPrice"),
		MaxPrice:   c.Query("maxPrice"),
		ProfitRate: c.Query("profitRate"),
	}

	err := form.valid()
	if err != nil {
		logger.Error("参数无效", logger.Err(err))
		render.Err400Msg(c, err.Error())
		return
	}

	bgp, _, err := form.calculateParams()
	if err != nil {
		logger.Error("calculateParams error", logger.Err(err), logger.Any("form", form))
		render.Err500(c, err)
		return
	}
	//averageProfit, averageProfitRate := model.CalculateProfit(form.Exchange, grids)
	//num := len(grids)
	//intervalSize := (grids[0].Price - grids[1].Price + grids[num-2].Price - grids[num-1].Price) / 2

	render.OK(c, bgp)
}

// GetMinMoney 根据交易所和品种获取最小投入资金
func GetMinMoney(c *gin.Context) {
	form := &autoGenerateForm{
		Exchange: c.Query("exchange"),
		Symbol:   c.Query("symbol"),
	}

	err := form.valid1()
	if err != nil {
		render.Err400Msg(c, err.Error())
		return
	}

	render.OK(c, gin.H{
		"minMoney": form.getMinMoney(),
	})
}

// GetGridParams 根据投入资金自动生成网格参数
func GetGridParams(c *gin.Context) {
	form := &autoGenerateForm{
		Exchange:  c.Query("exchange"),
		Symbol:    c.Query("symbol"),
		TotalSum:  str2Float64(c.Query("totalSum")),
		IsReverse: false,
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

	out := gin.H{
		"minPrice":          FloatRound(minPrice, pricePrecision),
		"maxPrice":          FloatRound(maxPrice, pricePrecision),
		"gridNum":           gridNum,
		"totalSum":          form.TotalSum,
		"averageProfitRate": averageProfitRate,
	}

	render.OK(c, out)
}

// GetBigGridParams 获取大网格参数
func GetBigGridParams(c *gin.Context) {
	form := &bigGridParamForm{
		Exchange:   c.Query("exchange"),
		Symbol:     c.Query("symbol"),
		IsAI:       c.Query("isAI"),
		MinPrice:   c.Query("minPrice"),
		ProfitRate: c.Query("profitRate"),
	}

	err := form.valid()
	if err != nil {
		logger.Error("参数无效", logger.Err(err))
		render.Err400Msg(c, err.Error())
		return
	}

	bgp, err := form.calculateParams()
	if err != nil {
		logger.Error("calculateParams error", logger.Err(err), logger.Any("form", form))
		render.Err500(c, err)
		return
	}

	render.OK(c, bgp)
}

// ListRunningStrategies 获取正在运行的策略列表
func ListRunningStrategies(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		//render.Err400Msg(c, "uid is empty")
		render.Err400Msg(c, "参数uid为空")
		return
	}

	form := &reqListForm{
		pageStr:  c.Query("page"),
		limitStr: c.Query("limit"),
		sort:     c.Query("sort"),
	}

	err := form.valid()
	if err != nil {
		logger.Error("参数错误", logger.Err(err), logger.Any("form", form))
		//render.Err400Msg(c, err.Error())
		render.Err400Msg(c, "参数错误")
		return
	}

	out, total, err := getRunningStrategies(uid, form)
	if err != nil {
		logger.Error("获取策略信息失败", logger.Err(err), logger.Any("form", form))
		render.Err500(c, "获取策略信息失败")
		return
	}

	render.OK(c, gin.H{"strategies": out, "total": total})
}

// GetStrategyDetail 获取策略详情
func GetStrategyDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		//render.Err400Msg(c, "id is empty")
		render.Err400Msg(c, "参数id为空")
		return
	}

	// 获取策略信息
	out, err := getStrategyDetail(id)
	if err != nil {
		logger.Error("获取策略信息失败", logger.Err(err), logger.String("id", id))
		render.Err500(c, "获取策略详情信息失败")
		return
	}

	render.OK(c, gin.H{"strategyDetail": out})
}

// GetStrategySimple 获取策略简要信息
func GetStrategySimple(c *gin.Context) {
	gsid := c.Param("id")
	if gsid == "" {
		render.Err400Msg(c, "策略id为空")
		return
	}

	query := bson.M{"_id": bson.ObjectIdHex(gsid)}
	gs, err := model.FindGridStrategy(query, bson.M{})
	if err != nil {
		logger.Error("model.FindGridStrategy error", logger.Err(err), logger.Any("query", query))
		render.Err500(c, "获取策略信息失败")
		return
	}
	if gs.AverageProfit == 0.0 {
		query = bson.M{"gsid": bson.ObjectIdHex(gsid)}
		gpo, err := model.FindGridPendingOrder(query, bson.M{})
		if err != nil {
			logger.Warn("model.FindGridPendingOrder error", logger.Err(err), logger.String("gsid", gsid))
		} else {
			gs.AverageProfit, gs.AverageProfitRate = model.CalculateProfit(gs.Exchange, gpo.Grids)
			query = bson.M{"_id": bson.ObjectIdHex(gsid)}
			update := bson.M{"$set": bson.M{"averageProfit": gs.AverageProfit, "averageProfitRate": gs.AverageProfitRate}}
			model.UpdateGridStrategy(query, update)
		}
	}
	el := model.GetExchangeLimitCache(model.GetKey(gs.Exchange, gs.Symbol))

	render.OK(c, gin.H{
		"exchange":          gs.Exchange,
		"symbol":            gs.Symbol,
		"anchorSymbol":      goex.GetAnchorCurrency(gs.Symbol),
		"totalSum":          gs.TotalSum,
		"minPrice":          FloatRound(gs.MinPrice, el.PricePrecision),
		"maxPrice":          FloatRound(gs.MaxPrice, el.PricePrecision),
		"gridNum":           gs.GridNum,
		"averageProfit":     gs.AverageProfit,
		"averageProfitRate": model.Float64ToStr(gs.AverageProfitRate*100) + "%",
		"type":              gs.Type,
	})
}

//GetStrategyTotal 获取一共有多少个网格策略在运行
func GetStrategyTotal(c *gin.Context) {
	gps := model.GetStrategyCaches()

	query := bson.M{"isRun": true}
	gss, err := model.FindGridStrategies(query, bson.M{}, 0, 1000000)
	if err != nil {
		logger.Error("FindGridStrategies error", logger.Any("query", query))
		render.Err500(c, "获取策略数量失败")
		return
	}

	out := gin.H{
		"realStrategyTotal":  len(gss),
		"cacheStrategyTotal": len(gps),
	}
	render.OK(c, out)
}

// GetStrategyCount 查看用户当前正在运行策略数量
func GetStrategyCount(c *gin.Context) {
	uid := c.Param("uid")

	n := getStrategyCount(uid)

	render.OK(c, n)
}
