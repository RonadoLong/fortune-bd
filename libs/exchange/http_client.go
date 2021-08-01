package exchange

import (
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Timeout:   time.Second * 5,
		Transport: &http.Transport{
			//Proxy: func(req *http.Request) (*url.URL, error) {
			//	return &url.URL{
			//		Scheme: "socks5",
			//		Host:   "192.168.123.172:1080"}, nil
			//},
		},
	}
)