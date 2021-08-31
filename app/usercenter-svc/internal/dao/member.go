package dao

import (
    "fortune-bd/app/usercenter-svc/internal/model"
    "fortune-bd/libs/logger"
)

func (d *Dao) GetMembersWithState(state int32) (wqMembers []*model.WqMember) {
    if err := d.db.Table(TABLE_WQ_MEMBER).Where("state = ?", state).Find(&wqMembers).Error; err != nil {
        logger.Warnf("GetMembersWithState has err %v", err)
        return
    }
    return
}

func (d *Dao) GetPaymentWithState(state int32) (wqPayment []*model.WqPayment) {
    if err := d.db.Table(TABLE_WQ_PAYMENT).Where("state = ?", state).Find(&wqPayment).Error; err != nil {
        logger.Warnf("GetPaymentWithState has err %v", err)
        return
    }
    return
}