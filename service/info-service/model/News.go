package model

type News struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	ThumbUrl string `json:"thumbUrl"`
	Author string `json:"author"`
	Avatar string `json:"avatar"`
	ReadCount int `json:"readCount"`
	CommentCount int `json:"commentCount"`
	LikeCount int `json:"likeCount"`
	Category string `json:"category"`
	ViewType int `json:"viewType"`
	IsRecommend int `json:"isRecommend"`
	Content string `json:"content"`
	CreateTime string `json:"createTime"`
} 