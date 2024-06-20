package htx_wss

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// 订阅订单更新
// ttps://www.htx.com/zh-cn/opend/newApiPages/?id=7ec49f15-7773-11ed-9966-0242ac110003
func SubSpotOrder(reciveHandle func(ReciveSpotOrderMsg), logHandle func(string), errHandle func(error)) {

	title := "SubSpotOrder"
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
		go logHandle(fmt.Sprintf("## connected %v", title))
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

		go logHandle(fmt.Sprintf("## %v send auth", title))
		go logHandle(fmt.Sprintf("## auth %s", authBuf))

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
			// fmt.Printf("message: %v\n", message)
			type Data struct {
				Symbol      string      `json:"symbol"`
				OrderStatus string      `json:"orderStatus"`
				OrderId     int         `json:"orderId"`
				Type        string      `json:"type"`
				TradeVolume json.Number `json:"tradeVolume"`
				TradePrice  json.Number `json:"tradePrice"`
				OrderSize   json.Number `json:"orderSize"`
				OrderValue  json.Number `json:"orderValue"`
				ExecAmt     json.Number `json:"execAmt"`
				OrderSource string      `json:"orderSource"`
				TradeTime   int64       `json:"tradeTime"`
				// SeqNum      int64       `json:"seqNum"`
			}
			type TickerRes struct {
				Data Data `json:"data"`
			}
			res := TickerRes{}
			json.Unmarshal([]byte(message), &res)
			symbol := strings.Replace(strings.ToUpper(res.Data.Symbol), "USDT", "", 1)
			if res.Data.OrderSource != "spot-api" {
				// 只处理spot-api的订单
				return
			}

			if res.Data.Type == "buy-market" {
				// 市价买单
				if res.Data.OrderStatus == "submitted" {
					// #     order = OrderStruct()
					// #     order.exchange = __exchange__
					// #     order.symbol = str(data['symbol']).upper().replace('USDT', '')
					// #     order.order_id = str(data['orderId'])
					// #     order.order_type = 'buy-market'
					// #     order.order_value = float(data['orderValue'])
					// #     order.status = 1
					// #     order.create_at = data['orderCreateTime']
					// #     return order.__dict__
				} else if res.Data.OrderStatus == "filled" {
					ordervalue, _ := res.Data.OrderValue.Float64()
					tradePrice, _ := res.Data.TradePrice.Float64()
					tradeValue, _ := res.Data.ExecAmt.Float64()
					tradeVolume, _ := res.Data.TradeVolume.Float64()
					msg := ReciveSpotOrderMsg{
						Exchange:    "htx",
						Symbol:      symbol,
						OrderId:     strconv.Itoa(res.Data.OrderId),
						OrderType:   res.Data.Type,
						OrderValue:  ordervalue,
						TradePrice:  tradePrice,
						TradeValue:  tradeValue,
						TradeVolume: tradeVolume,
						Status:      2,
						FilledAt:    res.Data.TradeTime,
					}
					go reciveHandle(msg)
				}
			} else if res.Data.Type == "sell-market" {
				// 市价卖单
				if res.Data.OrderStatus == "submitted" {
					// #     order = OrderStruct()
					// #     order.exchange = __exchange__
					// #     order.symbol = str(data['symbol']).upper().replace('USDT', '')
					// #     order.order_id = str(data['orderId'])
					// #     order.order_type = 'sell-market'
					// #     order.order_volume = float(data['orderSize'])
					// #     order.status = 1
					// #     order.create_at = data['orderCreateTime']
					// #     return order.__dict__
				} else if res.Data.OrderStatus == "filled" {
					orderSize, _ := res.Data.OrderSize.Float64()
					tradePrice, _ := res.Data.TradePrice.Float64()
					tradeValue, _ := res.Data.ExecAmt.Float64()
					tradeVolume, _ := res.Data.TradeVolume.Float64()
					msg := ReciveSpotOrderMsg{
						Exchange:    "htx",
						Symbol:      symbol,
						OrderId:     strconv.Itoa(res.Data.OrderId),
						OrderType:   res.Data.Type,
						OrderVolume: orderSize,
						TradePrice:  tradePrice,
						TradeVolume: tradeVolume,
						TradeValue:  tradePrice * tradeValue,
						Status:      2,
						FilledAt:    res.Data.TradeTime,
					}
					go reciveHandle(msg)
				} else {
					return
				}
			}
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
			go logHandle(fmt.Sprintf("## %v sub success", title))
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
