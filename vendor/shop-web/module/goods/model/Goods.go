package model

import "time"

type Goods struct {
	GoodsId     int64  `json:"goods_id"  gorm:"primary_key;size:20"`
	MerchantId  int64  `json:"merchant_id"`
	Code        string `json:"code"`
	CategoryId  int    `json:"category_id"`
	Title       string `json:"title"`
	EnTitle     string `json:"en_title"`
	SellPoint   string `json:"sell_point"`
	EnSellPoint string `json:"en_sell_point"`
	TagId       int    `json:"tag_id"`
	GoodsImages string `json:"goods_images"`
	HasActivity int    `json:"has_activity"`
	Status      int    `json:"status"`
}

type Product struct {
	ProductId   int64     `json:"product_id"`
	MerchantId  int64     `json:"merchant_id"`
	Code        string    `json:"code"`
	CategoryId  int       `json:"category_id"`
	Title       string    `json:"title"`
	EnTitle     string    `json:"en_title"`
	SellPoint   string    `json:"sell_point"`
	EnSellPoint string    `json:"en_sell_point"`
	Postage     int       `json:"postage"`
	GoodsImages string    `json:"goods_images"`
	MemberPrice int       `json:"memberPrice"`
	Price       int       `json:"price"`
	SoldCount   int       `json:"soldCount"`
	Stock       int       `json:"stock"`
	BuyMax      int       `json:"buyMax"`
	Commission  int       `json:"commission"`
	Integral    int       `json:"integral"`
	Status      int       `json:"status"`
	CreateTime  time.Time `json:"createTime"`
}
