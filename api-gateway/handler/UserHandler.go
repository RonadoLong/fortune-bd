package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"shop-micro/service/user-service/proto"
)

func Login(c *gin.Context) (string, bool) {

	var login = shop_srv_user.LoginReq{}

	if err := c.ShouldBindWith(&login, binding.JSON); err != nil {
		CreateErrorParams(c)
		return "", false
	}

	return "", true
}

func GetPhoneCode(c *gin.Context) {

}