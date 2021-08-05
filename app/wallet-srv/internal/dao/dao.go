package dao

import (
	"github.com/jinzhu/gorm"
	"wq-fotune-backend/app/wallet-srv/internal/model"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/dbclient"
)

type Dao struct {
	db *gorm.DB
}

func New() *Dao {
	return &Dao{
		db: dbclient.NewDB(env.DBDSN),
	}
}

func (d *Dao) CreateWallet(wallet *model.WqWallet) error {
	if err := d.db.Table(TableWqWallet).Create(wallet).Error; err != nil {
		logger.Warnf("用户创建钱包失败 %+v, %v", wallet, err)
		return err
	}
	return nil
}

func (d *Dao) UpdateWallet(wallet *model.WqWallet) error {
	if err := d.db.Table(TableWqWallet).Where("id=?", wallet.ID).Update(wallet).Error; err != nil {
		logger.Warnf("用户更新钱包失败 %+v, %v", wallet, err)
		return err
	}
	return nil
}

func (d *Dao) GetWalletByUserID(userId string) (*model.WqWallet, error) {
	wallet := &model.WqWallet{}
	if err := d.db.Table(TableWqWallet).Where("user_id = ?", userId).First(wallet).Error; err != nil {
		logger.Warnf("用户 %s 没有找到用户钱包的信息 %v", userId, err)
		return nil, err
	}
	return wallet, nil
}

func (d *Dao) GetWqCoinInfo() (*model.WqCoinInfo, error) {
	coinInfo := &model.WqCoinInfo{}
	if err := d.db.Table(TableWqCoinInfo).First(coinInfo).Error; err != nil {
		logger.Warnf("没有找到wqifc 币种的信息 %v", err)
		return nil, err
	}
	return coinInfo, nil
}

func (d *Dao) CreateTransferRecord(record *model.WqTransferRecord) error {
	if err := d.db.Table(TableWqTransferRecord).Create(record).Error; err != nil {
		logger.Warnf("用户创建划转记录失败 %+v, %v", record, err)
		return err
	}
	return nil
}

func (d *Dao) CreateWithdrawal(withdrawal *model.WqWithdrawal) error {
	if err := d.db.Table(TableWqWithdrawal).Create(withdrawal).Error; err != nil {
		logger.Warnf("创建提现记录失败 %+v", err)
		return err
	}
	return nil
}

func (d *Dao) CreateIfcGiftRecord(record *model.WqIfcGiftRecord) error {
	if err := d.db.Table(TableWqIfcGiftRecord).Create(record).Error; err != nil {
		logger.Warnf("创建ifc赠送记录失败 %v", err)
		return err
	}
	return nil
}

func (d *Dao) GetIfcGiftRecordByUid(userId string) (data []*model.WqIfcGiftRecord) {
	if err := d.db.Table(TableWqIfcGiftRecord).Where("user_id = ?", userId).Order("updated_at desc").Find(&data); err != nil {
		logger.Warnf("没找到赠送ifc数据 uid %s", userId)
		return
	}
	return
}

func (d *Dao) GetIfcGiftRecordBySql(userId, inUserId, exchange string) (data []*model.WqIfcGiftRecord) {
	query := "user_id = ? and in_user_id = ? and exchange = ?"
	if err := d.db.Table(TableWqIfcGiftRecord).Where(query, userId, inUserId, exchange).Find(&data); err != nil {
		return
	}
	return
}
