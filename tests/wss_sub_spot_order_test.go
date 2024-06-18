package htx_test

import (
	"encoding/json"
	"fmt"
	"testing"

	htx_wss "github.com/laoliu6668/esharp_htx_utils/wss"
)

func TestWssSubSpotOrder(t *testing.T) {
	htx_wss.SubSpotOrder(
		func(m htx_wss.ReciveSpotOrderMsg) {
			buf, _ := json.Marshal(m)
			fmt.Printf("m: %s\n", buf)
		},
		func(log string) {
			fmt.Printf("log: %v\n", log)
		},
		func(err error) {
			fmt.Printf("err: %v\n", err)
		})
	select {}
}
