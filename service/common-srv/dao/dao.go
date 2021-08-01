package dao

import (
	"github.com/jinzhu/gorm"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/dbclient"
	"wq-fotune-backend/service/common-srv/model"
)

type Dao struct {
	db *gorm.DB
}

func New() *Dao {
	return &Dao{
		db: dbclient.NewDB(env.DBDSN),
	}
}

func (d *Dao) GetCarousels() (carousels []*model.WqCommonCarousel) {
	if err := d.db.Table(TABLE_WQ_COMMON_CAROUSEL).Find(&carousels).Error; err != nil {
		logger.Warnf("GetCarousels err %s", err.Error)
		return
	}
	return
}

func (d *Dao) GetContact() (*model.WqCommonContact, error) {
	contract := &model.WqCommonContact{}
	if err := d.db.Table(TABLE_WQ_COMMON_CONTACT).First(&contract).Error; err != nil {
		logger.Warnf("GetContact err %s", err.Error)
		return nil, err
	}
	return contract, nil
}

func (d *Dao) GetAppVersion(platform string) (*model.WqAppVersion, error) {
	appVersion := &model.WqAppVersion{}
	if err := d.db.Table(TABLE_WQ_APPVERSION).Where("platform = ?", platform).First(&appVersion).Error; err != nil {
		logger.Warnf("GetAppVersion err %v", err)
		return nil, err
	}
	return appVersion, nil
}
