package htx_test

import (
	"encoding/json"
	"fmt"
	"testing"

	htx_wss "github.com/laoliu6668/esharp_htx_utils/wss"
)

func Test_SubSpotTicker(t *testing.T) {
	symbols := []string{"DOT", "BTC", "ETH", "DOGE", "XRP"}
	htx_wss.SubSpotTicker(symbols, func(m htx_wss.ReciveData, _ []byte) {
		buf, _ := json.Marshal(m.Ticker)
		fmt.Println(string(buf))
	}, func(err error) {
		fmt.Println(err)
	})

}

func Test_SubSwapTicker(t *testing.T) {
	symbols := []string{"DOT", "BTC", "ETH", "DOGE", "XRP"}
	htx_wss.SubSwapTicker(symbols, func(m htx_wss.ReciveData, _ []byte) {
		buf, _ := json.Marshal(m.Ticker)
		fmt.Println(string(buf))
	}, func(err error) {
		fmt.Println(err)
	})
}

func Test_JsonNumber(t *testing.T) {
	type JsonNumber struct {
		Key json.Number `json:"key"`
	}
	test := JsonNumber{}
	symbols := "{\"key\":\"123.00\"}"
	// symbols := "{\"key\":123}"
	// symbols := "{\"key\":123.00}"

	json.Unmarshal([]byte(symbols), &test)
	fmt.Println(test.Key)

	bf, _ := json.Marshal(&test)
	fmt.Printf("string(bf): %v\n", string(bf))
}
