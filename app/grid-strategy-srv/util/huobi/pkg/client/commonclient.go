package client

import (
	"encoding/json"
	"errors"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal/requestbuilder"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/getrequest"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/common"
)

// Responsible to get common information
type CommonClient struct {
	publicUrlBuilder *requestbuilder.PublicUrlBuilder
}

// Initializer
func (p *CommonClient) Init(host string) *CommonClient {
	p.publicUrlBuilder = new(requestbuilder.PublicUrlBuilder).Init(host)
	return p
}

func (p *CommonClient) GetSystemStatus() (string, error) {
	url := "https://status.huobigroup.com/api/v2/summary.json"
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return "", getErr
	}

	return getResp, nil
}

// Get all Supported Trading Symbol
// This endpoint returns all Huobi's supported trading symbol.
func (p *CommonClient) GetSymbols() ([]common.Symbol, error) {
	url := p.publicUrlBuilder.Build("/v1/common/symbols", nil)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := common.GetSymbolsResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Data != nil {
		return result.Data, nil
	}
	return nil, errors.New(getResp)
}

// Get all Supported Currencies
// This endpoint returns all Huobi's supported trading currencies.
func (p *CommonClient) GetCurrencys() ([]string, error) {
	url := p.publicUrlBuilder.Build("/v1/common/currencys", nil)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := common.GetCurrenciesResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)

	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Data != nil {
		return result.Data, nil
	}
	return nil, errors.New(getResp)
}

// APIv2 - Currency & Chains
// API user could query static reference information for each currency, as well as its corresponding chain(s). (Public Endpoint)
func (p *CommonClient) GetV2ReferenceCurrencies(optionalRequest getrequest.GetV2ReferenceCurrencies) ([]common.CurrencyChain, error) {
	request := new(getrequest.GetRequest).Init()
	if optionalRequest.Currency != "" {
		request.AddParam("currency", optionalRequest.Currency)
	}
	if optionalRequest.AuthorizedUser != "" {
		request.AddParam("authorizedUser", optionalRequest.AuthorizedUser)
	}

	url := p.publicUrlBuilder.Build("/v2/reference/currencies", request)

	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := common.GetV2ReferenceCurrenciesResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)

	if jsonErr != nil {
		return nil, jsonErr
	}

	if result.Code == 200 && result.Data != nil {
		return result.Data, nil
	}

	return nil, errors.New(result.Message)
}

// Get Current Timestamp
// This endpoint returns the current timestamp, i.e. the number of milliseconds that have elapsed since 00:00:00 UTC on 1 January 1970.
func (p *CommonClient) GetTimestamp() (int, error) {
	url := p.publicUrlBuilder.Build("/v1/common/timestamp", nil)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return 0, getErr
	}

	result := common.GetTimestampResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)

	if jsonErr != nil {
		return 0, jsonErr
	}
	if result.Status == "ok" && result.Data != 0 {
		return result.Data, nil
	}
	return 0, errors.New(getResp)
}
