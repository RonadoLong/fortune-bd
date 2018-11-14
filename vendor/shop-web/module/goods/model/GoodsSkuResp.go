package model

type GoodsSkuResp struct {

	SkuId        int64  `json:"skuId"  gorm:"primary_key;size:20"`
	ProductId      int64  `json:"productId"`
	//会员价
	MemberPrice  int `json:"memberPrice"`
	//'原价
	Price        int `json:"price"`
	//是否开启活动
	ActivityEnable string `json:"activityEnable"`
	//默认选中规格
	IsActive       int `json:"isActive"`
	//销量
	SoldCount    int `json:"soldCount"`
	//库存
	Stock  int `json:"stock"`
	LockStock int `json:"lockStock"`
	//0 为单品 1 为一种规格  2 为两种规格
	SkuType        int `json:"skuType"`
	//1.上架，0.下架
	Status         int `json:"status"`
	//最大购买数
	BuyNumMax int `json:"buyNumMax"`

	SkuPic string `json:"skuPic"`

	SkuPropertyRespList []GoodsSkuPropertyResp `json:"skuPropertyRespList"`
}


type GoodsSkuPropertyResp struct {
	Id int64 `json:"id"`
	SkuValue string `json:"skuValue"`
	SkuName string `json:"skuName"`
}