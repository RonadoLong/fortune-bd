package response

import "github.com/micro/go-micro/v2/errors"

func NewLoginPasswordOrPhoneErrMsg(id string) error {
	return errors.New(id, "手机号或者密码错误", ERROR_CODE_WRONG_PARAM)
}

func NewLoginReqFreqErrMsg(id string) error {
	return errors.New(id, "请求时间获取频繁 请间隔15s！", ERROR_CODE_Max_Req)
}

func NewLoginReqMaxErrMsg(id string) error {
	return errors.New(id, "今日调用次数已达上限！", ERROR_CODE_Max_Req)
}

func NewInternalServerErrMsg(id string) error {
	return errors.New(id, "内部服务错误！", ERROR_INTERNAL_SERVER)
}

func NewValidateCodeExpireErrMsg(id string) error {
	return errors.New(id, "验证码已失效！", ERROR_CODE_WRONG_PARAM)
}

func NewValidateCodeErrMsg(id string) error {
	return errors.New(id, "验证码填写错误！", ERROR_CODE_WRONG_PARAM)
}

func NewInvitationCodeErrMsg(id string) error {
	return errors.New(id, "邀请码无效！", ERROR_CODE_WRONG_PARAM)
}

func NewPhoneHasRegisterErrMsg(id string) error {
	return errors.New(id, "该手机号已被注册！", ERROR_CODE_CREATE)
}

func NewUserCreateErrMsg(id string) error {
	return errors.New(id, "用户创建失败！", ERROR_CODE_CREATE)
}

func NewUserNotFoundErrMsg(id string) error {
	return errors.New(id, "用户不存在！", ERROR_CODE_NOT_FOUND)
}

func NewUserSetPassErrMsg(id string) error {
	return errors.New(id, "修改用户密码失败！", ERROR_CODE_UPDATE)
}

func NewUserUpdateBaseErrMsg(id string) error {
	return errors.New(id, "修改用户信息失败！", ERROR_CODE_UPDATE)
}

func NewCarouselNotFoundErrMsg(id string) error {
	return errors.New(id, "图片没有找到！", ERROR_CODE_NOT_FOUND)
}

func NewContractNotFoundErrMsg(id string) error {
	return errors.New(id, "暂无联系方式！", ERROR_CODE_NOT_FOUND)
}

func NewExchangeNotFoundErrMsg(id string) error {
	return errors.New(id, "无交易所选项！", ERROR_CODE_NOT_FOUND)
}

func NewExchangeApiDuplicateErrMsg(id string) error {
	return errors.New(id, "同名apikey 已存在！", ERROR_CODE_CREATE)
}

func NewExchangeApiCreateErrMsg(id string) error {
	return errors.New(id, "apikey 创建失败！", ERROR_CODE_CREATE)
}

func NewCreateStrategyErrMsg(id, msg string) error {
	return errors.New(id, msg, ERROR_CODE_CREATE)
}

func NewExchangeApiCheckErrMsg(id string) error {
	return errors.New(id, "账户不可用,请填写真实账户！", ERROR_CODE_CREATE)
}

func NewExchangePassphraseNoneErrMsg(id string) error {
	return errors.New(id, "Passphrase 不可为空！", ERROR_CODE_WRONG_PARAM)
}

func NewExchangeApiListErrMsg(id string) error {
	return errors.New(id, "尚未绑定任何apikey！", ERROR_CODE_NOT_FOUND)
}

func NewUpdateExchangeApiErrMsg(id string) error {
	return errors.New(id, "更新apiKey失败!", ERROR_CODE_UPDATE)
}

func NewDeleteExchangeApiErrMsg(id string) error {
	return errors.New(id, "删除apiKey失败!", ERROR_CODE_DELETE)
}

func NewDeleteExchangeApiNotFoundErrMsg(id string) error {
	return errors.New(id, "该apiKey不存在!", ERROR_CODE_DELETE)
}

func NewTradeNotFoundErrMsg(id string) error {
	return errors.New(id, "暂无交易记录!", ERROR_CODE_NOT_FOUND)
}

func NewProfitNotFoundErrMsg(id string) error {
	return errors.New(id, "暂无盈亏记录!", ERROR_CODE_NOT_FOUND)
}

func NewUserStrategyNotFoundErrMsg(id string) error {
	return errors.New(id, "策略列表为空!", ERROR_CODE_NOT_FOUND)
}

func NewUserStrategyDetailErrMsg(id string) error {
	return errors.New(id, "查找该策略出错!", ERROR_CODE_NOT_FOUND)
}

func NewStrategyNotFoundErrMsg(id string) error {
	return errors.New(id, "暂无策略上架!", ERROR_CODE_NOT_FOUND)
}

func NewSetUserStrategyApiKeyErrMsg(id string) error {
	return errors.New(id, "更新apikey失败!", ERROR_CODE_UPDATE)
}

func NewSetApiKeySameErrMsg(id string) error {
	return errors.New(id, "新apikey与旧apikey相同!", ERROR_CODE_UPDATE)
}

func NewSetBalanceSameErrMsg(id string) error {
	return errors.New(id, "资金没有变更!", ERROR_CODE_UPDATE)
}

func NewSetBalanceErrMsg(id string) error {
	return errors.New(id, "资金更新失败!", ERROR_CODE_UPDATE)
}

func NewGetStrategyNotFoundErrMsg(id string) error {
	return errors.New(id, "策略不存在!", ERROR_CODE_NOT_FOUND)
}

func NewCreateUserStrategyErrMsg(id string) error {
	return errors.New(id, "购买策略失败!", ERROR_CODE_CREATE)
}

func NewEvaluationStrategyErrMsg(id string) error {
	return errors.New(id, "找不到订单绑定的策略!", ERROR_CODE_NOT_FOUND)
}

func NewExchangePosErrMsg(id string) error {
	return errors.New(id, "暂无数据,请检查您的apiKey是否失效!", ERROR_CODE_NOT_FOUND)
}

func NewDeleteApiHasStrategyErrMsg(id string) error {
	return errors.New(id, "该apiKey仍有正在运行的策略！", ERROR_CODE_NOT_FOUND)
}

func NewUpdateApiHasStrategyErrMsg(id string) error {
	return errors.New(id, "旧apiKey仍有正在运行的策略！", ERROR_CODE_UPDATE)
}

func NewExchangeApiExpireErrMsg(id string) error {
	return errors.New(id, "api账户已经失效或被删除", ERROR_CODE_NOT_FOUND)
}

func NewUserStrategyBalanceErrMsg(id string) error {
	return errors.New(id, "未设置初始资金", ERROR_CODE_NOT_FOUND)
}

func NewUserStrategyRunErrMsg(id string) error {
	return errors.New(id, "启动策略失败", ERROR_CODE_UPDATE)
}

func NewUserStrategyRunALReadyErrMsg(id string) error {
	return errors.New(id, "策略已经启动", ERROR_CODE_UPDATE)
}

func NewWqMembersNotFound(id string) error {
	return errors.New(id, "没有相关会员可购买", ERROR_CODE_NOT_FOUND)
}

func NewWqPayNotFound(id string) error {
	return errors.New(id, "无可用支付通道", ERROR_CODE_NOT_FOUND)
}

func NewDataNotFound(id, msg string) error {
	return errors.New(id, msg, ERROR_CODE_NOT_FOUND)
}

func NewInternalServerErrWithMsg(id, msg string) error {
	return errors.New(id, msg, ERROR_INTERNAL_SERVER)
}

func NewParamsErrWithMsg(id, msg string) error {
	return errors.New(id, msg, ERROR_CODE_WRONG_PARAM)
}
