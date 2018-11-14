package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	homePb "shop-micro/service/home-service/proto"
)

const HOME_NAME = "shop.srv.home"

func FindHomeHeadList(c *gin.Context) {
	homeService := homePb.NewHomeService(HOME_NAME, client.DefaultClient)
	homeNavListResp, err := homeService.FindHomeNav(c, &homePb.HomeNavListReq{})
	if err != nil {
		fmt.Printf("err %v \n", err)
		CreateErrorRequest(c)
		return
	}

	fmt.Printf("homeNavListResp %v \n", homeNavListResp)
	CreateSuccess(c, homeNavListResp)
}