package client

import (
	"context"
	"testing"
	fotune_srv_user "wq-fotune-backend/api/usercenter"
	"wq-fotune-backend/libs/env"
)

func TestNewUserClientClient(t *testing.T) {
	client := NewUserClient(env.EtcdAddr)
	ctx := context.Background()
	req := &fotune_srv_user.LoginReq{
		Phone:    "454545",
		Password: "455",
	}

	resp, err := client.Login(ctx, req)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(resp)
}
