package htx_apis

import (
	"encoding/json"
	"fmt"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
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

// SwapAccountDetail 是单个保证金币种的账户资产明细。
// 金额和比率均使用字符串，避免精度丢失。
type SwapAccountDetail struct {
	Currency              string `json:"currency"`
	Equity                string `json:"equity"`
	IsolatedEquity        string `json:"isolated_equity"`
	Available             string `json:"available"`
	WithdrawAvailable     string `json:"withdraw_available"`
	ProfitUnreal          string `json:"profit_unreal"`
	IsolatedProfitUnreal  string `json:"isolated_profit_unreal"`
	InitialMargin         string `json:"initial_margin"`
	MaintenanceMargin     string `json:"maintenance_margin"`
	MaintenanceMarginRate string `json:"maintenance_margin_rate"`
	InitialMarginRate     string `json:"initial_margin_rate"`
	AvailableMargin       string `json:"available_margin"`
	Voucher               string `json:"voucher"`
	VoucherValue          string `json:"voucher_value"`
	CreatedTime           int64  `json:"created_time"`
	UpdatedTime           int64  `json:"updated_time"`
}

// SwapAccountInfo 是 U 本位联合保证金账户余额信息。
type SwapAccountInfo struct {
	State                 string              `json:"state"`
	Equity                string              `json:"equity"`
	InitialMargin         string              `json:"initial_margin"`
	MaintenanceMargin     string              `json:"maintenance_margin"`
	MaintenanceMarginRate string              `json:"maintenance_margin_rate"`
	ProfitUnreal          string              `json:"profit_unreal"`
	AvailableMargin       string              `json:"available_margin"`
	VoucherValue          string              `json:"voucher_value"`
	CreatedTime           int64               `json:"created_time"`
	UpdatedTime           int64               `json:"updated_time"`
	Details               []SwapAccountDetail `json:"details"`
}

type ApiResponseSwapAccountInfo struct {
	htx.ApiResponseHBDMV5
	Data SwapAccountInfo `json:"data"`
}

// SwapPositionInfo 是当前合约持仓信息。
// 金额、价格、数量和比率均使用字符串，避免浮点精度丢失。
type SwapPositionInfo struct {
	ContractCode      string `json:"contract_code"`
	PositionSide      string `json:"position_side"`
	Direction         string `json:"direction"`
	OpenAvgPrice      string `json:"open_avg_price"`
	MarginMode        string `json:"margin_mode"`
	Volume            string `json:"volume"`
	Available         string `json:"available"`
	LeverRate         int64  `json:"lever_rate"`
	ADLRiskPercent    *int64 `json:"adl_risk_percent"`
	LiquidationPrice  string `json:"liquidation_price"`
	InitialMargin     string `json:"initial_margin"`
	MaintenanceMargin string `json:"maintenance_margin"`
	Margin            string `json:"margin"`
	ProfitUnreal      string `json:"profit_unreal"`
	ProfitRate        string `json:"profit_rate"`
	MarginRate        string `json:"margin_rate"`
	MarginCurrency    string `json:"margin_currency"`
	LastPrice         string `json:"last_price"`
	MarkPrice         string `json:"mark_price"`
	ContractType      string `json:"contract_type"`
	CreatedTime       string `json:"created_time"`
	UpdatedTime       string `json:"updated_time"`
}

type ApiResponseSwapPositionInfo struct {
	htx.ApiResponseHBDMV5
	Data []SwapPositionInfo `json:"data"`
}

// ## 获取当前持仓信息
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-1957f1fbee4
func GetSwapPositionInfo(symb string) (data []SwapPositionInfo, err error) {
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

	res := ApiResponseSwapPositionInfo{}
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

// SwapPositionModeData 是设置后的持仓模式。
// PositionMode 的可选值为 "single_side"（单向持仓）和 "dual_side"（双向持仓）。
type SwapPositionModeData struct {
	PositionMode string `json:"position_mode"`
}

type ApiResponseSwapPositionMode struct {
	htx.ApiResponseHBDMV5
	Data SwapPositionModeData `json:"data"`
}

// SetSwapPositionMode 设置 U 本位联合保证金账户的持仓模式。
// positionMode: "single_side" 为单向持仓，"dual_side" 为双向持仓。
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-1957f4ec40b
func SetSwapPositionMode(positionMode string) (data SwapPositionModeData, err error) {
	const symbol = "HTX SetSwapPositionMode"
	positionMode = strings.ToLower(strings.TrimSpace(positionMode))
	if positionMode != "single_side" && positionMode != "dual_side" {
		return data, fmt.Errorf("%s invalid position mode: %q", symbol, positionMode)
	}

	body, _, err := htx.ApiConfig.Post(gateway_hbdm, "/v5/position/mode", map[string]any{
		"position_mode": positionMode,
	})
	if err != nil {
		return data, fmt.Errorf("%s err: %w", symbol, err)
	}

	res := ApiResponseSwapPositionMode{}
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.UseNumber()
	if err := decoder.Decode(&res); err != nil {
		return data, fmt.Errorf("%s jsonDecodeErr: %w", symbol, err)
	}
	if !res.Success() {
		return data, fmt.Errorf("%s false:%v", symbol, res.Message)
	}

	return res.Data, nil
}
