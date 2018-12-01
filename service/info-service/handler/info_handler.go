package handler

import (
	"context"
	"log"
	pb "shop-micro/service/news-service/proto"
)

type NewsHandler struct {
	Repo *VideoRepository
}

func (ns *NewsHandler) GetVideoList(ctx context.Context, req *pb.VideoListReq, resp *pb.VideoListResp) error {
	log.Printf("req %v", req)
	videoResp, err := ns.Repo.FindVideosList(req)
	if err != nil {
		log.Printf("err %v", err)
	}
	resp.VideoResp = videoResp
	return nil
}
