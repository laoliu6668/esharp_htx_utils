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

type ApiResponseSwapV5Data struct {
	htx.ApiResponseHBDMV5
	Data map[string]any `json:"data"`
}

type ApiResponseSwapV5ListData struct {
	htx.ApiResponseHBDMV5
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

// ## 获取用户账户信息
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-195703a12d5
func GetSwapAccountInfo() (data map[string]any, err error) {
	const symbol = "HTX GetSwapAccountInfo"
	body, _, err := htx.ApiConfig.Get(gateway_hbdm, "/v5/account/balance", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ApiResponseSwapV5Data{}
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

// ## 获取当前持仓信息
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-1957f1fbee4
func GetSwapPositionInfo(symb string) (data []map[string]any, err error) {
	const symbol = "HTX GetSwapPositionInfo"
	params := map[string]any{}
	if contractCode := swapContractCode(symb); contractCode != "" {
		params["contract_code"] = contractCode
	}
	body, _, err := htx.ApiConfig.Get(gateway_hbdm, "/v5/trade/position/opens", params)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ApiResponseSwapV5ListData{}
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

func swapContractCode(symb string) string {
	contractCode := strings.ToUpper(symb)
	if contractCode == "" || strings.Contains(contractCode, "-") {
		return contractCode
	}
	return contractCode + "-USDT"
}

// ## 设置杠杆等级
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-1957f4a3b67
// marginMode: "cross" 为全仓，"isolated" 为逐仓。
// positionSide: 逐仓双向持仓时必传（"long" 或 "short"）；其他场景可传空字符串。
func SetSwapLeverageRate(symb string, marginMode string, positionSide string, leverRate int) (data map[string]any, err error) {
	const symbol = "HTX SetSwapLeverageRate"
	params := map[string]any{
		"contract_code": swapContractCode(symb),
		"margin_mode":   marginMode,
		"lever_rate":    leverRate,
	}
	if positionSide != "" {
		params["position_side"] = positionSide
	}
	body, _, err := htx.ApiConfig.Post(gateway_hbdm, "/v5/position/lever", params)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ApiResponseSwapV5Data{}
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

type SwapOrderV5Data struct {
	ClientOrderId string `json:"client_order_id"`
	OrderId       string `json:"order_id"`
}

type SwapOrderV5Res struct {
	htx.ApiResponseHBDMV5
	Data SwapOrderV5Data `json:"data"`
}

// ### 期货下单
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-1957dd521e6
// margin_mode: "cross":全仓 "isolated":逐仓
// volume 张数
// side: "buy":买 "sell":卖
// position_side: "long":多 "short":空 “both”:单向持仓，开平模式必填，买卖模式默认为both。
func SwapOrderV5(coin string, margin_mode string, volume int, side string, position_side string, order_price_type string) (data SwapOrderV5Data, err error) {
	var flag = fmt.Sprintf("HTX Swap%s%s", side, position_side)
	body, _, err := htx.ApiConfig.PostTimeout(gateway_hbdm, "/v5/trade/order", map[string]any{
		"contract_code": fmt.Sprintf("%s-USDT", strings.ToUpper(coin)),
		"margin_mode":   margin_mode,
		"volume":        fmt.Sprintf("%v", volume),
		"offset":        "open",
		"side":          side,
		"position_side": position_side,
		"type":          order_price_type, //  "market": 市价，"limit":限价, "post_only":只做maker
	}, time.Second)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	res := SwapOrderV5Res{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v %s", flag, res.Message, body)
		return
	}

	return res.Data, nil
}
