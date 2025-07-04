package htx_apis

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	htx "github.com/laoliu6668/esharp_htx_utils"

	"github.com/laoliu6668/esharp_htx_utils/util"
)

const gateway_hbdm = "api.hbdm.com"

// MODEL
type ApiResponseSwapData struct {
	htx.ApiResponseHBDM
	Data map[string]any `json:"data"`
}
type ApiResponseSwapListData struct {
	htx.ApiResponseHBDM
	Data []map[string]any `json:"data"`
}

// ## 获取合约交易对
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb72f34-77b5-11ed-9966-0242ac110003
func GetSwapSymbol() (data []map[string]any, err error) {
	const symbol = "HTX GetSwapSymbol"
	body, _, err := htx.ApiConfig.Get(gateway_hbdm, "/linear-swap-api/v1/swap_contract_info", map[string]any{
		"support_margin_mode": "all",
		"contract_type":       "swap",
		"business_type":       "swap",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}
	// util.WriteTestJsonFile(symbol, body)

	res := ApiResponseSwapListData{}
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
		var status json.Number = "1"
		if v["contract_status"] == status && v["trade_partition"] == "USDT" {
			data = append(data, v)
		}
	}
	return
}

// MODEL 获取用户持仓量限制
type ApiResponseSwapPositionLimit struct {
	htx.ApiResponseHBDM
	Data []map[string]any `json:"data"`
}

// ## 获取用户持仓量限制
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb7649f-77b5-11ed-9966-0242ac110003
func GetSwapPositionLimit(symb string) (data []map[string]any, err error) {
	const symbol = "HTX GetSwapPositionLimit"
	body, _, err := htx.ApiConfig.Post(gateway_hbdm, "/linear-swap-api/v1/swap_position_limit", map[string]any{
		"contract_code": fmt.Sprintf("%v-USDT", strings.ToUpper(symb)),
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}
	// util.WriteTestJsonFile(symbol, body)

	res := ApiResponseSwapListData{}
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
		if v["trade_partition"] == "USDT" {
			data = append(data, v)
		}
	}
	return
}

type SwapOrderLimitData struct {
	List []map[string]any `json:"list"`
}

type ApiResponseSwapOrderLimit struct {
	htx.ApiResponseHBDM
	Data SwapOrderLimitData `json:"data"`
}

// ## 获取用户下单量限制
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb76025-77b5-11ed-9966-0242ac110003
func GetSwapOrderLimit() (data []map[string]any, err error) {
	const symbol = "HTX GetSwapOrderLimit"
	body, _, err := htx.ApiConfig.Post(gateway_hbdm, "/linear-swap-api/v1/swap_order_limit", map[string]any{
		"order_price_type": "optimal_10",
		"contract_type":    "swap",
		"business_type":    "swap",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ApiResponseSwapOrderLimit{}
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
	for _, v := range res.Data.List {
		if v["trade_partition"] == "USDT" {
			data = append(data, v)
		}
	}
	return data, nil
}

// ## 获取合约资金费率
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb71b45-77b5-11ed-9966-0242ac110003
func GetSwapFundingRate() (data []map[string]any, err error) {
	const symbol = "HTX GetSwapFundingRate"
	body, _, err := htx.ApiConfig.Get(gateway_hbdm, "/linear-swap-api/v1/swap_batch_funding_rate", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}
	// util.WriteTestJsonFile(symbol, body)

	res := ApiResponseSwapListData{}
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
		if len(strings.Split(fmt.Sprintf("%s", v["contract_code"]), "-")) != 2 {
			continue
		}
		if v["trade_partition"] == "USDT" {
			data = append(data, v)
		}
	}
	return data, nil
}

// ## 【逐仓】获取用户账户信息
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb74886-77b5-11ed-9966-0242ac110003
func GetSwapAccountInfo(symb string) (data []map[string]any, err error) {
	const symbol = "HTX GetSwapAccountInfo"
	if symb != "" {
		symb = symb + "-USDT"
	}
	body, _, err := htx.ApiConfig.Post(gateway_hbdm, "/linear-swap-api/v1/swap_account_info", map[string]any{
		"contract_code": symb,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ApiResponseSwapListData{}
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

	return res.Data, nil
}

// ## 【逐仓】获取用户持仓信息
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb74886-77b5-11ed-9966-0242ac110003
func GetSwapPositionInfo(symb string) (data []map[string]any, err error) {
	const symbol = "HTX GetSwapPositionInfo"
	if symb != "" {
		symb = symb + "-USDT"
	}
	body, _, err := htx.ApiConfig.Post(gateway_hbdm, "/linear-swap-api/v1/swap_position_info", map[string]any{
		"contract_code": symb,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ApiResponseSwapListData{}
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

	return res.Data, nil
}

// ##【逐仓】获取用户账户和持仓信息
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb74886-77b5-11ed-9966-0242ac110003
func GetSwapAccountPositionInfo(symb string) (data []map[string]any, err error) {
	const symbol = "HTX GetSwapAccountPositionInfo"
	if symb != "" {
		symb = symb + "-USDT"
	}
	body, _, err := htx.ApiConfig.Post(gateway_hbdm, "/linear-swap-api/v1/swap_account_position_info", map[string]any{
		"contract_code": symb,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ApiResponseSwapListData{}
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

	return res.Data, nil
}

type SwapAccountTypeData struct {
	AccountType int `json:"account_type"`
}
type ApiResponseSwapAccountType struct {
	htx.ApiResponseHBDMV3
	Data SwapAccountTypeData `json:"data"`
}

// ## 获取用户账户类型
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb71825-77b5-11ed-9966-0242ac110003
func GetSwapAccountType() (accoutType int, err error) {
	const symbol = "HTX GetSwapAccountType"
	body, _, err := htx.ApiConfig.Get(gateway_hbdm, "/linear-swap-api/v3/swap_unified_account_type", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}
	res := ApiResponseSwapAccountType{}
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

	return res.Data.AccountType, nil
}

type AccountBalanceData struct {
	Balance string `json:"balance"`
}
type ApiResponseAccountBalance struct {
	htx.ApiResponseHBDM
	Data []AccountBalanceData `json:"data"`
}

// ## 获取账户总资产估值
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb74531-77b5-11ed-9966-0242ac110003
func GetSwapAccountBalance() (balance float64, err error) {
	const symbol = "HTX GetSwapAccountBalance"
	body, _, err := htx.ApiConfig.Post(gateway_hbdm, "/linear-swap-api/v1/swap_balance_valuation", map[string]any{
		"valuation_asset": "USDT",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	// fmt.Printf("string(body): %v\n", string(body))
	res := ApiResponseAccountBalance{}
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

	if len(res.Data) == 0 {
		err = fmt.Errorf("%s len err: zero len", symbol)
		return
	}

	return util.ParseFloat(res.Data[0].Balance, 0), nil
}

// 买入平空
func SwapBuyClose(symb string, volume int, lever_rate int) (orderId string, err error) {
	return SwapOrder(symb, volume, "buy", "close", lever_rate, "market")
}

// 卖出开空
func SwapSellOpen(symb string, volume int, lever_rate int) (orderId string, err error) {
	return SwapOrder(symb, volume, "sell", "open", lever_rate, "market")
}

// 买入开多
func SwapBuyOpen(symb string, volume int, lever_rate int) (data string, err error) {
	return SwapOrder(symb, volume, "buy", "open", lever_rate, "market")
}

// 卖出平多
func SwapSellClose(symb string, volume int, lever_rate int) (data string, err error) {
	return SwapOrder(symb, volume, "sell", "close", lever_rate, "market")
}

// ### 期货下单
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb77019-77b5-11ed-9966-0242ac110003
// volume 张数
func SwapOrder(symb string, volume int, direction string, offset string, lever_rate int, order_price_type string) (orderId string, err error) {
	// 买入平空
	var symbol = fmt.Sprintf("HTX Swap%s%s", direction, offset)
	body, _, err := htx.ApiConfig.PostTimeout(gateway_hbdm, "/linear-swap-api/v1/swap_order", map[string]any{
		"contract_code":    fmt.Sprintf("%s-USDT", strings.ToUpper(symb)),
		"volume":           volume,
		"direction":        direction,
		"offset":           offset,
		"lever_rate":       lever_rate,
		"order_price_type": order_price_type,
	}, time.Second)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	res := ApiResponseSwapData{}
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

	return fmt.Sprintf("%v", res.Data["order_id_str"]), nil
}
