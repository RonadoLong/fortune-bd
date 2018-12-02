package handler

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
	"shop-micro/helper"
	pb "shop-micro/service/info-service/proto"
	"shop-micro/shopproto/video"
)

type Video struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	ThumbUrl string `json:"thumbUrl"`
	Author string `json:"author"`
	Duration string `json:"duration"`
	ReadCount string `json:"readCount"`
	CommentCount string `json:"commentCount"`
	LikeCount int `json:"likeCount"`
	Category string `json:"category"`
	ViewType int `json:"viewType"`
	IsRecommend int `json:"isRecommend"`
	Content string `json:"content"`
	Tags string `json:"tags"`
	VideoDesc string `json:"videoDesc"`
	PusherInfo string `json:"pusherInfo"`
}
type InfoRepository struct {
	DB *gorm.DB
}

func (info *InfoRepository) FindVideosList(req *pb.InfoListReq) ([]*shop_srv_shopproto.Video, error) {

	var total int
	query := "status = 1"
	if req.Category != "" {
		query = fmt.Sprintf("category = '%s' and %s", req.Category, query)
	}

	err := info.DB.Table("video").Where(query).Count(&total).Error
	if err != nil {
		return nil, err
	}

	offset, pageSize := helper.GeneratorPage(int(req.PageNum), int(req.PageSize))
	if offset >= total {
		return nil, errors.New("offset > total")
	}

	var videoList []*Video
	err = info.DB.Table("video").
		Where(query).
		Order("`create_time` desc").Offset(offset).Limit(pageSize).
		Find(&videoList).Error
	if err != nil {
		return nil, err
	}

	var videoRespList []*shop_srv_shopproto.Video
	for _, video := range videoList {
		videoResp := shop_srv_shopproto.Video{}
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

		var pusherInfo shop_srv_shopproto.VideoPusherInfo
		err := jsoniter.Unmarshal([]byte(video.PusherInfo), &pusherInfo)
		if err == nil {
			videoResp.Author = pusherInfo.Name
			videoResp.Avatar = pusherInfo.Avatar
		}
		videoRespList = append(videoRespList, &videoResp)
	}

	return videoRespList, nil
}

func (info *InfoRepository) GetCategoryList(request *pb.Request, resp *pb.VideoCategorysResp) error {
	err := info.DB.Table("video_category").Order("`sort` asc").Where("`status` = 1 ").Find(&resp.VideoCategoryList).Error
	if err != nil {
		return err
	}
	return nil
}

func (info *InfoRepository) GetNewsCategoryList(request *pb.Request, resp *pb.NewsCategorysResp) error {
	err := info.DB.Table("news_category").Order("`sort` asc").Where("`status` = 1 ").Find(&resp.NewsCategoryList).Error
	if err != nil{
		return err
	}
	return nil
}

func (info *InfoRepository) GetNewsList(req *pb.InfoListReq, resp *pb.NewsListResp) error {

	query := " status = 1"
	if req.Category != "" {
		query = fmt.Sprintf(`category = '["%s"]' and %s`, req.Category, query)
	}

	total := 0
	if err := info.DB.Table("news").Where(query).Count(&total).Error; err != nil {
		return err
	}

	if total == 0 {
		return errors.New("no more content")
	}

	offset, pageSize := helper.GeneratorPage(int(req.PageNum), int(req.PageSize))
	if offset >= total {
		return errors.New("offset > total")
	}

	if err := info.DB.Table("news").Where(query).Order("`create_time` desc").Offset(offset).Limit(pageSize).Find(&resp.NewsList).Error; err != nil {
		return err
	}
	return nil
}


