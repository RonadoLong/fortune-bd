package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"log"
	resp2 "shop-micro/api-gateway/resp"
	videoPb "shop-micro/service/video-service/proto/video"
)

const NAME = "shop.srv.video"

func FindVideoList(ctx *gin.Context) {
	videoClient:= videoPb.NewVideoService(NAME, client.DefaultClient)
	resp, err := videoClient.GetVideoList(ctx, &videoPb.VideoListReq{})
	if err != nil {
		log.Fatalf("call GetGoodsList err %v \n", err)
	}
	log.Println(resp)
	resp2.CreateSuccess(ctx, resp.VideoResp)
}
