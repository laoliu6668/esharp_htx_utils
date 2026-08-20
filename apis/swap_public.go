package htx_apis

import (
	"encoding/json"
	"fmt"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
)

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
