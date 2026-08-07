package htx_apis

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util"
)

const gateway_huobiPro = "api.huobi.pro"

// ### 获取用户UID
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec491c9-7773-11ed-9966-0242ac110003
func GetUserId() (uid int, err error) {
	const symbol = "HTX GetUserId"
	body, _, err := htx.ApiConfig.Get(gateway_huobiPro, "/v2/user/uid", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := htx.ApiResponseIntData{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		return 0, fmt.Errorf("%s false:%v", symbol, res.Message)
	}
	return res.Data, nil
}

// # MODEL 获取用户账户
type ApiResponseAccountData struct {
	htx.ApiResponseV1
	Data []map[string]any `json:"data"`
}

// type AccountData struct {
// 	Id      int    `json:"id"`
// 	Type    string `json:"type"`
// 	Subtype string `json:"subtype"`
// 	State   string `json:"state"`
// }

// ### ！！！AccountData - > map[string]any
// ### 获取用户账户
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec40743-7773-11ed-9966-0242ac110003
func GetUserAccount() (data []map[string]any, err error) {
	const symbol = "HTX GetUserAccount"
	body, _, err := htx.ApiConfig.Get(gateway_huobiPro, "/v1/account/accounts", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}

	// fmt.Printf("string(body): %v\n", string(body))
	res := ApiResponseAccountData{}
	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()
	err = d.Decode(&res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		return nil, fmt.Errorf("%s false:%v", symbol, res.Message)
	}

	return res.Data, nil
}

// ### 现货账户向期货账户划转
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=10000095-77b7-11ed-9966-0242ac110003
func SpotToSwapTransfer(amount float64, symb string) (no int, err error) {
	const symbol = "HTX SpotToSwapTransfer"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v2/account/transfer", map[string]any{
		"currency":       "usdt",
		"amount":         amount,
		"from":           "spot",
		"to":             "linear-swap",
		"margin-account": strings.ToLower(symb) + "-usdt",
	})
	// fmt.Printf("string(body): %v\n", string(body))
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	res := htx.ApiResponseIntData{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		return 0, fmt.Errorf("%s false:%v", symbol, res.Message)
	}
	return res.Data, nil
}

// 期货账户向现货账户划转
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=10000095-77b7-11ed-9966-0242ac110003
func SwapToSpotTransfer(amount float64, symb string) (no int, err error) {
	const symbol = "HTX SwapToSpotTransfer"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v2/account/transfer", map[string]any{
		"currency":       "usdt",
		"amount":         amount,
		"from":           "linear-swap",
		"to":             "spot",
		"margin-account": strings.ToLower(symb) + "-usdt",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))

	res := htx.ApiResponseIntData{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		return 0, fmt.Errorf("%s false:%v", symbol, res.Message)
	}
	return res.Data, nil
}

// ### 现货账户向逐仓账户划转
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec42443-7773-11ed-9966-0242ac110003
func SpotToMarginTransfer(amount float64, symb string) (no int, err error) {
	const symbol = "HTX SpotToMarginTransfer"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v1/dw/transfer-in/margin", map[string]any{
		"currency": "usdt",
		"amount":   amount,
		"symbol":   strings.ToLower(symb) + "usdt",
	})
	// fmt.Printf("string(body): %v\n", string(body))
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	res := htx.ApiResponseIntData{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		return 0, fmt.Errorf("%s false:%v", symbol, res.Message)
	}
	return res.Data, nil
}

// ### 逐仓账户向现货账户划转
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec427c6-7773-11ed-9966-0242ac110003
func SwapToMarginTransfer(amount float64, symb string) (no int, err error) {
	const symbol = "HTX SpotToMarginTransfer"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v1/dw/transfer-out/margin", map[string]any{
		"currency": "usdt",
		"amount":   amount,
		"symbol":   strings.ToLower(symb) + "usdt",
	})
	// fmt.Printf("string(body): %v\n", string(body))
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	res := htx.ApiResponseIntData{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		return 0, fmt.Errorf("%s false:%v", symbol, res.Message)
	}
	return res.Data, nil
}

// ### U本位统一账户 向 现货账户划转
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb7be8a-77b5-11ed-9966-0242ac110003
func SwapToSpotTransferUnified(amount float64) (no int, err error) {
	const symbol = "HTX SwapToSpotTransferUnified"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v2/account/transfer ", map[string]any{
		"to":             "linear-swap",
		"from":           "spot",
		"currency":       "usdt",
		"amount":         amount,
		"margin-account": "USDT",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	res := htx.ApiResponseIntData{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		return 0, fmt.Errorf("%s false:%v %s", symbol, res.Message, string(body))
	}
	return res.Data, nil
}

// ### 现货账户 向 U本位统一账户划转
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb7be8a-77b5-11ed-9966-0242ac110003
func SpotToSwapTransferUnified(amount float64) (no int, err error) {
	const symbol = "HTX SpotToSwapTransferUnified"
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v2/account/transfer ", map[string]any{
		"from":           "linear-swap",
		"to":             "spot",
		"currency":       "usdt",
		"amount":         amount,
		"margin-account": "USDT",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	res := htx.ApiResponseIntData{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		return 0, fmt.Errorf("%s false:%v %s", symbol, res.Message, string(body))
	}
	return res.Data, nil
}

type AccountTotal struct {
	AccountBalanceUsdt string `json:"accountBalanceUsdt"`
}
type TotalData struct {
	ProfitAccountBalanceList []AccountTotal `json:"profitAccountBalanceList"`
}
type ApiResponseAccountTotal struct {
	Data    TotalData `json:"data"`
	Success bool      `json:"success"`
	Code    int       `json:"code"`
}

// ## 获取账户总资产估值
// https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec46584-7773-11ed-9966-0242ac110003
func GetAccountTotalValue() (balance float64, err error) {
	const symbol = "HTX GetAccountTotalValue"
	body, _, err := htx.ApiConfig.GetTimeout(gateway_huobiPro, "/v2/account/valuation", map[string]any{
		// "valuationCurrency": "BTC",
	}, time.Second*10)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ApiResponseAccountTotal{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success {
		err = fmt.Errorf("%s false:%v", symbol, res.Code)
		return
	}
	total := 0.0
	for _, v := range res.Data.ProfitAccountBalanceList {
		total += util.ParseFloat(v.AccountBalanceUsdt, 0)
	}

	return total, nil
}

// ## 账户类型更改
// https://www.htx.com/zh-cn/opend/newApiPages/?id=10000081-77b7-11ed-9966-0242ac110003
func SwitchAccountType() (err error) {
	const symbol = "HTX SwitchAccountType"
	body, _, err := htx.ApiConfig.PostTimeout(gateway_hbdm, "/linear-swap-api/v3/swap_switch_account_type", map[string]any{
		"account_type": 1,
	}, time.Second*10)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := htx.ApiResponseHBDMV3{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v %s \n resp: %s", symbol, res.Code, res.Message, string(body))
		return
	}

	return
}

type ResAccountType struct {
	htx.ApiResponseHBDMV3
	Data struct {
		AccountType int `json:"account_type"`
	} `json:"data"`
}

// ## 查询账户类型更改
// https://www.htx.com/zh-cn/opend/newApiPages/?id=10000079-77b7-11ed-9966-0242ac110003
func GetAccountType() (accountType int, err error) {
	const symbol = "HTX GetAccountType"
	body, _, err := htx.ApiConfig.GetTimeout(gateway_hbdm, "/linear-swap-api/v3/swap_unified_account_type", map[string]any{}, time.Second*10)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := ResAccountType{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v", symbol, res.Code)
		return
	}

	return res.Data.AccountType, nil
}

// ## 账户类型更改
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-1957dd1a995
// asset_mode: 0:单币种保证金模式（老）1:联合保证金模式 2:单币种保证金模式(新）
func SwitchSwapAccountTypeV5(asset_mode int) (err error) {
	const symbol = "HTX SwitchSwapAccountTypeV5"
	body, _, err := htx.ApiConfig.PostTimeout(gateway_hbdm, "/v5/account/asset_mode", map[string]any{
		"asset_mode": asset_mode,
	}, time.Second*10)
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		return
	}

	res := htx.ApiResponseHBDMV3{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false:%v %s \n resp: %s", symbol, res.Code, res.Message, string(body))
		return
	}

	return
}

// ### 发起万向划转
// doc: https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-19e1b02387c
// assetType[0] = fromAssetType, assetType[1] = toAssetType
func UniversalTransfer(fromAccountType, toAccountType, currency string, amount float64, assetType ...string) (err error) {
	const symbol = "HTX UniversalTransfer"

	var fromAssetType, toAssetType string
	if len(assetType) > 0 {
		fromAssetType = assetType[0]
	}
	if len(assetType) > 1 {
		toAssetType = assetType[1]
	}
	body, _, err := htx.ApiConfig.Post(gateway_huobiPro, "/v5/account/universal_transfer", map[string]any{
		"from_account_type": fromAccountType,
		"to_account_type":   toAccountType,
		"currency":          currency,
		"amount":            amount,
		"from_asset_type":   fromAssetType,
		"to_asset_type":     toAssetType,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	res := htx.ApiResponseV5{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if !res.Success() {
		err = fmt.Errorf("%s false: %v %s", symbol, res.Message, string(body))
		return
	}
	return
}
