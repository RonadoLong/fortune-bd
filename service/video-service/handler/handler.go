package handler

import (
	"context"
	pb "shop-micro/service/video-service/proto/video"
)

type VideoService struct {
	Repo *VideoRepository
}

func (vs *VideoService) GetVideoList(ctx context.Context, req *pb.VideoListReq, resp *pb.VideoListResp) error {
	resp.VideoResp, _ = vs.Repo.FindVideosList(req)
	return nil
}
