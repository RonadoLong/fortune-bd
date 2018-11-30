package handler

import (
	"context"
	"log"
	pb "shop-micro/service/video-service/proto"
)

type VideoService struct {
	Repo *VideoRepository
}

func (vs *VideoService) GetVideoList(ctx context.Context, req *pb.VideoListReq, resp *pb.VideoListResp) error {
	log.Printf("req %v", req)
	videoResp, err := vs.Repo.FindVideosList(req)
	if err != nil {
		log.Printf("err %v", err)
	}
	resp.VideoResp = videoResp
	return nil
}
