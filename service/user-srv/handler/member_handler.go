package handler

import (
	"context"
	"wq-fotune-backend/pkg/response"
	fotune_srv_user "wq-fotune-backend/service/user-srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
)

func (u *UserHandler) GetMembers(ctx context.Context, req *empty.Empty, resp *fotune_srv_user.GetMembersResp) error {
	wqMembers := u.dao.GetMembersWithState(1)
	if len(wqMembers) == 0 {
		return response.NewWqMembersNotFound(errID)
	}
	for _, value := range wqMembers {
		member := &fotune_srv_user.Member{
			Id:        value.ID,
			Name:      value.Name,
			Remark:    value.Remark,
			Price:     value.Price,
			OldPrice:  value.OldPrice,
			Duration:  value.Duration,
			State:     value.State,
			CreatedAt: value.CreatedAt.String(),
			UpdatedAt: value.UpdatedAt.String(),
		}
		resp.Members = append(resp.Members, member)
	}
	return nil
}

func (u *UserHandler) GetPaymentMethod(ctx context.Context, req *empty.Empty, resp *fotune_srv_user.GetPaymentMethodResp) error {
	wqPayment := u.dao.GetPaymentWithState(1)
	if len(wqPayment) == 0 {
		return response.NewWqPayNotFound(errID)
	}
	for _, value := range wqPayment {
		payment := &fotune_srv_user.Payment{
			Id:      value.ID,
			Name:    value.Name,
			Remark:  value.Remark,
			BitAddr: value.BitAddr,
			BitCode: value.BitCode,
			State:   value.State,
		}
		resp.Payments = append(resp.Payments, payment)
	}
	return nil
}
