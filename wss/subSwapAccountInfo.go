package htx_wss

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// 【统一账户联合保证金】资产变动数据（sub）
func SubSwapAccountInfoUnified(reciveHandle func(ReciveSwapAccountsMsg), logHandle func(string), errHandle func(error)) {
	flag := "SubSwapAccountInfoUnified"
	gateway := "api.hbdm.com"
	path := "/ws/v5/notification"
	requrl := fmt.Sprintf("wss://%s%s", gateway, path)
	proxyURL := ""
	if htx.UseProxy {
		logHandle(fmt.Sprintf("proxyUrl: %v\n", htx.ProxyUrl))
		proxyURL = fmt.Sprintf("http://%s", htx.ProxyUrl)
	}
	logHandle(fmt.Sprintf("requrl: %v\n", requrl))

	ws := websocketclient.New(requrl, proxyURL)
	ws.OnConnectError(errHandle)
	ws.OnDisconnected(errHandle)
	ws.OnSentError(func(err error) {
		errHandle(fmt.Errorf("send websocket message: %w", err))
	})
	ws.OnTextMessageReceived(func(message string) {
		body := []byte(message)
		var response struct {
			Op      string      `json:"op"`
			Topic   string      `json:"topic"`
			Type    string      `json:"type"`
			ErrCode int         `json:"err-code"`
			Event   string      `json:"event"`
			Ts      json.Number `json:"ts"`
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
			errHandle(fmt.Errorf("authentication failed: %s", body))
		case response.Op == "auth":
			if response.Type == "api" && response.ErrCode == 0 {
				// 订阅账户信息
				subAccountUpdateMp := map[string]any{
					"op":    "sub",
					"cid":   util.GetUUID32(),
					"topic": "account",
				}
				bf, _ := json.Marshal(subAccountUpdateMp)
				logHandle(fmt.Sprintf("sub: %v\n", string(bf)))
				ws.SendTextMessage(string(bf))
			}
			logHandle(fmt.Sprintf("subscribed: %s\n", body))
		case response.Op == "notify" && response.Topic == "account":
			logHandle(fmt.Sprintf("account update: %s\n", body))
			type Msg struct {
				Currency              string `json:"currency"`                // 币种
				Equity                string `json:"equity"`                  // 该币种资产权益
				Available             string `json:"available"`               // 可用余额
				WithdrawAvailable     string `json:"withdraw_available"`      // 可提余额
				ProfitUnreal          string `json:"profit_unreal"`           // 未实现盈亏
				InitialMargin         string `json:"initial_margin"`          // 初始保证金
				MaintenanceMargin     string `json:"maintenance_margin"`      // 维持保证金
				MaintenanceMarginRate string `json:"maintenance_margin_rate"` // 维持保证金率
				InitialMarginRate     string `json:"initial_margin_rate"`     // 初始保证金率
			}
			type TickerRes struct {
				Op    string `json:"op"`
				Ch    string `json:"ch"`
				Type  string `json:"type"`
				Event string `json:"event"`
				Data  struct {
					State   string `json:"state"`
					Details []Msg  `json:"details"`
				} `json:"data"`
			}
			res := TickerRes{}
			json.Unmarshal(body, &res)
			for _, v := range res.Data.Details {
				eq, _ := strconv.ParseFloat(v.Equity, 64)
				im, _ := strconv.ParseFloat(v.InitialMargin, 64)
				free := eq - im
				lock, _ := strconv.ParseFloat(v.InitialMargin, 64)
				lp, _ := strconv.ParseFloat(v.MaintenanceMargin, 64)
				rr, _ := strconv.ParseFloat(v.MaintenanceMarginRate, 64)
				go reciveHandle(ReciveSwapAccountsMsg{
					Symbol:      strings.ToUpper(v.Currency),
					FreeBalance: free,
					LockBalance: lock,
					LiquidPrice: lp,
					MarginRatio: math.Round(rr*100*100) / 100,
					UpdateAt:    htx.GetTimeFloat(),
				})
			}
		}
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
		body, err := json.Marshal(authentication)
		if err != nil {
			errHandle(fmt.Errorf("encode authentication request: %w", err))
			return
		}
		if err := ws.SendTextMessage(string(body)); err != nil {
			errHandle(fmt.Errorf("send authentication request: %w", err))
			return
		}
		logHandle("authentication request queued\n")
	})
	ws.OnBinaryMessageReceived(func(message []byte) {
		r, err := gzip.NewReader(bytes.NewReader(message))
		if err != nil {
			errHandle(fmt.Errorf("create gzip reader: %w", err))
			return
		}
		defer r.Close()
		body, err := io.ReadAll(r)
		if err != nil {
			errHandle(fmt.Errorf("read gzip message: %w", err))
			return
		}
		logHandle(fmt.Sprintf("received binary: %s\n", body))

	})
	ws.OnClose(func(code int, text string) {
		errHandle(fmt.Errorf("close: %v, %v", code, text))
	})
	ws.Connect()
}
