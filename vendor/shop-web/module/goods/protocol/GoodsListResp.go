package protocol

type GoodsListResp struct {
	ProductId     int64  `json:"productId"  gorm:"primary_key;size:20"`
	CategoryId  int    `json:"categoryId"`
	Title       string `json:"title"`
	SellPoint   string `json:"sellPoint"`
	Price       int    `json:"price"`
	MemberPrice int    `json:"memberPrice"`
	SoldCount   int    `json:"soldCount"`
	GoodsImages string `json:"goodsImages"`
	HasActivity int    `json:"hasActivity"`
	//商品类型 1: 单品 2：一种规格
	GoodsType int `json:"goodsType"`
	//0下架 1.上架 2.卖完
	Status int `json:"status"`
	//佣金
	Commission int `json:"commission"`
	//积分
	Integral int `json:"integral"`

	GoodsBanners string `json:"goodsBanners"`
	GoodsDetail  string `json:"goodsDetail"`
	GoodsDesc    string `json:"goodsDesc"`
}
