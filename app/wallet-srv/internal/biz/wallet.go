package biz

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"strings"
	"time"
	"wq-fotune-backend/api/response"
	"wq-fotune-backend/app/wallet-srv/cache"
	"wq-fotune-backend/app/wallet-srv/internal/model"
	"wq-fotune-backend/libs/exchangeclient"
	"wq-fotune-backend/libs/helper"
	"wq-fotune-backend/libs/logger"
)

const (
	IFC           = "IFC"
	USDT          = "USDT"
	TYPE_REGISTER = "register"
	TYPE_BIND_API = "api"
	TYPE_STRATEGY = "strategy"
)

func (w *WalletRepo) CreateWallet(userID string) error {
	oldWallet, err := w.GetWalletByUID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "1500") {
			return err
		}
	} else {
		return response.NewInternalServerErrWithMsg(ErrID, "钱包已经创建,用户id"+oldWallet.UserID)
	}

	subAccountId, err := w.binance.CreateSubAccount()
	if err != nil {
		logger.Warnf("用户 %s 创建子账户失败 %v", userID, err)
		return response.NewInternalServerErrMsg(ErrID)
	}
	err = w.binance.EnableSubAccountMargin(subAccountId)
	if err != nil {
		//记录下日至 先忽略
		logger.Warnf("user %s id-%s开启margin权限失败 %v", userID, subAccountId, err)
	}
	//创建子账户的apikey
	subApiResp, err := w.binance.CreateSubAccountApi(subAccountId, "true")
	if err != nil {
		logger.Warnf("用户 %s 创建子账户apiKey失败 %v", userID, err)
		return response.NewInternalServerErrMsg(ErrID)
	}
	//用子账户的api secret查询它的充值地址
	subBinance := exchangeclient.InitBinance(subApiResp.ApiKey, subApiResp.SecretKey)
	address, err := subBinance.GetAccountDepositAddress("USDT")
	//创建用户钱包
	walletModel := model.NewWqWalletModel(userID, subApiResp.ApiKey, subApiResp.SecretKey, subAccountId, address.Address)
	if err := w.dao.CreateWallet(walletModel); err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	return nil
}

func (w *WalletRepo) GetWalletByUID(userID string) (*model.WqWallet, error) {
	walletInfo, err := w.dao.GetWalletByUserID(userID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, response.NewDataNotFound(ErrID, "用户尚未创建钱包")
		}
		return nil, response.NewInternalServerErrMsg(ErrID)
	}
	return walletInfo, nil
}


func (w *WalletRepo) GetWqCoinInfo() (*model.WqCoinInfo, error) {
	info, err := w.dao.GetWqCoinInfo()
	if err != nil {
		return nil, response.NewDataNotFound(ErrID, "请管理员先添加币种兑换信息")
	}
	return info, nil
}

func (w *WalletRepo) Transfer(userID, fromCoin, toCoin string, fromCoinAmount float64) error {
	ifcInfo, err := w.dao.GetWqCoinInfo()
	if err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	walletInfo, err := w.GetWalletByUID(userID)
	if err != nil {
		return err
	}
	//coin, amount, amountBefore, amountAfter, toCoin, toCoinAmount, toCoinAmountAfter, txID string
	var amount = decimal.NewFromFloat(fromCoinAmount).String()
	var amountBefore string
	var amountAfter string
	var toCoinAmount string
	var toCoinAmountBefore string
	var toCoinAmountAfter string
	var record *model.WqTransferRecord

	ifcPrice, _ := decimal.NewFromString(ifcInfo.Price) //ifc 对 usdt 价值
	ifcBalance, _ := decimal.NewFromString(walletInfo.WqCoinBalance)
	if fromCoin == IFC { // 从ifc 转到usdt
		usdt := decimal.NewFromFloat(fromCoinAmount).Mul(ifcPrice).RoundBank(2)
		usdtString := usdt.String()
		//如果转账数量大于余额
		if decimal.NewFromFloat(fromCoinAmount).Cmp(ifcBalance) == 1 {
			return response.NewParamsErrWithMsg(ErrID, "ifc余额不足")
		}
		binance := exchangeclient.InitBinance(walletInfo.ApiKey, walletInfo.Secret)
		spotUsdt, err := binance.GetAccountSpotUsdt()
		if err != nil {
			logger.Warnf("用户%s 查询子账户资产错误 %v", userID, err)
		}
		toCoinAmountBefore = helper.Float64ToString(spotUsdt)
		//增加子账户usdt数量
		_, err = w.binance.ParentTransferToSubAccount(walletInfo.SubAccountID, "", "USDT", usdtString)
		if err != nil {
			if strings.Contains(err.Error(), "-9000") {
				return response.NewInternalServerErrWithMsg(ErrID, "ifortune账上余额不足,划转失败")
			}
			return response.NewInternalServerErrWithMsg(ErrID, fmt.Sprintf("划转失败 %s", err.Error()))
		}
		fromCoinAmountDecimal := decimal.NewFromFloat(fromCoinAmount)

		amountBefore = walletInfo.WqCoinBalance
		amountAfter = ifcBalance.Sub(fromCoinAmountDecimal).String()
		toCoinAmount = usdtString
		//添加记录
		spotUsdt, err = binance.GetAccountSpotUsdt()
		if err != nil {
			logger.Warnf("用户%s 查询子账户资产错误 %v", userID, err)
		}
		toCoinAmountAfter = helper.Float64ToString(spotUsdt)
		record = model.NewWqTransferRecord(userID, fromCoin, amount, amountBefore,
			amountAfter, toCoin, toCoinAmount, toCoinAmountBefore, toCoinAmountAfter, "")

		//更新ifc  balance
		walletInfo.WqCoinBalance = amountAfter

	}
	if fromCoin == USDT { // 从USDT 转到ifc
		ifc := decimal.NewFromFloat(fromCoinAmount).Div(ifcPrice).RoundBank(8)
		binance := exchangeclient.InitBinance(walletInfo.ApiKey, walletInfo.Secret)
		spotUsdt, err := binance.GetAccountSpotUsdt()
		if err != nil {
			logger.Warnf("用户%s 查询子账户资产错误 %v", userID, err)
		}
		amountBefore = helper.Float64ToString(spotUsdt)

		//减去子账户usdt数量
		_, err = w.binance.SubAccountTransferToParent(walletInfo.SubAccountID, "", "USDT", helper.Float64ToString(fromCoinAmount))
		if err != nil {
			return response.NewInternalServerErrWithMsg(ErrID, fmt.Sprintf("划转失败 %s", err.Error()))
		}
		spotUsdt, err = binance.GetAccountSpotUsdt()
		if err != nil {
			logger.Warnf("用户%s 查询子账户资产错误 %v", userID, err)
		}
		amountAfter = helper.Float64ToString(spotUsdt)
		toCoinAmount = ifc.String()
		toCoinAmountBefore = walletInfo.WqCoinBalance
		toCoinAmountAfter = ifcBalance.Add(ifc).String()
		record = model.NewWqTransferRecord(userID, fromCoin, amount, amountBefore, amountAfter, toCoin, toCoinAmount,
			toCoinAmountBefore, toCoinAmountAfter, "")
		//更新ifc  balance
		walletInfo.WqCoinBalance = toCoinAmountAfter
	}
	//TODO 记录表
	if err := w.dao.CreateTransferRecord(record); err != nil {
		logger.Warnf("创建划转记录失败 %+v", err)
	}
	walletInfo.UpdatedAt = time.Now()
	if err := w.dao.UpdateWallet(walletInfo); err != nil {
		logger.Warnf("用户 %s 更新钱包失败 %v", userID, err)
		return response.NewInternalServerErrWithMsg(ErrID, "余额更新失败")
	}
	return nil
}

func (w *WalletRepo) Withdrawal(userID, Coin, Addr string, Volume float64) error {
	wallet, err := w.GetWalletByUID(userID)
	if err != nil {
		return err
	}
	var withdrawal *model.WqWithdrawal
	if Coin == IFC {
		if Volume > helper.StringToFloat64(wallet.WqCoinBalance) {
			return response.NewParamsErrWithMsg(ErrID, "可提现数量不足")
		}
	}
	if Coin == USDT {
		binance := exchangeclient.InitBinance(wallet.ApiKey, wallet.Secret)
		usdt, err := binance.GetAccountSpotUsdt()
		if err != nil {
			return response.NewInternalServerErrMsg(ErrID)
		}
		if Volume > usdt {
			return response.NewParamsErrWithMsg(ErrID, "可提现数量不足")
		}
	}
	withdrawal = &model.WqWithdrawal{
		UserID:    userID,
		Coin:      Coin,
		CashAddr:  Addr,
		Amount:    helper.Float64ToString(Volume),
		State:     1,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	if err := w.dao.CreateWithdrawal(withdrawal); err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	return nil
}

func (w *WalletRepo) CreateWalletAtRunning() {
	resp, err := w.UserSrv.GetAllUserInfo(context.Background(), &empty.Empty{})
	if err != nil {
		logger.Warnf("获取所有用户失败 %+v", err)
		return
	}
	users := resp.UserInfo
	if len(users) == 0 {
		logger.Warnf("没有用户")
		return
	}
	for _, user := range users {
		_, err := w.GetWalletByUID(user.UserID)
		if err != nil && strings.Contains(err.Error(), "1404") {
			err := w.CreateWallet(user.UserID)
			if err != nil {
				logger.Warnf("用户id%s 创建用户钱包失败 %+v", user.UserID, err)
				continue
			}
			logger.Infof("创建钱包成功 用户id %s", user.UserID)
		}
	}
}

func (w *WalletRepo) AddIfcBalance(userID, inUserID, _type, exchange string, volumeIn float64) error {
	logger.Infof("增加用户ifc数量 uid %s volume %v", userID, volumeIn)
	wallet, err := w.GetWalletByUID(userID)
	if err != nil {
		return err
	}
	logger.Infof("用户原本ifc余额 %v", wallet.WqCoinBalance)
	oldIfcBlance, err := decimal.NewFromString(wallet.WqCoinBalance)
	volume := decimal.NewFromFloat(volumeIn)
	if err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	wallet.WqCoinBalance = oldIfcBlance.Add(volume).String()
	wallet.UpdatedAt = time.Now()
	if err = w.dao.UpdateWallet(wallet); err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	record := &model.WqIfcGiftRecord{
		UserID:    userID,
		InUserID:  inUserID,
		Volume:    volume.String(),
		Type:      _type,
		Exchange:  exchange,
		Before:    oldIfcBlance.String(),
		After:     wallet.WqCoinBalance,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	if err = w.dao.CreateIfcGiftRecord(record); err != nil {
		logger.Errorf("创建赠送记录失败 %v  %+v", err, record)
	}
	return nil
}

func (w *WalletRepo) GetIfcRecordByUid(uid string) []*model.WqIfcGiftRecord {
	return w.dao.GetIfcGiftRecordByUid(uid)
}

func (w *WalletRepo) GetIfcRecordByUidExchange(userMasterId, inUserID, exchange string) []*model.WqIfcGiftRecord {
	return w.dao.GetIfcGiftRecordBySql(userMasterId, inUserID, exchange)
}

// AddIfcByStrategyRunInfo 给当月首次启动策略加积分
func (w *WalletRepo) AddIfcByStrategyRunInfo(userMasterID, inUserID string, volume float64) {
	_, err := w.cacheService.GetUserStrategyRunInfo(userMasterID)
	if err != nil {
		if err == cache.KeyNotFound {
			//增加积分
			err := w.AddIfcBalance(userMasterID, inUserID, TYPE_STRATEGY, "", volume)
			if err != nil {
				logger.Errorf("增加积分失败 %v uid %s inUserId %s type %s volume %v", err.Error(),
					userMasterID, inUserID, TYPE_STRATEGY, volume)
				return
			}
			w.cacheService.CacheUserStrategyRunInfo(userMasterID)
			return
		}
		logger.Errorf("redis连接断开 增加积分失败 userMaster %s inUser %s volume %v", userMasterID, inUserID, volume)
		return
	}
	logger.Warnf("用户 uid%s inUserId %s 本月已经加过策略运行的积分", userMasterID, inUserID)
}
