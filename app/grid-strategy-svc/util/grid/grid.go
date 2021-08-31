package grid

import (
	"fmt"
	"math"
)

const (
	// ASGrid 等差分配网格
	ASGrid = "ASGrid"
	// GSGrid 等比分配网格
	GSGrid = "GSGrid"
)

// Grid 网格
type Grid struct {
	GID          int     `json:"gid" bson:"gid"`                   // 网格id
	Price        float64 `json:"price" bson:"price"`               // 挂单价格
	BuyQuantity  float64 `json:"buyQuantity" bson:"buyQuantity"`   // 买入数量
	SellQuantity float64 `json:"sellQuantity" bson:"sellQuantity"` // 卖出数量
	OrderID      string  `json:"orderId" bson:"orderId"`           // 订单号
	Side         string  `json:"side" bson:"side"`                 // 买卖方向，0:买入，1:卖出
}

// Generate 生成网格
func Generate(intervalType string, minPrice, maxPrice, totalSum float64, gridNum, pricePrecision, quantityPrecision int) ([]*Grid, error) {
	grids := []*Grid{}

	switch intervalType {
	case ASGrid:
		grids = arithmeticGrid(minPrice, maxPrice, totalSum, gridNum, pricePrecision, quantityPrecision)
	case GSGrid:
		grids = geometricSequenceGrid(minPrice, maxPrice, totalSum, gridNum, pricePrecision, quantityPrecision)
	default:
		return grids, fmt.Errorf("unknown grid type %s", intervalType)
	}

	return grids, nil
}

// Generate 生成等比序列的网格
func GenerateGS(minPrice, q, totalSum float64, gridNum, pricePrecision, quantityPrecision int) []*Grid {
	return geometricSequenceGrid2(minPrice, q, totalSum, gridNum, pricePrecision, quantityPrecision)
}

// 等差分配网格
func arithmeticGrid(minPrice, maxPrice, totalInvestment float64, gridNum int, pricePrecision int, quantityPrecision int) []*Grid {
	interval := totalInvestment / float64(gridNum) // 每个格子买卖资金
	d := (maxPrice - minPrice) / float64(gridNum)  // 公差

	grids := make([]*Grid, gridNum+1, gridNum+1) // 一共gridNum+1条线
	grids[0] = &Grid{                            // 第一个格子最大边缘
		GID:   0,
		Price: FloatRound(maxPrice, pricePrecision),
	}

	for i := 1; i <= gridNum; i++ {
		currentPrice := grids[i-1].Price - d
		buyQuantity := interval / currentPrice

		grids[i] = &Grid{
			GID:         i,
			Price:       FloatRound(currentPrice, pricePrecision),
			BuyQuantity: FloatRound(buyQuantity, quantityPrecision),
		}
		grids[i-1].SellQuantity = FloatRound(buyQuantity, quantityPrecision)
		//log.Printf("网格价格: %+v \n", grids[i])
	}

	return grids
}

// 等比分配网格
func geometricSequenceGrid(minPrice, maxPrice, totalSum float64, gridNum int, pricePrecision int, quantityPrecision int) []*Grid {
	interval := totalSum / float64(gridNum)              // 每个格子买卖金额
	d := math.Pow(maxPrice/minPrice, 1/float64(gridNum)) // 等比公差

	grids := make([]*Grid, gridNum+1, gridNum+1) // 一共gridNum+1条线
	grids[0] = &Grid{                            // 第一个格子最大边缘
		GID:   0,
		Price: FloatRound(maxPrice, pricePrecision),
	}

	for i := 1; i <= gridNum; i++ {
		currentPrice := grids[i-1].Price / d
		buyQuantity := interval / currentPrice

		grids[i] = &Grid{
			GID:         i,
			Price:       FloatRound(currentPrice, pricePrecision),
			BuyQuantity: FloatRound(buyQuantity, quantityPrecision),
		}
		grids[i-1].SellQuantity = FloatRound(buyQuantity, quantityPrecision)
	}

	return grids
}

// 倒叙等比分配网格2
func geometricSequenceGrid2(minPrice, q, totalSum float64, gridNum int, pricePrecision int, quantityPrecision int) []*Grid {
	interval := totalSum / float64(gridNum) // 每个格子买卖金额

	grids := make([]*Grid, gridNum+1, gridNum+1) // 一共gridNum+1条线
	grids[gridNum] = &Grid{                      // 第一个格子最大边缘
		GID:         gridNum,
		Price:       minPrice,
		BuyQuantity: FloatRound(interval/minPrice, quantityPrecision),
	}

	for i := gridNum - 1; i >= 0; i-- {
		currentPrice := grids[i+1].Price * q
		buyQuantity := interval / currentPrice

		grids[i] = &Grid{
			GID:          i,
			Price:        FloatRound(currentPrice, pricePrecision),
			BuyQuantity:  FloatRound(buyQuantity, quantityPrecision),
			SellQuantity: FloatRound(grids[i+1].BuyQuantity, quantityPrecision),
		}
		if i == 0 {
			grids[0].BuyQuantity = 0
		}
	}
	//PrintFormat(grids)
	return grids
}

// 顺序等比分配网格3
func geometricSequenceGrid3(minPrice, q, totalSum float64, gridNum int, pricePrecision int, quantityPrecision int) []*Grid {
	interval := totalSum / float64(gridNum) // 每个格子买卖金额

	grids := make([]*Grid, gridNum+1, gridNum+1) // 一共gridNum+1条线
	grids[0] = &Grid{                            // 第一个格子最大边缘
		GID:          0,
		Price:        minPrice,
		BuyQuantity:  FloatRound(interval/minPrice, quantityPrecision),
		SellQuantity: 0.0,
	}

	for i := 1; i <= gridNum; i++ {
		currentPrice := grids[i-1].Price * q
		buyQuantity := interval / currentPrice

		grids[i] = &Grid{
			GID:          i,
			Price:        FloatRound(currentPrice, pricePrecision),
			BuyQuantity:  FloatRound(buyQuantity, quantityPrecision),
			SellQuantity: FloatRound(grids[i-1].BuyQuantity, quantityPrecision),
		}
		if i == gridNum {
			grids[gridNum].BuyQuantity = 0
		}
	}
	//PrintFormat(grids)
	return grids
}

// FloatRound 截取小数位数，默认保留2位，四舍五入
func FloatRound(f float64, points ...int) float64 {
	size := 2
	if len(points) > 0 {
		size = points[0]
	}

	// 用到fmt.Sprintf，效率比较低
	//format := "%." + strconv.Itoa(size) + "f"
	//res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	//return res

	// 效率高
	shift := math.Pow(10, float64(size))
	fv := 0.00000000001 + f // 对浮点数产生.xxx999999999 计算不准进行处理
	return math.Floor(fv*shift+.4  ) / shift
}

// PrintFormat 格式化打印网格
func PrintFormat(grids []*Grid) {
	fmt.Printf("%15s %15s %15s %15s\n", "ID", "LatestPrice", "BuyQuantity", "SellQuantity")
	for _, v := range grids {
		fmt.Printf("%15v %15v %15v %15v\n", v.GID, v.Price, v.BuyQuantity, v.SellQuantity)
	}

	fmt.Println(grids[len(grids)-1].Price, "~", grids[0].Price)
}

// IsValidGrids 检查买卖单是否满足交易所最小订单量限制
func IsValidGrids(grids []*Grid, limitVolume float64) error {
	grid := grids[len(grids)-1]
	minimumOrderVolume := grid.BuyQuantity * grid.Price
	if minimumOrderVolume < limitVolume {
		return fmt.Errorf("minimum order volume=%v is less than %v\n", FloatRound(minimumOrderVolume, 8), limitVolume)
	}
	return nil
}

// CalculateProfit 计算网格平均每格利润
func CalculateProfit(grids []*Grid, feesRate float64) (float64, float64) {
	ids := []int{1, len(grids) - 1} // 计算一头一尾

	profits := []float64{}
	volumes := []float64{}
	for _, gridNum := range ids {
		grid := grids[gridNum]
		buyPrice := grid.Price
		buyFees := grid.BuyQuantity * grid.Price * feesRate

		grid = grids[gridNum-1]
		sellPrice := grid.Price
		sellFees := grid.SellQuantity * grid.Price * feesRate

		arbitrage := (sellPrice - buyPrice) * grid.SellQuantity
		profits = append(profits, arbitrage-buyFees-sellFees)
		volumes = append(volumes, grid.BuyQuantity*grid.Price+grid.SellQuantity*grid.Price)
	}

	averageProfit := 0.0
	if len(profits) == 2 {
		averageProfit = (profits[0] + profits[1]) / 2
	}

	averageVolume := 0.0
	if len(volumes) == 2 {
		averageVolume = (volumes[0] + volumes[1]) / 2
	}

	return FloatRound(averageProfit, 6), FloatRound(averageProfit/averageVolume, 4)
}

// 网格默认的最低价格和最高价格
func getMinMax(price float64, num int) (float64, float64) {
	k := 0.0
	if price > 15000 {
		k = 100
	} else if price > 14000 {
		k = 100
	} else if price > 13000 {
		k = 100
	} else if price > 12000 {
		k = 100
	} else if price > 11000 {
		k = 100
	} else if price > 10000 {
		k = 100
	} else if price > 5000 {
		k = 250
	} else if price > 100 {
		k = 100
	} else {
		k = 50
	}
	x := price / k
	min := price - 0.4*x*float64(num)
	max := price + 0.6*x*float64(num)

	return min, max
}

// GetMinMax 根据网格间隔和数量计算网格最小和最大价格
func GetMinMax(price float64, averageIntervalPrice float64, num int, isReverse bool) (float64, float64) {
	var min, max float64

	if !isReverse {
		min = price - 0.4*float64(num)*averageIntervalPrice
		max = price + 0.6*float64(num)*averageIntervalPrice
	} else {
		min = price - 0.7*float64(num)*averageIntervalPrice
		max = price + 0.3*float64(num)*averageIntervalPrice
	}

	return min, max
}
