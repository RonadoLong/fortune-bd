package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/mongo"
)

// GridTradeRecordCollectionName 表名
const GridTradeRecordCollectionName = "gridTradeRecord"

// GridTradeRecord 网格交易记录
type GridTradeRecord struct {
	mongo.PublicFields `bson:",inline"`

	GSID bson.ObjectId `json:"gsid" bson:"gsid"` // 网格策略id
	GID  int           `json:"gid" bson:"gid"`   // 网格编号

	OrderID       string  `json:"orderID" bson:"orderID"`             // 委托订单id，用户查询、取消订单
	ClientOrderID string  `json:"clientOrderID" bson:"clientOrderID"` // 用户自定义订单id
	OrderType     string  `json:"orderType" bson:"orderType"`         // 订单类型，limit:限价单，market:市价单
	Side          string  `json:"side" bson:"side"`                   // 买入卖出，buy:买入，sell:卖出
	Price         float64 `json:"price" bson:"price"`                 // 成交价格
	Quantity      float64 `json:"quantity" bson:"quantity"`           // 买卖数量
	Volume        float64 `json:"volume" bson:"volume"`               // 成交额
	Unit          string  `json:"unit" bson:"unit"`                   // 成交额单位
	Fees          float64 `json:"fees" bson:"fees"`                   // 买入卖出手续费

	OrderState string    `json:"orderState" bson:"orderState"` // 当前订单状态 submitted:委托中, canceled:取消, filled:已成交
	StateTime  time.Time `json:"stateTime" bson:"stateTime"`   // 订单状态变化时间

	IsStartUpOrder bool   `json:"isStartUpOrder" bson:"isStartUpOrder"` // 是否刚启动网格策略时的委托卖出订单，true:是，false:否，此字段只有Side=sell才有效
	BuyPrice       string `json:"buyPrice" bson:"buyPrice"`             // 成功卖出后的记录买入价格，此字段只有Side=sell才有值
	SellOrderID    string `json:"sellOrderID" bson:"sellOrderID"`       // 成功买入后的记录委托卖订单id，此字段只有Side=buy才有值

	Exchange string `json:"exchange" bson:"exchange"` // 交易所
	Symbol   string `json:"symbol" bson:"symbol"`     // 品种名称
}

// Insert 插入一条新的记录
func (object *GridTradeRecord) Insert() (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()
	object.SetFieldsValue()
	return mconn.Insert(GridTradeRecordCollectionName, object)
}

// FindGridTradeRecord 获取单条记录
func FindGridTradeRecord(selector bson.M, field bson.M) (*GridTradeRecord, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &GridTradeRecord{}
	return object, mconn.FindOne(GridTradeRecordCollectionName, object, selector, field)
}

// FindGridTradeRecords 获取多条记录
func FindGridTradeRecords(selector bson.M, field bson.M, page int, limit int, sort ...string) ([]*GridTradeRecord, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	// 默认从第一页开始
	if page < 0 {
		page = 0
	}
	if len(sort) == 0 || sort[0] == "" {
		sort = []string{"-_id"}
	}

	objects := []*GridTradeRecord{}
	return objects, mconn.FindAll(GridTradeRecordCollectionName, &objects, selector, field, page, limit, sort...)
}

// UpdateGridTradeRecord 更新单条记录
func UpdateGridTradeRecord(selector, update bson.M) (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateOne(GridTradeRecordCollectionName, selector, mongo.UpdatedTime(update))
}

// UpdateTradeRecords 更新多条记录
func UpdateTradeRecords(selector, update bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridTradeRecordCollectionName, selector, mongo.UpdatedTime(update))
}

// FindAndModifyGridTradeRecord 更新并返回最新记录
func FindAndModifyGridTradeRecord(selector bson.M, update bson.M) (*GridTradeRecord, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &GridTradeRecord{}
	return object, mconn.FindAndModify(GridTradeRecordCollectionName, object, selector, mongo.UpdatedTime(update))
}

// CountGridTradeRecords 统计数量，不包括删除记录
func CountGridTradeRecords(selector bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.Count(GridTradeRecordCollectionName, mongo.ExcludeDeleted(selector))
}

// DeleteGridTradeRecord 删除记录
func DeleteGridTradeRecord(selector bson.M) (int, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridTradeRecordCollectionName, selector, mongo.DeletedTime(bson.M{}))
}
