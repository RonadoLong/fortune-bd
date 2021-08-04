package v1

import (
	"fmt"
	"strings"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/app/grid-strategy-srv/model"
	"wq-fotune-backend/app/grid-strategy-srv/util/goex/binance"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/render"
)

// 同步交易所品种数量和交易额限制
func SyncExchangeLimit(c *gin.Context) {
	form := &exchangeLimitForm{}
	err := render.BindJSON(c, form)
	if err != nil {
		render.Err400Msg(c, "解析json错误, "+err.Error())
		return
	}

	err = form.valid()
	if err != nil {
		logger.Error("参数错误", logger.Err(err), logger.Any("form", form))
		render.Err400Msg(c, err.Error())
		return
	}

	symbols := form.filterSymbols()
	if len(symbols) == 0 {
		render.OK(c, fmt.Sprintf("成功添加0个交易对，一共%d交易对", len(form.Symbols)))
		return
	}

	// 区分不同交易所
	switch form.Exchange {
	case model.ExchangeHuobi:
		symbolInfos, err := huobi.GetSelectSymbols(symbols)
		if err != nil {
			logger.Error("获取交易所的品种信息失败", logger.Err(err), logger.Any("exchange", form.Exchange))
			render.Err500(c, fmt.Sprintf("获取交易所%s的品种信息失败", form.Exchange))
			return
		}

		if len(symbolInfos) == 0 {
			render.Err500(c, fmt.Sprintf("获取交易所%s的品种%v信息失败", form.Exchange, symbols))
			return
		}

		for _, v := range symbolInfos {
			el := &model.ExchangeLimit{
				Exchange:          form.Exchange,
				Symbol:            v.Symbol,
				ES:                model.GetKey(form.Exchange, v.Symbol),
				Currency:          strings.Replace(v.Symbol, form.AnchorCurrency, "", -1),
				AnchorCurrency:    form.AnchorCurrency,
				VolumeLimit:       v.MinOrderValue,
				QuantityLimit:     v.MinOrderAmt,
				PricePrecision:    v.PricePrecision,
				QuantityPrecision: v.AccountPrecision,
				LeverageRatio:     v.LeverageRatio,
			}
			if err = el.Insert(); err != nil {
				if !strings.Contains(err.Error(), "duplicate key") {
					logger.Error("insert data failed", logger.Err(err), logger.Any("exchangeLimit", el))
					continue
				}
			}
			model.SetExchangeLimitCache(el.ES, el.ToLimitValues()) // 更新缓存
			model.SetCurrencyPairCache(el.Symbol, el.Currency, el.AnchorCurrency)
		}

	case model.ExchangeBinance:
		sls, err := binance.GetSelectSymbols(symbols, env.ProxyAddr)
		if err != nil {
			logger.Error("获取交易所的品种信息失败", logger.Err(err), logger.Any("exchange", form.Exchange))
			render.Err500(c, fmt.Sprintf("获取交易所%s的品种信息失败", form.Exchange))
			return
		}

		if len(sls) == 0 {
			render.Err500(c, fmt.Sprintf("获取交易所%s的品种%v信息失败", form.Exchange, symbols))
			return
		}

		for _, sl := range sls {
			el := &model.ExchangeLimit{
				Exchange:          form.Exchange,
				Symbol:            sl.Symbol,
				ES:                model.GetKey(form.Exchange, sl.Symbol),
				Currency:          strings.Replace(sl.Symbol, form.AnchorCurrency, "", -1),
				AnchorCurrency:    form.AnchorCurrency,
				VolumeLimit:       sl.VolumeLimit,
				QuantityLimit:     sl.QuantityLimit,
				PricePrecision:    sl.PricePrecision,
				QuantityPrecision: sl.QuantityPrecision,
			}
			if err = el.Insert(); err != nil {
				if !strings.Contains(err.Error(), "duplicate key") {
					logger.Error("insert data failed", logger.Err(err), logger.Any("exchangeLimit", el))
					continue
				}
			}
			model.SetExchangeLimitCache(el.ES, el.ToLimitValues()) // 更新缓存
			model.SetCurrencyPairCache(el.Symbol, el.Currency, el.AnchorCurrency)
		}
	}

	if len(symbols) == len(form.Symbols) {
		render.OK(c, fmt.Sprintf("成功添加%d个交易对", len(symbols)))
		return
	}

	render.OK(c, fmt.Sprintf("成功添加%d个交易对，一共%d交易对, 部分品种添加失败的原因品种不是%s币本位", len(symbols), len(form.Symbols), form.AnchorCurrency))
}

// GetExchangeLimitAllCache 获取列表
func GetExchangeLimits(c *gin.Context) {
	form := &reqListForm{
		pageStr:  c.Query("page"),
		limitStr: c.Query("limit"),
		sort:     c.Query("sort"),
	}

	err := form.valid()
	if err != nil {
		logger.Error("参数错误", logger.Err(err), logger.Any("form", form))
		render.Err400Msg(c, err.Error())
		return
	}

	els, err := model.FindExchangeLimits(bson.M{}, bson.M{}, form.page, form.limit)
	if err != nil {
		logger.Error("获取数据失败", logger.Err(err), logger.String("page", form.pageStr), logger.String("limit", form.limitStr))
		render.Err500(c, err.Error())
		return
	}

	total, err := model.CountExchangeLimits(bson.M{})
	if err != nil {
		logger.Error("获取数据失败", logger.Err(err))
		render.Err500(c, err.Error())
		return
	}

	render.OK(c, gin.H{"exchangeLimits": els, "total": total})
}

// GetCurrencyPair 根据品种名称获取品种的交易对
func GetCurrencyPair(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		render.Err400Msg(c, "参数symbol为空")
		return
	}

	cp := model.GetCurrencyPairCache(symbol)
	if cp.CurrencyA.Symbol == "" || cp.CurrencyB.Symbol == "" {
		render.Err400Msg(c, "暂时不支持品种"+symbol)
		return
	}

	render.OK(c, gin.H{
		"currency":       strings.ToLower(cp.CurrencyA.Symbol),
		"anchorCurrency": strings.ToLower(cp.CurrencyB.Symbol),
	})
}

// UpdateCurrencyPairs 根据锚定币更新交易对
func UpdateCurrencyPairs(c *gin.Context) {
	anchorCurrency := c.Query("anchorCurrency")
	if anchorCurrency == "" {
		render.Err400Msg(c, "参数anchorCurrency为空")
		return
	}

	query := bson.M{"symbol": bson.M{"$regex": fmt.Sprintf(".*%s$", anchorCurrency)}}
	els, err := model.FindExchangeLimits(query, bson.M{}, 0, 10000)
	if err != nil {
		render.Err500(c, "更新失败，原因是读取数据失败")
		return
	}

	count := 0

	l := len(anchorCurrency)
	for _, el := range els {
		ls := len(el.Symbol)
		if ls > l+1 && el.Symbol[ls-l:] == anchorCurrency {
			if el.Currency != "" && el.AnchorCurrency != "" {
				continue
			}

			currency := strings.Replace(el.Symbol, anchorCurrency, "", -1)
			query := bson.M{"_id": el.ID}
			update := bson.M{
				"$set": bson.M{
					"anchorCurrency": anchorCurrency,
					"currency":       currency,
				}}
			err = model.UpdateExchangeLimit(query, update)
			if err != nil {
				fmt.Println(el.ID.Hex(), err.Error())
				continue
			}
			count++
			model.SetCurrencyPairCache(el.Symbol, currency, anchorCurrency)
		}
	}

	render.OK(c, fmt.Sprintf("成功更新%d个品种交易对字段", count))
}
