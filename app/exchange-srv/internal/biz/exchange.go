package biz

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
	pb "wq-fotune-backend/api/exchange"
	"wq-fotune-backend/api/protocol"
	"wq-fotune-backend/api/response"
	"wq-fotune-backend/app/exchange-srv/cache"
	"wq-fotune-backend/app/exchange-srv/client"
	"wq-fotune-backend/app/exchange-srv/internal/model"
	quoteCron "wq-fotune-backend/app/quote-srv/cron"
	"wq-fotune-backend/libs/encoding"
	"wq-fotune-backend/libs/exchange"
	"wq-fotune-backend/libs/exchangeclient"
	"wq-fotune-backend/libs/helper"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/utils"
)

func (e *ExOrderRepo) GetExchangeInfo() ([]*model.WqExchange, error) {
	exchangeList := e.dao.GetExchangeInfo()
	if len(exchangeList) == 0 {
		logger.Infof("ExChangeInfo 没有找到")
		return nil, response.NewExchangeNotFoundErrMsg(ErrID)
	}
	return exchangeList, nil
}

func (e *ExOrderRepo) checkIfApiValid(apiKey, secret, passphrase, exchange string) error {
	exchangeClient := e.GetExchangeClient(apiKey, secret, passphrase, exchange)
	if exchangeClient == nil {
		logger.Info("No exchange exchangeClient")
		return errors.New("exchange exchangeClient nil")
	}
	if err := exchangeClient.CheckIfApiValid(); err != nil {
		return response.NewExchangeApiCheckErrMsg(ErrID)
	}
	return nil
}

func (e *ExOrderRepo) AddExchangeApi(userID, apiKey, secret, passphrase string, exchangeID int64) error {
	ex, err := e.dao.GetExchangeById(exchangeID)
	if err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	if ex.Exchange == exchange.OKEX {
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
	err = e.checkIfApiValid(apiKey, secret, passphrase, ex.Exchange)
	if err != nil {
		log.Println(err)
		return err
	}

	secretCrypt, _ := encoding.AesEncrypt([]byte(secret))

	newAPI := &model.WqExchangeApi{
		UserID:       userID,
		ExchangeID:   exchangeID,
		ExchangeName: ex.Exchange,
		ApiKey:       apiKey,
		Secret:       hex.EncodeToString(secretCrypt),
		Passphrase:   passphrase,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
	if err := e.dao.AddExchangeApi(newAPI); err != nil {
		return response.NewExchangeApiCreateErrMsg(ErrID)
	}
	return nil
}

func (e *ExOrderRepo) GetExchangeAccountListFromCache(userID string) []byte {
	ret := cache.GetExchangeAccountList(userID)
	return ret
}

func (e *ExOrderRepo) SetExchangeAccountListCache(userID string, data []byte) {
	cache.CacheExchangeAccountList(userID, data)
}

func (e *ExOrderRepo) GetExchangeApiList(userId string) ([]*protocol.ExchangeApiResp, error) {
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

	rateUsdRmb, err := client.GetQuoteService().GetRate(context.Background(), &empty.Empty{})
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

func (e *ExOrderRepo) GetExchangeClient(apiKey, apiSecretKey, apiPassphrase, ex string) exchangeclient.ExClientI {
	if ex == exchange.HUOBI {
		return exchangeclient.InitHuobi(apiKey, apiSecretKey, true)
	}
	if ex == exchange.OKEX {
		return exchangeclient.InitOKEX(apiKey, apiSecretKey, apiPassphrase)
	}
	return exchangeclient.InitBinance(apiKey, apiSecretKey)
}

func (e *ExOrderRepo) GetTickWithExchange(ex, symbol string) (*quoteCron.Ticker, error) {
	if symbol == "USDT" {
		return &quoteCron.Ticker{Last: 1}, nil
	}
	if ex == exchange.OKEX {
		return cache.GetOKexQuote(fmt.Sprintf("%s%s", symbol, "-USDT"))
	}
	if ex == exchange.HUOBI {
		return cache.GetHuobiQuote(fmt.Sprintf("%s%s", symbol, "-USDT"))
	}
	if ex == exchange.BINANCE {
		return cache.GetBinanceQuote(fmt.Sprintf("%s%s", symbol, "-USDT"))
	}
	return nil, errors.New("exchange not valide")
}

func (e *ExOrderRepo) GetExchangePos(userId, ex string) ([]*pb.ExchangePos, error) {
	posList := make([]*pb.ExchangePos, 0)

	apiList := e.dao.GetExchangeApiListByUidAndPlatform(userId, ex) //查询数据库中okex的api信息
	for _, apiInfo := range apiList {
		secret, _ := hex.DecodeString(apiInfo.Secret)
		secretBytes, _ := encoding.AesDecrypt(secret)
		//获取对应的交易所客户端
		exchangeClient := e.GetExchangeClient(apiInfo.ApiKey, string(secretBytes), apiInfo.Passphrase, ex)
		if exchangeClient == nil {
			logger.Info("No exchange client")
			continue
		}
		err := utils.ReTryFunc(3, func() (bool, error) { //重试
			spotList, err := exchangeClient.GetAccountSpot() //调用现货接口
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
				tick, err := e.GetTickWithExchange(ex, key.Symbol)
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
				if ex == exchange.HUOBI {
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

	}
	return posList, nil
}

func (e *ExOrderRepo) UpdateExchangeApi(userID, apiKey, secret, passphrase string, exchangeID, apiID int64) error {
	ex, err := e.dao.GetExchangeById(exchangeID)
	if err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	if ex.Exchange == exchange.OKEX {
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
	err = e.checkIfApiValid(apiKey, secret, passphrase, ex.Exchange)
	if err != nil {
		log.Println(err)
		return err
	}
	apiInfo := &model.WqExchangeApi{
		ID:           apiID,
		UserID:       userID,
		ExchangeID:   exchangeID,
		ExchangeName: ex.Exchange,
		ApiKey:       apiKey,
		Secret:       hex.EncodeToString(secretCrypt),
		Passphrase:   passphrase,
		UpdatedAt:    time.Time{},
	}
	if err := e.dao.UpdateExchangeApi(apiInfo); err != nil {
		return response.NewUpdateExchangeApiErrMsg(ErrID)
	}

	return nil
}

func (e *ExOrderRepo) DeleteExchangeApi(userID string, apiId int64) error {
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

func (e *ExOrderRepo) GetApiKeyInfo(userID, apiKey string) (*model.WqExchangeApi, error) {
	return e.dao.GetExchangeApiByUidAndApi(userID, apiKey)
}
