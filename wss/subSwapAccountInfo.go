package htx_wss

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// 【逐仓】资产变动数据（sub）
func SubSwapAccountInfo(reciveHandle func(ReciveSwapAccountsMsg), logHandle func(string), errHandle func(error)) {

	flag := "SubSwapAccountInfo"
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
		logHandle(fmt.Sprintf("## connected %v\n", flag))
		// 发送鉴权消息
		mp["op"] = "auth"
		mp["type"] = "api"
		authBuf, _ := json.Marshal(mp)
		ws.SendTextMessage(string(authBuf))
		// fmt.Printf("AuthInfo: %v\n\n", string(authBuf))
		logHandle(fmt.Sprintf("AuthInfo: %v\n\n", string(authBuf)))
	})
	ws.OnBinaryMessageReceived(func(message []byte) {
		r, _ := gzip.NewReader(bytes.NewReader(message))
		buff, _ := io.ReadAll(r)

		type Msg struct {
			Op      string `json:"op"`
			Ch      string `json:"ch"`
			Type    string `json:"type"`
			Topic   string `json:"topic"`
			ErrCode int    `json:"err-code"`
			Event   string `json:"event"`
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
					"cid":   util.GetUUID32(),
					"topic": "accounts.*",
				}
				bf, _ := json.Marshal(subAccountUpdateMp)
				logHandle(fmt.Sprintf("subed: %v\n", string(bf)))
				ws.SendTextMessage(string(bf))
			}
		} else if msg.Op == "notify" && msg.Topic == "accounts" {
			if msg.Event == "init" {
				//初始推送（忽略）
				return
			} else if msg.Event == "snapshot" {

				type Msg struct {
					Symbol          string      `json:"symbol"`
					MarginAvailable json.Number `json:"margin_available"`  // 可用保金
					MarginFrozen    json.Number `json:"margin_frozen"`     // 冻结保金
					LiquidPrice     json.Number `json:"liquidation_price"` // 冻结保金
					RiskRate        json.Number `json:"risk_rate"`         // 风险率
				}
				type TickerRes struct {
					Op    string `json:"op"`
					Ch    string `json:"ch"`
					Type  string `json:"type"`
					Event string `json:"event"`
					Data  []Msg  `json:"data"`
				}
				res := TickerRes{}
				json.Unmarshal([]byte(string(buff)), &res)
				for _, v := range res.Data {
					free, _ := v.MarginAvailable.Float64()
					lock, _ := v.MarginFrozen.Float64()
					lp, _ := v.LiquidPrice.Float64()
					rr, _ := v.RiskRate.Float64()
					go reciveHandle(ReciveSwapAccountsMsg{
						Symbol:      strings.ToUpper(v.Symbol),
						FreeBalance: free,
						LockBalance: lock,
						LiquidPrice: lp,
						MarginRatio: math.Round(rr*100*100) / 100,
						UpdateAt:    htx.GetTimeFloat(),
					})
				}

			}

		}
	})
	ws.OnClose(func(code int, text string) {
		// fmt.Printf("close: %v, %v\n", code, text)
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
