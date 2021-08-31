package response

import "time"

type CarouselResp struct {
	ID        int       `json:"id"`
	Image     string    `json:"image"`
	ClickUrl  string    `json:"click_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppVersionResp struct {
	Id          int32  `json:"id"`
	HasUpdate   bool   `json:"has_update"`
	IsIgnorable bool   `json:"is_ignorable"`
	VersionCode int32  ` json:"version_code"`
	VersionName string `json:"version_name"`
	UpdateLog   string ` json:"update_log"`
	ApkUrl      string `json:"apk_url"`
	IosUrl      string `json:"ios_url"`
}

type RateRank struct {
	ID             int    `json:"id"`
	UserId         string `json:"user_id"`
	Avatar         string `json:"avatar"`
	Name           string `json:"name"`
	RateReturn     string `json:"rate_return"`
	RateReturnYear string `json:"rate_return_year"`
}
