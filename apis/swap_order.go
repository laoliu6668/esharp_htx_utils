package htx_apis

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	htx "github.com/laoliu6668/esharp_htx_utils"
)

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
func SwapOrderV5(coin string, margin_mode string, volume int, side string, position_side string, order_price_type string, channel_code string, client_order_id ...string) (data SwapOrderV5Data, err error) {
	var flag = fmt.Sprintf("HTX Swap%s%s", side, position_side)

	req := map[string]any{
		"contract_code": fmt.Sprintf("%s-USDT", strings.ToUpper(coin)),
		"margin_mode":   margin_mode,
		"volume":        fmt.Sprintf("%v", volume),
		"offset":        "open",
		"side":          side,
		"position_side": position_side,
		"type":          order_price_type, //  "market": 市价，"limit":限价, "post_only":只做maker
	}
	if channel_code != "" {
		req["channel_code"] = channel_code
	}
	if len(client_order_id) > 0 {
		req["client_order_id"] = client_order_id[0]
	}

	body, _, err := htx.ApiConfig.PostTimeout(gateway_hbdm, "/v5/trade/order", req, time.Second)
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
