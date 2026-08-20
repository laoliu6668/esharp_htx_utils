package htx_test

import (
	"fmt"
	"testing"

	htx_wss "github.com/laoliu6668/esharp_htx_utils/wss"
)

func TestWssSubSwapOrder(t *testing.T) {
	htx_wss.SubSwapOrderUnified(
		[]string{"*"},
		func(m htx_wss.ReciveSwapOrderMsg) {
			fmt.Printf("m: %+v\n", m)
		},
		func(log string) {
			fmt.Printf("log: %v\n", log)
		},
		func(err error) {
			fmt.Printf("err: %v\n", err)
		})
	select {}
}
