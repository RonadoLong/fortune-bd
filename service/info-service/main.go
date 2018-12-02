package main

import (
	"shop-micro/helper"
	"shop-micro/hystrix"
	"shop-micro/service/info-service/config"
	"shop-micro/service/info-service/handler"
	_ "shop-micro/service/info-service/subscriber"
	"time"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	pb "shop-micro/service/info-service/proto"
)

func main() {

	hystrix.Configure([]string{
		config.SRV_NAME + "infoHandler.GetCategoryList",
		config.SRV_NAME + "infoHandler.GetVideoList",
	})
	// New Service
	service := micro.NewService(
		micro.Name(config.SRV_NAME),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.WrapHandler(logWrapper),
		micro.WrapClient(hystrix.NewClientWrapper()),
	)

	// Initialise service
	service.Init()

	db, err := helper.CreateConnection()
	if err != nil {
		log.Fatalf("connect db err %v", err)
	}

	repo := &handler.InfoRepository{DB: db}
	infoHandler := handler.InfoHandler{Repo: repo}

	// Register Handler
	_ = pb.RegisterInfoHandler(service.Server(), &infoHandler)

	//// Register Struct as Subscriber
	//micro.RegisterSubscriber("shop.srv.video", service.Server(), new(subscriber.Example))
	//
	//// Register Function as Subscriber
	//micro.RegisterSubscriber("shop.srv.video", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
