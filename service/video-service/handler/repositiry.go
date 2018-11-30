package handler

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
	"shop-micro/commonUtils"
	"shop-micro/service/video-service/model"
	pb "shop-micro/service/video-service/proto"
)

type VideoRepository struct {
	DB *gorm.DB
}

func (vs *VideoRepository)FindVideosList(req *pb.VideoListReq) ([]*pb.VideotResp, error) {

	total, err := vs.FindVideoCount(req.Category)
	if err != nil {
		return nil, err
	}

	offset, pageSize := commonUtils.GeneratorPage(int(req.PageNum), int(req.PageSize))
	if offset >= total{
		return nil, errors.New("offset > total")
	}

	videoList,err := vs.FindVideosListByOffset(offset, pageSize, req.Category)
	if err != nil {
		return nil, err
	}

	var videoRespList []*pb.VideotResp
	for _,video := range videoList {
		videoResp := pb.VideotResp{}
		videoResp.Id = video.Id
		videoResp.Title = video.Title
		videoResp.ThumbUrl = video.ThumbUrl
		videoResp.Duration = video.Duration
		videoResp.ReadCount = video.ReadCount
		videoResp.CommentCount = video.CommentCount
		videoResp.LikeCount = int32(video.LikeCount)
		videoResp.Category = video.Category
		videoResp.Content = video.Content
		videoResp.Tags = video.Tags
		videoResp.VideoDesc = video.VideoDesc

		var pusherInfo model.PusherInfo
		err := jsoniter.Unmarshal([]byte(video.PusherInfo), &pusherInfo)
		if err == nil {
			videoResp.Author = pusherInfo.Name
			videoResp.Avatar = pusherInfo.Avatar
		}
		videoRespList = append(videoRespList, &videoResp)
	}

	return videoRespList, nil
}







