package v1

import (
	"errors"
	"strconv"
	"strings"
)

type exchangeLimitForm struct {
	Exchange       string   `json:"exchange"`       // 交易所
	Symbols        []string `json:"symbols"`        // 品种名称
	AnchorCurrency string   `json:"anchorCurrency"` // 锚定币
	currency       string   `json:"-"`
}

func (e *exchangeLimitForm) valid() error {
	if e.Exchange == "" {
		return errors.New("param exchange is empty")
	}
	if e.AnchorCurrency == "" {
		return errors.New("param anchorCurrency is empty")
	}
	e.AnchorCurrency = strings.ToLower(e.AnchorCurrency)
	if len(e.Symbols) == 0 {
		return errors.New("param symbol is empty")
	}

	return nil
}

func (e *exchangeLimitForm) filterSymbols() []string {
	symbols := []string{}

	l := len(e.AnchorCurrency)
	for _, symbol := range e.Symbols {
		ls := len(symbol)
		if ls > l+1 && symbol[ls-l:] == e.AnchorCurrency {
			symbols = append(symbols, strings.ToLower(symbol))
		}
	}

	return symbols
}

type reqListForm struct {
	page     int
	pageStr  string
	limit    int
	limitStr string
	sort     string
}

func (s *reqListForm) valid() error {
	if s.pageStr == "" {
		return errors.New("field page is empty")
	}
	if s.limitStr == "" {
		return errors.New("field limit is empty")
	}

	var err error
	s.page, err = strconv.Atoi(s.pageStr)
	if err != nil {
		return err
	}
	s.limit, err = strconv.Atoi(s.limitStr)
	if err != nil {
		return err
	}
	if s.limit < 1 {
		s.limit = 1
	} else if s.limit > 200 {
		s.limit = 200
	}

	return nil
}

func str2Float64(str string) float64 {
	if str == "" {
		return 0
	}

	isPercent := false
	l := len(str)
	if str[l-1] == '%' {
		isPercent = true
		str = strings.TrimRight(str, "%")
	}
	str = strings.Replace(str, ",", "", -1)

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}

	if isPercent {
		f = f / 100
	}

	return f
}
