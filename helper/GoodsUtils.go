package helper

//
//func TransferValue(goods model.Goods, language string) (model.GoodsResp, error) {
//
//	goodsResp := model.GoodsResp{}
//	goodsResp.GoodsId = goods.GoodsId
//	goodsResp.CategoryId = goods.CategoryId
//
//	if language == "EN"{
//		goodsResp.Title = goods.EnTitle
//		goodsResp.SellPoint = goods.EnSellPoint
//	}else {
//		goodsResp.Title = goods.Title
//		goodsResp.SellPoint = goods.SellPoint
//	}
//
//	goodsResp.GoodsImages = goods.GoodsImages
//	goodsResp.Status = goods.Status
//
//	goodsDetail, err := dao.FindGoodsDetailByGoodsId(goodsResp.GoodsId)
//	if err != nil{
//		return goodsResp, err
//	}
//	goodsResp.GoodsBanners = goodsDetail.GoodsBanners
//	goodsResp.GoodsDetail = goodsDetail.GoodsDetail
//	goodsResp.GoodsDesc = goodsDetail.GoodsDesc
//
//	goodsSkuList, _ := dao.FindGoodsSkuByGoodsId(goodsResp.GoodsId)
//	goodsSku := goodsSkuList[0]
//
//	if goodsSku.SkuType == 0{
//		//单品
//		goodsResp.GoodsType = goodsSku.SkuType
//		goodsResp.Price = goodsSku.Price
//		goodsResp.MemberPrice = goodsSku.MemberPrice
//	}
//
//	if goodsSku.SkuType >= 1{
//		var goodsSkuResp = model.GoodsSkuResp{}
//		for _,sku := range goodsSkuList {
//			if sku.IsActive == 1{
//				goodsResp.GoodsType = sku.SkuType
//				goodsResp.Price = sku.Price
//				goodsResp.MemberPrice = sku.MemberPrice
//			}
//			goodsSkuResp.Price = sku.Price
//			goodsSkuResp.SkuId = sku.SkuId
//			goodsSkuResp.ProductId = sku.GoodsId
//			goodsSkuResp.MemberPrice = sku.MemberPrice
//			goodsSkuResp.ActivityEnable = sku.ActivityEnable
//			goodsSkuResp.IsActive = sku.IsActive
//			goodsSkuResp.SkuType = sku.SkuType
//			goodsSkuResp.Stock= sku.Stock
//			goodsSkuResp.SoldCount = sku.SoldCount
//			goodsSkuResp.BuyNumMax = sku.BuyNumMax
//			goodsSkuResp.SkuPic = sku.SkuPic
//
//			goodsSkuPropertyList, err := dao.FindGoodsSkuPropertyBySkuId(goodsSkuResp.SkuId)
//			if err != nil {
//				return goodsResp,err
//			}
//
//			for _,goodsSkuProperty  := range goodsSkuPropertyList {
//				skuPropertyResp := model.GoodsSkuPropertyResp{}
//				if language == "EN" {
//					skuPropertyResp.Id = goodsSkuProperty.Id
//					skuPropertyResp.SkuName = goodsSkuProperty.EnSkuName
//					skuPropertyResp.SkuValue = goodsSkuProperty.EnSkuValue
//				}else {
//					skuPropertyResp.Id = goodsSkuProperty.Id
//					skuPropertyResp.SkuName = goodsSkuProperty.SkuName
//					skuPropertyResp.SkuValue = goodsSkuProperty.SkuValue
//				}
//				goodsSkuResp.SkuPropertyRespList = append(goodsSkuResp.SkuPropertyRespList,skuPropertyResp )
//			}
//
//			goodsResp.GoodsSkuList = append(goodsResp.GoodsSkuList, goodsSkuResp)
//
//		}
//	}
//
//	return goodsResp, nil
//}
//
//func TransferValues(goods model.Product, language string) (protocol.GoodsListResp, error) {
//
//	goodsListResp := protocol.GoodsListResp{}
//	goodsListResp.ProductId = goods.ProductId
//	goodsListResp.CategoryId = goods.CategoryId
//
//	if language == "EN"{
//		goodsListResp.Title = goods.EnTitle
//		goodsListResp.SellPoint = goods.EnSellPoint
//	}else {
//		goodsListResp.Title = goods.Title
//		goodsListResp.SellPoint = goods.SellPoint
//	}
//
//	goodsListResp.GoodsImages = goods.GoodsImages
//	goodsListResp.Status = goods.Status
//	goodsListResp.Commission = goods.Commission
//	goodsListResp.Integral = goods.Integral
//	goodsListResp.Price = goods.Price
//	goodsListResp.MemberPrice = goods.MemberPrice
//
//	goodsDetail, err := dao.FindGoodsDetailByGoodsId(goods.ProductId)
//	if err == nil{
//		goodsListResp.GoodsBanners = goodsDetail.GoodsBanners
//		goodsListResp.GoodsDetail = goodsDetail.GoodsDetail
//		goodsListResp.GoodsDesc = goodsDetail.GoodsDesc
//	}
//	return goodsListResp, nil
//}
//
//
//func TransferValueForList(goods model.Goods, language string) (protocol.GoodsListResp, error) {
//
//	goodsListResp := protocol.GoodsListResp{}
//	goodsListResp.ProductId = goods.GoodsId
//	goodsListResp.CategoryId = goods.CategoryId
//
//
//	if language == "EN"{
//		goodsListResp.Title = goods.EnTitle
//		goodsListResp.SellPoint = goods.EnSellPoint
//	}else {
//		goodsListResp.Title = goods.Title
//		goodsListResp.SellPoint = goods.SellPoint
//	}
//
//	goodsListResp.GoodsImages = goods.GoodsImages
//	goodsListResp.Status = goods.Status
//
//	//goodsSkuList, _ := dao.FindGoodsSkuByGoodsId(goodsListResp.GoodsId)
//	//goodsSku := goodsSkuList[0]
//
//	//if goodsSku.SkuType == 0{
//	//	//单品
//	//	goodsListResp.GoodsType = goodsSku.SkuType
//	//	goodsListResp.Price = goodsSku.Price
//	//	goodsListResp.MemberPrice = goodsSku.MemberPrice
//	//}
//	//
//	//if goodsSku.SkuType >= 1{
//	//	for _,sku := range goodsSkuList {
//	//		if sku.IsActive == 1{
//	//			goodsListResp.GoodsType = sku.SkuType
//	//			goodsListResp.Price = sku.Price
//	//			goodsListResp.MemberPrice = sku.MemberPrice
//	//			break
//	//		}
//	//	}
//	//}
//
//	return goodsListResp, nil
//}