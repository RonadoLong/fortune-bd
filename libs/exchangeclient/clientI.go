package exchangeclient

import (
	"fortune-bd/libs/env"
	"fortune-bd/libs/goex"
	"net/http"
	"net/url"
	"strings"
	"time"

)

type ExClientI interface {
	GetAccountSpot() (*goex.Account, error)
	GetAccountSwap() (*goex.Account, error)
	CheckIfApiValid() error
}

//todo 代理
var (
	client = &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   strings.Split(env.ProxyAddr, "//")[1],
				}, nil
			},
		},
	}
)
