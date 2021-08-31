package response

import "github.com/go-kratos/kratos/v2/errors"

func NewLoginPasswordOrPhoneErrMsg(id string) error {
	return errors.New(ERROR_CODE_WRONG_PARAM, "手机号或者密码错误", id)
}

func NewLoginReqFreqErrMsg(id string) error {
	return errors.New(ERROR_CODE_Max_Req, "请求时间获取频繁 请间隔15s！", id)
}

func NewLoginReqMaxErrMsg(id string) error {
	return errors.New(ERROR_CODE_Max_Req, "今日调用次数已达上限！", id)
}

func NewInternalServerErrMsg(id string) error {
	return errors.New(ERROR_INTERNAL_SERVER, "内部服务错误！", id)
}

func NewValidateCodeExpireErrMsg(id string) error {
	return errors.New(ERROR_CODE_WRONG_PARAM, "验证码已失效！", id)
}

func NewValidateCodeErrMsg(id string) error {
	return errors.New(ERROR_CODE_WRONG_PARAM, "验证码填写错误！", id)
}

func NewInvitationCodeErrMsg(id string) error {
	return errors.New(ERROR_CODE_WRONG_PARAM, "邀请码无效！", id)
}

func NewPhoneHasRegisterErrMsg(id string) error {
	return errors.New(ERROR_CODE_CREATE, "该手机号已被注册！", id)
}

func NewUserCreateErrMsg(id string) error {
	return errors.New(ERROR_CODE_CREATE, "用户创建失败！", id)
}

func NewUserNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "用户不存在！", id)
}

func NewUserSetPassErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "修改用户密码失败！", id)
}

func NewUserUpdateBaseErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "修改用户信息失败！", id)
}

func NewCarouselNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "图片没有找到！", id)
}

func NewContractNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "暂无联系方式！", id)
}

func NewExchangeNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "无交易所选项！", id)
}

func NewExchangeApiDuplicateErrMsg(id string) error {
	return errors.New(ERROR_CODE_CREATE, "同名apikey 已存在！", id)
}

func NewExchangeApiCreateErrMsg(id string) error {
	return errors.New(ERROR_CODE_CREATE, "apikey 创建失败！", id)
}

func NewCreateStrategyErrMsg(id, msg string) error {
	return errors.New(ERROR_CODE_CREATE, msg, id)
}

func NewExchangeApiCheckErrMsg(id string) error {
	return errors.New(ERROR_CODE_CREATE, "账户不可用,请填写真实账户！", id)
}

func NewExchangePassphraseNoneErrMsg(id string) error {
	return errors.New(ERROR_CODE_WRONG_PARAM, "Passphrase 不可为空！", id)
}

func NewExchangeApiListErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "尚未绑定任何apikey！", id)
}

func NewUpdateExchangeApiErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "更新apiKey失败!", id)
}

func NewDeleteExchangeApiErrMsg(id string) error {
	return errors.New(ERROR_CODE_DELETE, "删除apiKey失败!", id)
}

func NewDeleteExchangeApiNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_DELETE, "该apiKey不存在!", id)
}

func NewTradeNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "暂无交易记录!", id)
}

func NewProfitNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "暂无盈亏记录!", id)
}

func NewUserStrategyNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "策略列表为空!", id)
}

func NewUserStrategyDetailErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "查找该策略出错!", id)
}

func NewStrategyNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "暂无策略上架!", id)
}

func NewSetUserStrategyApiKeyErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "更新apikey失败!", id)
}

func NewSetApiKeySameErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "新apikey与旧apikey相同!", id)
}

func NewSetBalanceSameErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "资金没有变更!", id)
}

func NewSetBalanceErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "资金更新失败!", id)
}

func NewGetStrategyNotFoundErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "策略不存在!", id)
}

func NewCreateUserStrategyErrMsg(id string) error {
	return errors.New(ERROR_CODE_CREATE, "购买策略失败!", id)
}

func NewEvaluationStrategyErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "找不到订单绑定的策略!", id)
}

func NewExchangePosErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "暂无数据,请检查您的apiKey是否失效!", id)
}

func NewDeleteApiHasStrategyErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "该apiKey仍有正在运行的策略！", id)
}

func NewUpdateApiHasStrategyErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "旧apiKey仍有正在运行的策略！", id)
}

func NewExchangeApiExpireErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "api账户已经失效或被删除", id)
}

func NewUserStrategyBalanceErrMsg(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "未设置初始资金", id)
}

func NewUserStrategyRunErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "启动策略失败", id)
}

func NewUserStrategyRunALReadyErrMsg(id string) error {
	return errors.New(ERROR_CODE_UPDATE, "策略已经启动", id)
}

func NewWqMembersNotFound(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "没有相关会员可购买", id)
}

func NewWqPayNotFound(id string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, "无可用支付通道", id)
}

func NewDataNotFound(id, msg string) error {
	return errors.New(ERROR_CODE_NOT_FOUND, msg, id)
}

func NewInternalServerErrWithMsg(id, msg string) error {
	return errors.New(ERROR_INTERNAL_SERVER, msg, id)
}

func NewParamsErrWithMsg(id, msg string) error {
	return errors.New(ERROR_CODE_WRONG_PARAM, msg, id)
}
