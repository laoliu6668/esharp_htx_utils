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

// 【逐仓】订阅订单成交数据（sub）
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb706b4-77b5-11ed-9966-0242ac110003
func SubSwapOrder(reciveHandle func(ReciveSwapOrderMsg), logHandle func(string), errHandle func(error)) {

	flag := "SubSwapOrder"
	gateway := "api.hbdm.com"
	path := "/linear-swap-notification"

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
		logHandle(fmt.Sprintf("AuthInfo: %v\n\n", string(authBuf)))
	})
	ws.OnBinaryMessageReceived(func(message []byte) {
		r, _ := gzip.NewReader(bytes.NewReader(message))
		buff, _ := io.ReadAll(r)
		// fmt.Printf("buff: %s\n", buff)
		type Msg struct {
			Op      string `json:"op"`
			Ch      string `json:"ch"`
			Type    string `json:"type"`
			Topic   string `json:"topic"`
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
					"cid":   util.GetUUID32(),
					"topic": "orders.*",
				}
				bf, _ := json.Marshal(subAccountUpdateMp)
				logHandle(fmt.Sprintf("subed: %v\n", string(bf)))
				ws.SendTextMessage(string(bf))
			}
		} else if msg.Op == "notify" {
			type TickerRes struct {
				Status         int64       `json:"status"`
				OrderPriceType string      `json:"order_price_type"`
				Symbol         string      `json:"symbol"`
				OrderIdStr     string      `json:"order_id_str"`
				Direction      string      `json:"direction"`
				Offset         string      `json:"offset"`
				Volume         json.Number `json:"volume"`
				TradeVolume    json.Number `json:"trade_volume"`
				TradeAvgPrice  json.Number `json:"trade_avg_price"`
				TradeTurnover  json.Number `json:"trade_turnover"`
				CreatedAt      int64       `json:"created_at"`
			}
			res := TickerRes{}
			err := json.Unmarshal(buff, &res)
			if err != nil {
				errHandle(fmt.Errorf("decode: %v", err))
				return
			}
			// if res.OrderPriceType == "optimal_20" && res.Status == 6 {
			if res.Status == 6 {
				voluem, _ := res.Volume.Float64()
				trade_volume, _ := res.TradeVolume.Float64()
				tradeAvgPrice, _ := res.TradeAvgPrice.Float64()
				tradeTurnover, _ := res.TradeTurnover.Float64()

				ret := ReciveSwapOrderMsg{
					Exchange:    "htx",
					Symbol:      strings.ToUpper(res.Symbol),
					OrderId:     res.OrderIdStr,
					OrderType:   fmt.Sprintf("%s-%s", res.Direction, res.Offset),
					OrderVolume: int64(voluem),
					TradeVolume: int64(trade_volume),
					TradePrice:  tradeAvgPrice,
					TradeValue:  tradeTurnover,
					Status:      2,
					FilledAt:    res.CreatedAt,
				}
				go reciveHandle(ret)
			}

		}
	})
	ws.OnClose(func(code int, text string) {
		// fmt.Printf("close: %v, %v\n", code, text)
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
