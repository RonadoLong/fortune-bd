package client

import (
	"context"
	"github.com/micro/go-micro/client"
	"shop-micro/service/home-service/config"
	homePb "shop-micro/service/home-service/proto"
)

type HomeClient struct {
	client       homePb.HomeService
	serviceName string
}

func NewHomeClient() *HomeClient {
	c := homePb.NewHomeService(config.SRV_NAME, client.DefaultClient)
	return &HomeClient{
		client:           c,
		serviceName: config.SRV_NAME,
	}
}

func (h *HomeClient) FindHomeHeadList(ctx context.Context) (interface{}, error){
	homeNavListResp, err := h.client.FindHomeHeaders(ctx, nil)
	if err != nil {
		return nil, err
	}
	return homeNavListResp, nil
}