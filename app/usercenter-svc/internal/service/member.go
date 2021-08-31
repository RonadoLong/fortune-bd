package service

import (
	"context"
	"fortune-bd/api/response"
	pb "fortune-bd/api/usercenter/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetMembers 获取会员
func (s *UserService) GetMembers(ctx context.Context, req *emptypb.Empty) (*pb.GetMembersResp, error) {
	var resp = new(pb.GetMembersResp)
	wqMembers := s.dao.GetMembersWithState(1)
	if len(wqMembers) == 0 {
		return nil, response.NewWqMembersNotFound(errID)
	}
	for _, value := range wqMembers {
		member := &pb.Member{
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
	return resp, nil
}

// GetPaymentMethod 获取支付方法
func (s *UserService) GetPaymentMethod(ctx context.Context, req *emptypb.Empty) (*pb.GetPaymentMethodResp, error) {
	var resp = new(pb.GetPaymentMethodResp)
	wqPayment := s.dao.GetPaymentWithState(1)
	if len(wqPayment) == 0 {
		return nil,response.NewWqPayNotFound(errID)
	}
	for _, value := range wqPayment {
		payment := &pb.Payment{
			Id:      value.ID,
			Name:    value.Name,
			Remark:  value.Remark,
			BitAddr: value.BitAddr,
			BitCode: value.BitCode,
			State:   value.State,
		}
		resp.Payments = append(resp.Payments, payment)
	}
	return resp, nil
}
