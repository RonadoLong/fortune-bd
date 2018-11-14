package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	resp2 "shop-micro/api-gateway/resp"
	"shop-micro/commonUtils"
	videoPb "shop-micro/service/video-service/proto/video"
)

const NAME = "shop.srv.video"

func FindVideoList(c *gin.Context) {
	category := c.Param("category")

	if category == ""{
		resp2.CreateErrorParams(c)
		return
	}

	videoClient:= videoPb.NewVideoService(NAME, client.DefaultClient)
	offset, pageSize := commonUtils.GetOffset(c)
	req := &videoPb.VideoListReq{
		PageNum: int32(offset),
		PageSize: int32(pageSize),
		CategoryId: category,
	}

	resp, err := videoClient.GetVideoList(c, req)
	if err != nil {
		log.Fatalf("call GetGoodsList err %v \n", err)
		resp2.CreateErrorParams(c)
		return
	}

	resp2.CreateSuccess(c, resp.VideoResp)
}
