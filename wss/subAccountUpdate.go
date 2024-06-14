package htx_wss

import (
	"encoding/json"
	"fmt"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// 订阅现货账户变化
func SubAccountUpdate(reciveHandle func(ReciveBalanceMsg), logHandle func(string), errHandle func(error)) {

	gateway := "api.huobi.pro"
	path := "/ws/v2"
	mp := map[string]any{
		"accessKey":        htx.ApiConfig.AccessKey,
		"timestamp":        htx.UTCTimeNow(),
		"signatureMethod":  "HmacSHA256",
		"signatureVersion": "2.1",
	}
	mp["signature"] = htx.Signature("get", gateway, path, mp, htx.ApiConfig.SecretKey)
	requrl := fmt.Sprintf("wss://%s%s", gateway, path)
	proxyUrl := ""
	if htx.UseProxy {
		go logHandle(fmt.Sprintf("proxyUrl: %v\n", htx.ProxyUrl))
		proxyUrl = fmt.Sprintf("http://%s", htx.ProxyUrl)
	}
	go logHandle(fmt.Sprintf("requrl: %v\n", requrl))
	ws := websocketclient.New(requrl, proxyUrl)
	ws.OnConnectError(func(err error) {
		// fmt.Printf("err: %v\n", err)
		go errHandle(err)
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(err)
	})
	ws.OnConnected(func() {
		go logHandle("## connected SubAccountUpdate\n")
		// 发送鉴权消息
		mp["authType"] = "api"
		authMap := map[string]any{
			"action": "req",
			"ch":     "auth",
			"params": mp,
		}
		authBuf, _ := json.Marshal(authMap)
		// ws.SendBinaryMessage(authBuf)
		ws.SendTextMessage(string(authBuf))

		go logHandle("## SubAccountUpdate send auth\n")
		// fmt.Printf("authInfo: %v\n", string(authBuf))
	})

	ws.OnTextMessageReceived(func(message string) {
		type Msg struct {
			Action string         `json:"action"`
			Ch     string         `json:"ch"`
			Code   int            `json:"code"`
			Data   map[string]any `json:"data"`
		}
		msg := Msg{}
		err := json.Unmarshal([]byte(message), &msg)
		if err != nil {
			go errHandle(fmt.Errorf("decode: %v", err))
			return
		}
		if msg.Action == "ping" {
			type pingTs struct {
				Ts int64 `json:"ts"`
			}
			type pingRes struct {
				Action string `json:"action"`
				Data   pingTs `json:"data"`
			}
			pingRet := &pingRes{}
			json.Unmarshal([]byte(message), pingRet)
			pong := fmt.Sprintf(`{"action":"pong","data":{"ts":%d}}`, pingRet.Data.Ts)
			// 收到ping 回复pong
			ws.SendTextMessage(pong)
		} else if msg.Action == "push" && strings.Contains(msg.Ch, "accounts.update") {

			type Data struct {
				Currency    string      `json:"currency"`
				AccountId   int64       `json:"accountId"`
				Balance     json.Number `json:"balance"`
				Available   json.Number `json:"available"`
				AccountType string      `json:"accountType"`
				// SeqNum      int64       `json:"seqNum"`
			}
			fmt.Printf("msg: %v", message)
			type TickerRes struct {
				Data Data `json:"data"`
			}
			res := TickerRes{}
			json.Unmarshal([]byte(message), &res)
			if res.Data.AccountType == "trade" {
				a, _ := res.Data.Available.Float64()
				b, _ := res.Data.Balance.Float64()
				reciveHandle(ReciveBalanceMsg{
					Exchange:  htx.ExchangeName,
					Symbol:    strings.ToUpper(res.Data.Currency),
					AccountId: res.Data.AccountId,
					Available: a,
					Balance:   b,
				})
			}
		} else if msg.Action == "req" {
			if msg.Ch == "auth" && msg.Code == 200 {
				// 订阅账户信息
				subAccountUpdateMp := map[string]any{
					"action": "sub",
					"ch":     "accounts.update#0",
				}
				bf, _ := json.Marshal(subAccountUpdateMp)
				// fmt.Printf("sub: %v\n", string(bf))
				go logHandle(fmt.Sprintf("sub: %v\n", string(bf)))
				ws.SendTextMessage(string(bf))
			}
		} else if msg.Action == "sub" {
			go logHandle("## SubAccountUpdate sub success\n")
			// if msg.Code != 200 {
			// }
		} else {
			// fmt.Printf("unknown message: %v\n", string(message))
			go logHandle(fmt.Sprintf("unknown message: %v\n", string(message)))
		}
	})

	ws.OnClose(func(code int, text string) {
		// fmt.Printf("close: %v, %v\n", code, text)
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
