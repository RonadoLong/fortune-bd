package handler

import (
	"context"
	"log"
	pb "shop-micro/service/info-service/proto"
	"shop-micro/shopproto/video"
)

type InfoHandler struct {
	Repo *InfoRepository
}

func (infos *InfoHandler) GetVideoDetail(c context.Context, req *pb.VideoDetailReq, resp *shop_srv_shopproto.Video) error {
	return nil
}

func (infos *InfoHandler) GetNewsCategoryList(c context.Context, req *pb.Request, resp *pb.NewsCategorysResp) error {
	err := infos.Repo.GetNewsCategoryList(req, resp)
	if err != nil{
		log.Printf("GetNewsCategoryList err %v", err)
	}
	return nil
}

func (infos *InfoHandler) GetNewsList(c context.Context, req *pb.InfoListReq, resp *pb.NewsListResp) error {
	err := infos.Repo.GetNewsList(req, resp)
	if err != nil{
		log.Printf("GetNewsCategoryList err %v", err)
	}
	return nil
}

func (infos *InfoHandler) GetCategoryList(c context.Context,req *pb.Request, resp *pb.VideoCategorysResp) error {
	err := infos.Repo.GetCategoryList(req, resp)
	if err != nil{
		log.Printf("GetCategoryList err %v", err)
	}
	return nil
}

func (infos *InfoHandler) GetVideoList(ctx context.Context, req *pb.InfoListReq, resp *pb.VideoListResp) error {
	log.Printf("req %v", req)
	videoList, err := infos.Repo.FindVideosList(req)
	if err != nil {
		log.Printf("err %v", err)
	}
	resp.VideoList = videoList
	return nil
}
