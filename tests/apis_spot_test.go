package htx_test

import (
	"encoding/json"

	htx "github.com/laoliu6668/esharp_htx_utils/apis"
	"github.com/laoliu6668/esharp_htx_utils/util"

	"testing"
)

func Test_GetSpotSymbols(t *testing.T) {
	data, err := htx.GetSpotSymbols()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	b, _ := json.Marshal(data)
	util.WriteTestJsonFile("Test_GetSpotSymbols", b)
	t.Logf("data len: %v", len(data))
}

func Test_GetSpotAccountBalance(t *testing.T) {
	// data, err := htx.GetSpotAccountBalance(12222222)
	data, err := htx.GetSpotAccountBalance(61651347)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	b, _ := json.Marshal(data)
	util.WriteTestJsonFile("Test_GetSpotAccountBalance", b)

	t.Logf("data len: %v", len(data))
}

func Test_SpotBuyMarket(t *testing.T) {
	// 	no, err := htx.SpotBuyMarket("DOT", 9.0)
	no, err := htx.SpotBuyMarket("DOT", 10.0)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order no: %v", no)
}

func Test_SpotSellMarket(t *testing.T) {
	// no, err := htx.SpotSellMarket("DOT", 1.000064)
	no, err := htx.SpotSellMarket("DOT", 0.1)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order no: %v", no)
}
