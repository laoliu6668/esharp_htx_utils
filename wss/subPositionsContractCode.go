package htx_wss

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

func SubPositionsContractCode(reciveHandle func(ReciveAccountsMsg), errHandle func(error)) {

	gateway := "api.hbdm.com"
	path := "/linear-swap-notification"
	mp := map[string]any{
		"AccessKeyId":      htx.ApiConfig.AccessKey,
		"Timestamp":        htx.UTCTimeNow(),
		"SignatureMethod":  "HmacSHA256",
		"SignatureVersion": "2",
	}
	mp["Signature"] = htx.Signature("get", gateway, path, mp, htx.ApiConfig.SecretKey)
	requrl := fmt.Sprintf("wss://%s%s", gateway, path)
	proxyUrl := ""
	if htx.UseProxy {
		proxyUrl = fmt.Sprintf("http://%s", htx.ProxyUrl)
	}
	fmt.Printf("requrl: %v\n", requrl)
	fmt.Printf("proxyUrl: %v\n", proxyUrl)
	ws := websocketclient.New(requrl, proxyUrl)
	ws.OnConnectError(func(err error) {
		fmt.Printf("err: %v\n", err)
		errHandle(err)
	})
	ws.OnDisconnected(func(err error) {
		errHandle(err)
	})
	ws.OnConnected(func() {
		fmt.Println("\n## connected SubAccountUpdate")
		// 发送鉴权消息
		mp["op"] = "auth"
		mp["type"] = "api"
		authBuf, _ := json.Marshal(mp)
		ws.SendTextMessage(string(authBuf))
		fmt.Printf("AuthInfo: %v\n\n", string(authBuf))
	})
	ws.OnBinaryMessageReceived(func(message []byte) {
		r, _ := gzip.NewReader(bytes.NewReader(message))
		buff, _ := io.ReadAll(r)
		// fmt.Printf("message1: %s\n", buff)
		type Msg struct {
			Op      string `json:"op"`
			Ch      string `json:"ch"`
			Type    string `json:"type"`
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
				fmt.Printf("sub: %v\n", string(bf))
				ws.SendTextMessage(string(bf))
			}
		} else if msg.Op == "notify" && msg.Type == "positions" {
			type TickerRes struct {
				Op    string            `json:"op"`
				Ch    string            `json:"ch"`
				Type  string            `json:"type"`
				Event string            `json:"event"`
				Data  ReciveAccountsMsg `json:"data"`
			}
			res := TickerRes{}
			json.Unmarshal([]byte(string(message)), &res)

			reciveHandle(res.Data)
		}
	})

	ws.OnClose(func(code int, text string) {
		fmt.Printf("close: %v, %v\n", code, text)
		errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
