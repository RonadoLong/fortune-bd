package model

type GoodsSku struct {

	SkuId        int64  `json:"skuId"  gorm:"primary_key;size:20"`

	GoodsId      int64
	//会员价
	MemberPrice  int
	//'原价
	Price        int
	//进货价格
	LowPrice     int
	//活动价格
	ActivityPrice  int
	//折扣价格
	DiscountPrice  int
	//预售价格
	PreSalePrice   int
	//是否开启活动
	ActivityEnable string
	//默认选中规格
	IsActive       int
	SoldCount    int
	//库存
	Stock   int
	LockStock int
	//佣金
	Commission     int
	//积分
	Integral       int
	//0 为单品 1 为一种规格  2 为两种规格
	SkuType        int
	//1.上架，0.下架
	Status         int
	//最大购买数
	BuyNumMax int
	SkuPic string
}