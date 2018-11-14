package main

import (
	"shop-micro/service/video-service/database"
	"shop-micro/service/video-service/handler"
	_ "shop-micro/service/video-service/subscriber"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	pb "shop-micro/service/video-service/proto/video"
)

var videoService handler.VideoService

func init()  {
	db, err := database.CreateConnection()
	if err != nil {
		log.Fatalf("connect db err %v", err)
	}

	repo := &handler.VideoRepository{DB: db}
	videoService = handler.VideoService{Repo: repo}

}

func main() {



	// New Service
	service := micro.NewService(
		micro.Name("shop.srv.video"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

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
