package main

import (
	"github.com/micro/go-micro"
	"log"
	"shop-micro/helper"
	"shop-micro/hystrix"
	"shop-micro/service/home-service/config"
	"shop-micro/service/home-service/handler"
	pb "shop-micro/service/home-service/proto"
	"time"
)

func main() {

	hystrix.Configure([]string{
		config.SRV_NAME + ".HomeHandle.FindHomeHeaders",
		config.SRV_NAME + ".HomeHandle.FindHomeContents",
	})

	// New Service
	service := micro.NewService(
		micro.Name(config.SRV_NAME),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.WrapClient(hystrix.NewClientWrapper()),
	)

	redisPool := helper.CreateRedisPool()
	db, err := helper.CreateConnection()
	if err != nil {
		log.Printf("connect db error %v", err.Error())
		return
	}

	// Initialise service
	service.Init()

	repo := handler.HomeRepository{
		RedisPool: redisPool,
		DB:        db,
	}
	homeHandler := &handler.HomeHandle{
		Repo: &repo,
	}

	// Register Handler
	_ = pb.RegisterHomeServiceHandler(service.Server(), homeHandler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
