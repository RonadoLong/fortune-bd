package model

type ShoppingCart struct {
	Id          int64 `json:"id"`
	UserId      string
	ProductId   int64
	CheckStatus int
	SkuValues   string
	GoodsCount  int
	Status      int
}
