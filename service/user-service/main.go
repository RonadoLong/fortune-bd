package main

import (
	"github.com/labstack/gommon/log"
	"github.com/micro/go-micro"
	"shop-micro/service/user-service/handler"
	"shop-micro/service/user-service/proto"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

const SRV_NAME = "shop.srv.user"

func main() {

	userHandler, err := handler.NewUserHandler()
	if err != nil {
		log.Printf("NewUserHandler err %v", err)
		return
	}

	userService := micro.NewService(
		micro.Name(SRV_NAME),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
	)
	userService.Init()

	err = shop_srv_user.RegisterUserHandler(userService.Server(), userHandler)
	if err != nil {
		log.Printf("RegisterUserHandler err %v", err)
		return
	}

	if err := userService.Run(); err != nil {
		log.Printf("userService.Run err %v", err)
	}
}
