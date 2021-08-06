package dao

import (
	"github.com/jinzhu/gorm"
	"wq-fotune-backend/app/usercenter-srv/internal/model"
	"wq-fotune-backend/libs/cache"
	"wq-fotune-backend/libs/logger"
)

type Dao struct {
	db *gorm.DB
}

func New() *Dao {
	return &Dao{
		db: cache.Mysql(),
	}
}

// GetWqUserBaseByInCode 通过邀请码查询用户
func (d *Dao) GetWqUserBaseByInCode(inCode string) *model.WqUserBase {
	user := &model.WqUserBase{}
	if err := d.db.Table(TABLE_WQ_USER_BASE).Where("invitation_code = ?", inCode).First(user).Error; err != nil {
		logger.Warnf("GetWqUserBaseByInCode: %v invitation_code %s", err, inCode)
		return nil
	}
	return user
}

func (d *Dao) GetWqUserBaseByPhone(phone string) *model.WqUserBase {
	user := &model.WqUserBase{}
	if err := d.db.Table(TABLE_WQ_USER_BASE).Where("phone = ? and status = 1", phone).First(user).Error; err != nil {
		logger.Warnf("GetWqUserBaseByPhone: %v phone %s", err, phone)
		return nil
	}
	return user
}

func (d *Dao) GetWqUserBaseByUID(uid string) *model.WqUserBase {
	user := &model.WqUserBase{}
	if err := d.db.Table(TABLE_WQ_USER_BASE).Where("user_id = ? and status = 1", uid).First(user).Error; err != nil {
		logger.Warnf("GetWqUserBaseByUID: %v user_id %s", err, uid)
		return nil
	}
	return user
}

func (d *Dao) UpdateWqUserBaseByUID(user *model.WqUserBase) error {
	if err := d.db.Table(TABLE_WQ_USER_BASE).Where("user_id = ?", user.UserID).
		Update(user).Error; err != nil {
		logger.Warnf("UpdateWqUserBaseByUID: %v user_id %s", err, user.UserID)
		return err
	}
	return nil
}

func (d *Dao) CreateWqUserBase(db *gorm.DB, user *model.WqUserBase) error {
	if err := db.Table(TABLE_WQ_USER_BASE).Create(user).Error; err != nil {
		logger.Errorf("CreateWqUserInvite error %v ", err)
		return err
	}
	return nil
}

// 创建邀请码关联用户表
func (d *Dao) CreateWqUserInvite(db *gorm.DB, userMasterID, userInvitedID string) error {
	userInvite := &model.WqUserInvite{
		UserID:        userMasterID,
		InvitedUserID: userInvitedID,
	}
	if err := db.Table(TABLE_WQ_USER_INVITE).Create(userInvite).Error; err != nil {
		logger.Errorf("CreateWqUserInvite error %v userMasterID%s userInvitedID%s", err, userMasterID, userInvitedID)
		return err
	}
	return nil
}

func (d *Dao) GetUserMasterByInUserId(InUserId string) *model.WqUserInvite {
	data := &model.WqUserInvite{}
	if err := d.db.Table(TABLE_WQ_USER_INVITE).Where("in_user_id = ?", InUserId).First(data).Error; err != nil {
		logger.Warnf("查找邀请记录失败 InUserId %s err %v", InUserId, err)
		return nil
	}
	return data
}

func (d *Dao) GetAllUsers() (users []*model.WqUserBase) {
	if err := d.db.Table(TABLE_WQ_USER_BASE).Where("status=?", 1).Scan(&users).Error; err != nil {
		logger.Warnf("查找所有用户失败")
	}
	return
}
