package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	info "shop-micro/service/info-service/client"
)

const NAME = "shop.srv.video"

var (
	infoClient = info.NewInfoClient()
)

func FindVideoList(c *gin.Context) {
	videoList, err := infoClient.GetVideoList(c)
	if err != nil {
		log.Printf("FindVideoList err %v ", err)
		CreateError(c)
		return
	}
	if videoList == nil {
		CreateNotContent(c)
		return
	}
	CreateSuccess(c, videoList)
}


func FindVideoCategoryList(c *gin.Context) {
	videoCategoryList, err := infoClient.GetCategoryList(c)
	if err != nil {
		log.Printf("FindVideoList err %v ", err)
		CreateError(c)
		return
	}
	if videoCategoryList == nil {
		CreateNotContent(c)
		return
	}
	CreateSuccess(c, videoCategoryList)
}

func GetNewsList(c *gin.Context) {
	newsList, err := infoClient.GetNewsList(c)
	if err != nil {
		log.Printf("GetNewsList err %v ", err)
		CreateError(c)
		return
	}
	if newsList == nil {
		CreateNotContent(c)
		return
	}
	CreateSuccess(c, newsList)
}


func GetNewsCategoryList(c *gin.Context) {
	newsCategoryList, err := infoClient.GetNewsCategoryList(c)
	if err != nil {
		log.Printf("newsCategoryList err %v ", err)
		CreateError(c)
		return
	}
	if newsCategoryList == nil {
		CreateNotContent(c)
		return
	}
	CreateSuccess(c, newsCategoryList)
}