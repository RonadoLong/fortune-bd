package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/micro/go-micro/client"
	"shop-micro/commonUtils"
	videoPb "shop-micro/service/video-service/proto"
)

const NAME = "shop.srv.video"

func FindVideoList(c *gin.Context) {
	category := c.Param("category")
	if category == ""{
		CreateErrorParams(c)
		return
	}
	videoClient:= videoPb.NewVideoService(NAME, client.DefaultClient)
	offset, pageSize := commonUtils.GetOffset(c)
	req := &videoPb.VideoListReq{
		PageNum: int32(offset),
		PageSize: int32(pageSize),
		Category: category,
	}

	resp, err := videoClient.GetVideoList(c, req)
	if err != nil {
		log.Fatalf("call GetVideoList err %v \n", err)
		CreateErrorParams(c)
		return
	}
	CreateSuccess(c, resp.VideoResp)
}
