package service

import (
	"context"
	"fmt"
	"fortune-bd/api/constant"
	"fortune-bd/api/response"
	"fortune-bd/app/wallet-svc/internal/biz"
	"fortune-bd/libs/exchangeclient"
	"fortune-bd/libs/logger"
	"github.com/shopspring/decimal"
	"strings"

	userpb "fortune-bd/api/usercenter/v1"
	pb "fortune-bd/api/wallet/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var exchangeParam = map[string]string{constant.BINANCE: "", constant.HUOBI: ""}

type WalletService struct {
	pb.UnimplementedWalletServer
	walletSrv *biz.WalletRepo
}

func NewWalletService() *WalletService {
	return &WalletService{
		walletSrv: biz.NewWalletRepo(),
	}
}

func (s *WalletService) CreateWallet(ctx context.Context, req *pb.UidReq) (*emptypb.Empty, error) {
	err := s.walletSrv.CreateWallet(req.UserId)
	return nil, err
}

func (s *WalletService) Transfer(ctx context.Context, req *pb.TransferReq) (*emptypb.Empty, error) {
	if req.FromCoin == req.ToCoin {
		return nil, response.NewInternalServerErrWithMsg(biz.ErrID, "转入转出币种填写错误")
	}
	if req.FromCoin != biz.IFC && req.FromCoin != biz.USDT {
		return nil, response.NewInternalServerErrWithMsg(biz.ErrID, "转入转出币种填写错误")
	}
	if req.ToCoin != biz.IFC && req.ToCoin != biz.USDT {
		return nil, response.NewInternalServerErrWithMsg(biz.ErrID, "转入转出币种填写错误")
	}
	if req.FromCoinAmount <= 0.0 {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "数量填写错误")
	}
	err := s.walletSrv.Transfer(req.UserId, req.FromCoin, req.ToCoin, req.FromCoinAmount)
	return nil, err
}

func (s *WalletService) GetWalletIfc(ctx context.Context, req *pb.UidReq) (*pb.WalletBalanceResp, error) {
	var resp = new(pb.WalletBalanceResp)
	wallet, err := s.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return nil,  err
	}
	binance := exchangeclient.InitBinance(wallet.ApiKey, wallet.Secret)
	spot, err := binance.GetAccountSpot()
	if err != nil {
		if strings.Contains(err.Error(), "1022") {
			return nil, response.NewInternalServerErrWithMsg(biz.ErrID, "钱包密钥实效")
		}
		logger.Warnf("用户id %s 获取子账户财产失败 %+v", err)
		return nil, response.NewInternalServerErrMsg(biz.ErrID)
	}
	if len(spot.SubAccounts) == 0 {
		resp.Title = "usdt钱包"
		resp.Symbol = "usdt"
		resp.Total = "0"
		resp.Available = "0"
		return resp, nil
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
	return resp, nil
}

func (s *WalletService) GetWalletUsdt(ctx context.Context, req *pb.UidReq) (*pb.WalletBalanceResp, error) {
	var resp = new(pb.WalletBalanceResp)
	wallet, err := s.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return nil, err
	}
	resp.Symbol = "ifc"
	resp.Title = "ifc钱包"
	resp.Total = wallet.WqCoinBalance
	resp.Available = wallet.WqCoinBalance
	return resp, nil
}

func (s *WalletService) GetUsdtDepositAddr(ctx context.Context, req *pb.UidReq) (*pb.UsdtDepositAddrResp, error) {
	var resp = new(pb.UsdtDepositAddrResp)
	wallet, err := s.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return nil, err
	}
	resp.Address = wallet.DepositAddr
	return resp, nil
}

func (s *WalletService) ConvertCoinTips(ctx context.Context, req *pb.ConvertCoinTipsReq) (*pb.ConvertCoinTipsResp, error) {
	var resp = new(pb.ConvertCoinTipsResp)
	if req.From != biz.IFC && req.From != biz.USDT {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	if req.To != biz.IFC && req.To != biz.USDT {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	if req.From == req.To {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	wallet, err := s.walletSrv.GetWalletByUID(req.UserId)
	if err != nil {
		return nil, err
	}
	wqCoinInfo, err := s.walletSrv.GetWqCoinInfo()
	if err != nil {
		return nil, err
	}
	//ifc币种对usdt 的价格
	wqCoinUsdtPrice, _ := decimal.NewFromString(wqCoinInfo.Price)
	wqCoinBalance, _ := decimal.NewFromString(wallet.WqCoinBalance)

	if req.From == biz.IFC {
		usdtVolume := wqCoinBalance.Mul(wqCoinUsdtPrice) //算出可兑换多少usdt
		if usdtVolume.Equal(decimal.NewFromFloat(0)) {
			resp.Describe = fmt.Sprintf("ifc可用余额不足,可兑换 %s USDT, 当前比例 %d:%s", usdtVolume.RoundBank(2).String(), 1, wqCoinUsdtPrice.String())
			return resp, nil
		}
		resp.Describe = fmt.Sprintf("可兑换 %s USDT, 当前比例 %d:%s", usdtVolume.RoundBank(2).String(), 1, wqCoinUsdtPrice.String())
		return resp, nil
	}
	//else
	binance := exchangeclient.InitBinance(wallet.ApiKey, wallet.Secret)
	spot, err := binance.GetAccountSpot()
	if err != nil {
		logger.Warnf("用户id%s 钱包api-%s 获取财产%+v", req.UserId, wallet.ApiKey, err)
		return nil, response.NewInternalServerErrMsg(biz.ErrID)
	}
	if len(spot.SubAccounts) == 0 {
		resp.Describe = fmt.Sprintf("usdt可用余额不足, 可兑换 %s IFC, 当前比例 %s:%s", wqCoinUsdtPrice, 1)
		return resp, nil
	}
	for key, value := range spot.SubAccounts {
		if key.Symbol == biz.USDT {
			if value.Amount < 1 {
				resp.Describe = fmt.Sprintf("usdt可用余额不足, 可兑换 0 IFC, 当前比例 %s:%d", wqCoinUsdtPrice.String(), 1)
				return resp, nil
			}
			usdtBalance := decimal.NewFromFloat(value.Amount)
			ifcVolume := usdtBalance.Div(wqCoinUsdtPrice).RoundBank(8)
			resp.Describe = fmt.Sprintf("可兑换 %s IFC, 当前比例 %s:%d", ifcVolume, wqCoinUsdtPrice.String(), 1)
			return resp, nil
		}
	}
	return resp, nil
}

func (s *WalletService) ConvertCoin(ctx context.Context, req *pb.ConvertCoinReq) (*pb.ConvertCoinResp, error) {
	var resp = new(pb.ConvertCoinResp)
	if req.From != biz.IFC && req.From != biz.USDT {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	if req.To != biz.IFC && req.To != biz.USDT {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	if req.From == req.To {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "参数有误")
	}
	wqCoinInfo, err := s.walletSrv.GetWqCoinInfo()
	if err != nil {
		return nil, err
	}
	if req.Volume <= 0 {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "数量不能小于0")
	}
	wqCoinUsdtPrice, _ := decimal.NewFromString(wqCoinInfo.Price)
	if req.From == biz.IFC {
		resp.Describe = "可兑换usdt"
		resp.Volume, _ = decimal.NewFromFloat(req.Volume).Mul(wqCoinUsdtPrice).RoundBank(2).Float64()
		return resp,nil
	}
	resp.Describe = "可兑换ifc"
	resp.Volume, _ = decimal.NewFromFloat(req.Volume).Div(wqCoinUsdtPrice).RoundBank(8).Float64()
	return resp, nil
}

func (s *WalletService) Withdrawal(ctx context.Context, req *pb.WithdrawalReq) (*emptypb.Empty, error) {
	if req.Coin != biz.IFC && req.Coin != biz.USDT {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "不支持体现该币种"+req.Coin)
	}
	if req.Volume <= 0.0 {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "请输入正确数量")
	}
	err := s.walletSrv.Withdrawal(req.UserId, req.Coin, req.Address, req.Volume)
	return nil, err
}

func (s *WalletService) AddIfcBalance(ctx context.Context, req *pb.AddIfcBalanceReq) (*emptypb.Empty, error) {
	if req.Type != biz.TYPE_BIND_API && req.Type != biz.TYPE_REGISTER && req.Type != biz.TYPE_STRATEGY {
		return nil, response.NewParamsErrWithMsg(biz.ErrID, "type参数只能为register,api,strategy当中")
	}
	if req.Type == biz.TYPE_BIND_API {
		if _, ok := exchangeParam[req.Exchange]; !ok {
			return nil, response.NewParamsErrWithMsg(biz.ErrID, "交易所exchange只能为binance,huobi,okex")
		}
		record := s.walletSrv.GetIfcRecordByUidExchange(req.UserMasterId, req.InUserId, req.Exchange)
		if len(record) != 0 {
			logger.Warnf("已存在相同交易所关联的赠送记录,不需要重复赠送哦 %+v", req)
			return nil, nil
		}
	}
	err := s.walletSrv.AddIfcBalance(req.UserMasterId, req.InUserId, req.Type, req.Exchange, req.Volume)
	return nil, err
}

func (s *WalletService) GetTotalRebate(ctx context.Context, req *pb.GetTotalRebateReq) (*pb.GetTotalRebateResp, error) {
	var resp = new(pb.GetTotalRebateResp)
	data := s.walletSrv.GetIfcRecordByUid(req.UserId)
	if len(data) == 0 {
		return nil, response.NewDataNotFound(biz.ErrID, "暂无数据")
	}
	var total decimal.Decimal
	for _, v := range data {
		volume, _ := decimal.NewFromString(v.Volume)
		total = total.Add(volume)
		//被邀请用户
		inUser, err := s.walletSrv.UserSrv.GetUserInfo(context.Background(), &userpb.UserInfoReq{
			UserId: v.InUserID,
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
	return resp, nil
}

func (s *WalletService) StrategyRunNotify(ctx context.Context, req *pb.StrategyRunNotifyReq) (*emptypb.Empty, error) {
	info, err := s.walletSrv.UserSrv.GetUserMasterByInViteUser(context.Background(), &userpb.GetUserMasterReq{
		InviteUid: req.UserId,
	})
	if err != nil {
		logger.Warnf("没有找到邀请关系 不需要增加积分 %s", req.UserId)
		return nil, err
	}
	s.walletSrv.AddIfcByStrategyRunInfo(info.UserMasterId, req.UserId, 5.0)
	return nil, nil
}
