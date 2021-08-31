package router

import (
	"context"
	pb "fortune-bd/api/common/v1"
	"fortune-bd/api/response"
	"fortune-bd/app/common-svc/internal/service"
	"fortune-bd/libs/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	jsoniter "github.com/json-iterator/go"
)


var (
	commonService = service.NewCommonService()
)

func apiV1(group *gin.RouterGroup) {
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
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, oss)
}

func GetRateRank(c *gin.Context) {
	rank, err := commonService.GetUserRateRank(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	var dataResp []response.RateRank
	err = jsoniter.Unmarshal(rank.Data, &dataResp)
	if err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, dataResp)
}

func GetRateYearRank(c *gin.Context) {
	rank, err := commonService.GetUserRateYearRank(context.Background(),nil)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	var dataResp []response.RateRank
	err = jsoniter.Unmarshal(rank.Data, &dataResp)
	if err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, dataResp)
}

func GetCarousels(c *gin.Context) {
	list, err := commonService.Carousel(context.Background(), nil)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	var resp []response.CarouselResp
	if err := jsoniter.Unmarshal(list.Carousels, &resp); err != nil {
		logger.Infof("GetCarousels  err:%v", err.Error())
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, resp)
}

func GetContact(c *gin.Context) {
	contact, err := commonService.CustomerService(context.Background(), nil)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, &response.Contract{
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
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	resp := response.AppVersionResp{
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
