package htx_apis

import (
	"encoding/json"
	"fmt"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
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
