package htx_wss

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// 【通用】订阅资金费率推送(免鉴权)（sub）
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb7127e-77b5-11ed-9966-0242ac110003
func SubSwapFundingRate(reciveHandle func(ReciveSwapFundingRateMsg), logHandle func(string), errHandle func(error)) {
	title := "SubSwapFundingRate"
	gateway := "wss://api.hbdm.com/linear-swap-notification"
	proxyUrl := ""
	if htx.UseProxy {
		proxyUrl = fmt.Sprintf("http://%s", htx.ProxyUrl)
	}
	ws := websocketclient.New(gateway, proxyUrl)
	ws.OnConnectError(func(err error) {
		go errHandle(fmt.Errorf("OnConnectError: %v", err))
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(fmt.Errorf("disconnected: %v", err))
	})
	ws.OnConnected(func() {
		go logHandle(fmt.Sprintf("## connected %v", title))
		ws.SendTextMessage(fmt.Sprintf(`{"op":"sub","topic": "public.*.funding_rate", "cid": "%v"}`, util.GetUUID32()))
	})
	ws.OnBinaryMessageReceived(func(message []byte) {
		r, _ := gzip.NewReader(bytes.NewReader(message))
		buff, _ := io.ReadAll(r)
		type Msg struct {
			Op      string `json:"op"`
			Ch      string `json:"ch"`
			Type    string `json:"type"`
			Topic   string `json:"topic"`
			Event   string `json:"event"`
			ErrCode int    `json:"err-code"`
		}
		msg := Msg{}
		err := json.Unmarshal(buff, &msg)
		if err != nil {
			errHandle(fmt.Errorf("decode: %v", err))
			return
		}
		if msg.Op == "ping" {
			type pingRes struct {
				Op string `json:"op"`
				Ts int64  `json:"ts"`
			}
			pingRet := &pingRes{}
			json.Unmarshal(message, pingRet)
			pong := fmt.Sprintf(`{"op":"pong","ts":%d}`, pingRet.Ts)
			// 收到ping 回复pong
			ws.SendTextMessage(pong)
		} else if msg.Op == "auth" {
			if msg.Type == "api" && msg.ErrCode == 0 {
				// 订阅账户信息
				subAccountUpdateMp := map[string]any{
					"op":    "sub",
					"topic": "positions.*",
				}
				bf, _ := json.Marshal(subAccountUpdateMp)
				logHandle(fmt.Sprintf("sub: %v\n", string(bf)))
				// fmt.Printf("sub: %v\n", string(bf))
				ws.SendTextMessage(string(bf))
			}
		} else if msg.Op == "notify" && msg.Topic != "public.*.funding_rate" {
			// fmt.Printf("notify: %v\n", string(buff))
			type ReciveSwapFundingRate struct {
				Symbol      string      `json:"symbol"`
				FundingRate json.Number `json:"funding_rate"`
				FundingTime json.Number `json:"funding_time"` // 13位时间戳
			}
			type Msg struct {
				Data []ReciveSwapFundingRate `json:"data"`
			}
			res := Msg{}
			json.Unmarshal(buff, &res)
			// fmt.Printf("res.Data: %v\n", res.Data)
			if len(res.Data) > 0 {
				fmt.Printf("res.Data[0].FundingRate: %v\n", res.Data[0].FundingRate)
				fr, _ := res.Data[0].FundingRate.Float64()
				ft, _ := res.Data[0].FundingTime.Int64()
				reciveHandle(ReciveSwapFundingRateMsg{
					Symbol:      res.Data[0].Symbol,
					FundingRate: fr,
					FundingTime: ft / 1000,
					UpdateAt:    htx.GetTimeFloat(),
				})
			}
		}
	})

	ws.OnClose(func(code int, text string) {
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
