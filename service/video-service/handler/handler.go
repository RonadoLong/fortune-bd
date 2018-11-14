package handler

import (
	"context"
	"log"
	pb "shop-micro/service/video-service/proto/video"
)

type VideoService struct {
	Repo *VideoRepository
}

func (vs *VideoService) GetVideoList(ctx context.Context, req *pb.VideoListReq, resp *pb.VideoListResp) error {
	log.Printf("req %v", req)

	videotResps, err := vs.Repo.FindVideosList(req)
	log.Printf("videotResps %v", videotResps)

	if err != nil {
		log.Printf("err %v", err)
	}

	resp.VideoResp = videotResps
	return nil
}
