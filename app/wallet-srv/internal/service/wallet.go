package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/shopspring/decimal"
	"strings"
	"wq-fotune-backend/api/response"
	userPb "wq-fotune-backend/api/usercenter"
	pb "wq-fotune-backend/api/wallet"
	"wq-fotune-backend/app/wallet-srv/internal/biz"
	"wq-fotune-backend/libs/exchange"
	"wq-fotune-backend/libs/exchangeclient"
	"wq-fotune-backend/libs/logger"

)

type WalletService struct {

	walletSrv *biz.WalletRepo
}

func NewWalletHandler() *WalletService {
	return &WalletService{
		walletSrv: biz.NewWalletRepo(),
	}
}
func (w WalletService) GetUsdtDepositAddr(ctx context.Context, req *pb.UidReq, resp *pb.UsdtDepositAddrResp) error {
	wallet, err := w.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return err
	}
	resp.Address = wallet.DepositAddr
	return nil
}

func (w WalletService) CreateWallet(ctx context.Context, req *pb.UidReq, empty *empty.Empty) error {
	return w.walletSrv.CreateWallet(req.UserId)
}

func (w WalletService) Transfer(ctx context.Context, req *pb.TransferReq, e *empty.Empty) error {
	if req.FromCoin == req.ToCoin {
		return response.NewInternalServerErrWithMsg(biz.ErrID, "转入转出币种填写错误")
	}
	if req.FromCoin != biz.IFC && req.FromCoin != biz.USDT {
		return response.NewInternalServerErrWithMsg(biz.ErrID, "转入转出币种填写错误")
	}
	if req.ToCoin != biz.IFC && req.ToCoin != biz.USDT {
		return response.NewInternalServerErrWithMsg(biz.ErrID, "转入转出币种填写错误")
	}
	if req.FromCoinAmount <= 0.0 {
		return response.NewParamsErrWithMsg(biz.ErrID, "数量填写错误")
	}
	return w.walletSrv.Transfer(req.UserId, req.FromCoin, req.ToCoin, req.FromCoinAmount)
}

func (w WalletService) GetWalletUSDT(ctx context.Context, req *pb.UidReq, resp *pb.WalletBalanceResp) error {
	wallet, err := w.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return err
	}
	binance := exchangeclient.InitBinance(wallet.ApiKey, wallet.Secret)
	spot, err := binance.GetAccountSpot()
	if err != nil {
		if strings.Contains(err.Error(), "1022") {
			return response.NewInternalServerErrWithMsg(biz.ErrID, "钱包密钥实效")
		}
		logger.Warnf("用户id %s 获取子账户财产失败 %+v", err)
		return response.NewInternalServerErrMsg(biz.ErrID)
	}
	if len(spot.SubAccounts) == 0 {
		resp.Title = "usdt钱包"
		resp.Symbol = "usdt"
		resp.Total = "0"
		resp.Available = "0"
		return nil
	}
	for key, value := range spot.SubAccounts {
		if key.Symbol == biz.USDT {
			resp.Title = "usdt钱包"
			resp.Symbol = "usdt"
			resp.Total = decimal.NewFromFloat(value.Balance).String()
			resp.Available = decimal.NewFromFloat(value.Amount).String()
			break
		}
	}
	return nil
}

func (w WalletService) GetWalletIFC(ctx context.Context, req *pb.UidReq, resp *pb.WalletBalanceResp) error {
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

func (w WalletService) ConvertCoinTips(ctx context.Context, req *pb.ConvertCoinTipsReq, resp *pb.ConvertCoinTipsResp) error {
	if req.From != biz.IFC && req.From != biz.USDT {
		return response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	if req.To != biz.IFC && req.To != biz.USDT {
		return response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	if req.From == req.To {
		return response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
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

	if req.From == biz.IFC {
		usdtVolume := wqCoinBalance.Mul(wqCoinUsdtPrice) //算出可兑换多少usdt
		if usdtVolume.Equal(decimal.NewFromFloat(0)) {
			resp.Describe = fmt.Sprintf("ifc可用余额不足,可兑换 %s USDT, 当前比例 %d:%s", usdtVolume.RoundBank(2).String(), 1, wqCoinUsdtPrice.String())
			return nil
		}
		resp.Describe = fmt.Sprintf("可兑换 %s USDT, 当前比例 %d:%s", usdtVolume.RoundBank(2).String(), 1, wqCoinUsdtPrice.String())
		return nil
	}
	//else
	binance := exchangeclient.InitBinance(wallet.ApiKey, wallet.Secret)
	spot, err := binance.GetAccountSpot()
	if err != nil {
		logger.Warnf("用户id%s 钱包api-%s 获取财产%+v", req.UserId, wallet.ApiKey, err)
		return response.NewInternalServerErrMsg(biz.ErrID)
	}
	if len(spot.SubAccounts) == 0 {
		resp.Describe = fmt.Sprintf("usdt可用余额不足, 可兑换 %s IFC, 当前比例 %s:%s", wqCoinUsdtPrice, 1)
		return nil
	}
	for key, value := range spot.SubAccounts {
		if key.Symbol == biz.USDT {
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

func (w WalletService) ConvertCoin(ctx context.Context, req *pb.ConvertCoinReq, resp *pb.ConvertCoinResp) error {
	if req.From != biz.IFC && req.From != biz.USDT {
		return response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	if req.To != biz.IFC && req.To != biz.USDT {
		return response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	if req.From == req.To {
		return response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	wqCoinInfo, err := w.walletSrv.GetWqCoinInfo()
	if err != nil {
		return err
	}
	if req.Volume <= 0 {
		return response.NewParamsErrWithMsg(biz.ErrID, "数量不能小于0")
	}
	wqCoinUsdtPrice, _ := decimal.NewFromString(wqCoinInfo.Price)
	if req.From == biz.IFC {
		resp.Describe = "可兑换usdt"
		resp.Volume, _ = decimal.NewFromFloat(req.Volume).Mul(wqCoinUsdtPrice).RoundBank(2).Float64()
		return nil
	}
	resp.Describe = "可兑换ifc"
	resp.Volume, _ = decimal.NewFromFloat(req.Volume).Div(wqCoinUsdtPrice).RoundBank(8).Float64()
	return nil
}

func (w WalletService) Withdrawal(ctx context.Context, req *pb.WithdrawalReq, e *empty.Empty) error {
	if req.Coin != biz.IFC && req.Coin != biz.USDT {
		return response.NewParamsErrWithMsg(biz.ErrID, "不支持体现该币种"+req.Coin)
	}
	if req.Volume <= 0.0 {
		return response.NewParamsErrWithMsg(biz.ErrID, "请输入正确数量")
	}
	return w.walletSrv.Withdrawal(req.UserId, req.Coin, req.Address, req.Volume)
}

var exchangeParam = map[string]string{exchange.BINANCE: "", exchange.HUOBI: "", exchange.OKEX: ""}

func (w WalletService) AddIfcBalance(ctx context.Context, req *pb.AddIfcBalanceReq, e *empty.Empty) error {
	if req.Type != biz.TYPE_BIND_API && req.Type != biz.TYPE_REGISTER && req.Type != biz.TYPE_STRATEGY {
		return response.NewParamsErrWithMsg(biz.ErrID, "type参数只能为register,api,strategy当中")
	}
	if req.Type == biz.TYPE_BIND_API {
		if _, ok := exchangeParam[req.Exchange]; !ok {
			return response.NewParamsErrWithMsg(biz.ErrID, "交易所exchange只能为binance,huobi,okex")
		}
		record := w.walletSrv.GetIfcRecordByUidExchange(req.UserMasterId, req.InUserId, req.Exchange)
		if len(record) != 0 {
			logger.Warnf("已存在相同交易所关联的赠送记录,不需要重复赠送哦 %+v", req)
			return nil
		}
	}
	return w.walletSrv.AddIfcBalance(req.UserMasterId, req.InUserId, req.Type, req.Exchange, req.Volume)
}

func (w WalletService) GetTotalRebate(ctx context.Context, req *pb.GetTotalRebateReq, resp *pb.GetTotalRebateResp) error {
	data := w.walletSrv.GetIfcRecordByUid(req.UserId)
	if len(data) == 0 {
		return response.NewDataNotFound(biz.ErrID, "暂无数据")
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
		if v.Type == biz.TYPE_STRATEGY {
			msg = "邀请用户使用量化策略赠送"
		}
		if v.Type == biz.TYPE_REGISTER {
			msg = "邀请用户注册赠送"
		}
		if v.Type == biz.TYPE_BIND_API {
			msg = "邀请用户绑定api赠送"
		}
		date := v.UpdatedAt.Format("2006-01-02")
		resp.Record = append(resp.Record, &pb.IfcRecord{
			Phone:   phone,
			Volume:  v.Volume,
			TypeMsg: msg,
			Date:    date,
		})
	}
	resp.Total = total.String()
	return nil
}

func (w WalletService) StrategyRunNotify(ctx context.Context, req *pb.StrategyRunNotifyReq, e *empty.Empty) error {
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
