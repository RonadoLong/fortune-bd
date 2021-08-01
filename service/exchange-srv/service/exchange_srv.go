package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/shopspring/decimal"
	"log"
	"strings"
	"time"
	"wq-fotune-backend/api-gateway/protocol"
	apiBinance "wq-fotune-backend/libs/binance_client"
	"wq-fotune-backend/libs/exchange_clientI"
	"wq-fotune-backend/libs/helper"
	apiHuobi "wq-fotune-backend/libs/huobi_client"
	"wq-fotune-backend/libs/logger"
	api "wq-fotune-backend/libs/okex_client"
	"wq-fotune-backend/pkg/encoding"
	exchange_info "wq-fotune-backend/pkg/exchange-info"
	"wq-fotune-backend/pkg/response"
	"wq-fotune-backend/pkg/utils"
	"wq-fotune-backend/service/exchange-srv/client"
	"wq-fotune-backend/service/exchange-srv/model"
	pb "wq-fotune-backend/service/exchange-srv/proto"
	quoteCron "wq-fotune-backend/service/quote-srv/cron"
)

func (e *ExOrderService) GetExchangeInfo() ([]*model.WqExchange, error) {
	exchangeList := e.dao.GetExchangeInfo()
	if len(exchangeList) == 0 {
		logger.Infof("ExChangeInfo 没有找到")
		return nil, response.NewExchangeNotFoundErrMsg(ErrID)
	}
	return exchangeList, nil
}

func (e *ExOrderService) checkIfApiValid(apiKey, secret, passphrase, exchange string) error {
	client := e.GetExchangeClient(apiKey, secret, passphrase, exchange)
	if err := client.CheckIfApiValid(); err != nil {
		return response.NewExchangeApiCheckErrMsg(ErrID)
	}
	return nil
}

func (e *ExOrderService) AddExchangeApi(userID, apiKey, secret, passphrase string, exchangeID int64) error {
	exchange, err := e.dao.GetExchangeById(exchangeID)
	if err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	if exchange.Exchange == exchange_info.OKEX {
		if passphrase == "" {
			return response.NewExchangePassphraseNoneErrMsg(ErrID)
		}
	}
	oldAPI, _ := e.dao.GetExchangeApiByUidAndApi(userID, apiKey)
	if oldAPI != nil {
		return response.NewExchangeApiDuplicateErrMsg(ErrID)
	}

	log.Println("=================", exchangeID)
	//check
	err = e.checkIfApiValid(apiKey, secret, passphrase, exchange.Exchange)
	if err != nil {
		log.Println(err)
		return err
	}

	secretCrypt, _ := encoding.AesEncrypt([]byte(secret))

	newAPI := &model.WqExchangeApi{
		UserID:       userID,
		ExchangeID:   exchangeID,
		ExchangeName: exchange.Exchange,
		ApiKey:       apiKey,
		Secret:       hex.EncodeToString(secretCrypt),
		Passphrase:   passphrase,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
	if err := e.dao.AddExchangeApi(newAPI); err != nil {
		return response.NewExchangeApiCreateErrMsg(ErrID)
	}
	//赠送邀请ifc
	//resp, err := client.UserService.GetUserMasterByInViteUser(context.Background(), &userPb.GetUserMasterReq{
	//	InviteUid: userID,
	//})
	//if err != nil {
	//	logger.Warnf("找不到邀请数据 无需添加ifc userid %s", userID)
	//	return nil
	//}
	//_, err = client.WalletService.AddIfcBalance(context.Background(), &walletPb.AddIfcBalanceReq{
	//	UserMasterId: resp.UserMasterId,
	//	InUserId:     userID,
	//	Volume:       10,
	//	Type:         "api",
	//	Exchange:     exchange.Exchange,
	//})
	//if err != nil {
	//	logger.Warnf("添加apikey 增加发出邀请码用户的ifc失败 uid %s userMasterID %s  err %v", userID, resp.UserMasterId, err)
	//}
	return nil
}

func (e *ExOrderService) GetExchangeAccountListFromCache(userID string) []byte {
	ret := e.cacheService.GetExchangeAccountList(userID)
	return ret
}

func (e *ExOrderService) SetExchangeAccountListCache(userID string, data []byte) {
	e.cacheService.CacheExchangeAccountList(userID, data)
}

func (e *ExOrderService) GetExchangeApiList(userId string) ([]*protocol.ExchangeApiResp, error) {
	apiList := e.dao.GetExchangeApiListByUid(userId)
	if len(apiList) == 0 {
		return nil, response.NewExchangeApiListErrMsg(ErrID)
	}
	byteData, err := json.Marshal(apiList)
	if err != nil {
		logger.Errorf("GetExchangeApiList json.Marshal err %v ", err)
		return nil, response.NewInternalServerErrMsg(ErrID)
	}
	var apiResp []*protocol.ExchangeApiResp
	if err := json.Unmarshal(byteData, &apiResp); err != nil {
		logger.Errorf("GetExApiList json Unmarshal err %v", err)
		return nil, response.NewInternalServerErrMsg(ErrID)
	}
	for _, resp := range apiResp {
		resp.TotalUsdt = "0"
		resp.TotalRmb = "0"
		resp.BalanceDetail = make([]*pb.ExchangePos, 0)
		resp.UsdtBalance = "0"
		resp.BtcBalance = "0"
	}

	rateUsdRmb, err := client.QuoteService.GetRate(context.Background(), &empty.Empty{})
	if err != nil {
		logger.Warnf("quoteService调用GetRate 失败 %v", err)
	}
	logger.Infof("%+v", apiResp)
	for _, resp := range apiResp {
		if err != nil {
			return nil, response.NewInternalServerErrMsg(ErrID)
		}
		posResp, err := e.GetExchangePos(userId, resp.ExchangeName)
		if err != nil {
			logger.Infof("获取币种账户出错 %v", err)
			continue
		}
		balanceDe := decimal.NewFromFloat(0.0)
		for _, p := range posResp {
			totalUsdt, _ := decimal.NewFromString(p.TotalUsdt)
			balanceDe = balanceDe.Add(totalUsdt)
			if p.Symbol == "USDT" {
				resp.UsdtBalance = fmt.Sprintf("%.2f", helper.StringToFloat64(p.Available))
				continue
			}
			if p.Symbol == "BTC" {
				resp.BtcBalance = p.Available
			}
		}
		balance, _ := balanceDe.Round(2).Float64()
		resp.TotalUsdt = helper.Float64ToString(balance)
		rate := helper.StringToFloat64(rateUsdRmb.Rate)
		if rateUsdRmb != nil {
			resp.TotalRmb = helper.Float64ToString(utils.Keep2Decimal(balance * rate))
		}
		resp.BalanceDetail = posResp
	}
	return apiResp, nil
}

func (e *ExOrderService) GetExchangeClient(apiKey, apiSecretKey, apiPassphrase, exchange string) exchange_clientI.ClientI {
	if exchange == exchange_info.HUOBI {
		return apiHuobi.InitClient(apiKey, apiSecretKey, true)
	}
	if exchange == exchange_info.BINANCE {
		return api.InitClient(apiKey, apiSecretKey, apiPassphrase)
	}
	return apiBinance.InitClient(apiKey, apiSecretKey)
}

func (e *ExOrderService) GetTickWithExchange(exchange, symbol string) (*quoteCron.Ticker, error) {
	if symbol == "USDT" {
		return &quoteCron.Ticker{Last: 1}, nil
	}
	if exchange == exchange_info.OKEX {
		return e.cacheService.GetOKexQuote(fmt.Sprintf("%s%s", symbol, "-USDT"))
	}
	if exchange == exchange_info.HUOBI {
		return e.cacheService.GetHuobiQuote(fmt.Sprintf("%s%s", symbol, "-USDT"))
	}
	if exchange == exchange_info.BINANCE {
		return e.cacheService.GetBinanceQuote(fmt.Sprintf("%s%s", symbol, "-USDT"))
	}
	return nil, errors.New("exchange not valide")
}

func (e *ExOrderService) GetExchangePos(userId, exchange string) ([]*pb.ExchangePos, error) {
	posList := make([]*pb.ExchangePos, 0)

	apiList := e.dao.GetExchangeApiListByUidAndPlatform(userId, exchange) //查询数据库中okex的api信息
	for _, apiInfo := range apiList {
		secret, _ := hex.DecodeString(apiInfo.Secret)
		secretBytes, _ := encoding.AesDecrypt(secret)
		//获取对应的交易所客户端
		client := e.GetExchangeClient(apiInfo.ApiKey, string(secretBytes), apiInfo.Passphrase, exchange)

		err := utils.ReTryFunc(3, func() (bool, error) { //重试
			spotList, err := client.GetAccountSpot() //调用现货接口
			if err != nil {
				if err.Error() == "validation-format-error" {
					return false, nil
				}
				logger.Infof("GetExchangePos:GetAccount has err %v", err)
				if strings.Contains(err.Error(), "30006") {
					return true, errors.New("密钥失效或者不存在")
				}
				time.Sleep(time.Second * 1)
				return false, err
			}
			for key, value := range spotList.SubAccounts {
				if value.Balance == 0 && value.Amount == 0 { //过滤
					continue
				}
				price := 0.0
				tick, err := e.GetTickWithExchange(exchange, key.Symbol)
				if err == nil {
					price = tick.Last
				}
				pos := &pb.ExchangePos{
					Symbol:    key.Symbol,
					Balance:   decimal.NewFromFloat(value.Balance).String(),      //余额
					Available: decimal.NewFromFloat(value.Amount).String(),       //可用资金
					Frozen:    decimal.NewFromFloat(value.ForzenAmount).String(), //被冻结资金
					Price:     decimal.NewFromFloat(price).String(),
					TotalUsdt: decimal.NewFromFloat(price * value.Balance).Round(8).String(),
					Type:      "spot", //现货
				}
				balance := value.Balance
				if exchange == exchange_info.HUOBI {
					balance = value.Amount + value.ForzenAmount
				}
				pos.Balance = decimal.NewFromFloat(balance).Round(8).String()
				posList = append(posList, pos)
			}
			return false, nil
		})

		if err != nil {
			return nil, err
		}
		//time.Sleep(time.Second * 1)
		//err = utils.ReTryFunc(3, func() (bool, error) {
		//	swapList, err := client.GetAccountSwap() //查询永续合约
		//	if err != nil {
		//		if err.Error() == "validation-format-error" {
		//			return false, nil
		//		}
		//		logger.Infof("GetExchangePos:GetFutureUserinfo has err %v", err)
		//		time.Sleep(time.Second * 1)
		//		return false, err
		//	}
		//	for key, value := range swapList.SubAccounts {
		//		if value.Balance == 0 && value.Amount == 0 {
		//			continue
		//		}
		//		price := 0.0
		//		tick, err := e.cacheService.GetOKexQuote(fmt.Sprintf("%s%s", key.Symbol, "-USDT"))
		//		if err == nil {
		//			price = tick.Last
		//		}
		//		pos := &pb.ExchangePos{
		//			Symbol:    key.Symbol,
		//			Balance:   decimal.NewFromFloat(value.Balance).String(),
		//			Available: decimal.NewFromFloat(value.Amount).String(),
		//			Frozen:    decimal.NewFromFloat(value.ForzenAmount).String(),
		//			Price:     decimal.NewFromFloat(price).String(),
		//			TotalUsdt: decimal.NewFromFloat(price * value.Balance).Round(8).String(),
		//			Type:      "swap", //永续合约
		//		}
		//		posList = append(posList, pos)
		//	}
		//	return false, nil
		//})
		//if err != nil {
		//	return nil, err
		//}
	}
	return posList, nil
}

func (e *ExOrderService) UpdateExchangeApi(userID, apiKey, secret, passphrase string, exchangeID, apiID int64) error {
	exchange, err := e.dao.GetExchangeById(exchangeID)
	if err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	if exchange.Exchange == exchange_info.OKEX {
		if passphrase == "" {
			return response.NewExchangePassphraseNoneErrMsg(ErrID)
		}
	}
	if apiID == 0 {
		return response.NewUpdateExchangeApiErrMsg(ErrID)
	}
	oldAPi2, err := e.dao.GetExchangeApiByID(apiID)
	if err == nil {
		strategyList := e.dao.GetUserStrategyByApiKey(userID, oldAPi2.ApiKey)
		if len(strategyList) != 0 {
			return response.NewUpdateApiHasStrategyErrMsg(ErrID)
		}
	}
	oldAPI, _ := e.dao.GetExchangeApiByUidAndApi(userID, apiKey)
	if oldAPI != nil && oldAPI.ID != apiID {
		return response.NewExchangeApiDuplicateErrMsg(ErrID)
	}
	secretCrypt, _ := encoding.AesEncrypt([]byte(secret))

	//check
	err = e.checkIfApiValid(apiKey, secret, passphrase, exchange.Exchange)
	if err != nil {
		log.Println(err)
		return err
	}
	apiInfo := &model.WqExchangeApi{
		ID:           apiID,
		UserID:       userID,
		ExchangeID:   exchangeID,
		ExchangeName: exchange.Exchange,
		ApiKey:       apiKey,
		Secret:       hex.EncodeToString(secretCrypt),
		Passphrase:   passphrase,
		UpdatedAt:    time.Time{},
	}
	if err := e.dao.UpdateExchangeApi(apiInfo); err != nil {
		return response.NewUpdateExchangeApiErrMsg(ErrID)
	}
	//// todo 弃用 通过新的接口查
	//if err := e.dao.SetUserAllStrategyApi(userID, apiKey, exchange.Exchange); err != nil {
	//	logger.Warnf("SetUserAllStrategyApi has err %v", err)
	//}
	return nil
}

func (e *ExOrderService) DeleteExchangeApi(userID string, apiId int64) error {
	apiInfo, err := e.dao.GetExchangeApiByID(apiId)
	if err != nil {
		return response.NewDeleteExchangeApiNotFoundErrMsg(ErrID)
	}
	userStrategy := e.dao.GetUserStrategyByApiKey(userID, apiInfo.ApiKey)
	if len(userStrategy) != 0 {
		return response.NewDeleteApiHasStrategyErrMsg(ErrID)
	}
	if err := e.dao.DeleteExchangeApi(userID, apiId); err != nil {
		return response.NewDeleteExchangeApiErrMsg(ErrID)
	}
	return nil
}

func (e *ExOrderService) GetApiKeyInfo(userID, apiKey string) (*model.WqExchangeApi, error) {
	return e.dao.GetExchangeApiByUidAndApi(userID, apiKey)
}
