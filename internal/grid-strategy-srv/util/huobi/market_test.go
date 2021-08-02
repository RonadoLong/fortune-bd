// https://wq-grid-strategy/util/huobi

package huobi

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"wq-fotune-backend/internal/grid-strategy-srv/util/huobi/pkg/getrequest"

	"github.com/k0kubun/pp"
)

func TestGetSymbols(t *testing.T) {
	resp, err := GetSymbols()
	if err != nil {
		t.Error(err)
		return
	}

	symbols := "\""
	for _, v := range resp {
		v.Symbol = strings.ToLower(v.Symbol)
		if v.Symbol[len(v.Symbol)-4:] == "usdt" {
			symbols += v.Symbol + "\"" + ", \""
		}
	}

	//pp.Println(len(resp), resp)
	fmt.Println(symbols)
}

func TestGetLatestPrice(t *testing.T) {
	price, err := GetLatestPrice("btcusdt")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(price)
}

func TestGetKlineInfo(t *testing.T) {
	klines, err := GetKlines("ltcusdt", getrequest.MON1, 3)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(klines)
}

func TestGetLatestKline(t *testing.T) {
	lowPrice, highPrice, err := Get3CyclePriceRange("atomusdt", MON1)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(lowPrice, highPrice)
}

func TestGetAnchorCurrencySymbols(t *testing.T) {
	anchorCurrency := "eth"

	limitSymbols, err := GetAnchorCurrencySymbols(anchorCurrency)
	if err != nil {
		t.Error(err)
		return
	}

	symbols := []string{}
	for _, v := range limitSymbols {
		v.Symbol = strings.ToLower(v.Symbol)
		symbols = append(symbols, v.Symbol)

	}
	sort.Strings(symbols)

	symbolStr := "\""
	for _, symbol := range symbols {
		symbolStr += symbol + "\"" + ", \""
	}

	fmt.Println(symbolStr)
}
