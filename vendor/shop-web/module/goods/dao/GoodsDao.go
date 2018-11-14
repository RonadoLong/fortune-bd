package dao

import (
	"shop-web/common/dbutil"
	"shop-web/module/goods/model"
)

func FindGoodsCount(category int) (int,error){
	count := 0
	if err := dbutil.DB.Table("product").
		Where("category_id = ? and status = 1", category).
		Count(&count).Error; err != nil{

		return count, err
	}
	return count, nil
}

func FindGoodsAllCount()(int,error){
	count := 0
	if err := dbutil.DB.Table("product").Where("status = 1").
		Count(&count).Error; err != nil{
		return count, err
	}
	return count, nil
}

func FindAllGoodsListByOffset(pageNum int, pageSize int) ([]model.Product,error){
	var goodsList []model.Product
	if err := dbutil.DB.Table("product").
		Where("status = 1 ").
		Order("`create_time` ASC").Offset(pageNum).
		Limit(pageSize).
		Find(&goodsList).Error; err != nil{
		return goodsList, err
	}
	return goodsList, nil
}

func FindGoodsListByOffset(pageNum int, pageSize int, category int) ([]model.Product,error){
	var goodsList []model.Product

	if err := dbutil.DB.Table("product").
		Where("category_id = ? and status = 1 ", category).
		Order("`create_time` ASC").Offset(pageNum).
		Limit(pageSize).
		Find(&goodsList).Error; err != nil{

		return goodsList, err
	}
	return goodsList, nil
}


func FindGoodsByGoodsId(id int64) (model.Product,error){
	var goods model.Product
	if err := dbutil.DB.Table("product").
		Where("`status` = 1 and `product_id` = ? ", id).
		Find(&goods).Error; err != nil{

		return goods, err
	}
	return goods, nil
}

func FindGoodsDetailByGoodsId(id int64) ( model.GoodsDetail, error)  {
	var goodsDetail model.GoodsDetail
	err := dbutil.DB.Table("goods_detail").
		Where("`goods_id` = ? ", id).
		Find(&goodsDetail).Error

	if err != nil {
		return goodsDetail,err
	}
	return goodsDetail, nil
}

func FindGoodsSkuByGoodsId(id int64) ([]model.GoodsSku, error)  {
	var goodsSkuList []model.GoodsSku
	err := dbutil.DB.Table("goods_sku").
		Where("`product_id` = ? ", id).
		Find(&goodsSkuList).Error
	if err != nil {
		return goodsSkuList,err
	}
	return goodsSkuList, nil
}

func FindGoodsSkuByGoodsIdAndSkuId(skuId int64) (model.GoodsSku, error)  {
	var goodsSku model.GoodsSku
	err := dbutil.DB.Table("goods_sku").
		Where(" `sku_id` = ? ", skuId).
		Find(&goodsSku).Error

	if err != nil {
		return goodsSku,err
	}
	return goodsSku, nil
}

func FindGoodsSkuBySkuId(product_id int64) (model.Product, error)  {
	var goods model.Product
	err := dbutil.DB.Table("product").
		Where(" `product_id` = ? ", product_id).
		Find(&goods).Error

	if err != nil {
		return goods,err
	}
	return goods, nil
}


func FindGoodsSkuPropertyBySkuId(id int64) ([]model.GoodsSkuProperty, error)  {
	var goodsSkuPropertys  []model.GoodsSkuProperty
	err := dbutil.DB.Table("goods_sku_property").
		Where("`sku_id` = ? ", id).
		Find(&goodsSkuPropertys).Error

	if err != nil {
		return goodsSkuPropertys,err
	}
	return goodsSkuPropertys, nil
}