package htx_wss

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// 【逐仓】持仓变动更新数据（sub）
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb70c13-77b5-11ed-9966-0242ac110003
func SubSwapPositionInfo(reciveHandle func(ReciveSwapPositionMsg), logHandle func(string), errHandle func(error)) {
	title := "SubPositionsContractCode"
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
		logHandle(fmt.Sprintf("proxyUrl: %v\n", proxyUrl))
		proxyUrl = fmt.Sprintf("http://%s", htx.ProxyUrl)
	}
	logHandle(fmt.Sprintf("requrl: %v\n", requrl))
	// fmt.Printf("requrl: %v\n", requrl)
	ws := websocketclient.New(requrl, proxyUrl)
	ws.OnConnectError(func(err error) {
		fmt.Printf("err: %v\n", err)
		errHandle(err)
	})
	ws.OnDisconnected(func(err error) {
		errHandle(err)
	})
	ws.OnConnected(func() {
		logHandle(fmt.Sprintf("connected %s", title))
		// fmt.Println("\n## connected SubAccountUpdate")
		// 发送鉴权消息
		mp["op"] = "auth"
		mp["type"] = "api"
		authBuf, _ := json.Marshal(mp)
		ws.SendTextMessage(string(authBuf))
		logHandle(fmt.Sprintf("AuthInfo: %v\n", string(authBuf)))
		// fmt.Printf("AuthInfo: %v\n\n", string(authBuf))
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
		} else if msg.Op == "notify" && msg.Topic == "positions" {
			if msg.Event == "init" {
				// 初始推送（忽略）
				return
			} else if msg.Event == "snapshot" {
				type Data struct {
					Symbol    string      `json:"symbol"`
					Direction string      `json:"direction"` // buy or sell
					Volume    json.Number `json:"volume"`    // 持仓张数
					UpdateAt  float64     `json:"update_at"` // 更新时间
				}
				type Msg struct {
					Data []Data `json:"data"`
				}
				res := Msg{}
				err := json.Unmarshal([]byte(string(buff)), &res)
				if err != nil {
					errHandle(fmt.Errorf("decode: %v", err))
					return
				}
				ret := ReciveSwapPositionMsg{}
				for _, v := range res.Data {
					ret.Symbol = strings.ToUpper(v.Symbol)
					if v.Direction == "buy" {
						f, _ := v.Volume.Float64()
						ret.BuyVolume = int64(f)
					} else if v.Direction == "sell" {
						f, _ := v.Volume.Float64()
						ret.SellVolume = int64(f)
					}
				}
				ret.UpdateAt = htx.GetTimeFloat()
			}
		}
	})

	ws.OnClose(func(code int, text string) {
		fmt.Printf("close: %v, %v\n", code, text)
		errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
