package handler

import (
	"fmt"
	"shop-micro/helper"
	"shop-micro/service/home-service/proto"
	"testing"
)

func TestHomeRepository(t *testing.T) {
	redisPool := helper.GetRedisPool("0.0.0.0:6379","")
	db, err := helper.GetDbByHost("root", "123456", "0.0.0.0")
	if err != nil {
		fmt.Printf("connect db error %v \n", err.Error())
		return
	}
	homeRepository := HomeRepository{
		DB:        db,
		RedisPool: redisPool,
	}

	resp := shop_srv_home.HomeHeadersResp{}
	err = homeRepository.FindHomeNav(&shop_srv_home.HomeHeaderReq{}, &resp)
	fmt.Printf("resp %v \n", resp.HomeCourseList)

}