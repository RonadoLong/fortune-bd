package handler

import (
	"shop-micro/service/video-service/model"
)

func (vs *VideoRepository) FindVideosListByOffset( pageNum int, pageSize int, category string) ([]*model.Video, error) {
	var videoList []*model.Video
	if err := vs.DB.Table("video").
		Where("category = ? and status = 1", category).
		Order("`create_time` desc").Offset(pageNum).Limit(pageSize).
		Find(&videoList).Error; err != nil{
		return videoList, err
	}
	return videoList, nil
}


func (vs *VideoRepository)FindVideoCount(category string) (int,error){
	count := 0
	if err := vs.DB.Table("video").Where("category = ? and status = 1", category).Count(&count).Error; err != nil{
		return count, err
	}
	return count, nil
}