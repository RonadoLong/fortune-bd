package handler

import (
	"github.com/jinzhu/gorm"
	"shop-micro/service/user-service/proto"
)

const (
	WECHAT = "wechat"
	FACEBOOK = "facebook"
	PHONE = "phone"
)

type userRepository struct {
	DB *gorm.DB
}

func (repository *userRepository) login(req *shop_srv_user.LoginReq, resp *shop_srv_user.UserResp) error {

	if req.Type == WECHAT || req.Type == FACEBOOK {
		
	} else if req.Type == PHONE {

	}
	return nil
}


