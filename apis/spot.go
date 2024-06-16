package htx_apis

import (
	"encoding/json"
	"fmt"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
)

type ApiResponseListData struct {
	Status  string           `json:"status"`
	Message string           `json:"err_msg"`
	Data    []map[string]any `json:"data"`
}

func (r *ApiResponseListData) Success() bool {
	return r.Status == "ok"
}

// ### 获取现货交易对
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec456a5-7773-11ed-9966-0242ac110003
func GetSpotSymbols() (data []map[string]any, err error) {
	const symbol = "HTX GetSpotSymbols"
	body, _, err := htx.ApiConfig.Get(gateway_huobiPro, "/v1/settings/common/market-symbols", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}

	res := ApiResponseListData{}
	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()
	err = d.Decode(&res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v", symbol, res.Message)
		return
	}
	// 过滤
	data = []map[string]any{}
	for _, v := range res.Data {
		if v["state"] == "online" && v["qc"] == "usdt" {
			data = append(data, v)
		}
	}
	return data, nil
}

// MODEL 获取现货账户余额
type AssetBalance struct {
	Currency  string `json:"currency"`
	Type      string `json:"type"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Debt      string `json:"debt"`
	SeqNum    string `json:"seq-num"`
}

type AccountSpotData struct {
	Id    int    `json:"id"`
	Type  string `json:"type"`
	State string `json:"state"`
	List  []AssetBalance
}

type ApiResponseSpotAccountBalance struct {
	htx.ApiResponseV1
	Data AccountSpotData `json:"data"`
}

// ### 获取现货账户余额
// https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec40922-7773-11ed-9966-0242ac110003
func GetSpotAccountBalance(account_id int) (data []map[string]any, err error) {
	const symbol = "HTX GetSpotAccountBalance"
	body, _, err := htx.ApiConfig.Get(gateway_huobiPro, fmt.Sprintf("/v1/account/accounts/%d/balance", account_id), nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}

	// util.WriteTestJsonFile(symbol, body)

	res := ApiResponseSpotAccountBalance{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v", symbol, res.Message)
		return
	}
	// 过滤
	data = []map[string]any{}
	tmp := map[string]map[string]any{}
	for _, v := range res.Data.List {
		syb := strings.ToUpper(v.Currency)
		typ := v.Type
		if _, ok := tmp[syb]; !ok {
			tmp[syb] = map[string]any{
				"symbol": syb,
			}
		}
		tmp[syb][typ] = v.Balance
	}
	for _, v := range tmp {
		data = append(data, v)
	}
	return data, nil
}

// ### 现货下单
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec44bc5-7773-11ed-9966-0242ac110003
type ApiResponseV1String struct {
	// 定义响应结构体
	Data string `json:"data"`
	htx.ApiResponseV1
}

func SpotBuyMarket(symb string, amount float64) (data string, err error) {
	// 市价买入
	const symbol = "HTX SpotBuyMarket"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v1/order/orders/place", map[string]any{
		"account-id": htx.ApiConfig.AccountId,
		"symbol":     fmt.Sprintf("%susdt", strings.ToLower(symb)),
		"type":       "buy-market",
		"amount":     amount,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := ApiResponseV1String{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v", symbol, res.Message)
		return
	}

	return res.Data, nil
}

func SpotSellMarket(symb string, volume float64) (data string, err error) {
	// 市价卖出
	const symbol = "HTX SpotSellMarket"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v1/order/orders/place", map[string]any{
		"account-id": htx.ApiConfig.AccountId,
		"symbol":     fmt.Sprintf("%susdt", strings.ToLower(symb)),
		"type":       "sell-market",
		"amount":     volume,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := ApiResponseV1String{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v", symbol, res.Message)
		return
	}

	return res.Data, nil
}

type ApiResponseV1No struct {
	// 定义响应结构体
	Data int64 `json:"data"`
	htx.ApiResponseV1
}
type ApiResponseV1List struct {
	// 定义响应结构体
	Data []map[string]any `json:"data"`
	htx.ApiResponseV1
}

// 查询借币币息率及额度（逐仓）
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec43494-7773-11ed-9966-0242ac110003
func GetSpotMarginLoanInfo(symbols []string) (data []map[string]any, err error) {
	const flag = "HTX GetSpotMarginLoanInfo"
	list := []string{}
	for _, symbol := range symbols {
		list = append(list, fmt.Sprintf("%susdt", strings.ToLower(symbol)))
	}
	body, _, err := htx.ApiConfig.Get(gateway_huobiPro, "/v1/margin/loan-info", map[string]any{
		"symbols": strings.Join(list, ","),
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	fmt.Printf("string(body): %v\n", string(body))
	res := ApiResponseV1List{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v", flag, res.Message)
		return
	}
	return res.Data, nil

}

//	(申请借币（逐仓）)
//
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec438b9-7773-11ed-9966-0242ac110003
func SpotBorrow(symbol string, amount float64) (data int64, err error) {
	const flag = "HTX SpotBorrow"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v1/margin/orders", map[string]any{
		"symbol":   fmt.Sprintf("%susdt", strings.ToLower(symbol)),
		"currency": strings.ToLower(symbol),
		"amount":   amount,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	fmt.Printf("string(body): %v\n", string(body))
	res := ApiResponseV1No{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v", flag, res.Message)
		return
	}
	return res.Data, nil
}

//	(归还借币（逐仓）)
//
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec438b9-7773-11ed-9966-0242ac110003
func SpotReturn(orderId int64, amount float64) (data int64, err error) {
	const flag = "HTX SpotBorrow"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, fmt.Sprintf("/v1/margin/orders/%d/repay", orderId), map[string]any{
		"amount": amount,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	fmt.Printf("string(body): %v\n", string(body))
	res := ApiResponseV1No{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v", flag, res.Message)
		return
	}
	return res.Data, nil
}
