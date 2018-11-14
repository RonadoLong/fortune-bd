package main

import (
	"fmt"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"shop-micro/commonUtils"
	"shop-micro/service/home-service/handler"
	pb "shop-micro/service/home-service/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("shop.srv.home"),
		micro.Version("latest"),
	)

	redisPool := commonUtils.CreateRedisPool()
	db, err := commonUtils.CreateConnection()
	if err != nil {
		fmt.Printf("connect db error %v", err.Error())
		return
	}

	// Initialise service
	service.Init()

	repo := handler.HomeRepository{
		RedisPool: redisPool,
		DB: db,
	}
	homeHandler := &handler.HomeHandle{
		Repo: &repo,
	}

	// Register Handler
	pb.RegisterHomeServiceHandler(service.Server(), homeHandler)
	//
	//// Register Struct as Subscriber
	//micro.RegisterSubscriber("shop.srv.home-service", service.Server(), new(subscriber.Example))
	//
	//// Register Function as Subscriber
	//micro.RegisterSubscriber("shop.srv.home-service", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
