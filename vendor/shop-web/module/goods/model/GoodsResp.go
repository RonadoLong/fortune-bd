package model

type GoodsResp struct {
	GoodsId int64 `json:"goodsId"  gorm:"primary_key;size:20"`
	CategoryId int `json:"categoryId"`
	Title string `json:"title"`
	SellPoint string `json:"sellPoint"`
	Price int `json:"price"`
	MemberPrice  int `json:"memberPrice"`
	SoldCount int `json:"soldCount"`
	GoodsImages string `json:"goodsImages"`
	HasActivity int `json:"hasActivity"`
	//商品类型 1: 单品 2：一种规格
	GoodsType int `json:"goodsType"`
	//0下架 1.上架 2.卖完
	Status int `json:"status"`
	GoodsBanners string `json:"goodsBanners"`
	GoodsDetail  string `json:"goodsDetail"`
	GoodsDesc    string `json:"goodsDesc"`

	GoodsSkuList []GoodsSkuResp `json:"goodsSkuList"`
}