package handler

import (
	"context"
	"shop-micro/service/user-service/proto"
	"testing"
	"time"
)

func TestUserHandler_Login(t *testing.T) {
	handler, err := NewUserHandler()
	if err != nil {
		t.Log(err)
	}

		time.Sleep(time.Millisecond * 500)
		req := &shop_srv_user.LoginReq{Type:FACEBOOK}

		resp := &shop_srv_user.UserResp{RealName:"ssss"}
		_ = handler.Login(context.Background(), req, resp)

	//for {
	//
	//	time.Sleep(time.Millisecond * 500)
	//	req := &shop_srv_user.LoginReq{Type:FACEBOOK}
	//
	//	resp := &shop_srv_user.UserResp{RealName:"ssss"}
	//	_ = handler.Login(context.Background(), req, resp)
	//
	//}

}