package model

type GoodsStockFlow struct {
	Id int64
	SkuId int64
	OrderId string
	StockBefore int
	StockAfter int
	StockChange int
	CheckStatus int
	LockStockBefore int
	LockStockAfter int
	LockStockChange int
}