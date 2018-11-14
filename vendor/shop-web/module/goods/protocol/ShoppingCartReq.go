package protocol

type ShoppingCartReq struct {
	SkuId      int64  `json:"skuId"`
	ProductId  int64  `json:"productId"`
	SkuValues  string `json:"skuValues"`
	GoodsCount int    `json:"goodsCount"`
}
