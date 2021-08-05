package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"strings"
	"time"
	pb "wq-fotune-backend/api/common"
	"wq-fotune-backend/app/common-srv/cache"
	"wq-fotune-backend/app/common-srv/internal/dao"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/oss"
	"wq-fotune-backend/pkg/response"
)

const (
	errID = "common"
)

type CommonService struct {
	dao      *dao.Dao
	cacheSrv *cache.Service
}

func (c *CommonService) PushProfitImageOss(ctx context.Context, req *pb.PushImageReq, resp *pb.ImageResp) error {
	client, err := oss.NewOssClient(oss.AccessKeyId, oss.AccessKeySecret, oss.Endpoint, "ifortune")
	if err != nil {
		return response.NewInternalServerErrMsg(errID)
	}
	nano := time.Now().UnixNano()
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(req.ImageBase64))
	objectName := fmt.Sprintf("%s/image/%v.jpeg", env.RunMode, nano)
	err = client.PushFileWithIOReader(objectName, decoder)
	if err != nil {
		if strings.Contains(err.Error(), "illegal base64") {
			return response.NewInternalServerErrWithMsg(errID, "base64格式不正确")
		}
		return response.NewInternalServerErrMsg(errID)
	}
	resp.ImageUrl = fmt.Sprintf("https://%s.%s/%s", "ifortune", oss.EndPointNoHttp, objectName)
	return nil
}

func (c *CommonService) GetUserRateRank(ctx context.Context, e *empty.Empty, resp *pb.UserRateRankResp) error {
	data, err := c.cacheSrv.GetData("rateReturnSort")
	if err != nil {
		logger.Warnf("暂无排名数据 %v", err)
		return response.NewDataNotFound(errID, "暂无排名数据")
	}
	resp.Data = data
	return nil
}

func (c *CommonService) GetUserRateYearRank(ctx context.Context, e *empty.Empty, resp *pb.UserRateRankResp) error {
	data, err := c.cacheSrv.GetData("rateReturnYearSort")
	if err != nil {
		logger.Warnf("暂无排名数据 %v", err)
		return response.NewDataNotFound(errID, "暂无排名数据")
	}
	resp.Data = data
	return nil
}

func (c *CommonService) GetAppVersion(ctx context.Context, req *pb.VersionReq, resp *pb.AppVersion) error {
	appVersion, err := c.dao.GetAppVersion(req.Platform)
	if err != nil {
		return response.NewDataNotFound(errID, "获取app版本号出错")
	}
	log.Printf("%+v", appVersion)
	resp.HasUpdate = appVersion.HasUpdate
	resp.IsIgnorable = appVersion.IsIgnorable
	resp.VersionCode = appVersion.VersionCode
	resp.VersionName = appVersion.VersionName
	resp.UpdateLog = appVersion.UpdateLog
	resp.ApkUrl = appVersion.ApkUrl
	resp.Id = appVersion.ID
	return nil
}


func (c *CommonService) Carousel(ctx context.Context, req *empty.Empty, resp *pb.CarouselList) error {
	carousels := c.dao.GetCarousels()
	if len(carousels) == 0 {
		return response.NewCarouselNotFoundErrMsg(errID)
	}
	byteData, _ := json.Marshal(carousels)
	resp.Carousels = byteData
	return nil
}

func (c *CommonService) CustomerService(ctx context.Context, req *empty.Empty, resp *pb.ContractAddr) error {
	contact, err := c.dao.GetContact()
	if err != nil {
		return response.NewContractNotFoundErrMsg(errID)
	}
	resp.Image = contact.Image
	resp.Content = contact.Content
	resp.Contact = contact.Contact
	return nil
}
