package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/shopspring/decimal"
	"strings"
	apiBinance "wq-fotune-backend/libs/binance_client"
	"wq-fotune-backend/libs/logger"
	exchange_info "wq-fotune-backend/pkg/exchange-info"
	"wq-fotune-backend/pkg/response"
	userPb "wq-fotune-backend/internal/usercenter-srv/proto"
	fotune_srv_wallet "wq-fotune-backend/internal/wallet-srv/proto"
	"wq-fotune-backend/internal/wallet-srv/service"
)

type WalletHandler struct {
	walletSrv *service.WalletService
}

func NewWalletHandler() *WalletHandler {
	return &WalletHandler{
		walletSrv: service.NewWalletService(),
	}
}
func (w WalletHandler) GetUsdtDepositAddr(ctx context.Context, req *fotune_srv_wallet.UidReq, resp *fotune_srv_wallet.UsdtDepositAddrResp) error {
	wallet, err := w.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return err
	}
	resp.Address = wallet.DepositAddr
	return nil
}

func (w WalletHandler) CreateWallet(ctx context.Context, req *fotune_srv_wallet.UidReq, empty *empty.Empty) error {
	return w.walletSrv.CreateWallet(req.UserId)
}

func (w WalletHandler) Transfer(ctx context.Context, req *fotune_srv_wallet.TransferReq, e *empty.Empty) error {
	if req.FromCoin == req.ToCoin {
		return response.NewInternalServerErrWithMsg(service.ErrID, "转入转出币种填写错误")
	}
	if req.FromCoin != service.IFC && req.FromCoin != service.USDT {
		return response.NewInternalServerErrWithMsg(service.ErrID, "转入转出币种填写错误")
	}
	if req.ToCoin != service.IFC && req.ToCoin != service.USDT {
		return response.NewInternalServerErrWithMsg(service.ErrID, "转入转出币种填写错误")
	}
	if req.FromCoinAmount <= 0.0 {
		return response.NewParamsErrWithMsg(service.ErrID, "数量填写错误")
	}
	return w.walletSrv.Transfer(req.UserId, req.FromCoin, req.ToCoin, req.FromCoinAmount)
}

func (w WalletHandler) GetWalletUSDT(ctx context.Context, req *fotune_srv_wallet.UidReq, resp *fotune_srv_wallet.WalletBalanceResp) error {
	wallet, err := w.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return err
	}
	binance := apiBinance.InitClient(wallet.ApiKey, wallet.Secret)
	spot, err := binance.GetAccountSpot()
	if err != nil {
		if strings.Contains(err.Error(), "1022") {
			return response.NewInternalServerErrWithMsg(service.ErrID, "钱包密钥实效")
		}
		logger.Warnf("用户id %s 获取子账户财产失败 %+v", err)
		return response.NewInternalServerErrMsg(service.ErrID)
	}
	if len(spot.SubAccounts) == 0 {
		resp.Title = "usdt钱包"
		resp.Symbol = "usdt"
		resp.Total = "0"
		resp.Available = "0"
		return nil
	}
	for key, value := range spot.SubAccounts {
		if key.Symbol == service.USDT {
			resp.Title = "usdt钱包"
			resp.Symbol = "usdt"
			resp.Total = decimal.NewFromFloat(value.Balance).String()
			resp.Available = decimal.NewFromFloat(value.Amount).String()
			break
		}
	}
	return nil
}

func (w WalletHandler) GetWalletIFC(ctx context.Context, req *fotune_srv_wallet.UidReq, resp *fotune_srv_wallet.WalletBalanceResp) error {
	wallet, err := w.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return err
	}
	resp.Symbol = "ifc"
	resp.Title = "ifc钱包"
	resp.Total = wallet.WqCoinBalance
	resp.Available = wallet.WqCoinBalance
	return nil
}

func (w WalletHandler) ConvertCoinTips(ctx context.Context, req *fotune_srv_wallet.ConvertCoinTipsReq, resp *fotune_srv_wallet.ConvertCoinTipsResp) error {
	if req.From != service.IFC && req.From != service.USDT {
		return response.NewParamsErrWithMsg(service.ErrID, "参数有误")
	}
	if req.To != service.IFC && req.To != service.USDT {
		return response.NewParamsErrWithMsg(service.ErrID, "参数有误")
	}
	if req.From == req.To {
		return response.NewParamsErrWithMsg(service.ErrID, "参数有误")
	}
	wallet, err := w.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return err
	}
	wqCoinInfo, err := w.walletSrv.GetWqCoinInfo()
	if err != nil {
		return err
	}
	//ifc币种对usdt 的价格
	wqCoinUsdtPrice, _ := decimal.NewFromString(wqCoinInfo.Price)
	wqCoinBalance, _ := decimal.NewFromString(wallet.WqCoinBalance)

	if req.From == service.IFC {
		usdtVolume := wqCoinBalance.Mul(wqCoinUsdtPrice) //算出可兑换多少usdt
		if usdtVolume.Equal(decimal.NewFromFloat(0)) {
			resp.Describe = fmt.Sprintf("ifc可用余额不足,可兑换 %s USDT, 当前比例 %d:%s", usdtVolume.RoundBank(2).String(), 1, wqCoinUsdtPrice.String())
			return nil
		}
		resp.Describe = fmt.Sprintf("可兑换 %s USDT, 当前比例 %d:%s", usdtVolume.RoundBank(2).String(), 1, wqCoinUsdtPrice.String())
		return nil
	}
	//else
	binance := apiBinance.InitClient(wallet.ApiKey, wallet.Secret)
	spot, err := binance.GetAccountSpot()
	if err != nil {
		logger.Warnf("用户id%s 钱包api-%s 获取财产%+v", req.UserId, wallet.ApiKey, err)
		return response.NewInternalServerErrMsg(service.ErrID)
	}
	if len(spot.SubAccounts) == 0 {
		resp.Describe = fmt.Sprintf("usdt可用余额不足, 可兑换 %s IFC, 当前比例 %s:%s", wqCoinUsdtPrice, 1)
		return nil
	}
	for key, value := range spot.SubAccounts {
		if key.Symbol == service.USDT {
			if value.Amount < 1 {
				resp.Describe = fmt.Sprintf("usdt可用余额不足, 可兑换 0 IFC, 当前比例 %s:%d", wqCoinUsdtPrice.String(), 1)
				return nil
			}
			usdtBlance := decimal.NewFromFloat(value.Amount)
			ifcVolume := usdtBlance.Div(wqCoinUsdtPrice).RoundBank(8)
			resp.Describe = fmt.Sprintf("可兑换 %s IFC, 当前比例 %s:%d", ifcVolume, wqCoinUsdtPrice.String(), 1)
			return nil
		}
	}
	return nil
}

func (w WalletHandler) ConvertCoin(ctx context.Context, req *fotune_srv_wallet.ConvertCoinReq, resp *fotune_srv_wallet.ConvertCoinResp) error {
	if req.From != service.IFC && req.From != service.USDT {
		return response.NewParamsErrWithMsg(service.ErrID, "参数有误")
	}
	if req.To != service.IFC && req.To != service.USDT {
		return response.NewParamsErrWithMsg(service.ErrID, "参数有误")
	}
	if req.From == req.To {
		return response.NewParamsErrWithMsg(service.ErrID, "参数有误")
	}
	wqCoinInfo, err := w.walletSrv.GetWqCoinInfo()
	if err != nil {
		return err
	}
	if req.Volume <= 0 {
		return response.NewParamsErrWithMsg(service.ErrID, "数量不能小于0")
	}
	wqCoinUsdtPrice, _ := decimal.NewFromString(wqCoinInfo.Price)
	if req.From == service.IFC {
		resp.Describe = "可兑换usdt"
		resp.Volume, _ = decimal.NewFromFloat(req.Volume).Mul(wqCoinUsdtPrice).RoundBank(2).Float64()
		return nil
	}
	resp.Describe = "可兑换ifc"
	resp.Volume, _ = decimal.NewFromFloat(req.Volume).Div(wqCoinUsdtPrice).RoundBank(8).Float64()
	return nil
}

func (w WalletHandler) Withdrawal(ctx context.Context, req *fotune_srv_wallet.WithdrawalReq, e *empty.Empty) error {
	if req.Coin != service.IFC && req.Coin != service.USDT {
		return response.NewParamsErrWithMsg(service.ErrID, "不支持体现该币种"+req.Coin)
	}
	if req.Volume <= 0.0 {
		return response.NewParamsErrWithMsg(service.ErrID, "请输入正确数量")
	}
	return w.walletSrv.Withdrawal(req.UserId, req.Coin, req.Address, req.Volume)
}

var exchangeParam = map[string]string{exchange_info.BINANCE: "", exchange_info.HUOBI: "", exchange_info.OKEX: ""}

func (w WalletHandler) AddIfcBalance(ctx context.Context, req *fotune_srv_wallet.AddIfcBalanceReq, e *empty.Empty) error {
	if req.Type != service.TYPE_BIND_API && req.Type != service.TYPE_REGISTER && req.Type != service.TYPE_STRATEGY {
		return response.NewParamsErrWithMsg(service.ErrID, "type参数只能为register,api,strategy当中")
	}
	if req.Type == service.TYPE_BIND_API {
		if _, ok := exchangeParam[req.Exchange]; !ok {
			return response.NewParamsErrWithMsg(service.ErrID, "交易所exchange只能为binance,huobi,okex")
		}
		record := w.walletSrv.GetIfcRecordByUidExchange(req.UserMasterId, req.InUserId, req.Exchange)
		if len(record) != 0 {
			logger.Warnf("已存在相同交易所关联的赠送记录,不需要重复赠送哦 %+v", req)
			return nil
		}
	}
	return w.walletSrv.AddIfcBalance(req.UserMasterId, req.InUserId, req.Type, req.Exchange, req.Volume)
}

func (w WalletHandler) GetTotalRebate(ctx context.Context, req *fotune_srv_wallet.GetTotalRebateReq, resp *fotune_srv_wallet.GetTotalRebateResp) error {
	data := w.walletSrv.GetIfcRecordByUid(req.UserId)
	if len(data) == 0 {
		return response.NewDataNotFound(service.ErrID, "暂无数据")
	}
	var total decimal.Decimal
	for _, v := range data {
		volume, _ := decimal.NewFromString(v.Volume)
		total = total.Add(volume)
		//被邀请用户
		inUser, err := w.walletSrv.UserSrv.GetUserInfo(context.Background(), &userPb.UserInfoReq{
			UserID: v.InUserID,
		})
		phone := ""
		if err != nil {
			phone = "00000000000"
		} else {
			phone = inUser.Phone
		}
		msg := ""
		if v.Type == service.TYPE_STRATEGY {
			msg = "邀请用户使用量化策略赠送"
		}
		if v.Type == service.TYPE_REGISTER {
			msg = "邀请用户注册赠送"
		}
		if v.Type == service.TYPE_BIND_API {
			msg = "邀请用户绑定api赠送"
		}
		date := v.UpdatedAt.Format("2006-01-02")
		resp.Record = append(resp.Record, &fotune_srv_wallet.IfcRecord{
			Phone:   phone,
			Volume:  v.Volume,
			TypeMsg: msg,
			Date:    date,
		})
	}
	resp.Total = total.String()
	return nil
}

func (w WalletHandler) StrategyRunNotify(ctx context.Context, req *fotune_srv_wallet.StrategyRunNotifyReq, e *empty.Empty) error {
	info, err := w.walletSrv.UserSrv.GetUserMasterByInViteUser(context.Background(), &userPb.GetUserMasterReq{
		InviteUid: req.UserId,
	})
	if err != nil {
		logger.Warnf("没有找到邀请关系 不需要增加积分 %s", req.UserId)
		return nil
	}
	w.walletSrv.AddIfcByStrategyRunInfo(info.UserMasterId, req.UserId, 5.0)
	return nil
}
