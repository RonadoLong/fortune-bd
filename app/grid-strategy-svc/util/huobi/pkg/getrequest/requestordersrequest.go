package getrequest

type RequestOrdersRequest struct {
	Op        string `json:"op"`
	Topic     string `json:"topic"`
	ClientId  string `json:"cid"`
	AccountId int    `json:"account-id"`
	Symbol    string `json:"symbol"`
	Types     string `json:"types"`
	States    string `json:"states"`
	StartDate string `json:"start-date"`
	EndDate   string `json:"end-date"`
	From      string `json:"from"`
	Direct    string `json:"direct"`
	Size      string `json:"size"`
}
