package response

import "errors"

type LoginReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type ValidateCodeReq struct {
	Phone string `json:"phone"`
}

type CheckValidateCodeReq struct {
	Phone        string `json:"phone"`
	ValidateCode string `json:"validate_code"`
}

func (req *CheckValidateCodeReq) CheckNotNull() error {
	switch "" {
	case req.Phone:
		return errors.New("手机号不能为空")
	case req.ValidateCode:
		return errors.New("验证码不能为空")
	}
	return nil
}

type RegisterReq struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	InvitationCode  string `json:"invitation_code"`
	ValidateCode    string `json:"validate_code"`
}

func (req *RegisterReq) CheckBaseParam() error {
	switch "" {
	case req.Phone:
		return errors.New("手机号不能为空")
	case req.Password:
		return errors.New("密码不能为空")
	case req.ConfirmPassword:
		return errors.New("确认密码不能为空")
	case req.ValidateCode:
		return errors.New("验证码不能为空")
	}
	if req.Password != req.ConfirmPassword {
		return errors.New("两次输入密码不一致")
	}
	return nil
}

// UpdateUserBaseReq 更新用户基础信息
type UpdateUserBaseReq struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (req *UpdateUserBaseReq) CheckNotNull() error {
	NotNullPlease := "请勿填空值, 至少修改一项用户基础信息"
	Null := ""
	if Null == req.Name && Null == req.Avatar {
		return errors.New(NotNullPlease)
	}
	return nil
}

type ChangePasswordReq struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (req *ChangePasswordReq) CheckPassword() error {
	switch "" {
	case req.Password:
		return errors.New("输入密码不能为空")
	case req.ConfirmPassword:
		return errors.New("确认密码不能为空")
	}
	if req.Password != req.ConfirmPassword {
		return errors.New("两次输入密码不一致")
	}
	return nil
}

type ForgetPasswordReq struct {
	Phone        string `json:"phone"`
	ValidateCode string `json:"validate_code"`
	ChangePasswordReq
}

func (req *ForgetPasswordReq) CheckPhoneVCode() error {
	switch "" {
	case req.Phone:
		return errors.New("手机号不能为空")
	case req.ValidateCode:
		return errors.New("验证码不能为空")
	}
	return nil
}
