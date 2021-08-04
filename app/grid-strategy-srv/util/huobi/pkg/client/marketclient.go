package client

import (
	"encoding/json"
	"errors"
	"strconv"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal/requestbuilder"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/getrequest"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/market"
)

// Responsible to get market information
type MarketClient struct {
	publicUrlBuilder *requestbuilder.PublicUrlBuilder
}

// Initializer
func (p *MarketClient) Init(host string) *MarketClient {
	p.publicUrlBuilder = new(requestbuilder.PublicUrlBuilder).Init(host)
	return p
}

// Retrieves all klines in a specific range.
func (client *MarketClient) GetCandlestick(symbol string, optionalRequest getrequest.GetCandlestickOptionalRequest) ([]market.Candlestick, error) {

	request := new(getrequest.GetRequest).Init()
	request.AddParam("symbol", symbol)
	if optionalRequest.Period != "" {
		request.AddParam("period", optionalRequest.Period)
	}
	if optionalRequest.Size != 0 {
		request.AddParam("size", strconv.Itoa(optionalRequest.Size))
	}

	url := client.publicUrlBuilder.Build("/market/history/kline", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := market.GetCandlestickResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Data != nil {

		return result.Data, nil
	}
	return nil, errors.New(getResp)
}

// Retrieves the latest ticker with some important 24h aggregated market data.
func (client *MarketClient) GetLast24hCandlestickAskBid(symbol string) (*market.CandlestickAskBid, error) {

	request := new(getrequest.GetRequest).Init()
	request.AddParam("symbol", symbol)

	url := client.publicUrlBuilder.Build("/market/detail/merged", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := market.GetLast24hCandlestickAskBidResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Tick != nil {

		return result.Tick, nil
	}

	return nil, errors.New(getResp)
}

// Retrieve the latest tickers for all supported pairs.
func (client *MarketClient) GetAllSymbolsLast24hCandlesticksAskBid() ([]market.SymbolCandlestick, error) {

	request := new(getrequest.GetRequest).Init()

	url := client.publicUrlBuilder.Build("/market/tickers", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := market.GetAllSymbolsLast24hCandlesticksAskBidResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Data != nil {

		return result.Data, nil
	}

	return nil, errors.New(getResp)
}

// Retrieves the current order book of a specific pair
func (client *MarketClient) GetDepth(symbol string, step string, optionalRequest getrequest.GetDepthOptionalRequest) (*market.Depth, error) {

	request := new(getrequest.GetRequest).Init()
	request.AddParam("symbol", symbol)
	request.AddParam("type", step)
	if optionalRequest.Size != 0 {
		request.AddParam("depth", strconv.Itoa(optionalRequest.Size))
	}

	url := client.publicUrlBuilder.Build("/market/depth", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := market.GetDepthResponse{}

	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Tick != nil {

		return result.Tick, nil
	}

	return nil, errors.New(getResp)
}

// Retrieves the latest trade with its price, volume, and direction.
func (client *MarketClient) GetLatestTrade(symbol string) (*market.TradeTick, error) {

	request := new(getrequest.GetRequest).Init()
	request.AddParam("symbol", symbol)

	url := client.publicUrlBuilder.Build("/market/trade", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := market.GetLatestTradeResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Tick != nil {

		return result.Tick, nil
	}

	return nil, errors.New(getResp)
}

// Retrieves the most recent trades with their price, volume, and direction.
func (client *MarketClient) GetHistoricalTrade(symbol string, optionalRequest getrequest.GetHistoricalTradeOptionalRequest) ([]market.TradeTick, error) {

	request := new(getrequest.GetRequest).Init()
	request.AddParam("symbol", symbol)
	if optionalRequest.Size != 0 {
		request.AddParam("size", strconv.Itoa(optionalRequest.Size))
	}

	url := client.publicUrlBuilder.Build("/market/history/trade", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := market.GetHistoricalTradeResponse{}

	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Data != nil {

		return result.Data, nil
	}

	return nil, errors.New(getResp)
}

// Retrieves the summary of trading in the market for the last 24 hours.
func (client *MarketClient) GetLast24hCandlestick(symbol string) (*market.Candlestick, error) {

	request := new(getrequest.GetRequest).Init()
	request.AddParam("symbol", symbol)

	url := client.publicUrlBuilder.Build("/market/detail", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := market.GetLast24hCandlestick{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if result.Status == "ok" && result.Tick != nil {

		return result.Tick, nil
	}

	return nil, errors.New(getResp)
}
