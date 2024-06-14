package htx_wss

import "encoding/json"

type Values struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}
type Ticker struct {
	Exchange string  `json:"exchange"`
	Symbol   string  `json:"symbol"`
	Buy      Values  `json:"buy"`
	Sell     Values  `json:"sell"`
	UpdateAt float64 `json:"update_at"`
}
type ReciveData struct {
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
	Ticker   Ticker `json:"ticker"`
}

type ReciveBalanceMsg struct {
	Exchange  string  `json:"exchange"`
	Symbol    string  `json:"symbol"`
	Available float64 `json:"available"`
	AccountId int64   `json:"accountId"`
	Balance   float64 `json:"balance"`
}

type ReciveAccountsMsg struct {
	Symbol        string      `json:"symbol"`
	ContractCode  string      `json:"contract_code"`
	MarginBalance json.Number `json:"margin_balance"`
	MarginStatic  json.Number `json:"margin_static"`
}
