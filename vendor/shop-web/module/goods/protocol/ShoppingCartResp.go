package protocol

type ShoppingCartResp struct {
	Id int64 `json:"id"`
	SkuId int64 `json:"skuId"`
	ProductId int64 `json:"productId"`
	SkuValues string `json:"skuValues"`
	GoodsCount int `json:"goodsCount"`
	GoodsImages string `json:"goodsImages"`
	Title string `json:"goodsTitle"`
	SellPoint string `json:"sellPoint"`
	StockNumber int `json:"stockNumber"`
	CheckStatus bool `json:"checkStatus"`
	Price int `json:"price"`
	MemberPrice int `json:"memberPrice"`
	Status int `json:"status"`
}
