package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	home "shop-micro/service/home-service/client"
)

func FindHomeHeadList(c *gin.Context) {
	homeClient := home.NewHomeClient()
	resp, err := homeClient.FindHomeHeadList(c)
	if err != nil {
		log.Printf("home err %v",err)
	}
	CreateSuccess(c, resp)
}