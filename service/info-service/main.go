package main

import (
	"shop-micro/commonUtils"
	"shop-micro/service/news-service/config"
	"shop-micro/service/news-service/handler"
	_ "shop-micro/service/news-service/subscriber"
	"time"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	pb "shop-micro/service/news-service/proto"
)

func main() {

	// New Service
	service := micro.NewService(
		micro.Name(config.SRV_NAME),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.WrapHandler(logWrapper),
	)

	// Initialise service
	service.Init()

	db, err := commonUtils.CreateConnection()
	if err != nil {
		log.Fatalf("connect db err %v", err)
	}

	repo := &handler.VideoRepository{DB: db}
	videoService := handler.VideoService{Repo: repo}

	// Register Handler
	pb.RegisterVideoHandler(service.Server(), &videoService)

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
