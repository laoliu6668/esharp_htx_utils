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
	"github.com/laoliu6668/esharp_htx_utils/util/websocketclient"
)

// reciveHandle:并发 logHandle:并发 errHandle:并发
func SubSwapTicker(symbols []string, reciveHandle func(Ticker), logHandle func(string), errHandle func(error)) {
	gateway := "wss://api.hbdm.com/linear-swap-ws"
	proxyUrl := ""
	if htx.UseProxy {
		proxyUrl = fmt.Sprintf("http://%s", htx.ProxyUrl)
	}
	ws := websocketclient.New(gateway, proxyUrl)
	ws.OnConnectError(func(err error) {
		go errHandle(fmt.Errorf("OnConnectError: %v", err))
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(fmt.Errorf("disconnected: %v", err))
	})
	ws.OnConnected(func() {
		// go logHandle("## connected SubSwapTicker")
		for _, symbol := range symbols {
			ws.SendTextMessage(fmt.Sprintf(`{"sub": "market.%s-USDT.bbo", "id": "id%v"}`, strings.ToUpper(symbol), time.Now().Unix()))
		}
		go logHandle(fmt.Sprintf("Sub: %v\n", strings.Join(symbols, "、")))
	})
	ws.OnBinaryMessageReceived(func(message []byte) {
		r, _ := gzip.NewReader(bytes.NewReader(message))
		buff, _ := io.ReadAll(r)
		mp := map[string]any{}
		d := json.NewDecoder(strings.NewReader(string(buff)))
		d.UseNumber()
		err := d.Decode(&mp)
		if err != nil {
			go errHandle(fmt.Errorf("decode: %v", err))
			return
		}
		if _, ok := mp["ping"]; ok {
			// 收到ping 回复pong
			timestamp, _ := mp["ping"].(json.Number).Int64()
			ws.SendTextMessage(fmt.Sprintf(`{"pong":%d}`, timestamp))
		} else if _, ok := mp["ch"]; ok {
			type Tick struct {
				Mrid    int64     `json:"mrid"`
				ID      int       `json:"id"`
				Bid     []float64 `json:"bid"`
				Ask     []float64 `json:"ask"`
				Ts      int64     `json:"ts"`
				Version int64     `json:"version"`
				Ch      string    `json:"ch"`
			}
			type Res struct {
				Ch   string `json:"ch"`
				Ts   int64  `json:"ts"`
				Tick Tick   `json:"tick"`
			}
			res := Res{}
			json.Unmarshal([]byte(string(buff)), &res)
			symbolArr := strings.Split(res.Ch, ".")
			if len(symbolArr) > 1 {
				symbolCp := strings.ToUpper(symbolArr[1])
				symbolCpArr := strings.Split(symbolCp, "-")
				if len(symbolCpArr) > 0 {
					res.Ch = symbolCpArr[0]
				}
			}
			buyPrice := 0.0
			buySize := 0.0
			sellPrice := 0.0
			sellSize := 0.0
			if len(res.Tick.Bid) >= 2 {
				buyPrice = res.Tick.Bid[0]
				buySize = res.Tick.Bid[1]
			}
			if len(res.Tick.Ask) >= 2 {
				sellPrice = res.Tick.Ask[0]
				sellSize = res.Tick.Ask[1]
			}

			go reciveHandle(Ticker{
				Exchange: htx.ExchangeName,
				Symbol:   res.Ch,
				Buy:      Values{Price: buyPrice, Size: buySize},
				Sell:     Values{Price: sellPrice, Size: sellSize},
				UpdateAt: htx.GetTimeFloat(),
			})
		} else if _, ok := mp["subbed"]; ok {
			// go logHandle(fmt.Sprintf("subbed: %v", string(buff)))
		} else {
			go logHandle(fmt.Sprintf("unknown message: %v", string(buff)))
		}
	})

	ws.OnClose(func(code int, text string) {
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
