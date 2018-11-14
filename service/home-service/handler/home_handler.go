package handler

import (
	"context"
	pb "shop-micro/service/home-service/proto"
)

type HomeHandle struct{
	Repo *HomeRepository
}

func (h *HomeHandle) FindHomeNav(ctx context.Context, req *pb.HomeNavListReq, resp *pb.HomeNavListResp) (error) {
	h.Repo.FindHomeNav(req, resp)
	return nil
}

func (h *HomeHandle) FindHomeList(ctx context.Context, req *pb.HomeContentListReq, resp *pb.HomeContentListResp,) ( error) {
	return nil
}


