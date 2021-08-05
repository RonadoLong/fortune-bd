package dao

import (
	"github.com/jinzhu/gorm"
	"wq-fotune-backend/app/exchange-srv/internal/model"
	"wq-fotune-backend/libs/logger"
)

func (d *Dao) GetExchangeInfo() (exchangeList []*model.WqExchange) {
	if err := d.db.Table(TABLE_WQ_EXCHANGE).Where("status = 1").Find(&exchangeList).Error; err != nil {
		logger.Warnf("GetExchangeInfo has err %v", err)
		return
	}
	return
}

func (d *Dao) GetExchangeById(id int64) (*model.WqExchange, error) {
	exchange := &model.WqExchange{}
	if err := d.db.Table(TABLE_WQ_EXCHANGE).Where("id = ?", id).First(exchange).Error; err != nil {
		logger.Warnf("GetExchangeInfo has err %v  exchangeId %v", err, id)
		return nil, err
	}
	return exchange, nil
}

func (d *Dao) AddExchangeApi(exApi *model.WqExchangeApi) error {
	if err := d.db.Table(TABLE_WQ_EXCHANGE_API).Create(&exApi).Error; err != nil {
		logger.Warnf("AddExchangeApi has err %v uid %s", err, exApi.UserID)
		return err
	}
	return nil
}

func (d *Dao) GetExchangeApiByUidAndApi(uid, apiKey string) (*model.WqExchangeApi, error) {
	exApi := &model.WqExchangeApi{}
	if err := d.db.Table(TABLE_WQ_EXCHANGE_API).Where("user_id=? and api_key= ?", uid, apiKey).
		First(exApi).Error; err != nil {
		logger.Warnf("GetExchangeApiByUID has err %v %v", err, uid)
		return nil, err
	}
	return exApi, nil
}

func (d *Dao) GetExchangeApiByUidAndExID(uid string, exchangeId int64) (*model.WqExchangeApi, error) {
	exApi := &model.WqExchangeApi{}
	if err := d.db.Table(TABLE_WQ_EXCHANGE_API).Where("user_id=? and exchange_id= ?", uid, exchangeId).
		First(exApi).Error; err != nil {
		logger.Warnf("GetExchangeApiByUidAndExID has err %v %v", err, uid)
		return nil, err
	}
	return exApi, nil
}

func (d *Dao) GetExchangeApiListByUid(uid string) (ApiList []*model.WqExchangeApi) {
	if err := d.db.Table(TABLE_WQ_EXCHANGE_API).Select("id, user_id,exchange_name, exchange_id, api_key, created_at, updated_at").
		Where("user_id = ?", uid).Order("exchange_id desc").Find(&ApiList).Error; err != nil {
		logger.Warnf("GetExchangeApiListByUid has err %v", err.Error())
		return
	}
	return
}

func (d *Dao) GetExchangeApiListByUidAndPlatform(uid, exName string) (ApiList []*model.WqExchangeApi) {
	if err := d.db.Table(TABLE_WQ_EXCHANGE_API).Select("id, user_id, exchange_name, exchange_id, passphrase, api_key, secret, created_at, updated_at").
		Where("user_id = ? and exchange_name = ?", uid, exName).Find(&ApiList).Error; err != nil {
		logger.Warnf("GetExchangeApiListByUidAndPlatform has err %v", err.Error())
		return
	}
	return
}

func (d *Dao) UpdateExchangeApi(api *model.WqExchangeApi) error {
	db := d.db.Table(TABLE_WQ_EXCHANGE_API).Where("user_id = ? and id = ?", api.UserID, api.ID).
		Update(api)

	if db.Error != nil {
		logger.Errorf("UpdateExchangeApi has err %v ", db.Error)
		return db.Error
	}

	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *Dao) GetExchangeApiByID(id int64) (*model.WqExchangeApi, error) {
	api := &model.WqExchangeApi{}
	if err := d.db.Table(TABLE_WQ_EXCHANGE_API).Where("id = ?", id).First(api).Error; err != nil {
		logger.Warnf("GetExchangeApiByID has err %v", err)
		return nil, err
	}
	return api, nil
}

func (d *Dao) DeleteExchangeApi(uid string, id int64) error {
	db := d.db.Table(TABLE_WQ_EXCHANGE_API).Delete(&model.WqExchangeApi{}, "user_id=? and id = ?", uid, id)
	if err := db.Error; err != nil {
		logger.Errorf("DeleteExchangeApi has err %v", err)
		return err
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
