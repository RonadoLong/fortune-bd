package main

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/micro/go-web"
	"shop-micro/api-gateway/handler"
)

const SRV_NAME = "shop.srv.apigateway"

func main() {
	// optionally setup command line usage
	//cmd.Init(
	//	cmd.Name(SRV_NAME),
	//	cmd.Version("latest"),
	//)

	webSrv := web.NewService(
		web.Name(SRV_NAME),
		web.Address(":20050"),
	)

	router := handler.ClientEngine()
	router.Use(gin.Logger())
	webSrv.Handle("/", router)

	//server := &http.Server{
	//	Addr:           ":20050",
	//	Handler:        router,
	//	ReadTimeout:    15 * time.Second,
	//	WriteTimeout:   15 * time.Second,
	//	IdleTimeout:    120 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}

	if err := webSrv.Init(); err != nil {
		log.Fatal(err)
	}


	if err := webSrv.Run(); err != nil {
		log.Fatalf("run err %v ", err)
	}
	// Run server
	//server.ListenAndServe()
}
