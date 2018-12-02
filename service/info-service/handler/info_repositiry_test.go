package handler

import (
	"fmt"
	"shop-micro/helper"
	"shop-micro/service/info-service/proto"
	"testing"
)

func TestInfoRepository(t *testing.T) {
	db, err := helper.GetDbByHost("root", "123456", "0.0.0.0")
	if err != nil {
		fmt.Printf("connect db error %v \n", err.Error())
		return
	}
	repository := InfoRepository{DB: db}
	//resp := shop_srv_info.VideoCategorysResp{}
	//
	//err = repository.GetCategoryList(&shop_srv_info.Request{}, &resp)
	//if err != nil {
	//	fmt.Printf("GetCategoryListerror %v \n", err.Error())
	//	return
	//}
	//
	//for _, cate := range resp.VideoCategoryList {
	//	fmt.Printf("resp %s \n", cate.Title)
	//	req := &shop_srv_info.InfoListReq{
	//		Category: cate.Title,
	//	}
	//	videos, err := repository.FindVideosList(req)
	//	fmt.Printf("videos %s  err = %v\n", videos, err)
	//}
	resp := &shop_srv_info.NewsCategorysResp{}
	err = repository.GetNewsCategoryList(&shop_srv_info.Request{}, resp)
	if err != nil {
		fmt.Printf("GetCategoryList error %v \n", err.Error())
		return
	}

	for _, cate := range resp.NewsCategoryList {
		fmt.Printf("resp %s \n", cate.Title)
		req := &shop_srv_info.InfoListReq{
			Category: cate.Title,
		}
		listResp := &shop_srv_info.NewsListResp{}
		err := repository.GetNewsList(req, listResp)
		fmt.Printf("videos %s  err = %v\n", listResp, err)
	}

}