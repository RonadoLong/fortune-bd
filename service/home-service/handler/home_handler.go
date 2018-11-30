package handler

import (
	"context"
	"log"
	pb "shop-micro/service/home-service/proto"
)

type HomeHandle struct{
	Repo *HomeRepository
}

func (h *HomeHandle) FindHomeHeaders(c context.Context, req *pb.HomeHeaderReq, resp *pb.HomeHeadersResp) error {
	log.Printf("FindHomeHeaders %v", req)
	return h.Repo.FindHomeNav(req, resp)
}

func (h *HomeHandle) FindHomeContents(context.Context, *pb.HomeContentsReq, *pb.HomeContentsResp) error {
	return nil
}



