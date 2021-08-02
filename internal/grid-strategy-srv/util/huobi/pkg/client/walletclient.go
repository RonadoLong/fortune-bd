package client

import (
	"encoding/json"
	"errors"
	"strconv"
	"wq-fotune-backend/internal/grid-strategy-srv/util/huobi/internal"
	"wq-fotune-backend/internal/grid-strategy-srv/util/huobi/internal/requestbuilder"
	"wq-fotune-backend/internal/grid-strategy-srv/util/huobi/pkg/getrequest"
	"wq-fotune-backend/internal/grid-strategy-srv/util/huobi/pkg/postrequest"
	"wq-fotune-backend/internal/grid-strategy-srv/util/huobi/pkg/response/wallet"
)

// Responsible to operate wallet
type WalletClient struct {
	privateUrlBuilder *requestbuilder.PrivateUrlBuilder
}

// Initializer
func (p *WalletClient) Init(accessKey string, secretKey string, host string) *WalletClient {
	p.privateUrlBuilder = new(requestbuilder.PrivateUrlBuilder).Init(accessKey, secretKey, host)
	return p
}

// Get deposit address of corresponding chain, for a specific crypto currency (except IOTA)
func (p *WalletClient) GetDepositAddress(currency string) ([]wallet.DepositAddress, error) {
	request := new(getrequest.GetRequest).Init()

	request.AddParam("currency", currency)

	url := p.privateUrlBuilder.Build("GET", "/v2/account/deposit/address", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := wallet.GetDepositAddressResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Code == 200 && result.Data != nil {
		return result.Data, nil
	}
	return nil, errors.New(getResp)
}

// Parent user query sub user deposit address of corresponding chain, for a specific crypto currency (except IOTA)
func (p *WalletClient) GetSubUserDepositAddress(subUid int64, currency string) ([]wallet.DepositAddress, error) {
	request := new(getrequest.GetRequest).Init()
	request.AddParam("subUid", strconv.FormatInt(subUid, 10))
	request.AddParam("currency", currency)

	url := p.privateUrlBuilder.Build("GET", "/v2/sub-user/deposit-address", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := wallet.GetDepositAddressResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Code == 200 && result.Data != nil {
		return result.Data, nil
	}
	return nil, errors.New(getResp)
}

// Query withdraw quota for currencies
func (p *WalletClient) GetWithdrawQuota(currency string) (*wallet.WithdrawQuota, error) {
	request := new(getrequest.GetRequest).Init()

	request.AddParam("currency", currency)

	url := p.privateUrlBuilder.Build("GET", "/v2/account/withdraw/quota", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := wallet.GetWithdrawQuotaResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Code == 200 && result.Data != nil {
		return result.Data, nil
	}
	return nil, errors.New(getResp)
}

// Withdraw from spot trading account to an external address.
func (p *WalletClient) CreateWithdraw(request postrequest.CreateWithdrawRequest) (int64, error) {
	postBody, jsonErr := postrequest.ToJson(request)

	url := p.privateUrlBuilder.Build("POST", "/v1/dw/withdraw/api/create", nil)
	postResp, postErr := internal.HttpPost(url, postBody)
	if postErr != nil {
		return 0, postErr
	}

	result := wallet.CreateWithdrawResponse{}
	jsonErr = json.Unmarshal([]byte(postResp), &result)
	if jsonErr != nil {
		return 0, jsonErr
	}

	if result.Status == "ok" && result.Data != 0 {
		return result.Data, nil
	}
	return 0, errors.New(postResp)
}

// Cancels a previously created withdraw request by its transfer id.
func (p *WalletClient) CancelWithdraw(withdrawId int64) (int64, error) {

	url := p.privateUrlBuilder.Build("POST", "/v1/dw/withdraw-virtual/"+strconv.FormatInt(withdrawId, 10)+"}/cancel", nil)
	postResp, postErr := internal.HttpPost(url, "")
	if postErr != nil {
		return 0, postErr
	}
	result := wallet.CancelWithdrawResponse{}
	jsonErr := json.Unmarshal([]byte(postResp), &result)
	if jsonErr != nil {
		return 0, jsonErr
	}

	if result.Status == "ok" && result.Data != 0 {
		return result.Data, nil
	}
	return 0, errors.New(postResp)

}

// Returns all existed withdraws and deposits and return their latest status.
func (p *WalletClient) QueryDepositWithdraw(depositOrWithdraw string, optionalRequest getrequest.QueryDepositWithdrawOptionalRequest) ([]wallet.DepositWithdraw, error) {
	request := new(getrequest.GetRequest).Init()

	request.AddParam("type", depositOrWithdraw)

	if optionalRequest.Currency != "" {
		request.AddParam("currency", optionalRequest.Currency)
	}
	if optionalRequest.From != "" {
		request.AddParam("from", optionalRequest.From)
	}
	if optionalRequest.Direct != "" {
		request.AddParam("direct", optionalRequest.Direct)
	}
	if optionalRequest.Size != "" {
		request.AddParam("size", optionalRequest.Size)
	}

	url := p.privateUrlBuilder.Build("GET", "/v1/query/deposit-withdraw", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := wallet.QueryDepositWithdrawResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Data != nil {
		return result.Data, nil
	}
	return nil, errors.New(getResp)
}

// Parent user query sub user deposits history
func (p *WalletClient) QuerySubUserDepositHistory(subUid int64, optionalRequest getrequest.QuerySubUserDepositHistoryOptionalRequest) ([]wallet.DepositHistory, error) {
	request := new(getrequest.GetRequest).Init()

	request.AddParam("subUid", strconv.FormatInt(subUid, 10))

	if optionalRequest.Currency != "" {
		request.AddParam("currency", optionalRequest.Currency)
	}
	if optionalRequest.StartTime != 0 {
		request.AddParam("startTime", strconv.FormatInt(optionalRequest.StartTime, 10))
	}
	if optionalRequest.EndTime != 0 {
		request.AddParam("endTime", strconv.FormatInt(optionalRequest.EndTime, 10))
	}
	if optionalRequest.Sort != "" {
		request.AddParam("sort", optionalRequest.Sort)
	}
	if optionalRequest.Limit != "" {
		request.AddParam("limit", optionalRequest.Limit)
	}
	if optionalRequest.FromId != 0 {
		request.AddParam("fromId", strconv.FormatInt(optionalRequest.FromId, 10))
	}

	url := p.privateUrlBuilder.Build("GET", "/v2/sub-user/query-deposit", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := wallet.QuerySubUserDepositHistoryResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Code == 200 && result.Data != nil {
		return result.Data, nil
	}
	return nil, errors.New(getResp)
}
