package model

type GoodsDetail struct {
	DetailId     int64  `json:"detailId"  gorm:"primary_key;size:20"`
	GoodsId      int64  `json:"goods_id"`
	GoodsBanners string `json:"goods_banners"`
	GoodsDetail  string `json:"goods_detail"`
	GoodsDesc    string `json:"goods_desc"`
}