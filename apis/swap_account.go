package htx_apis

import (
	"encoding/json"
	"fmt"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util"
)

// ## 获取用户账户信息
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-195703a12d5
func GetSwapAccountInfo() (data SwapAccountInfo, err error) {
	const symbol = "HTX GetSwapAccountInfo"
	body, _, err := htx.ApiConfig.Get(gateway_hbdm, "/v5/account/balance", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ApiResponseSwapAccountInfo{}
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
