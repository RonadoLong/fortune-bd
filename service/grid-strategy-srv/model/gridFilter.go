package model

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"wq-fotune-backend/service/grid-strategy-srv/util/grid"

	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mongo"
)

// GridFilterCollectionName 表名
const GridFilterCollectionName = "gridFilter"

// GridFilter 交易所限制
type GridFilter struct {
	mongo.PublicFields `bson:",inline"`

	Exchange string  `json:"exchange" bson:"exchange"` // 交易所
	Symbol   string  `json:"symbol" bson:"symbol"`     // 目标品种
	Fees     float64 `json:"fees" bson:"fees"`         // 手续费

	TotalSum float64 `json:"totalSum" bson:"totalSum"` // 投入金额

	GridNum              int     `json:"gridNum" bson:"gridNum"`                           // 网格数量
	MinPrice             float64 `json:"minPrice" bson:"minPrice"`                         // 最小价格
	CurrentPrice         float64 `json:"currentPrice" bson:"currentPrice"`                 // 当前价格
	MaxPrice             float64 `json:"maxPrice" bson:"maxPrice"`                         // 最大价格
	PriceDifference      float64 `json:"priceDifference" bson:"priceDifference"`           // 最大和最小价格差
	AverageIntervalPrice float64 `json:"averageIntervalPrice" bson:"averageIntervalPrice"` // 平均每格价格间距

	AverageProfit     float64 `json:"averageProfit" bson:"averageProfit"`         // 平均每格利润
	AverageProfitRate float64 `json:"averageProfitRate" bson:"averageProfitRate"` // 平均每格利润率

	UniqueVal string `json:"uniqueVal" bson:"uniqueVal"` // 唯一值，避免重复，由交易所、品种、网格数、投入金额组成
}

// Insert 插入一条新的记录
func (object *GridFilter) Insert() (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()
	object.SetFieldsValue()

	return mconn.Insert(GridFilterCollectionName, object)
}

// FindGridFilter 获取单条记录
func FindGridFilter(selector bson.M, field bson.M) (*GridFilter, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &GridFilter{}
	return object, mconn.FindOne(GridFilterCollectionName, object, selector, field)
}

// FindGridFilters 获取多条记录
func FindGridFilters(selector bson.M, field bson.M, page int, limit int, sort ...string) ([]*GridFilter, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	// 默认从第一页开始
	if page < 0 {
		page = 0
	}
	if len(sort) == 0 || sort[0] == "" {
		sort = []string{"-_id"}
	}

	objects := []*GridFilter{}
	return objects, mconn.FindAll(GridFilterCollectionName, &objects, selector, field, page, limit, sort...)
}

// UpdateGridFilter 更新单条记录
func UpdateGridFilter(selector, update bson.M) (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateOne(GridFilterCollectionName, selector, mongo.UpdatedTime(update))
}

// UpdateGridFilters 更新多条记录
func UpdateGridFilters(selector, update bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridFilterCollectionName, selector, mongo.UpdatedTime(update))
}

// FindAndModifyGridFilter 更新并返回最新记录
func FindAndModifyGridFilter(selector bson.M, update bson.M) (*GridFilter, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &GridFilter{}
	return object, mconn.FindAndModify(GridFilterCollectionName, object, selector, mongo.UpdatedTime(update))
}

// CountGridFilters 统计数量，不包括删除记录
func CountGridFilters(selector bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.Count(GridFilterCollectionName, mongo.ExcludeDeleted(selector))
}

// DeleteGridFilter 删除记录
func DeleteGridFilter(selector bson.M) (int, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridFilterCollectionName, selector, mongo.DeletedTime(bson.M{}))
}

// -------------------------------------------------------------------------------------------------

type ValueRange struct {
	Start float64 `json:"start"` // 起始值
	End   float64 `json:"end"`   // 结束值
	Step  float64 `json:"step"`  // 步长
}

type SymbolParams struct {
	TotalSum      *ValueRange `json:"totalSum"`      // 投入资金参数
	LatestPrice   *ValueRange `json:"latestPrice"`   // 价格参数
	IntervalPrice *ValueRange `json:"intervalPrice"` // 价格间隔参数
}

type CalculateGrid struct {
	Exchange             string        `json:"exchange"`             // 交易所
	Symbol               string        `json:"symbol"`               // 品种
	TargetProfitRate     float64       `json:"targetProfitRate"`     // 最小目标收益率
	PriceDifferenceLimit float64       `json:"priceDifferenceLimit"` // 最小价格差限制
	ParamsRange          *SymbolParams `json:"paramsRange"`          // 品种参数范围
	IsReverse            bool          `json:"isReverse"`            // 是否为反向
}

func (c *CalculateGrid) Done(limitVolume float64, fees float64) *GridFilter {
	exchange := c.Exchange
	symbol := c.Symbol
	targetProfitRate := 0.0005
	if c.TargetProfitRate > targetProfitRate {
		targetProfitRate = c.TargetProfitRate
	}

	//feesMap := map[string]float64{"huobi": 0.002, "binance": 0.001}
	//limitVolumeMap := map[string]float64{"huobi": 5, "binance": 10}
	//fees := feesMap[exchange]

	minPrice, maxPrice, priceDifference := 0.0, 0.0, 0.0
	logger.Info("CalculateGrid params", logger.Any("CalculateGrid", c))

	for latestPrice := c.ParamsRange.LatestPrice.Start; latestPrice <= c.ParamsRange.LatestPrice.End; latestPrice += c.ParamsRange.LatestPrice.Step { // 价格轮询，范围和步长根据品种来设定
		c.setPriceDifferenceLimit(latestPrice)                                                                                          // 设置最小间距
		for totalSum := c.ParamsRange.TotalSum.Start; totalSum <= c.ParamsRange.TotalSum.End; totalSum += c.ParamsRange.TotalSum.Step { // 投资金额轮询，范围50~20000，步长为10
			gfs := gridFilters{}
			maxGridNum := 100
			if limitVolume > 0.0 {
				maxGridNum := int(totalSum / limitVolume)
				if maxGridNum > 100 {
					maxGridNum = 100
				}
			}
			for gridNum := 5; gridNum <= maxGridNum; gridNum += 1 { // 网格数量轮询，范围5~100，步长为1
				c.setIntervalPrice(latestPrice)                                                                                                                               // 设置网格间隔范围
				for intervalPrice := c.ParamsRange.IntervalPrice.Start; intervalPrice <= c.ParamsRange.IntervalPrice.End; intervalPrice += c.ParamsRange.IntervalPrice.Step { // 网格价格间隔轮询，范围和步长根据品种设定

					minPrice, maxPrice = grid.GetMinMax(latestPrice, intervalPrice, gridNum, c.IsReverse)
					priceDifference = maxPrice - minPrice
					if priceDifference < c.PriceDifferenceLimit { // 判断最大最小价格差是否低于限制值
						continue
					}

					grids, _ := grid.Generate(grid.GSGrid, minPrice, maxPrice, totalSum, gridNum, 8, 6)
					//if err := grid.IsValidGrids(grids, limitVolumeMap[exchange]); err != nil {
					if err := grid.IsValidGrids(grids, limitVolume); err != nil {
						continue
					}

					averageProfit, averageProfitRate := grid.CalculateProfit(grids, fees)
					if averageProfit < 0.0 {
						continue
					}

					if averageProfitRate > targetProfitRate { // 目标网格的平均利润率
						uniqueVal := fmt.Sprintf("%s.%v.%v", exchange, totalSum, latestPrice)

						gf := &GridFilter{
							Exchange: exchange,
							Symbol:   symbol,
							Fees:     fees,

							TotalSum: totalSum,

							GridNum:              gridNum,
							MinPrice:             minPrice,
							CurrentPrice:         latestPrice,
							MaxPrice:             maxPrice,
							PriceDifference:      maxPrice - minPrice,
							AverageIntervalPrice: intervalPrice,

							AverageProfit:     averageProfit,
							AverageProfitRate: averageProfitRate,

							UniqueVal: uniqueVal,
						}
						gfs = append(gfs, gf)
						//fmt.Println(gf.TotalSum, gf.GridNum, gf.AverageIntervalPrice, gf.MinPrice, gf.CurrentPrice, gf.MaxPrice, gf.PriceDifference, fmt.Sprintf("%v%%", gf.AverageProfitRate*100), gf.AverageProfit)
						//data <- gf
						break // 找到就下一个循环
					}
				}
			}

			// 在相同的投入资金和最新价格上，筛选间隔最小的网格参数
			if len(gfs) > 0 {
				return getBestGrid(gfs)
			}
		}
	}

	return nil
}

func (c *CalculateGrid) Done2() *GridFilter {
	exchange := c.Exchange
	symbol := c.Symbol
	targetProfitRate := 0.0005
	if c.TargetProfitRate > targetProfitRate {
		targetProfitRate = c.TargetProfitRate
	}

	feesMap := map[string]float64{"huobi": 0.002, "binance": 0.001}
	limitVolumeMap := map[string]float64{"huobi": 5, "binance": 10}
	fees := feesMap[exchange]

	minPrice, maxPrice, priceDifference := 0.0, 0.0, 0.0
	logger.Info("CalculateGrid params", logger.Any("CalculateGrid", c))

	for latestPrice := c.ParamsRange.LatestPrice.Start; latestPrice <= c.ParamsRange.LatestPrice.End; latestPrice += c.ParamsRange.LatestPrice.Step { // 价格轮询，范围和步长根据品种来设定
		c.setPriceDifferenceLimit(latestPrice)                                                                                          // 设置最小间距
		for totalSum := c.ParamsRange.TotalSum.Start; totalSum <= c.ParamsRange.TotalSum.End; totalSum += c.ParamsRange.TotalSum.Step { // 投资金额轮询，范围50~20000，步长为10
			gfs := gridFilters{}
			for gridNum := 100; gridNum <= 200; gridNum += 1 { // 网格数量轮询，范围5~100，步长为1
				c.setIntervalPrice(latestPrice)                                                                                                                               // 设置网格间隔范围
				for intervalPrice := c.ParamsRange.IntervalPrice.Start; intervalPrice <= c.ParamsRange.IntervalPrice.End; intervalPrice += c.ParamsRange.IntervalPrice.Step { // 网格价格间隔轮询，范围和步长根据品种设定

					minPrice, maxPrice = grid.GetMinMax(latestPrice, intervalPrice, gridNum, c.IsReverse)
					priceDifference = maxPrice - minPrice
					if priceDifference < c.PriceDifferenceLimit { // 判断最大最小价格差是否低于限制值
						continue
					}

					grids, _ := grid.Generate(grid.GSGrid, minPrice, maxPrice, totalSum, gridNum, 8, 6)
					if err := grid.IsValidGrids(grids, limitVolumeMap[exchange]); err != nil {
						continue
					}

					averageProfit, averageProfitRate := grid.CalculateProfit(grids, fees)
					if averageProfit < 0.0 {
						continue
					}

					if averageProfitRate > targetProfitRate { // 目标网格的平均利润率
						uniqueVal := fmt.Sprintf("%s.%v.%v", exchange, totalSum, latestPrice)

						gf := &GridFilter{
							Exchange: exchange,
							Symbol:   symbol,
							Fees:     fees,

							TotalSum: totalSum,

							GridNum:              gridNum,
							MinPrice:             minPrice,
							CurrentPrice:         latestPrice,
							MaxPrice:             maxPrice,
							PriceDifference:      maxPrice - minPrice,
							AverageIntervalPrice: intervalPrice,

							AverageProfit:     averageProfit,
							AverageProfitRate: averageProfitRate,

							UniqueVal: uniqueVal,
						}
						gfs = append(gfs, gf)
						//fmt.Println(gf.TotalSum, gf.GridNum, gf.AverageIntervalPrice, gf.MinPrice, gf.CurrentPrice, gf.MaxPrice, gf.PriceDifference, fmt.Sprintf("%v%%", gf.AverageProfitRate*100), gf.AverageProfit)
						//data <- gf
						break // 找到就下一个循环
					}
				}
			}

			// 在相同的投入资金和最新价格上，筛选间隔最小的网格参数
			if len(gfs) > 0 {
				return getBestGrid(gfs)
			}
		}
	}

	return nil
}

func (c *CalculateGrid) DoneAndSave(data chan *GridFilter) {
	exchange := c.Exchange
	symbol := c.Symbol
	targetProfitRate := 0.0005
	if c.TargetProfitRate > targetProfitRate {
		targetProfitRate = c.TargetProfitRate
	}

	feesMap := map[string]float64{"huobi": 0.002, "binance": 0.001}
	limitVolumeMap := map[string]float64{"huobi": 5, "binance": 10}
	fees := feesMap[exchange]

	count := 0
	minPrice, maxPrice, priceDifference := 0.0, 0.0, 0.0

	for latestPrice := c.ParamsRange.LatestPrice.Start; latestPrice <= c.ParamsRange.LatestPrice.End; latestPrice += c.ParamsRange.LatestPrice.Step { // 价格轮询，范围和步长根据品种来设定
		c.setPriceDifferenceLimit(latestPrice)                                                                                          // 设置最小间距
		for totalSum := c.ParamsRange.TotalSum.Start; totalSum <= c.ParamsRange.TotalSum.End; totalSum += c.ParamsRange.TotalSum.Step { // 投资金额轮询，范围50~20000，步长为10

			gfs := gridFilters{}

			for gridNum := 5; gridNum <= 60; gridNum += 1 { // 网格数量轮询，范围5~100，步长为1
				//gridNum := getGridNum(c.Symbol, totalSum) // 通过枚举获取

				for intervalPrice := c.ParamsRange.IntervalPrice.Start; intervalPrice <= c.ParamsRange.IntervalPrice.End; intervalPrice += c.ParamsRange.IntervalPrice.Step { // 网格价格间隔轮询，范围和步长根据品种设定
					count++
					if count%100000 == 0 {
						logger.Infof("count=%d  totalSum=%v, gridNum=%d, intervalPrice=%v, latestPrice=%v", count, totalSum, gridNum, intervalPrice, latestPrice)
					}

					minPrice, maxPrice = grid.GetMinMax(latestPrice, intervalPrice, gridNum, c.IsReverse)
					priceDifference = maxPrice - minPrice
					if priceDifference < c.PriceDifferenceLimit { // 判断最大最小价格差是否低于限制值
						continue
					}

					grids, _ := grid.Generate(grid.GSGrid, minPrice, maxPrice, totalSum, gridNum, 8, 6)
					if err := grid.IsValidGrids(grids, limitVolumeMap[exchange]); err != nil {
						continue
					}

					averageProfit, averageProfitRate := grid.CalculateProfit(grids, fees)
					if averageProfit < 0.0 {
						continue
					}

					if averageProfitRate > targetProfitRate { // 目标网格的平均利润率
						uniqueVal := fmt.Sprintf("%s.%v.%v", exchange, totalSum, latestPrice)

						gf := &GridFilter{
							Exchange: exchange,
							Symbol:   symbol,
							Fees:     fees,

							TotalSum: totalSum,

							GridNum:              gridNum,
							MinPrice:             minPrice,
							CurrentPrice:         latestPrice,
							MaxPrice:             maxPrice,
							PriceDifference:      maxPrice - minPrice,
							AverageIntervalPrice: intervalPrice,

							AverageProfit:     averageProfit,
							AverageProfitRate: averageProfitRate,

							UniqueVal: uniqueVal,
						}
						gfs = append(gfs, gf)
						//fmt.Println(gf.TotalSum, gf.GridNum, gf.AverageIntervalPrice, gf.MinPrice, gf.CurrentPrice, gf.MaxPrice, gf.PriceDifference, fmt.Sprintf("%v%%", gf.AverageProfitRate*100), gf.AverageProfit)
						//data <- gf
						break // 找到就下一个循环
					}
				}
			}

			// 在相同的投入资金和最新价格上，筛选间隔最小的网格参数
			if len(gfs) > 0 {
				data <- getBestGrid(gfs)
			}
		}
	}

	logger.Infof("calculate count %d", count)
}

func (c *CalculateGrid) setPriceDifferenceLimit(latestPrice float64) {
	if latestPrice >= 200 {
		c.PriceDifferenceLimit = 30
	} else {
		c.PriceDifferenceLimit = latestPrice * 0.11
	}
}

func (c *CalculateGrid) setIntervalPrice(latestPrice float64) {
	if c.ParamsRange.IntervalPrice == nil {
		c.ParamsRange.IntervalPrice = &ValueRange{}
	}
	c.ParamsRange.IntervalPrice.Start = latestPrice * 0.001
	c.ParamsRange.IntervalPrice.End = latestPrice * 0.1
	c.ParamsRange.IntervalPrice.Step = latestPrice * 0.001
}

// 从mongodb读取数据写入csv
func mgo2csv(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	wr := csv.NewWriter(f)
	wr.Write([]string{"exchange", "symbol", "totalSum", "gridNum", "minPrice", "currentPrice", "maxPrice", "priceDifference", "averageIntervalPrice", "averageProfit", "averageProfitRate"})
	wr.Flush()

	count, _ := CountGridFilters(bson.M{})
	limit := 100
	page := count / 100
	fmt.Println(count, page, limit)
	for i := 0; i <= page; i++ {
		gfs, err := FindGridFilters(bson.M{}, bson.M{}, i, limit, "-totalSum")
		if err != nil {
			continue
		}

		for _, v := range gfs {
			wr.Write([]string{
				v.Exchange,
				v.Symbol,
				toString(v.TotalSum),
				toString(v.GridNum),
				toString(v.MinPrice),
				toString(v.CurrentPrice),
				toString(v.MaxPrice),
				toString(v.PriceDifference),
				toString(v.AverageIntervalPrice),
				toString(v.AverageProfit),
				toString(v.AverageProfitRate),
			})
		}
		wr.Flush()
	}

	return nil
}

func toString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// Save2Mgo 保存到数据库
func Save2Mgo(data chan *GridFilter) {
	for {
		select {
		case v, canUse := <-data:
			if !canUse {
				return
			}
			v.Insert()
		}
	}
}

// Save2CSV 通过管道保存到csv文件
func Save2CSV(file string, data chan *GridFilter) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
		return
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	wr := csv.NewWriter(f)
	wr.Write([]string{"totalSum", "gridNum", "averageIntervalPrice", "priceDifference", "minPrice", "currentPrice", "maxPrice", "averageProfit", "averageProfitRate"})
	wr.Flush()

	count := 0
	for {
		select {
		case v, canUse := <-data:
			if !canUse {
				wr.Flush()
				return
			}

			wr.Write([]string{
				toString(v.TotalSum),
				toString(v.GridNum),
				toString(v.AverageIntervalPrice),
				toString(v.PriceDifference),
				toString(v.MinPrice),
				toString(v.CurrentPrice),
				toString(v.MaxPrice),
				toString(v.AverageProfit),
				toString(v.AverageProfitRate*100) + "%",
			})

			count++
			if count%100 == 0 {
				wr.Flush()
			}
		}
	}
}

// -------------------------------------------------------------------------------------------------

type gridFilters []*GridFilter

func (g gridFilters) Len() int { // 重写 Len() 方法
	return len(g)
}
func (g gridFilters) Swap(i, j int) { // 重写 Swap() 方法
	g[i], g[j] = g[j], g[i]
}
func (g gridFilters) Less(i, j int) bool { // 重写 Less() 方法
	return g[j].AverageIntervalPrice > g[i].AverageIntervalPrice
}

// 从一组数据中获取最好的网格
func getBestGrid(gfs gridFilters) *GridFilter {
	//if len(gfs) == 0 {
	//	return nil
	//}

	sort.Sort(gfs)

	refVal := 0.0
	stopKey := 0
	for k, v := range gfs {
		if k == 0 {
			refVal = v.AverageIntervalPrice
			stopKey = k
			continue
		}

		if v.AverageIntervalPrice == refVal {
			stopKey = k
			continue
		} else {
			break
		}
	}

	if stopKey > 0 {
		stopKey = stopKey / 2
	}

	return gfs[stopKey]
}

// -------------------------------------------------------------------------------------------------

// GetDefaultGridNum 获取默认的网格数
func GetDefaultGridNum(totalSum float64) int {
	return btcusdtMatch(totalSum)
}

func getGridNum(symbol string, totalSum float64) int {
	switch symbol {
	case "btcusdt":
		return btcusdtMatch(totalSum)
	case "ethusdt":
		return ethusdtMatch(totalSum)
	}

	return 0
}

func ethusdtMatch(totalSum float64) int {
	switch true {
	case totalSum >= 6000:
		return int(totalSum) / 200
	case totalSum >= 300:
		return 30
	case totalSum >= 200:
		return int(totalSum) / 10
	case totalSum >= 100:
		return int(totalSum) / 10
	case totalSum >= 90:
		return 9
	case totalSum >= 80:
		return 8
	case totalSum >= 70:
		return 7
	case totalSum >= 60:
		return 6
	case totalSum > 50:
		return 5
	}

	return 0
}

func btcusdtMatch(totalSum float64) int {
	switch true {
	case totalSum >= 7000:
		return int(totalSum) / 200
		//return 35
	case totalSum >= 390:
		return 30
	case totalSum >= 380:
		return 29
	case totalSum >= 370:
		return 28
	case totalSum >= 360:
		return 27
	case totalSum >= 340:
		return 26
	case totalSum >= 330:
		return 25
	case totalSum >= 320:
		return 24
	case totalSum >= 300:
		return 23
	case totalSum >= 290:
		return 22
	case totalSum >= 280:
		return 21
	case totalSum >= 260:
		return 20
	case totalSum >= 250:
		return 19
	case totalSum >= 240:
		return 18
	case totalSum >= 230:
		return 17
	case totalSum >= 210:
		return 16
	case totalSum >= 200:
		return 15
	case totalSum >= 190:
		return 14
	case totalSum >= 170:
		return 13
	case totalSum >= 160:
		return 12
	case totalSum >= 150:
		return 11
	case totalSum >= 130:
		return 10
	case totalSum >= 120:
		return 9
	case totalSum >= 110:
		return 8
	case totalSum >= 100:
		return 7
	case totalSum >= 90:
		return 6
	case totalSum > 50:
		return 5
	}

	return 0
}
