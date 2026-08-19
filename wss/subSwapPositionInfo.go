package htx_wss

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// 【逐仓】持仓变动更新数据（sub）
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb70c13-77b5-11ed-9966-0242ac110003
func SubSwapPositionInfo(symbols []string, reciveHandle func(ReciveSwapPositionMsg), logHandle func(string), errHandle func(error)) {
	title := "SubPositionsContractCode"
	gateway := "api.hbdm.com"
	path := "/linear-swap-notification"

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
		// 发送鉴权消息
		mp := map[string]any{
			"AccessKeyId":      htx.ApiConfig.AccessKey,
			"Timestamp":        htx.UTCTimeNow(),
			"SignatureMethod":  "HmacSHA256",
			"SignatureVersion": "2",
		}
		mp["Signature"] = htx.Signature("get", gateway, path, mp, htx.ApiConfig.SecretKey)
		mp["op"] = "auth"
		mp["type"] = "api"
		authBuf, _ := json.Marshal(mp)
		ws.SendTextMessage(string(authBuf))
		logHandle(fmt.Sprintf("AuthInfo: %v\n", string(authBuf)))
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

				for _, s := range symbols {
					subAccountUpdateMp := map[string]any{
						"op":    "sub",
						"topic": fmt.Sprintf("positions.%s-USDT", strings.ToUpper(s)),
					}
					bf, _ := json.Marshal(subAccountUpdateMp)
					ws.SendTextMessage(string(bf))
				}
				logHandle(fmt.Sprintf("sub: %s\n", strings.Join(symbols, ",")))
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

					UpdateAt float64 `json:"update_at"` // 更新时间
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
				reciveHandle(ret)
			}
		}
	})

	ws.OnClose(func(code int, text string) {
		fmt.Printf("close: %v, %v\n", code, text)
		errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}

// 【统一账户联合保证金】持仓变动更新数据（sub）
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-19582def41f
func SubSwapPositionInfoUnified(symbols []string, reciveHandle func(ReciveSwapPositionMsg), logHandle func(string), errHandle func(error)) {
	const (
		flag    = "SubSwapPositionInfoUnified"
		gateway = "api.hbdm.com"
		path    = "/ws/v5/notification"
	)

	reqURL := fmt.Sprintf("wss://%s%s", gateway, path)
	proxyURL := ""
	if htx.UseProxy {
		proxyURL = fmt.Sprintf("http://%s", htx.ProxyUrl)
		logHandle(fmt.Sprintf("proxyUrl: %v\n", htx.ProxyUrl))
	}
	logHandle(fmt.Sprintf("requrl: %v\n", reqURL))

	ws := websocketclient.New(reqURL, proxyURL)
	ws.OnConnectError(errHandle)
	ws.OnDisconnected(errHandle)
	ws.OnSentError(func(err error) {
		errHandle(fmt.Errorf("send websocket message: %w", err))
	})

	handleMessage := func(body []byte) {
		var response struct {
			Op           string      `json:"op"`
			Topic        string      `json:"topic"`
			Type         string      `json:"type"`
			Event        string      `json:"event"`
			ErrCode      int         `json:"err-code"`
			ErrMessage   string      `json:"err-msg"`
			ContractCode string      `json:"contract_code"`
			Ts           json.Number `json:"ts"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			errHandle(fmt.Errorf("decode websocket response: %w", err))
			return
		}

		switch {
		case response.Op == "ping":
			if err := ws.SendTextMessage(fmt.Sprintf(`{"op":"pong","ts":%s}`, response.Ts)); err != nil {
				errHandle(fmt.Errorf("send pong: %w", err))
			}
		case response.Op == "auth" && response.ErrCode != 0:
			errHandle(fmt.Errorf("authentication failed: %s", string(body)))
		case response.Op == "auth":
			if response.Type != "api" {
				return
			}
			for _, symbol := range symbols {
				contractCode := strings.ToUpper(strings.TrimSpace(symbol))
				if contractCode == "" {
					continue
				}
				if !strings.Contains(contractCode, "-") {
					contractCode += "-USDT"
				}
				subscription := map[string]any{
					"op":            "sub",
					"cid":           util.GetUUID32(),
					"topic":         "positions",
					"contract_code": contractCode,
				}
				request, err := json.Marshal(subscription)
				if err != nil {
					errHandle(fmt.Errorf("encode position subscription: %w", err))
					continue
				}
				if err := ws.SendTextMessage(string(request)); err != nil {
					errHandle(fmt.Errorf("send position subscription: %w", err))
				}
			}
			logHandle(fmt.Sprintf("sub: %s\n", strings.Join(symbols, ",")))
		case response.Op == "sub" && response.ErrCode != 0:
			errHandle(fmt.Errorf("position subscription failed for %s: %s", response.ContractCode, response.ErrMessage))
		case response.Op == "notify" && response.Topic == "positions":
			var notification struct {
				Data []struct {
					Symbol       string      `json:"symbol"`
					ContractCode string      `json:"contract_code"`
					PositionSide string      `json:"position_side"`
					Direction    string      `json:"direction"`
					Volume       json.Number `json:"volume"`
				} `json:"data"`
			}
			if err := json.Unmarshal(body, &notification); err != nil {
				errHandle(fmt.Errorf("decode position notification: %w", err))
				return
			}

			positions := make(map[string]ReciveSwapPositionMsg)
			for _, position := range notification.Data {
				symbol := strings.ToUpper(position.Symbol)
				if symbol == "" {
					symbol = strings.ToUpper(strings.Split(position.ContractCode, "-")[0])
				}
				volume, err := position.Volume.Int64()
				if err != nil {
					errHandle(fmt.Errorf("parse position volume for %s: %w", position.ContractCode, err))
					continue
				}
				update := positions[symbol]
				update.Symbol = symbol
				update.UpdateAt = htx.GetTimeFloat()
				side := strings.ToLower(position.PositionSide)
				if side == "long" || strings.ToLower(position.Direction) == "buy" {
					update.BuyVolume += volume
				} else if side == "short" || strings.ToLower(position.Direction) == "sell" {
					update.SellVolume += volume
				}
				positions[symbol] = update
			}
			for _, position := range positions {
				reciveHandle(position)
			}
		}
	}

	ws.OnTextMessageReceived(func(message string) {
		handleMessage([]byte(message))
	})
	ws.OnBinaryMessageReceived(func(message []byte) {
		reader, err := gzip.NewReader(bytes.NewReader(message))
		if err != nil {
			errHandle(fmt.Errorf("create gzip reader: %w", err))
			return
		}
		defer reader.Close()
		body, err := io.ReadAll(reader)
		if err != nil {
			errHandle(fmt.Errorf("read gzip message: %w", err))
			return
		}
		handleMessage(body)
	})
	ws.OnConnected(func() {
		logHandle(fmt.Sprintf("## connected %v\n", flag))
		authentication := map[string]any{
			"AccessKeyId":      htx.ApiConfig.AccessKey,
			"Timestamp":        htx.UTCTimeNow(),
			"SignatureMethod":  "HmacSHA256",
			"SignatureVersion": "2",
		}
		authentication["Signature"] = htx.Signature("GET", gateway, path, authentication, htx.ApiConfig.SecretKey)
		authentication["op"] = "auth"
		authentication["type"] = "api"
		request, err := json.Marshal(authentication)
		if err != nil {
			errHandle(fmt.Errorf("encode authentication request: %w", err))
			return
		}
		if err := ws.SendTextMessage(string(request)); err != nil {
			errHandle(fmt.Errorf("send authentication request: %w", err))
		}
	})
	ws.OnClose(func(code int, text string) {
		errHandle(fmt.Errorf("close: %v, %v", code, text))
	})
	ws.Connect()
}
