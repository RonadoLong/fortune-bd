package main

import (
	"shop-micro/service/video-service/handler"
	_ "shop-micro/service/video-service/subscriber"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	pb "shop-micro/service/video-service/proto/video"
)

func main() {

	db, err := handler.CreateConnection()
	if err != nil {
		log.Fatal(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("shop.srv.video"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	repo := &handler.VideoRepository{db}
	video := handler.VideoService{repo}
	// Register Handler
	pb.RegisterVideoHandler(service.Server(), &video)

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
