package service

import (
	"context"
	"errors"
	"time"
	"wq-fotune-backend/api/response"
	fotune_srv_user "wq-fotune-backend/api/usercenter"
	walletPb "wq-fotune-backend/api/wallet"
	"wq-fotune-backend/app/usercenter-srv/internal/model"
	"wq-fotune-backend/libs/bcrypt2"
	"wq-fotune-backend/libs/limitReq"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/message"

	validate_code "wq-fotune-backend/libs/validate-code"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserService) GetAllUserInfo(ctx context.Context, req *empty.Empty, resp *fotune_srv_user.AllUserInfoResp) error {
	users := u.dao.GetAllUsers()
	for _, user := range users {
		lastLogin, _ := ptypes.TimestampProto(user.LastLogin)
		resp.UserInfo = append(resp.UserInfo, &fotune_srv_user.LoginResp{
			UserID:         user.UserID,
			InvitationCode: user.InvitationCode,
			Name:           user.Name,
			Avatar:         user.Avatar,
			Phone:          user.Phone,
			LoginCount:     int32(user.LoginCount),
			LastLogin:      lastLogin,
		})
	}
	return nil
}

func (u *UserService) GetUserInfo(ctx context.Context, req *fotune_srv_user.UserInfoReq, resp *fotune_srv_user.LoginResp) error {
	user := u.dao.GetWqUserBaseByUID(req.GetUserID())
	if user == nil {
		return response.NewUserNotFoundErrMsg(errID)
	}
	resp.UserID = user.UserID
	resp.Name = user.Name
	resp.Avatar = user.Avatar
	resp.Phone = user.Phone
	resp.InvitationCode = user.InvitationCode
	resp.LoginCount = int32(user.LoginCount)
	resp.LastLogin, _ = ptypes.TimestampProto(user.LastLogin)
	return nil
}

func (u *UserService) SendValidateCode(ctx context.Context, req *fotune_srv_user.ValidateCodeReq, resp *fotune_srv_user.ValidateCodeResp) error {
	ipCount := limitReq.GetReqCount(req.Phone)
	if ipCount > 0 {
		return response.NewLoginReqFreqErrMsg(errID)
	}
	if err := limitReq.SetReqCount(req.Phone, ipCount+1); err != nil {
		logger.Errorf("limitReq.SetReqCount %v %s", err, req.Phone)
		return response.NewInternalServerErrMsg(errID)
	}
	count, err := validate_code.CheckCount(req.Phone)
	if err != nil { // 禁止调用次数超出
		return response.NewLoginReqMaxErrMsg(errID)
	}
	code := validate_code.GenValidateCode(6)
	logger.Infof("手机验证码 %s", code)
	timeout := 10 * time.Minute
	if err := validate_code.SaveValidateCode(req.Phone, code, count+1, timeout); err != nil {
		logger.Errorf("SendValidateCode 保存手机验证码到redis错误 :%v", err)
		return response.NewInternalServerErrMsg(errID)
	}
	err = message.SendMsg(req.Phone, code)
	if err != nil {
		logger.Errorf("send code server err :%v", err)
		return response.NewInternalServerErrMsg(errID)
	}
	resp.Code = code
	return nil
}

func (u *UserService) Login(ctx context.Context, req *fotune_srv_user.LoginReq, resp *fotune_srv_user.LoginResp) error {
	user := u.dao.GetWqUserBaseByPhone(req.Phone)
	if user == nil {
		return response.NewLoginPasswordOrPhoneErrMsg(errID)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return response.NewLoginPasswordOrPhoneErrMsg(errID)
	}
	resp.UserID = user.UserID
	resp.Phone = user.Phone
	resp.Name = user.Name
	resp.Avatar = user.Avatar
	resp.InvitationCode = user.InvitationCode
	resp.LastLogin, _ = ptypes.TimestampProto(user.LastLogin)
	resp.LoginCount = int32(user.LoginCount)
	return nil
}

func (u *UserService) Register(ctx context.Context, req *fotune_srv_user.RegisterReq, resp *empty.Empty) error {

	vcode, err := validate_code.GetValidateCode(req.Phone)
	if err != nil {
		logger.Warnf("GetValidateCode 查验证码失败 %v", err)
		return response.NewValidateCodeExpireErrMsg(errID)
	}
	if req.ValidateCode != vcode {
		return response.NewValidateCodeErrMsg(errID)
	}
	var userMaster *model.WqUserBase
	if req.InvitationCode != "" {
		//主人用户 发出邀请码的用户
		userMaster = u.dao.GetWqUserBaseByInCode(req.InvitationCode)
		if userMaster == nil {
			return response.NewInvitationCodeErrMsg(errID)
		}
	}
	if user := u.dao.GetWqUserBaseByPhone(req.Phone); user != nil {
		return response.NewPhoneHasRegisterErrMsg(errID)
	}
	dbclt := u.dao.BeginTran()
	user := model.NewWqUserBase(req.GetPhone(), req.GetPassword())
	if err := u.dao.CreateWqUserBase(dbclt, user); err != nil {
		dbclt.Rollback()
		return response.NewUserCreateErrMsg(errID)
	}
	//创建钱包
	_, err = u.walletSrv.CreateWallet(context.Background(), &walletPb.UidReq{UserId: user.UserID})
	if err != nil {
		logger.Warnf("用户%s 注册时,创建钱包失败 %v", user.UserID, err)
	}
	if userMaster == nil {
		dbclt.Commit()
		return nil
	}
	if err := u.dao.CreateWqUserInvite(dbclt, userMaster.UserID, user.UserID); err != nil {
		dbclt.Rollback()
		return response.NewUserCreateErrMsg(errID)
	}
	//给发出邀请的用户增加ifc
	if err := u.AddIfcBalance(userMaster.UserID, user.UserID, "", "register", 5.0); err != nil {
		logger.Warnf("注册时增加邀请码用户的ifc失败 uid %s userMasterId %s err %v", user.UserID, userMaster.UserID, err)
	}
	dbclt.Commit()
	//delete validateCode
	validate_code.DeleteValidateCode(req.Phone)
	return nil
}

func (u *UserService) AddIfcBalance(userMasterId, inUserID, exchange, _type string, volume float64) error {
	_, err := u.walletSrv.AddIfcBalance(context.Background(), &walletPb.AddIfcBalanceReq{
		UserMasterId: userMasterId,
		InUserId:     inUserID,
		Volume:       volume,   //手数
		Type:         _type,    //register api  strategy
		Exchange:     exchange, //注册这里留空
	})
	return err
}

func (u *UserService) ResetPassword(ctx context.Context, req *fotune_srv_user.ChangePasswordReq, resp *empty.Empty) error {
	user := u.dao.GetWqUserBaseByUID(req.UserID)
	if user == nil {
		return response.NewUserNotFoundErrMsg(errID)
	}
	userUpdateField := &model.WqUserBase{UserID: user.UserID, Password: bcrypt2.CryptPassword(req.Password),
		UpdatedAt: time.Now()}
	if err := u.dao.UpdateWqUserBaseByUID(userUpdateField); err != nil {
		return response.NewUserSetPassErrMsg(errID)
	}
	return nil
}

func (u *UserService) GetUserInfoByPhone(phone string) (*model.WqUserBase, error) {
	user := u.dao.GetWqUserBaseByPhone(phone)
	if user == nil {
		return nil, errors.New("该手机没有绑定用户")
	}
	return user, nil
}

func (u *UserService) ForgetPassword(ctx context.Context, req *fotune_srv_user.ForgetPasswordReq, resp *empty.Empty) error {
	vcode, err := validate_code.GetValidateCode(req.Phone)
	if err != nil {
		logger.Warnf("GetValidateCode 查验证码失败 %v key %s", err)
		return response.NewValidateCodeExpireErrMsg(errID)
	}
	if req.ValidateCode != vcode {
		return response.NewValidateCodeErrMsg(errID)
	}
	user, err := u.GetUserInfoByPhone(req.Phone)
	if err != nil {
		return response.NewUserNotFoundErrMsg(errID)
	}

	changePassReq := &fotune_srv_user.ChangePasswordReq{
		UserID:          user.UserID,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}
	if err := u.ResetPassword(context.Background(), changePassReq, nil); err != nil {
		return response.NewUserSetPassErrMsg(errID)
	}
	validate_code.DeleteValidateCode(req.Phone)
	return nil
}

func (u *UserService) UpdateUser(ctx context.Context, req *fotune_srv_user.UpdateUserReq, resp *empty.Empty) error {
	user := u.dao.GetWqUserBaseByUID(req.GetUserId())
	if user == nil {
		return response.NewUserNotFoundErrMsg(errID)
	}
	switch {
	case req.Avatar != "":
		user.Avatar = req.Avatar
	case req.Name != "":
		user.Name = req.Name
	}
	user.UpdatedAt = time.Now()
	if err := u.dao.UpdateWqUserBaseByUID(user); err != nil {
		return response.NewUserUpdateBaseErrMsg(errID)
	}
	return nil
}

func (u *UserService) GetUserMasterByInViteUser(ctx context.Context, req *fotune_srv_user.GetUserMasterReq, resp *fotune_srv_user.UserMasterResp) error {
	data := u.dao.GetUserMasterByInUserId(req.InviteUid)
	if data == nil {
		return response.NewDataNotFound(errID, "没有找到邀请数据")
	}
	resp.UserMasterId = data.UserID
	return nil
}
