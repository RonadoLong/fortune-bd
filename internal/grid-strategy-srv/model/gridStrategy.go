package model

import (
	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/mongo"
)

// GridStrategyCollectionName 表名
const GridStrategyCollectionName = "gridStrategy"

// GridStrategy 网格策略
type GridStrategy struct {
	mongo.PublicFields `bson:",inline"`

	UID  string `json:"uid" bson:"uid"`   // 用户id
	Type int    `json:"type" bson:"type"` // 0:网格交易，1:杠杆网格，2:借贷网格，3:反向网格，4:反向杠杆网格，5:无线网格

	// 交易币网格信息
	ApiKey   string `json:"apiKey" bson:"apiKey"`     // 交易所apikey
	Exchange string `json:"exchange" bson:"exchange"` // 交易所
	Symbol   string `json:"symbol" bson:"symbol"`     // 交易的品种

	GridIntervalType  string  `json:"gridIntervalType" bson:"gridIntervalType"`   // 网格间隔类型，ASGrid:等差, GSGrid:等比
	GridDesc          string  `json:"gridDesc" bson:"gridDesc"`                   // 网格说明，例如xxx震荡上涨
	MinPrice          float64 `json:"minPrice" bson:"minPrice"`                   // 网格最低价格
	MaxPrice          float64 `json:"maxPrice" bson:"maxPrice"`                   // 网格最高价格
	GridNum           int     `json:"gridNum" bson:"gridNum"`                     // 网格数量
	TotalSum          float64 `json:"totalSum" bson:"totalSum"`                   // 网格投资总额
	AverageProfit     float64 `json:"averageProfit" bson:"averageProfit"`         // 平均每次套利的利润
	AverageProfitRate float64 `json:"averageProfitRate" bson:"averageProfitRate"` // 平均每次套利的利润率(低买高卖)

	StopProfitPrice float64 `json:"stopProfitPrice" bson:"stopProfitPrice"` // 止盈价格
	StopLossPrice   float64 `json:"stopLossPrice" bson:"stopLossPrice"`     // 止损价格
	BasisPrice      float64 `json:"basisPrice" bson:"basisPrice"`           // 网格基准价格，也就是当前币的价格
	EntryPrice      float64 `json:"entryPrice" bson:"entryPrice"`           // 入场价，启动网格策略时记录的第一个基准价格，趋势网格使用
	ResetPriceCount int     `json:"resetPriceCount" bson:"resetPriceCount"` // 重新设置网格价格次数

	StartupMinPrice float64 `json:"startupMinPrice" bson:"startupMinPrice"` // 启动网格策略时最小价格
	StartupMaxPrice float64 `json:"startupMaxPrice" bson:"startupMaxPrice"` // 启动网格策略时最大价格
	IntervalSize    float64 `json:"intervalSize" bson:"intervalSize"`       // 启动网格策略时间隔，如果等比类型表示公比，等差类型表示公差

	// 锚定币信息
	AnchorSymbol string `json:"anchorSymbol" bson:"anchorSymbol"` // 锚定币的品种名称，例如USDT

	BuyCoinQuantity float64 `json:"buyCoinQuantity" bson:"buyCoinQuantity"` // 启动网格策略时以市价单买入币的数量
	GridBaseNO      int     `json:"gridBaseNO" bson:"gridBaseNO"`           // 启动网格策略时网格基准线编号

	IsRun bool `json:"isRun" bson:"isRun"` // 网格策略是否运行中
}

// Insert 插入一条新的记录
func (object *GridStrategy) Insert() (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()
	object.SetFieldsValue()
	return mconn.Insert(GridStrategyCollectionName, object)
}

// FindGridStrategy 获取单条记录
func FindGridStrategy(selector bson.M, field bson.M) (*GridStrategy, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &GridStrategy{}
	return object, mconn.FindOne(GridStrategyCollectionName, object, selector, field)
}

// FindGridStrategies 获取多条记录
func FindGridStrategies(selector bson.M, field bson.M, page int, limit int, sort ...string) ([]*GridStrategy, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	// 默认从第一页开始
	if page < 0 {
		page = 0
	}
	if len(sort) == 0 || sort[0] == "" {
		sort = []string{"-_id"}
	}

	objects := []*GridStrategy{}
	return objects, mconn.FindAll(GridStrategyCollectionName, &objects, selector, field, page, limit, sort...)
}

// UpdateGridStrategy 更新单条记录
func UpdateGridStrategy(selector, update bson.M) (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateOne(GridStrategyCollectionName, selector, mongo.UpdatedTime(update))
}

// UpdateGridStrategies 更新多条记录
func UpdateGridStrategies(selector, update bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridStrategyCollectionName, selector, mongo.UpdatedTime(update))
}

// FindAndModifyGridStrategy 更新并返回最新记录
func FindAndModifyGridStrategy(selector bson.M, update bson.M) (*GridStrategy, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &GridStrategy{}
	return object, mconn.FindAndModify(GridStrategyCollectionName, object, selector, mongo.UpdatedTime(update))
}

// CountGridStrategies 统计数量，不包括删除记录
func CountGridStrategies(selector bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.Count(GridStrategyCollectionName, mongo.ExcludeDeleted(selector))
}

// DeleteGridStrategy 删除记录
func DeleteGridStrategy(selector bson.M) (int, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridStrategyCollectionName, selector, mongo.DeletedTime(bson.M{}))
}
