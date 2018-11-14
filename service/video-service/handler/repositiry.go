package handler

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"shop-micro/commonUtils"
	pb "shop-micro/service/video-service/proto/video"
	"shop-web/module/news/newsModel"
)

type Repository interface {
	FindVideosList(req *pb.VideoListReq) (*pb.VideoListResp, error)
}

type VideoRepository struct {
	DB *gorm.DB
}

func (vs *VideoRepository)FindVideosList(req *pb.VideoListReq) ([]*pb.VideotResp, error) {

	total, err := vs.FindVideoCount(req.CategoryId)
	if err != nil {
		return nil, err
	}

	offset, pageSize := commonUtils.GeneratorPage(int(req.PageNum), int(req.PageSize))
	if offset >= total{
		return nil, nil
	}

	videoList,err := vs.FindVideosListByOffset(offset, pageSize, req.CategoryId)
	if err != nil {
		return nil, nil
	}

	videoRespList := []pb.VideotResp{}
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

		var pusherInfo newsModel.PusherInfo
		err := json.Unmarshal([]byte(video.PusherInfo), &pusherInfo)
		if err == nil {
			videoResp.Author = pusherInfo.Name
			videoResp.Avatar = pusherInfo.Avatar
		}
		videoRespList = append(videoRespList, videoResp)
	}

	return nil, nil
}







