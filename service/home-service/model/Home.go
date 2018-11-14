package model


type HomeNav struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	EnTitle string `json:"enTitle"`
	ImgUrl string `json:"imgUrl"`
	Types int `json:"types"`
	ClassId int `json:"classId"`
}

type HomeCarousel struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	EnTitle string `json:"enTitle"`
	ImgUrl string `json:"imgUrl"`
	Url string `json:"url"`
}

type HomeNews struct {
	Id int64 `json:"id"`
	Type int `json:"type"`
	Status string `json:"status"`
}

type HomeGoods struct {
	Id int64 `json:"id"`
	Type int `json:"type"`
	Status string `json:"status"`
}
