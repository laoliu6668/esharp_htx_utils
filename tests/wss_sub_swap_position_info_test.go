package htx_test

import (
	"fmt"
	"testing"

	htx_wss "github.com/laoliu6668/esharp_htx_utils/wss"
)

func TestWssSubSwapPositionInfo(t *testing.T) {
	sub()
	select {}
}

func sub() {
	htx_wss.SubSwapPositionInfo(
		func(m htx_wss.ReciveSwapPositionMsg) {
			fmt.Printf("m: %v\n", m)
		},
		func(log string) {
			fmt.Printf("log: %v\n", log)
		},
		func(err error) {
			fmt.Printf("err: %v\n", err)
		})
}
