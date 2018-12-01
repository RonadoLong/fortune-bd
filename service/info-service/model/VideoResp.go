package model



type VideoResp struct {

	Id int64 `json:"id"`
	Title string `json:"title"`
	ThumbUrl string `json:"thumbUrl"`
	Avatar string `json:"avatar"`
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
}