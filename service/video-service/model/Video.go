package model



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