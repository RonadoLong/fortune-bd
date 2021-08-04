package model

import "time"

type WqCommonCarousel struct {
	ID        int       `gorm:"column:id" json:"id"`
	Image     string    `gorm:"column:image" json:"image"`
	ClickUrl  string    `gorm:"column:click_url" json:"click_url"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type WqCommonContact struct {
	ID        int       `gorm:"column:id" json:"id"`
	Image     string    `gorm:"column:image" json:"image"`
	Content   string    `gorm:"column:content" json:"content"`
	Contact   string    `gorm:"column:contact" json:"contact"`
	Platform  string    `gorm:"column:platform" json:"platform"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type WqAppVersion struct {
	ID          int32     `gorm:"column:id" json:"id"`
	HasUpdate   bool      `gorm:"column:has_update" json:"hasUpdate"`
	IsIgnorable bool      `gorm:"column:is_ignorable" json:"isIgnorable"`
	VersionCode int32     `gorm:"column:version_code" json:"versionCode"`
	VersionName string    `gorm:"column:version_name" json:"versionName"`
	UpdateLog   string    `gorm:"column:update_log" json:"updateLog"`
	ApkUrl      string    `gorm:"column:apk_url" json:"apkUrl"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	Platform    string    `gorm:"column:platform" json:"platform"`
}
