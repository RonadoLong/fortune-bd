package client

import (
	"github.com/micro/go-micro/client"
	"shop-micro/service/news-service/config"
	pb "shop-micro/service/news-service/proto"
)

type VideoClient struct {
	client *pb.VideoService
	servicename string
}

func NewVideoClient() *VideoClient {
	videoService := pb.NewVideoService(config.SRV_NAME, client.DefaultClient)
	return &VideoClient{
		client:&videoService,
		servicename:config.SRV_NAME,
	}
}

