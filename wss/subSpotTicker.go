package htx_wss

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	htx "github.com/laoliu6668/esharp_htx_utils"
	"github.com/laoliu6668/esharp_htx_utils/util"
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// reciveHandle:并发 logHandle:并发 errHandle:并发
func SubSpotTicker(symbols []string, reciveHandle func(Ticker), logHandle func(string), errHandle func(error)) {
	gateway := "wss://api.huobi.pro/wss"
	proxyUrl := ""
	if htx.UseProxy {
		proxyUrl = fmt.Sprintf("http://%s", htx.ProxyUrl)
	}
	ws := websocketclient.New(gateway, proxyUrl)
	ws.OnConnectError(func(err error) {
		fmt.Printf("err: %v\n", err)
		go errHandle(err)
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(err)
	})
	ws.OnSentError(func(err error) {
		go errHandle(fmt.Errorf("OnSentError: %v", err))
	})
	ws.OnConnected(func() {
		// go logHandle("## connected SubSpotTicker")
		for _, symbol := range symbols {
			ws.SendTextMessage(fmt.Sprintf(`{"sub": "market.%susdt.bbo", "id": "id%v"}`, strings.ToLower(symbol), time.Now().Unix()))
		}
		go logHandle(fmt.Sprintf("Sub: %v\n", strings.Join(symbols, "、")))
	})
	ws.OnBinaryMessageReceived(func(message []byte) {
		r, err := gzip.NewReader(bytes.NewReader(message))
		if err != nil {
			go errHandle(fmt.Errorf("gzip.NewReader: %v", err))
			return
		}
		buff, err := io.ReadAll(r)
		if err != nil {
			go errHandle(fmt.Errorf("io.ReadAll: %v", err))
			return
		}

		mp := map[string]any{}
		d := json.NewDecoder(strings.NewReader(string(buff)))
		d.UseNumber()
		err = d.Decode(&mp)
		if err != nil {
			go errHandle(fmt.Errorf("decode: %v", err))
			return
		}
		if _, ok := mp["ping"]; ok {
			// 收到ping 回复pong
			timestamp := util.ParseInt(mp["ping"], 0)
			ws.SendTextMessage(fmt.Sprintf(`{"pong":%d}`, timestamp))
		} else if _, ok := mp["ch"]; ok {
			type Tick struct {
				SeqID     int64   `json:"seqId"`
				Ask       float64 `json:"ask"`
				AskSize   float64 `json:"askSize"`
				Bid       float64 `json:"bid"`
				BidSize   float64 `json:"bidSize"`
				QuoteTime int64   `json:"quoteTime"`
				Symbol    string  `json:"symbol"`
			}

			type TickerRes struct {
				Ch   string `json:"ch"`
				Ts   int64  `json:"ts"`
				Tick Tick   `json:"tick"`
			}

			res := TickerRes{}
			json.Unmarshal([]byte(string(buff)), &res)
			symbolArr := strings.Split(res.Ch, ".")
			if len(symbolArr) > 1 {
				res.Ch = strings.Replace(strings.ToUpper(symbolArr[1]), "USDT", "", 1)
			}
			// gzip
			go reciveHandle(Ticker{
				Exchange: htx.ExchangeName,
				Symbol:   res.Ch,
				Buy:      Values{Price: res.Tick.Bid, Size: res.Tick.BidSize},
				Sell:     Values{Price: res.Tick.Ask, Size: res.Tick.AskSize},
				UpdateAt: htx.GetTimeFloat(),
			})
		} else if _, ok := mp["subbed"]; ok {
			// go logHandle(fmt.Sprintf("subbed: %v", string(buff)))
		} else {
			go logHandle(fmt.Sprintf("unknown message: %v", string(buff)))
		}
	})

	ws.OnClose(func(code int, text string) {
		// fmt.Printf("close: %v, %v\n", code, text)
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
