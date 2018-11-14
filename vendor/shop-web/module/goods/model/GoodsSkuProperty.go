package model

type GoodsSkuProperty struct {
	Id int64 `json:"id"`
	SkuId int64 `json:"sku_id"`
	SkuValue string `json:"sku_value"`
	EnSkuValue string `json:"en_sku_value"`
	SkuName string `json:"sku_name"`
	EnSkuName string `json:"en_sku_name"`
}