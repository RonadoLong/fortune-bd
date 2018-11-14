package main

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/micro/go-web"
	"shop-micro/api-gateway/handler"
)

const SRV_NAME = "shop.srv.apigateway"

func main() {

	webSrv := web.NewService(
		web.Name(SRV_NAME),
		web.Address(":20050"),
	)

	router := handler.ClientEngine()
	router.Use(gin.Logger())
	webSrv.Handle("/", router)
	if err := webSrv.Init(); err != nil {
		log.Fatal(err)
	}

	if err := webSrv.Run(); err != nil {
		log.Fatalf("run err %v ", err)
	}

}
