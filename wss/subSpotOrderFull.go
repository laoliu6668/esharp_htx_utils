package htx_wss

import (
	"encoding/json"
	"fmt"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// 订阅订单更新
// https://www.htx.com/zh-cn/opend/newApiPages/?id=7ec49f15-7773-11ed-9966-0242ac110003
// reciveHandle:并发 logHandle:并发 errHandle:并发
func SubSpotOrderFull(reciveHandle func(map[string]any, []byte), logHandle func(string), errHandle func(error)) {
	gateway := "api.huobi.pro"
	path := "/ws/v2"

	requrl := fmt.Sprintf("wss://%s%s", gateway, path)
	proxyUrl := ""
	if htx.UseProxy {
		go logHandle(fmt.Sprintf("proxyUrl: %v\n", htx.ProxyUrl))
		proxyUrl = fmt.Sprintf("http://%s", htx.ProxyUrl)
	}
	go logHandle(fmt.Sprintf("requrl: %v\n", requrl))
	ws := websocketclient.New(requrl, proxyUrl)
	ws.OnConnectError(func(err error) {
		go errHandle(err)
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(err)
	})
	ws.OnConnected(func() {
		// 发送鉴权消息
		mp := map[string]any{
			"accessKey":        htx.ApiConfig.AccessKey,
			"timestamp":        htx.UTCTimeNow(),
			"signatureMethod":  "HmacSHA256",
			"signatureVersion": "2.1",
		}
		mp["signature"] = htx.Signature("get", gateway, path, mp, htx.ApiConfig.SecretKey)
		mp["authType"] = "api"
		authMap := map[string]any{
			"action": "req",
			"ch":     "auth",
			"params": mp,
		}
		authBuf, _ := json.Marshal(authMap)
		ws.SendTextMessage(string(authBuf))
		go logHandle(fmt.Sprintf("auth %s", authBuf))

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
		} else if msg.Action == "push" && strings.Contains(msg.Ch, "orders#*") {
			go reciveHandle(msg.Data, []byte(message))
		} else if msg.Action == "req" {
			if msg.Ch == "auth" && msg.Code == 200 {
				// 订阅账户信息
				subAccountUpdateMp := map[string]any{
					"action": "sub",
					"ch":     "orders#*",
				}
				bf, _ := json.Marshal(subAccountUpdateMp)
				go logHandle(fmt.Sprintf("sub: %v\n", string(bf)))
				ws.SendTextMessage(string(bf))
			}
		} else if msg.Action == "sub" {
			go logHandle(fmt.Sprintf("%v sub success", "orders#*"))
		} else {
			go logHandle(fmt.Sprintf("unknown message: %v\n", string(message)))
		}
	})

	ws.OnClose(func(code int, text string) {
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
