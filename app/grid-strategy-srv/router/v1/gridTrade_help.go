package v1

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/app/grid-strategy-srv/model"
	"wq-fotune-backend/app/grid-strategy-srv/util/goex"
	"wq-fotune-backend/app/grid-strategy-srv/util/goex/binance"
	"wq-fotune-backend/app/grid-strategy-srv/util/grid"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi"

	"sort"

	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/gohttp"
	"github.com/zhufuyi/pkg/krand"
	"github.com/zhufuyi/pkg/logger"
)

type gridTradeForm struct {
	UID       string `json:"uid"`    // 用户id
	ApiKey    string `json:"apiKey"` // 交易所账号apikey
	secretKey string `json:"-"`      // 交易所账号访问密钥

	Type int `json:"type"` // 0:网格交易，1:趋势网格，2:无限网格，3:反向网格

	Exchange     string `json:"exchange"`     // 交易所
	Symbol       string `json:"symbol"`       // 交易品种
	AnchorSymbol string `json:"anchorSymbol"` // 锚定币的品种名称，例如USDT

	GridIntervalType string  `json:"gridIntervalType"` // 网格间隔类型，ASGrid:等差, GSGrid:等比
	MinPrice         float64 `json:"minPrice"`         // 网格最低价格
	MaxPrice         float64 `json:"maxPrice"`         // 网格最高价格
	TotalSum         float64 `json:"totalSum"`         // 投资总额
	GridNum          int     `json:"gridNum"`          // 网格数量
	GridBasisNO      int     `json:"-"`                // 网格基准线编号

	StopProfitPrice float64 `json:"stopProfitPrice"` // 止盈价格
	StopLossPrice   float64 `json:"stopLossPrice"`   // 止损价格
	BasisPrice      float64 `json:"basisPrice"`      // 网格开单价格，当价格小于等于基准价格时买入，否则卖出，用来统计需要的资金和已买入的币，当基准价格为0时，使用当前币价格
	LatestPrice     float64 `json:"-"`               // 当前最新价格

	exchangeAccount  goex.Accounter `json:"-"` // 交易所账号
	gridStrategyID   string         `json:"-"` // 网格策略id
	buyCoinQuantity  float64        `json:"-"` // 启动网格时以市价单买入币的数量
	sellCoinQuantity float64        `json:"-"` // 启动网格时以市价单卖出币的数量
}

func (g *gridTradeForm) valid() error {
	switch "" {
	case g.UID:
		//return errors.New("field uid is empty")
		return errors.New("参数uid为空")
	case g.ApiKey:
		//return errors.New("field apiKey is empty")
		return errors.New("参数apiKey为空")
	case g.Exchange:
		//return errors.New("field exchange is empty")
		return errors.New("参数exchange为空")
	case g.Symbol:
		//return errors.New("field symbol is empty")
		return errors.New("参数symbol为空")
	//case g.AnchorSymbol:
	//	//return errors.New("field anchorSymbol is empty")
	//	return errors.New("参数anchorSymbol为空")
	case g.GridIntervalType:
		g.GridIntervalType = grid.GSGrid
	}

	if g.Exchange != model.ExchangeHuobi && g.Exchange != model.ExchangeBinance {
		return fmt.Errorf("暂时不支持%s交易所", g.Exchange)
	}
	// 判断该品种网格策略是否已经存在
	//key := model.StrategyCacheKey(g.UID, g.Exchange, g.Symbol)
	//if _, ok := model.GetStrategyCache(key); ok {
	//	return fmt.Errorf("grid strategy already exists, only one grid strategy is allowed to run per symbol(%s)", g.Symbol)
	//	return fmt.Errorf("同一个品种%s下只允许运行一个网格策略", g.Symbol)
	//}

	switch g.Type {
	case model.GridTypeNormal, model.GridTypeTrend, model.GridTypeReverse:
		if g.GridNum < 5 || g.GridNum > 100 {
			return errors.New("网格数设置必须在5 ~ 100范围")
		}
	case model.GridTypeInfinite:
		g.GridIntervalType = grid.GSGrid
		if g.GridNum < 5 || g.GridNum > 200 {
			return errors.New("网格数设置必须在5 ~ 200范围")
		}
	default:
		return errors.New("不支持未知的网格类型")
	}
	//if g.Type == model.GridTypeNormal || g.Type == model.GridTypeTrend || g.Type == model.GridTypeReverse {
	//	if g.GridNum < 5 || g.GridNum > 100 {
	//		//return errors.New("gridNum value must be  in the range of 5 to 60")
	//		return errors.New("网格数设置必须在5 ~ 100范围")
	//	}
	//} else { // 无限网格最多200个网格
	//	if g.GridNum < 5 || g.GridNum > 200 {
	//		//return errors.New("gridNum value must be  in the range of 5 to 60")
	//		return errors.New("网格数设置必须在5 ~ 200范围")
	//	}
	//}
	//if g.Type == model.GridTypeInfinite { // 无限网格使用等比数列
	//	g.GridIntervalType = grid.GSGrid
	//}

	if g.MinPrice <= 0.0 {
		//return errors.New("field buyLimit is illegality")
		return errors.New("网格最低价格不能小于0")
	}
	if g.MaxPrice <= 0.0 {
		//return errors.New("field sellLimit is illegality")
		return errors.New("网格最高价格不能小于0")
	}
	if g.MinPrice >= g.MaxPrice {
		//return fmt.Errorf("minPrice=%v cannot be greater than maxPrice=%v", g.MinPrice, g.MaxPrice)
		return fmt.Errorf("网格参数最低价格(%v)不能超过最高价格(%v)", g.MinPrice, g.MaxPrice)
	}

	el := model.GetExchangeLimitCache(model.GetKey(g.Exchange, g.Symbol))
	if el.VolumeLimit <= 0.0 {
		//return fmt.Errorf("not supported symbol (%s %s) for the moment, need add it", g.Exchange, g.Symbol)
		return fmt.Errorf("暂时不支持%s交易所的%s品种", g.Exchange, g.Symbol)
	}

	if g.BasisPrice > 0.0 && (g.BasisPrice < g.MinPrice || g.BasisPrice > g.MaxPrice) {
		//return fmt.Errorf("basisPrice(%v) cannot be outside the range of minPrice(%v) and maxPrice(%v)", g.BasisPrice, g.MinPrice, g.MaxPrice)
		return fmt.Errorf("设置网格价格范围(%v ~ %v)已经超出最新价格%v", g.MinPrice, g.MaxPrice, g.BasisPrice)
	}

	if g.StopProfitPrice > 0.0 && g.StopProfitPrice <= g.MaxPrice {
		//return fmt.Errorf("stopProfitPrice(%v) must be greater than maxPrice(%v)", g.StopProfitPrice, g.MaxPrice)
		return fmt.Errorf("止盈价格%v不能低于网格最高价格%v", g.StopProfitPrice, g.MaxPrice)
	}

	if g.StopLossPrice > 0.0 && g.StopLossPrice >= g.MinPrice {
		//return fmt.Errorf("stopLossPrice(%v) cannot be greater than minPrice(%v)", g.StopLossPrice, g.MinPrice)
		return fmt.Errorf("止损价格%v不能低于网格最高价格%v", g.StopLossPrice, g.MinPrice)
	}

	latestPrice, err := model.GetLatestPrice(g.Exchange, g.Symbol)
	if err != nil {
		return err
	}
	g.LatestPrice = latestPrice

	// 如果basisPrice为空，从交易所获取最新价格
	if g.BasisPrice <= 0.0 {
		g.BasisPrice = latestPrice
	}

	if g.LatestPrice < g.MinPrice || g.LatestPrice > g.MaxPrice {
		//return fmt.Errorf("basis price(%f) is outside the range of minPrice(%f) and maxPrice(%f), grid is illegality", g.BasisPrice, g.MinPrice, g.MaxPrice)
		return fmt.Errorf("当前最新价格(%v)超出网格范围(%v ~ %v)", g.LatestPrice, g.MinPrice, g.MaxPrice)
	}

	g.AnchorSymbol = goex.GetAnchorCurrency(g.Symbol)
	var minLimit, maxLimit float64
	if g.Type == model.GridTypeReverse { // 反向网格
		minLimit, maxLimit, err = getReverseMinAndMaxTotalSum(g.Exchange, g.Symbol)
		if err != nil {
			return errors.New("获取投资金额范围失败")
		}
	} else {
		minLimit, maxLimit = getMinAndMaxTotalSum(g.Exchange, g.AnchorSymbol, g.LatestPrice, el.VolumeLimit)
	}
	if g.TotalSum < minLimit || g.TotalSum > maxLimit {
		return fmt.Errorf("投资金额设置必须在%v ~ %v范围", FloatRound(minLimit, el.PricePrecision), maxLimit)
	}

	if g.Type != model.GridTypeReverse {
		if g.AnchorSymbol != "usdt" {
			g.TotalSum = FloatRound(g.TotalSum, el.GetVolumePrecision())
		}
	}
	g.gridStrategyID = bson.NewObjectId().Hex() // 网格策略id

	n := getStrategyCount(g.UID)
	if n > 20 {
		return errors.New("已经运行的机器人数达到最大限制20个，启动新机器人失败")
	}

	return nil
}

func (g *gridTradeForm) valid2() error {
	switch "" {
	case g.Exchange:
		//return errors.New("field exchange is empty")
		return errors.New("参数exchange为空")
	case g.Symbol:
		//return errors.New("field symbol is empty")
		return errors.New("参数symbol为空")
	}

	el := model.GetExchangeLimitCache(model.GetKey(g.Exchange, g.Symbol))
	if el.QuantityLimit <= 0.0 {
		//return fmt.Errorf("not supported symbol (%s %s) for the moment, need add it", g.Exchange, g.Symbol)
		return fmt.Errorf("暂时不支持%s交易所的%s品种", g.Exchange, g.Symbol)
	}

	if g.GridIntervalType == "" {
		g.GridIntervalType = grid.GSGrid
	}

	// 如果没有设置网格，使用默认网格数
	if g.GridNum < 5 || g.GridNum > 100 {
		g.GridNum = env.GridNum
	}

	var latestPrice float64
	var err error
	// 区分不同交易所
	switch g.Exchange {
	case model.ExchangeHuobi:
		latestPrice, err = huobi.GetLatestPrice(g.Symbol)
		if err != nil {
			return err
		}
	case model.ExchangeBinance:
		// 根据品种获取当前价格
		latestPrice, err = binance.GetLatestPrice(g.Symbol, env.ProxyAddr)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("暂时不支持%s交易所", g.Exchange)
	}

	g.LatestPrice = latestPrice
	minPrice, maxPrice := getMinMax(g.LatestPrice, g.GridNum)
	// 默认30%是挂买单
	if g.MinPrice <= 0.0 {
		g.MinPrice = minPrice
	}
	// 默认70%是挂卖单
	if g.MaxPrice <= 0.0 {
		g.MaxPrice = maxPrice
	}
	if g.MinPrice >= g.MaxPrice {
		//return fmt.Errorf("minPrice=%v cannot be greater than maxPrice=%v", g.MinPrice, g.MaxPrice)
		return fmt.Errorf("网格最低价格%v不能超过最高%v", g.MinPrice, g.MaxPrice)
	}

	return nil
}

// 初始化交易所账号
func (g *gridTradeForm) initExchangeAccount() error {
	accessKey, secretKey, account, err := model.InitExchangeAccount(g.UID, g.Exchange, g.ApiKey)
	if err != nil {
		return err
	}
	g.exchangeAccount = account
	g.ApiKey = accessKey
	g.secretKey = secretKey

	return nil
}

func (g *gridTradeForm) desensitize() gridTradeForm {
	gtf := *g
	gtf.secretKey = ""
	gtf.exchangeAccount = nil
	return gtf
}

// 初始化网格
func (g *gridTradeForm) initGrid() ([]*grid.Grid, error) {
	el := model.GetExchangeLimitCache(model.GetKey(g.Exchange, g.Symbol))

	grids, err := grid.Generate(
		g.GridIntervalType,
		g.MinPrice,
		g.MaxPrice,
		g.TotalSum,
		g.GridNum,
		el.PricePrecision,
		el.QuantityPrecision,
	)
	if err != nil {
		return grids, err
	}

	err = g.checkGrids(grids) // 检查网格

	return grids, err
}

// 判断网格是否满足交易所的最小金额和数量限制
func (g *gridTradeForm) checkGrids(grids []*grid.Grid) error {
	el := model.GetExchangeLimitCache(model.GetKey(g.Exchange, g.Symbol))

	for k, v := range grids {
		// 判断买
		if k > 0 {
			if v.BuyQuantity < el.QuantityLimit {
				//return fmt.Errorf("buyQuantity(%v) must be greater than %v", v.BuyQuantity, el.QuantityLimit)
				logger.Warnf("每格最小交易数量不能小于%v，当前值为%v", el.QuantityLimit, v.BuyQuantity)
				return errors.New("网格参数不合法，请增大投资金额或减小网格数量")
			}
			volume := v.BuyQuantity * v.Price
			if volume < el.VolumeLimit {
				//return fmt.Errorf("buy volume(%v) must be greater than %v", volume, el.VolumeLimit)
				logger.Warnf("每格最小交易金额不能小于%v，当前值为%v", el.VolumeLimit, volume)
				return errors.New("网格参数不合法，请增大投资金额或减小网格数量")
			}
		}

		// 判断卖
		if k < g.GridNum {
			if v.SellQuantity < el.QuantityLimit {
				//return fmt.Errorf("sellQuantity(%v) must be greater than %v", v.SellQuantity, el.QuantityLimit)
				logger.Warnf("每格最小交易数量不能小于%v，当前值为%v", el.QuantityLimit, v.SellQuantity)
				return errors.New("网格参数不合法，请增大投资金额或减小网格数量")
			}
			volume := v.SellQuantity * v.Price
			if volume < el.VolumeLimit {
				//return fmt.Errorf("sell volume(%v) must be greater than %v", volume, el.VolumeLimit)
				logger.Warnf("每格最小交易金额不能小于%v，当前值为%v", el.VolumeLimit, volume)
				return errors.New("网格参数不合法，请增大投资金额或减小网格数量")
			}
		}
	}

	_, averageProfitRate := model.CalculateProfit(g.Exchange, grids)
	if averageProfitRate < 0.0005 || averageProfitRate > 0.05 {
		logger.Warn("网格参数不合法，每格纯利润率范围0.05%~5%，当前值为" + model.Float64ToStr(averageProfitRate*100) + "%")
		return errors.New("网格参数不合法，请增大投资金额或减小网格数量")
	}

	return nil
}

// 检测账号下持仓量和资金余额是否能够满足网格
func (g *gridTradeForm) checkAccountBalance(grids []*grid.Grid) (float64, error) {
	basisPrice := g.BasisPrice
	latestPrice := g.LatestPrice
	gridBasisNO := 0
	needBuyCoinQuantity := 0.0
	needMoney := 0.0
	err := error(nil)

	buyCoin := 0.0
	for k, v := range grids {
		if v.Price > basisPrice { // 大于基准线价格，需要卖出币的数量，也就是账号下必须已经持币数量
			gridBasisNO = k // 网格编号是有序的，最后一个大于basisPrice对应编号
			needBuyCoinQuantity += v.SellQuantity
		} else { // 小于等于基准线价格，统计委托挂单需要的金额
			needMoney += v.Price * v.BuyQuantity
			buyCoin += v.BuyQuantity
		}
	}
	needBuyCoinQuantity -= grids[gridBasisNO].SellQuantity // 减去和当前价格相近的卖单数量，不需要挂卖单

	// 区分不同交易所
	switch g.Exchange {
	case model.ExchangeHuobi:
		needBuyCoinQuantity = needBuyCoinQuantity/(1-huobi.FilledFees) + buyCoin*huobi.FilledFees // 扣除手续费之后的持仓数量
	case model.ExchangeBinance:
		needBuyCoinQuantity = needBuyCoinQuantity/(1-binance.FilledFees) + buyCoin*binance.FilledFees // 扣除手续费之后的持仓数量
	}

	if g.exchangeAccount == nil {
		return needBuyCoinQuantity, errors.New("not found account")
	}

	el := model.GetExchangeLimitCache(model.GetKey(g.Exchange, g.Symbol))

	// 查询账户基准币的余额，判断余额是否满足网格所需的金额
	anchorCurrencyBalance, err := g.exchangeAccount.GetCurrencyBalance(g.AnchorSymbol)
	if err != nil {
		return needBuyCoinQuantity, fmt.Errorf("get currency balance error, err=%s", err.Error())
	}
	needTotalMoney := needMoney + needBuyCoinQuantity*latestPrice
	logger.Info("check account balance", logger.Float64("gridNeedTotalMoney", needTotalMoney), logger.Float64(fmt.Sprintf("%sBalance", g.AnchorSymbol), anchorCurrencyBalance), logger.Int("quantityPrecision", el.QuantityPrecision))
	if anchorCurrencyBalance-needTotalMoney < 0.0 {
		return needBuyCoinQuantity, fmt.Errorf("%s balance(%f) is less than grid need money(%f)", g.AnchorSymbol, anchorCurrencyBalance, needTotalMoney)
	}

	g.GridBasisNO = gridBasisNO
	g.buyCoinQuantity = FloatRound(needBuyCoinQuantity, el.QuantityPrecision)

	return g.buyCoinQuantity, nil
}

func (g *gridTradeForm) calculateGridNeedMoney() (float64, error) {
	el := model.GetExchangeLimitCache(model.GetKey(g.Exchange, g.Symbol))
	g.TotalSum = float64(g.GridNum+1) * el.VolumeLimit // 给定最小价格

	// 初始化网格
	grids, err := g.initGrid()
	if err != nil {
		return 0, err
	}

	needMoney := 0.0
	moneyTmp := 0.0

	for _, v := range grids {
		if v.Price > g.LatestPrice {
			moneyTmp = g.LatestPrice * v.SellQuantity
			needMoney += moneyTmp
		} else {
			moneyTmp = v.Price * v.BuyQuantity
			needMoney += moneyTmp
		}
	}

	return FloatRound(needMoney+1, el.GetVolumePrecision()), nil
}

// 买入市价单
func (g *gridTradeForm) placeMarketOrder(side string, needBuyCoinQuantity float64) (*model.TradeRecord, error) {
	marketOrderRecord := &model.TradeRecord{}

	if needBuyCoinQuantity > 0.0 {
		if g.exchangeAccount == nil {
			return marketOrderRecord, errors.New("not found account")
		}

		el := model.GetExchangeLimitCache(model.GetKey(g.Exchange, g.Symbol))

		//side := "buy"
		//clientOrderID := fmt.Sprintf("mo_%s_%s", g.gridStrategyID, krand.String(krand.R_NUM|krand.R_LOWER, 4))
		clientOrderID := model.GenerateClientOrderID(g.Exchange, model.PrefixIDMob, g.gridStrategyID)
		volume := 0.0
		// 区分不同交易所
		switch g.Exchange {
		case model.ExchangeHuobi:
			volume = FloatRound(needBuyCoinQuantity*g.LatestPrice, el.PricePrecision) // 火币需要转换为总额
		case model.ExchangeBinance:
			volume = needBuyCoinQuantity
		}
		//coinQuantity := fmt.Sprintf("%v", volume)
		coinQuantity := model.Float64ToStr(volume, el.QuantityPrecision)

		orderID, err := g.exchangeAccount.PlaceMarketOrder(side, g.Symbol, coinQuantity, clientOrderID)
		if err != nil {
			return marketOrderRecord, err
		}

		logger.Info("place market order success",
			logger.String("side", side),
			logger.String("symbol", g.Symbol),
			logger.Float64("price", FloatRound(g.LatestPrice, el.PricePrecision)),
			logger.String("coinQuantity", coinQuantity),
			logger.String("clientOrderID", clientOrderID),
			logger.String("orderID", orderID),
		)

		marketOrderRecord = &model.TradeRecord{
			OrderID:       orderID,
			ClientOrderID: clientOrderID,
			OrderType:     "market",
			Side:          side,
			Price:         g.LatestPrice,
			Quantity:      needBuyCoinQuantity,
			Volume:        needBuyCoinQuantity * g.LatestPrice,
			Unit:          g.AnchorSymbol,
			OrderState:    "filled",
		}
	} else {
		logger.Info("no need to place market order")
	}

	return marketOrderRecord, nil
}

// 买入、卖出网格限价单
func (g *gridTradeForm) placeLimitOrder(grids []*grid.Grid, handleKeys map[int]bool) ([]*model.TradeRecord, error) {
	records := []*model.TradeRecord{}

	if g.exchangeAccount == nil {
		return records, errors.New("not found account")
	}

	var side, price, amount, clientOrderID, orderID string
	var err error
	for k, grid := range grids {
		if !handleKeys[k] {
			continue
		}
		if k == g.GridBasisNO { // 忽略接近当前价格的卖单
			continue
		}
		quantity, volume := 0.0, 0.0
		if k > g.GridBasisNO {
			// 买入限价单
			side = "buy"
			clientOrderID = model.NewGridClientOrderID(g.Exchange, side, g.gridStrategyID, k)
			amount = fmt.Sprintf("%v", grid.BuyQuantity)
			quantity = grid.BuyQuantity
			volume = grid.Price * grid.BuyQuantity
		} else {
			// 卖出限价单
			side = "sell"
			clientOrderID = model.NewGridClientOrderID(g.Exchange, side, g.gridStrategyID, k)
			amount = fmt.Sprintf("%v", grid.SellQuantity)
			quantity = grid.SellQuantity
			volume = grid.Price * grid.SellQuantity
		}

		//price = fmt.Sprintf("%v", grid.Price)
		price = strconv.FormatFloat(grid.Price, 'f', -1, 64)

		orderID, err = g.exchangeAccount.PlaceLimitOrder(side, g.Symbol, price, amount, clientOrderID)
		if err != nil {
			logger.Error("placeLimitOrder error", logger.Err(err), logger.String("param", fmt.Sprintf("%v, %v, %v, %v, %v", side, g.Symbol, price, amount, clientOrderID)))
			continue
		}

		grid.OrderID = orderID
		grid.Side = side
		logger.Info("place limit order success",
			logger.String("side", side),
			logger.String("symbol", g.Symbol),
			logger.String("price", price),
			logger.String("amount", amount),
			logger.String("clientOrderID", clientOrderID),
			logger.String("orderID", orderID),
		)

		records = append(records, &model.TradeRecord{
			GID:           k,
			OrderID:       orderID,
			ClientOrderID: clientOrderID,
			OrderType:     "limit",
			Side:          side,
			Price:         grid.Price,
			Quantity:      quantity,
			Volume:        volume,
			Unit:          g.AnchorSymbol,
			OrderState:    huobi.OrderStateSubmitted,
		})

		time.Sleep(time.Millisecond * 50) // 延时下单
	}

	return records, nil
}

// 买入、卖出网格限价单
func (g *gridTradeForm) placeLimitOrderAndSave(gp *model.GridProcess, grids []*grid.Grid, handleKeys map[int]bool) error {
	time.Sleep(time.Second) // 延时下单

	if g.exchangeAccount == nil {
		return errors.New("not found account")
	}

	var side, price, amount, clientOrderID, orderID string
	var err error
	for k, grid := range grids {
		if !handleKeys[k] {
			continue
		}
		if k == g.GridBasisNO { // 忽略接近当前价格的卖单
			continue
		}
		quantity, volume := 0.0, 0.0
		if k > g.GridBasisNO {
			// 买入限价单
			side = "buy"
			clientOrderID = model.NewGridClientOrderID(g.Exchange, side, g.gridStrategyID, k)
			amount = fmt.Sprintf("%v", grid.BuyQuantity)
			quantity = grid.BuyQuantity
			volume = grid.Price * grid.BuyQuantity
		} else {
			// 卖出限价单
			side = "sell"
			clientOrderID = model.NewGridClientOrderID(g.Exchange, side, g.gridStrategyID, k)
			amount = fmt.Sprintf("%v", grid.SellQuantity)
			quantity = grid.SellQuantity
			volume = grid.Price * grid.SellQuantity
		}

		//price = fmt.Sprintf("%v", grid.Price)
		price = strconv.FormatFloat(grid.Price, 'f', -1, 64)

		orderID, err = g.exchangeAccount.PlaceLimitOrder(side, g.Symbol, price, amount, clientOrderID)
		if err != nil {
			logger.Error("placeLimitOrder error", logger.Err(err), logger.String("param", fmt.Sprintf("%v, %v, %v, %v, %v", side, g.Symbol, price, amount, clientOrderID)))
			continue
		}

		grid.OrderID = orderID
		grid.Side = side
		gp.Grids[k].Side = side
		gp.Grids[k].OrderID = orderID

		isStartUpOrder := false
		if side == "sell" {
			isStartUpOrder = true
		}

		// 更新订单记录
		gtr := &model.GridTradeRecord{
			GSID:          bson.ObjectIdHex(g.gridStrategyID),
			GID:           k,
			OrderID:       orderID,
			ClientOrderID: clientOrderID,
			OrderType:     "limit",
			Side:          side,
			Price:         grid.Price,
			Quantity:      quantity,
			Volume:        volume,
			Unit:          g.AnchorSymbol,
			OrderState:    huobi.OrderStateSubmitted,

			StateTime:      time.Now().Local(),
			IsStartUpOrder: isStartUpOrder,
			Exchange:       g.Exchange,
			Symbol:         g.Symbol,
		}
		err = gtr.Insert()
		if err != nil {
			logger.Warn("gtr.Insert error", logger.Err(err), logger.Any("gtr", gtr))
			time.Sleep(time.Second)
			continue
		}

		// 更新网格参数
		query := bson.M{"gsid": gtr.GSID}
		update := bson.M{
			"$set": bson.M{
				fmt.Sprintf("grids.%d.side", k):    gtr.Side,
				fmt.Sprintf("grids.%d.orderId", k): gtr.OrderID,
			},
		}
		err := model.UpdateGridPendingOrder(query, update)
		if err != nil {
			logger.Warn("UpdateGridPendingOrder error", logger.Err(err), logger.Any("query", query), logger.Any("update", update))
			time.Sleep(time.Second)
			continue
		}

		logger.Info("place delay limit order success",
			logger.String("side", side),
			logger.String("symbol", g.Symbol),
			logger.String("price", price),
			logger.String("amount", amount),
			logger.String("clientOrderID", clientOrderID),
			logger.String("orderID", orderID),
		)
		time.Sleep(time.Second) // 延时下单
	}

	return nil
}

// 分批挂单，最近当前价格的先挂单
func splitGridKeys(grids []*grid.Grid, basisNO int) (map[int]bool, map[int]bool) {
	firstKeys, delayKeys := map[int]bool{}, map[int]bool{}

	for key := range grids {
		if math.Abs(float64(key-basisNO)) <= 5 {
			firstKeys[key] = true
		} else {
			delayKeys[key] = true
		}
	}

	return firstKeys, delayKeys
}

func getMinAndMaxTotalSum(exchange string, anchorSymbol string, latestPrice float64, volumeLimit float64) (float64, float64) {
	minLimit, maxLimit := 55.0, 0.0

	if anchorSymbol == "usdt" {
		if latestPrice < 1 {
			maxLimit = 10000.0
		} else {
			maxLimit = 20000.0
		}
	} else {
		minLimit = volumeLimit * 10.02
		// 获取当前锚地币的usdt计价价格
		anPrice, err := model.GetLatestPrice(exchange, anchorSymbol+"usdt")
		if err != nil || anPrice == 0.0 {
			if latestPrice < 0.01 {
				maxLimit = volumeLimit * 10000
			} else {
				maxLimit = volumeLimit * 5000
			}
		} else {
			maxLimit = 20000 / anPrice
		}

		size := countPoint(fmt.Sprintf("%v", volumeLimit))
		minLimit = FloatRound(minLimit, size)
		maxLimit = FloatRound(maxLimit, size)
	}

	return minLimit, maxLimit
}

func countPoint(str string) int {
	ss := strings.Split(str, ".")
	if len(ss) == 2 {
		return len(ss[1])
	}

	return 0
}

// -------------------------------------------------------------------------------------------------

func (g *gridTradeForm) toGridStrategy(grids []*grid.Grid) *model.GridStrategy {
	averageProfit, averageProfitRate := model.CalculateProfit(g.Exchange, grids)
	num := len(grids)
	intervalSize := (grids[0].Price - grids[1].Price + grids[num-2].Price - grids[num-1].Price) / 2

	gs := &model.GridStrategy{
		UID:    g.UID,
		Type:   g.Type,
		ApiKey: g.ApiKey,

		Exchange: g.Exchange,
		Symbol:   g.Symbol,

		GridIntervalType:  g.GridIntervalType,
		MinPrice:          g.MinPrice,
		MaxPrice:          g.MaxPrice,
		GridNum:           g.GridNum,
		AverageProfit:     averageProfit,
		AverageProfitRate: averageProfitRate,

		StopProfitPrice: g.StopProfitPrice,
		StopLossPrice:   g.StopLossPrice,
		BasisPrice:      g.BasisPrice,
		EntryPrice:      g.BasisPrice,

		StartupMinPrice: g.MinPrice,
		StartupMaxPrice: g.MaxPrice,
		IntervalSize:    intervalSize,

		AnchorSymbol: g.AnchorSymbol,
		TotalSum:     g.TotalSum,

		BuyCoinQuantity: g.buyCoinQuantity,
		GridBaseNO:      g.GridBasisNO,

		IsRun: true,
	}

	gs.ID = bson.ObjectIdHex(g.gridStrategyID)

	return gs
}

func (g *gridTradeForm) toGridPendingOrder(gs *model.GridStrategy, grids []*grid.Grid) *model.GridPendingOrder {
	gpo := &model.GridPendingOrder{
		GSID: gs.ID,

		Grids:       grids,
		BasisGridNO: g.GridBasisNO,
		//EachGridMoney: gs.EachGridMoney,

		Exchange: g.Exchange,
		Symbol:   g.Symbol,
		//Position: ,
		//BuyCost: ,

		// 锚定币持仓分布
		AnchorSymbol: g.AnchorSymbol,
		//AnchorSymbolPosition: ,
	}

	gpo.ID = bson.NewObjectId()

	return gpo
}

func (g *gridTradeForm) toGridTradeRecord(gs *model.GridStrategy, records []*model.TradeRecord) []*model.GridTradeRecord {
	gtrs := []*model.GridTradeRecord{}

	now := time.Now()
	isStartUpOrder := false
	for _, v := range records {

		if v.Side == "sell" {
			isStartUpOrder = true
		}

		gtr := &model.GridTradeRecord{
			GSID: gs.ID,
			GID:  v.GID,

			OrderID:       v.OrderID,
			ClientOrderID: v.ClientOrderID,
			OrderType:     v.OrderType,
			Side:          v.Side,
			Price:         v.Price,
			Quantity:      v.Quantity,
			Volume:        v.Volume,
			Unit:          v.Unit,
			//Fees : ,

			OrderState: v.OrderState,
			StateTime:  now,

			IsStartUpOrder: isStartUpOrder,
			//BuyOrderID: ,
			//GridPeerBuyOrderID  : ,

			Exchange: g.Exchange,
			Symbol:   g.Symbol,
		}

		gtr.ID = bson.NewObjectId()
		gtrs = append(gtrs, gtr)
	}

	return gtrs
}

func (g *gridTradeForm) toGridProcess(grids []*grid.Grid, basisGridNO int, gsid string) *model.GridProcess {
	ctx, cancel := context.WithCancel(context.Background())

	return &model.GridProcess{
		Gsid:            gsid,
		UID:             g.UID,
		Exchange:        g.Exchange,
		Symbol:          g.Symbol,
		Type:            g.Type,
		AccessKey:       g.ApiKey,
		SecretKey:       g.secretKey,
		AnchorSymbol:    g.AnchorSymbol,
		Grids:           grids,
		BasisGridNO:     basisGridNO,
		ExchangeAccount: g.exchangeAccount,

		Ctx:    ctx,
		Cancel: cancel,
	}
}

// -------------------------------------------------------------------------------------------------

// FloatRound 截取小数位数，默认保留2位，四舍五入
func FloatRound(f float64, points ...int) float64 {
	size := 2
	if len(points) > 0 {
		size = points[0]
	}

	format := "%." + strconv.Itoa(size) + "f"
	if size < 1 {
		format = "%." + "f"
	}

	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

// -------------------------------------------------------------------------------------------------

// 待保存的网格信息
type saveStrategyInfo struct {
	strategy     *model.GridStrategy
	pendingOrder *model.GridPendingOrder
	tradeRecords []*model.GridTradeRecord
}

type errInfo struct {
	data   interface{}
	errMsg string
}

func (s *saveStrategyInfo) insert() map[string]map[string]*errInfo {
	errMaps := map[string]map[string]*errInfo{} // 记录保存数据失败的数据

	err := s.strategy.Insert()
	if err != nil {
		errMaps[model.GridStrategyCollectionName] = map[string]*errInfo{
			"reSave": {s.strategy, err.Error()},
		}
	}

	err = s.pendingOrder.Insert()
	if err != nil {
		errMaps[model.GridPendingOrderCollectionName] = map[string]*errInfo{
			"reSave": {s.pendingOrder, err.Error()},
		}
	}

	errRecords := map[string]*errInfo{}
	for _, v := range s.tradeRecords {
		err = v.Insert()
		if err != nil {
			errRecords[v.ID.Hex()] = &errInfo{v, err.Error()}
		}
	}
	if len(errRecords) > 0 {
		errMaps[model.GridTradeRecordCollectionName] = errRecords
	}

	return errMaps
}

func (g *gridTradeForm) saveInfo(grids []*grid.Grid, records []*model.TradeRecord) string {
	strategy := g.toGridStrategy(grids)
	pendingOrder := g.toGridPendingOrder(strategy, grids)
	tradeRecords := g.toGridTradeRecord(strategy, records)

	ssi := &saveStrategyInfo{
		strategy:     strategy,
		pendingOrder: pendingOrder,
		tradeRecords: tradeRecords,
	}

	errMaps := ssi.insert()
	if len(errMaps) == 0 {
		logger.Info("save (gridStrategy, pendingOrder, tradeRecords) data success", logger.String("gsid", strategy.ID.Hex()))
		return strategy.ID.Hex()
	}

	// 如果有保存失败，重新写入，尽可能保持一致
	go func() {
		logger.Warn("insert some data failed, insert data again", logger.Any("errMsg", errMaps))
		var err error

		for i := 0; i < 10; i++ {
			for key, errInfo := range errMaps {
				switch key {
				case model.GridStrategyCollectionName:
					err = errInfo["reSave"].data.(*model.GridStrategy).Insert()
					if err != nil {
						logger.Error(fmt.Sprintf("insert data to %s error", key), logger.Err(err), logger.Any("data", errInfo["reSave"].data.(*model.GridStrategy)))
						errInfo["reSave"].errMsg = err.Error()
						continue
					}
					delete(errMaps, key)

				case model.GridPendingOrderCollectionName:
					err = errInfo["reSave"].data.(*model.GridPendingOrder).Insert()
					if err == nil {
						logger.Error(fmt.Sprintf("insert data to %s error", key), logger.Err(err), logger.Any("data", errInfo["reSave"].data.(*model.GridPendingOrder)))
						errInfo["reSave"].errMsg = err.Error()
						continue
					}
					delete(errMaps, key)

				case model.GridTradeRecordCollectionName:
					for k, v := range errInfo {
						data := v.data.(*model.GridTradeRecord)
						err = data.Insert()
						if err != nil {
							logger.Error(fmt.Sprintf("insert data to %s error", key), logger.Err(err), logger.Any("data", v.data.(*model.GridTradeRecord)))
							v.errMsg = err.Error()
							continue
						} else {
							// 写入成功，删除key
							delete(errInfo, k)
							if len(errInfo) == 0 {
								delete(errMaps, key)
							} else {
								errMaps[key] = errInfo
							}
						}
					}
				}
			}

			time.Sleep(2 * time.Second)
		}

		if len(errMaps) == 0 {
			logger.Error("insert data failed after repeated", logger.Any("data", errMaps))
		}
	}()

	return strategy.ID.Hex()
}

// -------------------------------------------------------------------------------------------------
type autoGenerateForm struct {
	Exchange      string             `json:"exchange"`      // 交易所
	Symbol        string             `json:"symbol"`        // 品种
	TotalSum      float64            `json:"totalSum"`      // 投资金额
	MinProfitRate float64            `json:"minProfitRate"` // 每个最小盈利率
	IsReverse     bool               `json:"isReverse"`     // 是否为反向网格
	esl           *model.LimitValues `json:"-"`             // 交易所品种的限制值
	latestPrice   float64            `json:"-"`             // 品种的最新价格
}

func (a *autoGenerateForm) valid1() error {
	switch "" {
	case a.Exchange:
		return errors.New("参数exchange为空")
	case a.Symbol:
		return errors.New("参数symbol为空")
	}

	key := model.GetKey(a.Exchange, a.Symbol)
	a.esl = model.GetExchangeLimitCache(key)
	if a.esl.VolumeLimit <= 0.0 {
		return fmt.Errorf("暂时不支持%s交易所的%s品种", a.Exchange, a.Symbol)
	}

	return nil
}

func (a *autoGenerateForm) valid2() error {
	if err := a.valid1(); err != nil {
		return err
	}

	var err error
	var latestPrice float64
	// 区分不同交易所
	switch a.Exchange {
	case model.ExchangeHuobi:
		latestPrice, err = huobi.GetLatestPrice(a.Symbol)
		if err != nil {
			return err
		}
	case model.ExchangeBinance:
		// 根据品种获取当前价格
		latestPrice, err = binance.GetLatestPrice(a.Symbol, env.ProxyAddr)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("暂时不支持%s交易所", a.Exchange)
	}
	a.latestPrice = latestPrice

	if a.IsReverse {
		a.TotalSum *= latestPrice // 反向网格转换为计价币
	}

	anchorCurrency := goex.GetAnchorCurrency(a.Symbol)
	minLimit, maxLimit := getMinAndMaxTotalSum(a.Exchange, anchorCurrency, latestPrice, a.esl.VolumeLimit)
	if a.TotalSum < minLimit || a.TotalSum > maxLimit {
		if a.IsReverse {
			minLimitSize := model.Float64ToStr(minLimit/a.latestPrice, a.esl.PricePrecision)
			maxLimitSize := model.Float64ToStr(maxLimit/a.latestPrice, a.esl.PricePrecision)
			return fmt.Errorf("投资金额设置范围必须在%s ~ %s", minLimitSize, maxLimitSize)
		}
		return fmt.Errorf("投资金额设置范围必须在%v ~ %v", FloatRound(minLimit, a.esl.PricePrecision), FloatRound(maxLimit, a.esl.PricePrecision))
	}

	if a.MinProfitRate <= 0.0 {
		a.MinProfitRate = 0.001
	}

	return nil
}

// 获取最小投资金额
func (a *autoGenerateForm) getMinMoney() float64 {
	defaultGridNum := 5

	if a.esl == nil {
		key := model.GetKey(a.Exchange, a.Symbol)
		a.esl = model.GetExchangeLimitCache(key)
	}

	volumeLimit := a.esl.VolumeLimit * float64(defaultGridNum)
	anchorCurrency := goex.GetAnchorCurrency(a.Symbol)
	if anchorCurrency == "usdt" {
		if volumeLimit < 50 {
			volumeLimit = 50
		}
		return FloatRound(volumeLimit + 5)
	}

	price, _ := model.GetLatestPrice(a.Exchange, anchorCurrency+"usdt")
	if price <= 0 {
		return FloatRound(a.esl.VolumeLimit*40, a.esl.GetVolumePrecision())
	}

	return FloatRound(55/price, a.esl.GetVolumePrecision())
}

func (a *autoGenerateForm) getBestGrid() (*model.GridFilter, error) {
	totalSum := a.TotalSum
	if goex.GetAnchorCurrency(a.Symbol) != "usdt" {
		totalSum *= 0.8
	} else {
		totalSum *= 0.98
	}

	cg := &model.CalculateGrid{
		Exchange:         a.Exchange,
		Symbol:           a.Symbol,
		TargetProfitRate: a.MinProfitRate,
		ParamsRange: &model.SymbolParams{
			TotalSum:    &model.ValueRange{totalSum, totalSum, 1},
			LatestPrice: &model.ValueRange{a.latestPrice, a.latestPrice, 1},
		},
		IsReverse: a.IsReverse,
	}

	limitVolume, fees := model.GetExchangeRule(a.Exchange, a.Symbol)
	gf := cg.Done(limitVolume, fees)
	if gf == nil {
		return nil, errors.New("not found")
	}

	return gf, nil
}

// -------------------------------------------------------------------------------------------------

type cancelGridStrategyForm struct {
	UID             string `json:"uid"`             // 用户id
	Exchange        string `json:"exchange"`        // 交易所
	Symbol          string `json:"symbol"`          // 品种
	Gsid            string `json:"gsid"`            // 网格策略id
	IsClosePosition bool   `json:"isClosePosition"` // 是否平仓，true:是，false：否

	apiKey          string         `json:"-"`
	exchangeAccount goex.Accounter `json:"-"` // 交易所账号
}

func (c *cancelGridStrategyForm) valid() error {
	switch "" {
	case c.UID:
		//return errors.New("field uid is empty")
		return errors.New("参数uid为空")
	case c.Exchange:
		return errors.New("参数exchange为空")
	case c.Symbol:
		return errors.New("参数symbol为空")
	case c.Gsid:
		return errors.New("参数gsid为空")
	}

	query := bson.M{"_id": bson.ObjectIdHex(c.Gsid)}
	gs, err := model.FindGridStrategy(query, bson.M{})
	if err != nil {
		return errors.New("获取策略信息失败")
	}
	if gs.IsRun == false {
		//return errors.New("grid strategy had stopped")
		return errors.New("网格策略已经在停止状态")
	}
	c.apiKey = gs.ApiKey

	return nil
}

// 初始化交易所账号
func (c *cancelGridStrategyForm) initExchangeAccount() error {
	_, _, account, err := model.InitExchangeAccount(c.UID, c.Exchange, c.apiKey)
	if err != nil {
		return err
	}
	c.exchangeAccount = account

	return nil
}

// 取消在交易所的委托
func (c *cancelGridStrategyForm) cancelCommissionOrder() ([]string, error) {
	if c.exchangeAccount == nil {
		return []string{}, errors.New("not found account")
	}

	// 获取所有委托订单
	query := bson.M{"gsid": bson.ObjectIdHex(c.Gsid), "orderState": huobi.OrderStateSubmitted}
	field := bson.M{"orderID": true, "_id": true, "symbol": true}
	gtrs, err := model.FindGridTradeRecords(query, field, 0, 100)
	if err != nil {
		return nil, err
	}

	successOrderIDs := []string{}
	failedOrderIDs := []string{}
	for _, v := range gtrs {
		err = c.exchangeAccount.CancelOrder(v.OrderID, v.Symbol)
		if err != nil {
			isCancelel := false
			// 区分不同交易所
			switch c.Exchange {
			case model.ExchangeHuobi:
				if strings.Contains(err.Error(), "order-orderstate-error") {
					isCancelel = true
				}
			case model.ExchangeBinance:
				if strings.Contains(err.Error(), "Unknown order sent") {
					isCancelel = true
				}
			}

			if isCancelel {
				query = bson.M{"_id": v.ID}
				update := bson.M{"$set": bson.M{"orderState": huobi.OrderStateFilled, "stateTime": time.Now()}}
				model.UpdateGridTradeRecord(query, update)
			} else {
				failedOrderIDs = append(failedOrderIDs, v.OrderID)
				logger.Warn("cancel order failed", logger.Err(err), logger.String("clientOrderID", v.ClientOrderID), logger.String("orderID", v.OrderID))
				continue
			}
		}

		successOrderIDs = append(successOrderIDs, v.OrderID)
		time.Sleep(time.Millisecond * 50)
	}

	if len(failedOrderIDs) > 0 {
		logger.Error("取消失败的订单", logger.String("uid", c.UID), logger.Any("failedOrderIDs", failedOrderIDs))
	}

	return successOrderIDs, nil
}

// 平仓
func (c *cancelGridStrategyForm) closePosition() error {
	needClosePosition, err := model.CalculatePosition(c.Exchange, c.Symbol, c.Gsid, c.exchangeAccount)
	if err != nil {
		logger.Warn("model.CalculatePosition error", logger.Err(err))
		return err
	}

	// 平仓
	tradeRecord, err := c.placeMarketOrder("sell", needClosePosition)
	if err != nil {
		return err
	}

	if tradeRecord.Exchange != "" {
		err = tradeRecord.Insert()
		if err != nil {
			logger.Error("tradeRecord.Insert error", logger.Err(err), logger.Any("tradeRecord", tradeRecord))
			return err
		}
	}

	return nil
}

// 买卖市价单
func (c *cancelGridStrategyForm) placeMarketOrder(side string, placeCoinQuantity float64) (*model.GridTradeRecord, error) {
	marketOrderRecord := &model.GridTradeRecord{}

	if placeCoinQuantity > 0.0 {
		if c.exchangeAccount == nil {
			return marketOrderRecord, errors.New("not found account")
		}

		el := model.GetExchangeLimitCache(model.GetKey(c.Exchange, c.Symbol))

		//side := "sell"
		clientOrderID := model.GenerateClientOrderID(c.Exchange, model.PrefixIDMos, string(krand.String(krand.R_All, 6)))
		coinQuantity := model.RoundOffToStr(placeCoinQuantity, el.QuantityPrecision)
		placeCoinQuantity = str2Float64(coinQuantity)
		latestPrice := 0.0
		feesRate := 0.0
		err := error(nil)

		// 区分不同交易所
		switch c.Exchange {
		case model.ExchangeHuobi:
			latestPrice, err = huobi.GetLatestPrice(c.Symbol)
			if err != nil {
				logger.Warn("huobi.GetLatestPrice error", logger.Err(err), logger.String("symbol", c.Symbol))
			}
			feesRate = huobi.FilledFees

		case model.ExchangeBinance:
			latestPrice, err = binance.GetLatestPrice(c.Symbol, env.ProxyAddr)
			if err != nil {
				logger.Warn("binance.GetLatestPrice error", logger.Err(err), logger.String("symbol", c.Symbol))
			}
			feesRate = binance.FilledFees
		}

		if placeCoinQuantity*latestPrice < el.VolumeLimit {
			return &model.GridTradeRecord{}, errors.New("下市价单的volume太小")
		}

		orderID, err := c.exchangeAccount.PlaceMarketOrder(side, c.Symbol, coinQuantity, clientOrderID)
		if err != nil {
			logger.Warn("PlaceMarketOrder error", logger.Err(err), logger.String("params", fmt.Sprintf("%s, %s, %v, %s", side, c.Symbol, coinQuantity, clientOrderID)))
			return marketOrderRecord, err
		}

		logger.Info("place market order success",
			logger.String("side", side),
			logger.String("symbol", c.Symbol),
			logger.Float64("price", FloatRound(latestPrice, el.PricePrecision)),
			logger.String("coinQuantity", coinQuantity),
			logger.String("clientOrderID", clientOrderID),
			logger.String("orderID", orderID),
		)

		//if volume == placeCoinQuantity {
		//	volume = FloatRound(placeCoinQuantity * latestPrice)
		//}
		//volume := FloatRound(placeCoinQuantity*latestPrice, el.GetVolumePrecision())
		volume := FloatRound(placeCoinQuantity*latestPrice, el.PricePrecision)
		anchorSymbol := goex.GetAnchorCurrency(c.Symbol)

		marketOrderRecord = &model.GridTradeRecord{
			GSID: bson.ObjectIdHex(c.Gsid),
			GID:  0,

			OrderID:       orderID,
			ClientOrderID: clientOrderID,
			OrderType:     "market",
			Side:          side,
			Price:         latestPrice,
			Quantity:      placeCoinQuantity,
			Volume:        volume,
			Unit:          anchorSymbol,
			Fees:          volume * feesRate,

			OrderState: "filled",
			StateTime:  time.Now(),

			IsStartUpOrder: false,

			Exchange: c.Exchange,
			Symbol:   c.Symbol,
		}
	} else {
		logger.Info("no need to place market order")
	}

	return marketOrderRecord, nil
}

func (c *cancelGridStrategyForm) getCurrencyPosition(totalPosition float64) float64 {
	needClosePosition := totalPosition

	currency, _ := goex.SplitSymbol(c.Symbol)
	if currency != "" {
		currencyBalance, err := c.exchangeAccount.GetCurrencyBalance(currency)
		if err != nil {
			logger.Warn("GetCurrencyBalance error", logger.Err(err), logger.String("currency", currency))
		} else {
			if currencyBalance <= totalPosition {
				needClosePosition = currencyBalance * (1 - 0.002)
			} else {
				needClosePosition = totalPosition * (1 - 0.002)
			}
		}
	}

	el := model.GetExchangeLimitCache(model.GetKey(c.Exchange, c.Symbol))

	return FloatRound(needClosePosition, el.QuantityPrecision)
}

// 更新订单状态
func updateOrderStatus(orderIDs []string) []string {
	failedOrderIDs := []string{}

	update := bson.M{
		"$set": bson.M{"orderState": huobi.OrderStateCanceled, "stateTime": time.Now()},
	}
	for _, orderID := range orderIDs {
		query := bson.M{"orderID": orderID}
		err := model.UpdateGridTradeRecord(query, update)
		if err != nil {
			failedOrderIDs = append(failedOrderIDs, orderID)
			logger.Error("updateGridTradeRecord error", logger.Err(err), logger.String("orderID", orderID), logger.Any("update", update))
			continue
		}
	}

	return failedOrderIDs
}

// 更新策略运行状态
func updateStrategyOrder(gsid string) error {
	query := bson.M{"_id": bson.ObjectIdHex(gsid)}
	update := bson.M{
		"$set": bson.M{"isRun": false},
	}
	return model.UpdateGridStrategy(query, update)
}

// 网格默认的最低价格和最高价格
func getMinMax(price float64, num int) (float64, float64) {
	k := 0.0
	if price > 5000 {
		k = 250
	} else if price > 100 {
		k = 100
	} else {
		k = 50
	}
	x := price / k
	min := price - 0.3*x*float64(num)
	max := price + 0.7*x*float64(num)

	return min, max
}

// -------------------------------------------------------------------------------------------------

type strategyOut struct {
	ID          string `json:"id"`          // 策略id
	Name        string `json:"name"`        // 网格类型名称
	Exchange    string `json:"exchange"`    // 交易所
	Symbol      string `json:"symbol"`      // 交易的品种
	StartupTime string `json:"startupTime"` // 策略启动时间
	Type        int    `json:"type"`        // 策略类型，0:网格交易，1:无限网格

	TotalSum     string `json:"totalSum"`     // 投资本金
	AnchorSymbol string `json:"anchorSymbol"` // 锚定币的品种名称，例如USDT
	TradeCount   int    `json:"tradeCount"`   // 交易次数

	TotalProfit       string `json:"totalProfit"` // 总利润
	RateReturn        string `json:"rateReturn"`
	RealizedRevenue   string `json:"realizedRevenue"`   // 已实现收益
	UnrealizedRevenue string `json:"unrealizedRevenue"` // 未实现收益
	AnnualReturn      string `json:"annualReturn"`      // 年化收益率

	IsRun bool `json:"isRun"` // 策略是否运行中
}

// 获取运行的策略
func getRunningStrategies(uid string, form *reqListForm) ([]*strategyOut, int, error) {
	gos := []*strategyOut{}
	total := 0

	query := bson.M{"uid": uid, "isRun": true}
	gss, err := model.FindGridStrategies(query, bson.M{}, form.page, form.limit, form.sort)
	if err != nil {
		return gos, total, err
	}

	total, err = model.CountGridStrategies(query)
	if err != nil {
		return gos, total, err
	}

	// 获取统计信息
	for _, gs := range gss {
		// 和交易所一致计算方式
		query = bson.M{"gsid": gs.ID, "orderState": huobi.OrderStateFilled}
		count, _ := model.CountGridTradeRecords(query)

		resp, err := getStatisticalInfo(gs.UID, gs.ID.Hex())
		if err != nil {
			logger.Error("getStatisticalInfo error", logger.Err(err), logger.String("uid", gs.UID), logger.String("apikey", gs.ApiKey))
		}

		gos = append(gos, &strategyOut{
			ID:          gs.ID.Hex(),
			Name:        model.GetStrategyTypeCache(gs.Type).Name,
			Exchange:    gs.Exchange,
			Symbol:      gs.Symbol,
			StartupTime: gs.CreatedAt.Local().Format(model.DateTimeUTC),
			Type:        gs.Type,

			TotalSum:     fmt.Sprintf("%v", gs.TotalSum),
			AnchorSymbol: gs.AnchorSymbol,
			TradeCount:   count,

			RateReturn:      resp.Data.RateReturn,
			TotalProfit:     resp.Data.RealizeProfit,
			RealizedRevenue: resp.Data.RealizeProfit,
			AnnualReturn:    resp.Data.RateReturnYear,
			IsRun:           true,
		})
	}

	return gos, total, nil
}

// 策略详情信息
type detailOut struct {
	ID          string `json:"id"`          // 策略id
	Name        string `json:"name"`        // 网格类型名称
	Exchange    string `json:"exchange"`    // 交易所
	Symbol      string `json:"symbol"`      // 交易的品种
	StartupTime string `json:"startupTime"` // 策略启动时间
	Type        int    `json:"type"`        // 策略类型，0:网格交易，1:无限网格

	TotalSum     string `json:"totalSum"`                         // 投资本金
	AnchorSymbol string `json:"anchorSymbol" bson:"anchorSymbol"` // 锚定币的品种名称，例如USDT

	TotalProfit       string `json:"totalProfit"`       // 总利润
	RealizedRevenue   string `json:"realizedRevenue"`   // 已实现收益
	UnrealizedRevenue string `json:"unrealizedRevenue"` // 未实现收益
	RateReturn        string `json:"rateReturn"`
	AnnualReturn      string `json:"annualReturn"` // 年化收益率
	IsRun             bool   `json:"isRun"`        // 策略是否运行中

	TradeCount            int        `json:"tradeCount"`            // 交易数量
	BuyCount              int        `json:"buyCount"`              // 已经成功买入数量
	SellCount             int        `json:"sellCount"`             // 已经成功卖出数量
	SymbolPositions       float64    `json:"symbolPositions"`       // 交易品种持仓数量
	AnchorSymbolPositions float64    `json:"anchorSymbolPositions"` // 锚定币持仓数量
	EvaDailyList          []EvaDaily `json:"evaDailyList"`
	CommissionBuyOrders   []*QP      `json:"commissionBuyOrders"`  // 正在委托买订单列表
	CommissionSellOrders  []*QP      `json:"commissionSellOrders"` // 正在委托卖订单列表
}

// 买入卖出数量和价格
type QP struct {
	Quantity float64 `json:"quantity"` // 数量
	Price    float64 `json:"price"`    // 价格
}

type ascendingQPs []*QP

// Len slice长度
func (a ascendingQPs) Len() int {
	return len(a)
}

// Less 判断大小(升序)
func (a ascendingQPs) Less(i, j int) bool {
	return a[i].Price < a[j].Price
}

// Swap 交换元素
func (a ascendingQPs) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type reverseQPs []*QP

// Len slice长度
func (r reverseQPs) Len() int {
	return len(r)
}

// Less 判断大小(倒序)
func (r reverseQPs) Less(i, j int) bool {
	return r[i].Price > r[j].Price
}

// Swap 交换元素
func (r reverseQPs) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type EvaDaily struct {
	Date        string `json:"date"`
	ProfitDaily string `json:"profit_daily"`
}

type statisticalInfoResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		RealizeProfit  string     `json:"realize_profit"`   // 已实现收益
		RateReturnYear string     `json:"rate_return_year"` // 年化收益率
		RateReturn     string     `json:"rate_return"`
		EvaDailyList   []EvaDaily `json:"eva_daily_list"`
	} `json:"data"`
}

func getStatisticalInfo(uid string, strategyID string) (*statisticalInfoResp, error) {
	ek := &statisticalInfoResp{}
	url := env.StatisticalInfoURL + fmt.Sprintf("/%s/%s", uid, strategyID)
	err := gohttp.GetJSON(ek, url, nil)
	if err != nil {
		return ek, err
	}
	if ek.Code != 0 {
		return ek, errors.New(ek.Msg)
	}

	return ek, nil
}

// 获取详情
func getStrategyDetail(id string) (*detailOut, error) {
	query := bson.M{"_id": bson.ObjectIdHex(id)}
	gs, err := model.FindGridStrategy(query, bson.M{})
	if err != nil {
		return nil, err
	}

	// 获取统计信息
	resp, err := getStatisticalInfo(gs.UID, gs.ID.Hex())
	if err != nil {
		logger.Error("getStatisticalInfo error", logger.Err(err), logger.String("uid", gs.UID), logger.String("apikey", gs.ApiKey))
	}

	so := &detailOut{
		ID:          gs.ID.Hex(),
		Name:        model.GetStrategyTypeCache(gs.Type).Name,
		Exchange:    gs.Exchange,
		Symbol:      gs.Symbol,
		StartupTime: gs.CreatedAt.Local().Format(model.DateTimeUTC),
		Type:        gs.Type,

		TotalSum:     fmt.Sprintf("%v", gs.TotalSum),
		AnchorSymbol: gs.AnchorSymbol,

		TotalProfit:     resp.Data.RealizeProfit,
		RealizedRevenue: resp.Data.RealizeProfit,
		RateReturn:      resp.Data.RateReturn,
		AnnualReturn:    resp.Data.RateReturnYear,
		EvaDailyList:    resp.Data.EvaDailyList,

		IsRun: gs.IsRun,
	}

	// 获取成功买入的数量
	query = bson.M{"gsid": bson.ObjectIdHex(id), "side": "buy", "orderState": huobi.OrderStateFilled}
	so.BuyCount, _ = model.CountGridTradeRecords(query)

	// 获取成功卖出的数量
	query = bson.M{"gsid": bson.ObjectIdHex(id), "side": "sell", "orderState": huobi.OrderStateFilled}
	so.SellCount, _ = model.CountGridTradeRecords(query)

	// 交易总数量
	so.TradeCount = so.BuyCount + so.SellCount

	// 获取交易品种和锚定币的持仓数量
	query = bson.M{"gsid": bson.ObjectIdHex(id), "orderState": huobi.OrderStateSubmitted}
	field := bson.M{"side": true, "quantity": true, "price": true}
	gtrs, err := model.FindGridTradeRecords(query, field, 0, 200)
	if err != nil {
		return nil, err
	}

	useMoney := 0.0
	buyOrders := []*QP{}
	sellOrders := []*QP{}
	for _, gtr := range gtrs {
		//so.SymbolPositions += gtr.Quantity

		if gtr.Side == "buy" {
			buyOrders = append(buyOrders, &QP{gtr.Quantity, gtr.Price})
			useMoney -= gtr.Price * gtr.Quantity
		} else if gtr.Side == "sell" {
			so.SymbolPositions += gtr.Quantity
			sellOrders = append(sellOrders, &QP{gtr.Quantity, gtr.Price})
		}
	}
	// 计算锚定币的持参数
	query = bson.M{"gsid": bson.ObjectIdHex(id), "orderState": huobi.OrderStateFilled}
	gtrs2, err := model.FindGridTradeRecords(query, bson.M{}, 0, 200)
	if err != nil {
		return nil, err
	}
	for _, gtr := range gtrs2 {
		if gtr.Side == "buy" {
			useMoney -= gtr.Price * gtr.Quantity
		} else if gtr.Side == "sell" {
			useMoney += gtr.Price * gtr.Quantity
		}
	}

	el := model.GetExchangeLimitCache(model.GetKey(so.Exchange, so.Symbol))
	if el.VolumeLimit <= 0.0 {
		el.QuantityPrecision = 6
	}
	so.SymbolPositions = FloatRound(so.SymbolPositions, el.QuantityPrecision)
	so.AnchorSymbolPositions = FloatRound(str2Float64(so.TotalSum)+useMoney, el.QuantityPrecision)

	so.CommissionBuyOrders = buyOrders
	so.CommissionSellOrders = sellOrders

	// 设置无限网格的买和卖最大显示数量
	if gs.Type == model.GridTypeInfinite {
		dsplNum := 25
		if len(buyOrders) > dsplNum {
			bos := reverseQPs(buyOrders)
			sort.Sort(bos)
			so.CommissionBuyOrders = bos[0:dsplNum]
		}
		if len(sellOrders) > dsplNum {
			sos := ascendingQPs(sellOrders)
			sort.Sort(sos)
			so.CommissionSellOrders = sos[0:dsplNum]
		}
	}

	return so, nil
}

// -------------------------------------------------------------------------------------------------

// 更新网格策略表单
type updateGridStrategyForm struct {
	Exchange        string  `json:"exchange"`        // 交易所
	Symbol          string  `json:"symbol"`          // 交易品种
	Gsid            string  `json:"gsid"`            // 网格策略id
	MinPrice        float64 `json:"minPrice"`        // 网格最小价格
	MaxPrice        float64 `json:"maxPrice"`        // 网格最大价格
	IsClosePosition bool    `json:"isClosePosition"` // 是否平仓

	latestPrice     float64             `json:"-"` // 最新价格
	exchangeAccount goex.Accounter      `json:"-"` // 交易所账号
	gs              *model.GridStrategy `json:"-"` // 网格策略
}

func (u *updateGridStrategyForm) valid() error {
	switch "" {
	case u.Exchange:
		return errors.New("参数exchange为空")
	case u.Symbol:
		return errors.New("参数symbol为空")
	case u.Gsid:
		return errors.New("参数gsid为空")
	}

	if !bson.IsObjectIdHex(u.Gsid) {
		return errors.New("参数gsid无效")
	}

	latestPrice, err := model.GetLatestPrice(u.Exchange, u.Symbol)
	if err != nil {
		return err
	}
	u.latestPrice = latestPrice

	if u.MinPrice <= 0.0 {
		return errors.New("网格最低价格不能小于0")
	}
	if u.MaxPrice <= 0.0 {
		return errors.New("网格最高价格不能小于0")
	}
	if u.MinPrice >= u.MaxPrice {
		return fmt.Errorf("网格参数最低价格(%v)不能超过最高价格(%v)", u.MinPrice, u.MaxPrice)
	}
	if u.latestPrice > 0.0 && (u.latestPrice < u.MinPrice || u.latestPrice > u.MaxPrice) {
		return fmt.Errorf("设置网格价格范围(%v ~ %v)不合法，设置范围必须包括当前最新价格%v在内", u.MinPrice, u.MaxPrice, u.latestPrice)
	}

	return nil
}

// 计算修改的网格是否为盈利的网格
func (u *updateGridStrategyForm) genNewGrid() ([]*grid.Grid, error) {
	query := bson.M{"_id": bson.ObjectIdHex(u.Gsid)}
	gs, err := model.FindGridStrategy(query, bson.M{})
	if err != nil {
		logger.Warn("model.FindGridStrategy", logger.Err(err), logger.String("_id", u.Gsid))
		return nil, errors.New("获取网格策略信息失败")
	}
	u.gs = gs

	if gs.Type != model.GridTypeNormal {
		return nil, errors.New("只有网格交易类型机器人才允许修改参数")
	}

	if !gs.IsRun {
		return nil, errors.New("网格策略在停止状态，禁止修改参数")
	}

	if time.Now().Local().Unix()-gs.UpdatedAt.Local().Unix() < 86400 {
		return nil, errors.New("每24小时只能修改一次网格参数")
	}

	return model.GenNewGrid(gs.Exchange, gs.Symbol, gs.TotalSum, gs.GridNum, u.MinPrice, u.MaxPrice)
}

// 初始化交易所账号
func (u *updateGridStrategyForm) initExchangeAccount() error {
	if u.gs == nil {
		u.gs, _ = model.FindGridStrategy(bson.M{"_id": bson.ObjectIdHex(u.Gsid)}, bson.M{})
	}
	_, _, account, err := model.InitExchangeAccount(u.gs.UID, u.Exchange, u.gs.ApiKey)
	if err != nil {
		return err
	}
	u.exchangeAccount = account

	return nil
}

type bigGridParamForm struct {
	Exchange        string  `json:"exchange"`   // 交易所
	Symbol          string  `json:"symbol"`     // 品种
	MinPrice        string  `json:"minPrice"`   // 最小价格
	ProfitRate      string  `json:"profitRate"` // 每格收益率
	IsAI            string  `json:"isAI"`       // 是否开启AI获取参数
	minPriceFloat   float64 `json:"-"`
	profitRateFloat float64 `json:"-"`
	isAI            bool    `json:"-"`
}

func (b *bigGridParamForm) valid() error {
	switch "" {
	case b.Exchange:
		return errors.New("参数exchange为空")
	case b.Symbol:
		return errors.New("参数symbol为空")
	case b.IsAI:
		return errors.New("参数IsAI为空")
	}

	b.IsAI = strings.ToLower(b.IsAI)
	if b.IsAI == "true" {
		b.isAI = true
	} else if b.IsAI == "false" {
		b.isAI = false
	} else {
		return errors.New("参数isAI不合法")
	}

	if !b.isAI { // 手动设置
		switch "" {
		case b.MinPrice:
			return errors.New("参数minPrice为空")
		case b.ProfitRate:
			return errors.New("参数profitRate为空")
		}

		if strings.Contains(b.ProfitRate, "%") {
			b.profitRateFloat = str2Float64(strings.Replace(b.ProfitRate, "%", "", -1)) / 100
		} else {
			b.profitRateFloat = str2Float64(b.ProfitRate)
		}
		if b.profitRateFloat < 0.001 || b.profitRateFloat > 0.05 {
			return errors.New("每格收益率必须在" + "0.1% ~ 5%范围")
		}

		b.minPriceFloat = str2Float64(b.MinPrice)
		latestPrice, err := model.GetLatestPrice(b.Exchange, b.Symbol)
		if err != nil {
			logger.Warn("model.GetLatestPrice failed", logger.Err(err), logger.String("exchange", b.Exchange), logger.String("symbol", b.Symbol))
			return errors.New("获取最新价格失败")
		}
		if b.minPriceFloat > latestPrice {
			return fmt.Errorf("最低价格(%s)不能大于最新价格(%v)", b.MinPrice, latestPrice)
		}
	}

	return nil
}

type bigGridParams struct {
	MinPrice          float64 `json:"minPrice"`          // 最小价格
	MaxPrice          float64 `json:"maxPrice"`          // 最高价格
	GSQ               float64 `json:"gsq"`               // 等比数列公比
	GridNum           int     `json:"gridNum"`           // 网格数量
	MinTotalSum       float64 `json:"minTotalSum"`       // 最小投资金额
	ProfitRate        string  `json:"profitRate"`        // 收益率
	AverageProfit     float64 `json:"averageProfit"`     // 平均收益
	AverageProfitRate float64 `json:"averageProfitRate"` // 平均收益率
	AverageInterval   float64 `json:"averageInterval"`   // 平均间隔
}

func (b *bigGridParamForm) calculateParams() (*bigGridParams, error) {
	bgp := &bigGridParams{}
	var lowerPrice, highPrice, feesRate float64
	var err error

	// 判断不同交易所
	switch b.Exchange {
	case model.ExchangeHuobi:
		feesRate = huobi.FilledFees
		lowerPrice, highPrice, err = huobi.Get3CyclePriceRange(b.Symbol, huobi.MON1)
		if err != nil {
			logger.Error("huobi.Get3CyclePriceRange error", logger.Err(err), logger.String("symbol", b.Symbol))
			return bgp, err
		}
	case model.ExchangeBinance:
		feesRate = binance.FilledFees
		lowerPrice, highPrice, err = binance.Get3CyclePriceRange(b.Symbol, goex.KLINE_PERIOD_1MONTH, env.ProxyAddr)
		if err != nil {
			logger.Error("huobi.Get3CyclePriceRange error", logger.Err(err), logger.String("symbol", b.Symbol))
			return bgp, err
		}
	default:
		return bgp, errors.New("暂时不支持交易所" + b.Exchange)
	}

	el := model.GetExchangeLimitCache(model.GetKey(b.Exchange, b.Symbol))
	if b.isAI {
		b.profitRateFloat = 0.004
		b.minPriceFloat = FloatRound(lowerPrice*0.98, el.PricePrecision)
	}

	maxPrice := highPrice * 1.02 // 最高价格
	QRate := math.Log10(maxPrice / b.minPriceFloat)

	count := 0
	defer func() {
		logger.Infof("total count = %d", count)
	}()

	for minTotalSum := 20 * el.VolumeLimit; minTotalSum <= 2000*el.VolumeLimit; minTotalSum += el.VolumeLimit {
		lastNum := 0
		for minQ := 1.2; minQ > b.profitRateFloat+1; minQ -= 0.0001 {
			num := int(QRate / math.Log10(minQ)) // 网格数量
			if num < 20 || lastNum == num || num > 180 {
				continue
			}
			if lastNum != num {
				lastNum = num
			}
			count++

			grids := grid.GenerateGS(b.minPriceFloat, minQ, minTotalSum, num, el.PricePrecision, el.QuantityPrecision)
			if grids[num-1].Price*grids[num-1].BuyQuantity <= el.VolumeLimit {
				continue
			}
			profit, profitRate := grid.CalculateProfit(grids, feesRate)
			if profitRate <= b.profitRateFloat {
				//grid.PrintFormat(grids)
				//fmt.Println(lowerPrice, c.MinPrice, minQ, grids[0].Price, num, minTotalSum, profitRate, c.profitRateFloat)

				averageInterval := (grids[0].Price - grids[1].Price + grids[num-2].Price - grids[num-1].Price) / 2
				if minTotalSum/float64(num) <= el.VolumeLimit {
					minTotalSum *= 1.02
				}

				return &bigGridParams{
					MinPrice: model.FloatRound(b.minPriceFloat, el.PricePrecision),
					MaxPrice: model.FloatRound(grids[0].Price, el.PricePrecision),
					GSQ:      FloatRound(minQ, 4),
					GridNum:  num,
					//MinTotalSum: model.FloatRound(minTotalSum, el.GetVolumePrecision()),
					MinTotalSum:       model.FloatRound(minTotalSum, el.PricePrecision),
					ProfitRate:        model.Float64ToStr(profitRate*100) + "%",
					AverageProfit:     profit,
					AverageProfitRate: profitRate,
					AverageInterval:   FloatRound(averageInterval, el.PricePrecision),
				}, nil
			}
		}
	}

	if b.isAI {
		return bgp, errors.New("计算失败，暂时不建议做无限网格")
	}
	return bgp, errors.New("设置最低价格太小，请重新调整参数")
}

// 获取网格数量
func getStrategyCount(uid string) int {
	query := bson.M{"uid": uid, "isRun": true}

	n, err := model.CountGridStrategies(query)
	if err != nil {
		logger.Error("model.CountGridStrategies error", logger.Err(err), logger.Any("query", query))
		return -1
	}

	return n
}

func strategyStartUpNotify(uid string, gsid string) error {
	params := map[string]interface{}{
		"user_id": uid,
	}

	resp := &gohttp.JSONResponse{}
	url := env.NotifyStrategyStartUpURL
	err := gohttp.PostJSON(resp, url, params)
	if err != nil {
		logger.Warn("notifyStrategyStartUp error", logger.Err(err))
		return err
	}
	if resp.Code != 0 {
		err = fmt.Errorf("code=%d, msg=%s", resp.Code, resp.Msg)
		logger.Error("notifyStrategyStartUp error", logger.Err(err), logger.String("url", url), logger.Any("params", params))
		return err
	}

	logger.Info("notifyStrategyStartUp success", logger.String("uid", uid), logger.String("gsid", gsid))

	return nil
}

type calculateGridParamForm struct {
	Exchange        string  `json:"exchange"`   // 交易所
	Symbol          string  `json:"symbol"`     // 品种
	MinPrice        string  `json:"minPrice"`   // 最小价格
	MaxPrice        string  `json:"maxPrice"`   // 最大价格
	ProfitRate      string  `json:"profitRate"` // 每格收益率
	minPriceFloat   float64 `json:"-"`
	maxPriceFloat   float64 `json:"-"`
	profitRateFloat float64 `json:"-"`
}

func (c *calculateGridParamForm) valid() error {
	switch "" {
	case c.Exchange:
		return errors.New("参数exchange为空")
	case c.Symbol:
		return errors.New("参数symbol为空")
	case c.MinPrice:
		return errors.New("参数minPrice为空")
	case c.MaxPrice:
		return errors.New("参数maxPrice为空")
	case c.ProfitRate:
		return errors.New("参数profitRate为空")
	}

	c.minPriceFloat = str2Float64(c.MinPrice)
	if c.minPriceFloat < 0 {
		return errors.New("参数minPrice不合法")
	}
	c.maxPriceFloat = str2Float64(c.MaxPrice)
	if c.maxPriceFloat < c.minPriceFloat {
		return errors.New("参数maxPrice不合法")
	}

	if strings.Contains(c.ProfitRate, "%") {
		c.profitRateFloat = str2Float64(strings.Replace(c.ProfitRate, "%", "", -1)) / 100
	} else {
		c.profitRateFloat = str2Float64(c.ProfitRate)
	}
	if c.profitRateFloat < 0.001 || c.profitRateFloat > 0.05 {
		return errors.New("每格收益率必须在" + "0.1% ~ 5%范围")
	}

	return nil
}

func (c *calculateGridParamForm) calculateParams() (*bigGridParams, []*grid.Grid, error) {
	bgp := &bigGridParams{}

	feesRate := model.GetExchangeFees(c.Exchange)
	el := model.GetExchangeLimitCache(model.GetKey(c.Exchange, c.Symbol))

	//minQ := c.profitRateFloat + 1 // 公比
	QRate := math.Log10(c.maxPriceFloat / c.minPriceFloat)
	//num := int(QRate / math.Log10(minQ))         // 网格数量
	//minTotalSum := 10 * el.VolumeLimit // 最小投入资金

	// 计算需要最低投资金额和网格数
	var grids []*grid.Grid
	for minTotalSum := 10 * el.VolumeLimit; minTotalSum <= 2000*el.VolumeLimit; minTotalSum += el.VolumeLimit {
		for minQ := 1.2; minQ > c.profitRateFloat+1; minQ -= 0.0001 {
			num := int(QRate / math.Log10(minQ)) // 网格数量
			if num < 20 {
				continue
			}
			grids = grid.GenerateGS(c.minPriceFloat, minQ, minTotalSum, num, el.PricePrecision, el.QuantityPrecision)
			if grids[num-1].Price*grids[num-1].BuyQuantity <= el.VolumeLimit {
				continue
			}
			profit, profitRate := grid.CalculateProfit(grids, feesRate)
			if profitRate <= c.profitRateFloat {
				//grid.PrintFormat(grids)
				//fmt.Println(lowerPrice, c.MinPrice, minQ, grids[0].Price, num, minTotalSum, profitRate, c.profitRateFloat)
				if num > 100 {
					return bgp, grids, errors.New("资金利用率过低，请重新调整参数")
				}
				averageInterval := (grids[0].Price - grids[1].Price + grids[num-2].Price - grids[num-1].Price) / 2
				if minTotalSum/float64(num) <= el.VolumeLimit {
					minTotalSum *= 1.02
				}

				return &bigGridParams{
					MinPrice: model.FloatRound(c.minPriceFloat, el.PricePrecision),
					MaxPrice: model.FloatRound(grids[0].Price, el.PricePrecision),
					GSQ:      FloatRound(minQ, 4),
					GridNum:  num,
					//MinTotalSum: model.FloatRound(minTotalSum, el.GetVolumePrecision()),
					MinTotalSum:       model.FloatRound(minTotalSum, el.PricePrecision),
					ProfitRate:        model.Float64ToStr(profitRate*100) + "%",
					AverageProfit:     profit,
					AverageProfitRate: profitRate,
					AverageInterval:   FloatRound(averageInterval, el.PricePrecision),
				}, grids, nil
			}
		}
	}

	return bgp, grids, errors.New("设置最低价格太小，请重新调整参数")
}
