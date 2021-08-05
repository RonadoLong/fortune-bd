package protocol

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
	HasUpdate   bool   `json:"hasUpdate"`
	IsIgnorable bool   `json:"isIgnorable"`
	VersionCode int32  ` json:"versionCode"`
	VersionName string `json:"versionName"`
	UpdateLog   string ` json:"updateLog"`
	ApkUrl      string `json:"apkUrl"`
	IosUrl      string `json:"iosUrl"`
}

type RateRank struct {
	ID             int    `json:"id"`
	UserId         string `json:"user_id"`
	Avatar         string `json:"avatar"`
	Name           string `json:"name"`
	RateReturn     string `json:"rate_return"`
	RateReturnYear string `json:"rate_return_year"`
}
