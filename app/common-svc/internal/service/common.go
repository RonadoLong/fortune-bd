package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"fortune-bd/api/response"
	"fortune-bd/app/common-svc/cache"
	"fortune-bd/app/common-svc/internal/dao"
	"fortune-bd/libs/env"
	"fortune-bd/libs/logger"
	"fortune-bd/libs/oss"
	jsoniter "github.com/json-iterator/go"
	"log"
	"strings"
	"time"

	pb "fortune-bd/api/common/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	errID = "common"
)


type CommonService struct {
	pb.UnimplementedCommonServer
	dao      *dao.Dao
}

func NewCommonService() *CommonService {
	return &CommonService{
		dao: dao.New(),
	}
}

func (s *CommonService) Carousel(ctx context.Context, req *emptypb.Empty) (*pb.CarouselList, error) {
	carousels := s.dao.GetCarousels()
	if len(carousels) == 0 {
		return nil, response.NewCarouselNotFoundErrMsg(errID)
	}
	byteData, _ := jsoniter.Marshal(carousels)
	return &pb.CarouselList{Carousels: byteData}, nil
}
func (s *CommonService) CustomerService(ctx context.Context, req *emptypb.Empty) (*pb.ContractAddr, error) {
	contact, err := s.dao.GetContact()
	if err != nil {
		return nil, response.NewContractNotFoundErrMsg(errID)
	}
	var resp = &pb.ContractAddr{}
	resp.Image = contact.Image
	resp.Content = contact.Content
	resp.Contact = contact.Contact
	return resp, nil
}

func (s *CommonService) GetAppVersion(ctx context.Context, req *pb.VersionReq) (*pb.AppVersion, error) {
	appVersion, err := s.dao.GetAppVersion(req.Platform)
	if err != nil {
		return nil, response.NewDataNotFound(errID, "获取app版本号出错")
	}
	var resp = new(pb.AppVersion)
	log.Printf("%+v", appVersion)
	resp.HasUpdate = appVersion.HasUpdate
	resp.IsIgnorable = appVersion.IsIgnorable
	resp.VersionCode = appVersion.VersionCode
	resp.VersionName = appVersion.VersionName
	resp.UpdateLog = appVersion.UpdateLog
	resp.ApkUrl = appVersion.ApkUrl
	resp.Id = appVersion.ID
	return resp, nil
}
func (s *CommonService) GetUserRateRank(ctx context.Context, req *emptypb.Empty) (*pb.UserRateRankResp, error) {
	data, err := cache.GetData("rateReturnSort")
	if err != nil {
		logger.Warnf("暂无排名数据 %v", err)
		return nil, response.NewDataNotFound(errID, "暂无排名数据")
	}
	return &pb.UserRateRankResp{Data: data}, nil
}

func (s *CommonService) GetUserRateYearRank(ctx context.Context, req *emptypb.Empty) (*pb.UserRateRankResp, error) {
	data, err := cache.GetData("rateReturnYearSort")
	if err != nil {
		logger.Warnf("暂无排名数据 %v", err)
		return nil, response.NewDataNotFound(errID, "暂无排名数据")
	}
	return &pb.UserRateRankResp{Data: data}, nil
}

func (s *CommonService) PushProfitImageOss(ctx context.Context, req *pb.PushImageReq) (*pb.ImageResp, error) {
	client, err := oss.NewOssClient(oss.AccessKeyId, oss.AccessKeySecret, oss.Endpoint, "ifortune")
	if err != nil {
		return nil, response.NewInternalServerErrMsg(errID)
	}
	var resp = new(pb.ImageResp)
	nano := time.Now().UnixNano()
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(req.ImageBase64))
	objectName := fmt.Sprintf("%s/image/%v.jpeg", env.RunMode, nano)
	err = client.PushFileWithIOReader(objectName, decoder)
	if err != nil {
		if strings.Contains(err.Error(), "illegal base64") {
			return nil, response.NewInternalServerErrWithMsg(errID, "base64格式不正确")
		}
		return nil, response.NewInternalServerErrMsg(errID)
	}
	resp.ImageUrl = fmt.Sprintf("https://%s.%s/%s", "ifortune", oss.EndPointNoHttp, objectName)
	return resp, nil
}
