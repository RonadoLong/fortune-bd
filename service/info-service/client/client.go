package client

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"shop-micro/helper"
	"shop-micro/service/info-service/config"
	pb "shop-micro/service/info-service/proto"
)

type InfoClient struct {
	client pb.InfoService
	servicename string
}

func NewInfoClient() *InfoClient {
	videoService := pb.NewInfoService(config.SRV_NAME, client.DefaultClient)
	return &InfoClient{
		client: videoService,
		servicename: config.SRV_NAME,
	}
}

func (c *InfoClient) GetCategoryList(ctx context.Context) (interface{}, error) {
	return c.client.GetCategoryList(ctx, nil)
}

func (c *InfoClient) GetVideoList(ctx *gin.Context) (interface{}, error){
	category := ctx.Param("category")
	offset, pageSize := helper.GetOffset(ctx)

	req := &pb.InfoListReq{
		PageNum: int32(offset),
		PageSize: int32(pageSize),
		Category: category,
	}
	return c.client.GetVideoList(ctx,req)
}

func (c *InfoClient) GetNewsCategoryList(ctx context.Context) (interface{}, error) {
	return c.client.GetNewsCategoryList(ctx, nil)
}

func (c *InfoClient) GetNewsList(ctx *gin.Context) (interface{}, error){
	category := ctx.Param("category")
	offset, pageSize := helper.GetOffset(ctx)

	req := &pb.InfoListReq{
		PageNum: int32(offset),
		PageSize: int32(pageSize),
		Category: category,
	}
	return c.client.GetNewsList(ctx,req)
}