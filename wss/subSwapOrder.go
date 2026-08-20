package htx_wss

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
	"github.com/shopspring/decimal"
)

// 【逐仓】订阅订单成交数据（sub）
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb706b4-77b5-11ed-9966-0242ac110003
// reciveHandle:并发 logHandle:并发 errHandle:并发
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
			go errHandle(fmt.Errorf("decode: %v", err))
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
				go logHandle(fmt.Sprintf("subed: %v\n", string(bf)))
				ws.SendTextMessage(string(bf))
			}
		} else if msg.Op == "notify" {
			type TickerRes struct {
				Status         int64       `json:"status"`
				OrderPriceType string      `json:"order_price_type"`
				Symbol         string      `json:"symbol"`
				OrderIdStr     string      `json:"order_id_str"`
				Direction      string      `json:"direction"`
				OrderSource    string      `json:"order_source"`
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
				go errHandle(fmt.Errorf("decode: %v", err))
				return
			}
			if res.OrderSource != "api" {
				return
			}
			// if res.OrderPriceType == "optimal_20" && res.Status == 6 {
			if res.Status == 6 {
				volume, _ := decimal.NewFromString(res.Volume.String())               // 下单张数d
				trade_volume, _ := decimal.NewFromString(res.TradeVolume.String())    // 成交张数d
				tradeAvgPrice, _ := decimal.NewFromString(res.TradeAvgPrice.String()) // 成交均价d
				tradeTurnover, _ := decimal.NewFromString(res.TradeTurnover.String()) // 成交额d
				size := tradeTurnover.Div(tradeAvgPrice).Div(trade_volume)            // 面值d
				orderVolume, _ := volume.Mul(size).Float64()                          // 下单数量
				tradeVolume, _ := trade_volume.Mul(size).Float64()                    // 成交数量
				tradePrice, _ := tradeAvgPrice.Float64()                              // 成交价格
				tradeValuie, _ := tradeTurnover.Float64()                             // 成交额
				// 面值    3105/3105/ 0.01
				ret := ReciveSwapOrderMsg{
					Exchange:    "htx",
					Symbol:      strings.ToUpper(res.Symbol),
					OrderId:     res.OrderIdStr,
					OrderType:   fmt.Sprintf("%s-%s", res.Direction, res.Offset),
					OrderVolume: orderVolume,
					TradeVolume: tradeVolume,
					TradePrice:  tradePrice,
					TradeValue:  tradeValuie,
					Status:      2,
					CreateAt:    res.CreatedAt,
					FilledAt:    time.Now().UnixMilli(),
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

// 【统一账户联合保证金】订阅订单成交数据（sub）
// https://www.htx.com/zh-cn/opend/newApiPages/?id=8cb89359-77b5-11ed-9966-19582def41f
// symbols 支持 BTC、BTC-USDT 或 *；仅回调 API 来源的完全成交订单。
func SubSwapOrderUnified(symbols []string, reciveHandle func(ReciveSwapOrderMsg), logHandle func(string), errHandle func(error)) {
	const (
		flag    = "SubSwapOrderUnified"
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
			errHandle(fmt.Errorf("authentication failed: %s", body))
		case response.Op == "auth":
			if response.Type != "api" {
				return
			}
			for _, symbol := range symbols {
				contractCode := strings.ToUpper(strings.TrimSpace(symbol))
				if contractCode == "" {
					continue
				}
				if contractCode != "*" && !strings.Contains(contractCode, "-") {
					contractCode += "-USDT"
				}
				subscription := map[string]any{
					"op":            "sub",
					"cid":           util.GetUUID32(),
					"topic":         "orders",
					"contract_code": contractCode,
				}
				request, err := json.Marshal(subscription)
				if err != nil {
					errHandle(fmt.Errorf("encode order subscription: %w", err))
					continue
				}
				if err := ws.SendTextMessage(string(request)); err != nil {
					errHandle(fmt.Errorf("send order subscription: %w", err))
				}
			}
			logHandle(fmt.Sprintf("sub: %s\n", strings.Join(symbols, ",")))
		case response.Op == "sub" && response.ErrCode != 0:
			errHandle(fmt.Errorf("order subscription failed for %s: %s", response.ContractCode, response.ErrMessage))
		case response.Op == "notify" && response.Topic == "orders":
			var notification struct {
				Data struct {
					ContractCode  string `json:"contract_code"`
					OrderId       string `json:"order_id"`
					OrderSource   string `json:"order_source"`
					Side          string `json:"side"`          // sell | buy
					PositionSide  string `json:"position_side"` // short | long
					State         string `json:"state"`         // filled
					Volume        string `json:"volume"`
					TradeVolume   string `json:"trade_volume"`
					TradeAvgPrice string `json:"trade_avg_price"`
					TradeTurnover string `json:"trade_turnover"`
					CreatedTime   string `json:"created_time"`
					UpdateTime    string `json:"updated_time"`
				} `json:"data"`
			}
			if err := json.Unmarshal(body, &notification); err != nil {
				errHandle(fmt.Errorf("decode order notification: %w", err))
				errHandle(fmt.Errorf("decode order notification:body %+v", string(body)))
				return
			}

			order := notification.Data
			if order.State != "filled" {
				return
			}
			volume, err := decimal.NewFromString(order.Volume)
			if err != nil {
				errHandle(fmt.Errorf("parse order volume for %s: %w", order.OrderId, err))
				return
			}
			tradeVolume, err := decimal.NewFromString(order.TradeVolume)
			if err != nil {
				errHandle(fmt.Errorf("parse trade volume for %s: %w", order.OrderId, err))
				return
			}
			if tradeVolume.IsZero() {
				errHandle(fmt.Errorf("trade volume is zero for %s", order.OrderId))
				return
			}
			tradePrice, err := decimal.NewFromString(order.TradeAvgPrice)
			if err != nil {
				errHandle(fmt.Errorf("parse trade average price for %s: %w", order.OrderId, err))
				return
			}
			if tradePrice.IsZero() {
				errHandle(fmt.Errorf("trade average price is zero for %s", order.OrderId))
				return
			}
			tradeTurnover, err := decimal.NewFromString(order.TradeTurnover)
			if err != nil {
				errHandle(fmt.Errorf("parse trade turnover for %s: %w", order.OrderId, err))
				return
			}

			contractSize := tradeTurnover.Div(tradePrice).Div(tradeVolume)
			orderVolume, _ := volume.Mul(contractSize).Float64()
			filledVolume, _ := tradeVolume.Mul(contractSize).Float64()
			filledPrice, _ := tradePrice.Float64()
			filledValue, _ := tradeTurnover.Float64()
			filledAtStr := order.UpdateTime
			if filledAtStr == "" {
				filledAtStr = strconv.FormatInt(time.Now().UnixMilli(), 10)
			}

			direction := strings.ToLower(order.Side)
			offset := "open"
			if order.PositionSide == "short" && direction == "buy" {
				offset = "close"
			} else if order.PositionSide == "long" && direction == "sell" {
				offset = "close"
			}

			filledAt, _ := strconv.ParseInt(filledAtStr, 10, 64)
			createdAt, _ := strconv.ParseInt(order.CreatedTime, 10, 64)
			reciveHandle(ReciveSwapOrderMsg{
				Exchange:    "htx",
				Symbol:      strings.ReplaceAll(strings.ToUpper(order.ContractCode), "-USDT", ""),
				OrderId:     order.OrderId,
				OrderType:   fmt.Sprintf("%s-%s", direction, offset),
				OrderVolume: orderVolume,
				TradeVolume: filledVolume,
				TradePrice:  filledPrice,
				TradeValue:  filledValue,
				Status:      2,
				CreateAt:    createdAt,
				FilledAt:    filledAt,
			})
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
