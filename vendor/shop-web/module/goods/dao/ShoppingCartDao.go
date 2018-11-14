package dao

import (
	"shop-web/module/goods/model"
	"shop-web/common/dbutil"
)

func AddCart(cart model.ShoppingCart)  error {
	if err := dbutil.DB.Table("shopping_cart").Create(&cart).Error; err != nil {
		return err
	}
	return nil
}

func FindCartList(userId string) ([]model.ShoppingCart, error) {
	var cartList []model.ShoppingCart
	if err := dbutil.DB.Table("shopping_cart").
		Where("`user_id` = ? and `status` = 1", userId).
		Find(&cartList).Error; err != nil{

		return nil,err
	}
	return cartList, nil
}

func DelCart(cartId int64, userId string) (error)  {
	return dbutil.DB.Exec(" update `shopping_cart` set `status` = 1-`status` where `user_id` = ? and `id` = ? ", userId, cartId).Error
}

func HasSameSku(skuId int64, userId string) (model.ShoppingCart)  {

	var cart model.ShoppingCart
	err :=  dbutil.DB.Table("shopping_cart").Where("`user_id` = ? and `product_id` = ? ", userId, skuId).Find(&cart).Error
	if err != nil || cart.Id != 0{
		return cart
	}
	return cart
}

func UpdateSkuCount(id int64, goodsCount int) (error)  {
	return dbutil.DB.Exec(" update `shopping_cart` set `goods_count` = `goods_count` + ?  where `id` = ? ", goodsCount, id).Error
}