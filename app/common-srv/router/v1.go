package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/v2/errors"
	pb "wq-fotune-backend/api/common"
	"wq-fotune-backend/api/protocol"
	"wq-fotune-backend/app/common-srv/client"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/response"
)


var (
	commonService pb.CommonService
)

func apiV1(group *gin.RouterGroup) {
	commonService = client.NewCommonClient(env.EtcdAddr)
	group.GET("/carousels", GetCarousels)
	group.GET("/contact", GetContact)
	group.GET("/appVersion/:platform", GetAppVersion)
	group.GET("/userRateRank", GetRateRank)
	group.GET("/userRateYearRank", GetRateYearRank)
	group.POST("/pushProfitImage", pushProfitImage)
}

func pushProfitImage(c *gin.Context) {
	var req pb.PushImageReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	oss, err := commonService.PushProfitImageOss(context.Background(), &req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, oss)
}

func GetRateRank(c *gin.Context) {
	rank, err := commonService.GetUserRateRank(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	var dataResp []protocol.RateRank
	err = jsoniter.Unmarshal(rank.Data, &dataResp)
	if err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, dataResp)
}

func GetRateYearRank(c *gin.Context) {
	rank, err := commonService.GetUserRateYearRank(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	var dataResp []protocol.RateRank
	err = jsoniter.Unmarshal(rank.Data, &dataResp)
	if err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, dataResp)
}

func GetCarousels(c *gin.Context) {
	list, err := commonService.Carousel(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	var resp []protocol.CarouselResp
	if err := jsoniter.Unmarshal(list.Carousels, &resp); err != nil {
		logger.Infof("GetCarousels  err:%v", err.Error())
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, resp)
}

func GetContact(c *gin.Context) {
	contact, err := commonService.CustomerService(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, &protocol.Contract{
		Image:   contact.Image,
		Content: contact.Content,
		Contact: contact.Contact,
	})
}

func GetAppVersion(c *gin.Context) {
	platform := c.Param("platform")
	version, err := commonService.GetAppVersion(context.Background(), &pb.VersionReq{Platform: platform})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	resp := protocol.AppVersionResp{
		Id:          version.Id,
		HasUpdate:   version.HasUpdate,
		IsIgnorable: version.IsIgnorable,
		VersionCode: version.VersionCode,
		VersionName: version.VersionName,
		UpdateLog:   version.UpdateLog,
		ApkUrl:      version.ApkUrl,
	}
	response.NewSuccess(c, resp)
}
