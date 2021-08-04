package session

import (
	"sync"
	"wq-fotune-backend/app/forward-offer-srv/srv/model"
)

var (
	mutex      = new(sync.Mutex)
	ClientPool = make(map[string]*Client)
)

// RegisterClientOrGetClient 注册登录的交易账户
func RegisterClientOrGetClient(info *model.ExchangeInfo) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	exitC := isExit(info.APIKey)
	if exitC != nil {
		return exitC
	}
	client := initClient(info.APIKey, info.SecretKey, info.EcPass)
	if client != nil {
		registerClient(client)
		client.SubscriptExchangeEvent()
		return client
	}
	return nil
}

func isExit(key string) *Client {
	val, ok := ClientPool[key]
	if ok && val != nil {
		return val
	}
	return nil
}

func registerClient(client *Client) {
	ClientPool[client.ApiClient.GetApiKey()] = client
}

func delClient(apiKey string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(ClientPool, apiKey)
}
