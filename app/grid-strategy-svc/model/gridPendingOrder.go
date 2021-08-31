package model

import (
	"fortune-bd/app/grid-strategy-svc/util/grid"
	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/mongo"
)

// GridPendingOrderCollectionName 表名
const GridPendingOrderCollectionName = "gridPendingOrder"

// GridPendingOrder 网格挂单详情
type GridPendingOrder struct {
	mongo.PublicFields `bson:",inline"`

	GSID bson.ObjectId `json:"gsid" bson:"gsid"` // 网格策略id

	Grids       []*grid.Grid `json:"grids" bson:"grids"`             // 网格
	BasisGridNO int          `json:"basisGridNO" bson:"basisGridNO"` // 网格基准线编号
	//EachGridMoney float64      `json:"eachGridMoney" bson:"eachGridMoney"` // 每格买卖金额

	Exchange string `json:"exchange" bson:"exchange"` // 交易所
	// 交易品种的持仓分布
	Symbol string `json:"symbol" bson:"symbol"` // 品种
	//Position float64 `json:"position" bson:"position"` // 交易品种的持仓数量
	//BuyCost  float64 `json:"buyCost" bson:"buyCost"`   // 买入成本

	// 锚定币持仓分布
	AnchorSymbol         string `json:"anchorSymbol" bson:"anchorSymbol"`                 // 锚定币的品种名称，例如USDT
	AnchorSymbolPosition string `json:"anchorSymbolPosition" bson:"anchorSymbolPosition"` // 锚定币的持仓数量
}

// Grid 网格
//type Grid struct {
//	GID          int     `json:"gid" bson:"gid"`                   // 网格id
//	Price        float64 `json:"price" bson:"price"`               // 挂单价格
//	BuyQuantity  float64 `json:"buyQuantity" bson:"buyQuantity"`   // 买入数量
//	SellQuantity float64 `json:"sellQuantity" bson:"sellQuantity"` // 卖出数量
//	OrderID      string  `json:"orderId" bson:"orderId"`           // 订单号
//	Side         string  `json:"side" bson:"side"`                 // 买卖方向，0:买入，1:卖出
//}

// Insert 插入一条新的记录
func (object *GridPendingOrder) Insert() (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()
	object.SetFieldsValue()
	return mconn.Insert(GridPendingOrderCollectionName, object)
}

// FindGridPendingOrder 获取单条记录
func FindGridPendingOrder(selector bson.M, field bson.M) (*GridPendingOrder, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &GridPendingOrder{}
	return object, mconn.FindOne(GridPendingOrderCollectionName, object, selector, field)
}

// FindGridPendingOrders 获取多条记录
func FindGridPendingOrders(selector bson.M, field bson.M, page int, limit int, sort ...string) ([]*GridPendingOrder, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	// 默认从第一页开始
	if page < 0 {
		page = 0
	}
	if len(sort) == 0 || sort[0] == "" {
		sort = []string{"-_id"}
	}

	objects := []*GridPendingOrder{}
	return objects, mconn.FindAll(GridPendingOrderCollectionName, &objects, selector, field, page, limit, sort...)
}

// UpdateGridPendingOrder 更新单条记录
func UpdateGridPendingOrder(selector, update bson.M) (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateOne(GridPendingOrderCollectionName, selector, mongo.UpdatedTime(update))
}

// UpdateGridPendingOrders 更新多条记录
func UpdateGridPendingOrders(selector, update bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridPendingOrderCollectionName, selector, mongo.UpdatedTime(update))
}

// FindAndModifyGridPendingOrder 更新并返回最新记录
func FindAndModifyGridPendingOrder(selector bson.M, update bson.M) (*GridPendingOrder, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &GridPendingOrder{}
	return object, mconn.FindAndModify(GridPendingOrderCollectionName, object, selector, mongo.UpdatedTime(update))
}

// CountGridPendingOrders 统计数量，不包括删除记录
func CountGridPendingOrders(selector bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.Count(GridPendingOrderCollectionName, mongo.ExcludeDeleted(selector))
}

// DeleteGridPendingOrder 删除记录
func DeleteGridPendingOrder(selector bson.M) (int, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridPendingOrderCollectionName, selector, mongo.DeletedTime(bson.M{}))
}
