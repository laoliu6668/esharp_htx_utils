package htx_test

import (
	"encoding/json"

	htx "github.com/laoliu6668/esharp_htx_utils/apis"
	"github.com/laoliu6668/esharp_htx_utils/util"

	"testing"
)

func Test_GetSwapSymbol(t *testing.T) {
	data, err := htx.GetSwapSymbol()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	b, _ := json.Marshal(data)
	util.WriteTestJsonFile("Test_GetSwapSymbol", b)
	t.Logf("data len: %v", len(data))
}

func Test_GetSwapPositionLimit(t *testing.T) {
	data, err := htx.GetSwapPositionLimit("ETH")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	b, _ := json.Marshal(data)
	util.WriteTestJsonFile("Test_GetSwapPositionLimit", b)
	t.Logf("data len: %v", len(data))
}

func Test_GetSwapOrderLimit(t *testing.T) {
	data, err := htx.GetSwapOrderLimit()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	b, _ := json.Marshal(data)
	util.WriteTestJsonFile("Test_GetSwapOrderLimit", b)
	t.Logf("data len: %v", len(data))
}

func Test_GetSwapFundingRate(t *testing.T) {
	data, err := htx.GetSwapFundingRate()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	b, _ := json.Marshal(data)
	util.WriteTestJsonFile("Test_GetSwapFundingRate", b)
	t.Logf("data len: %v", len(data))
}

func Test_GetSwapAccountPositionInfo(t *testing.T) {
	data, err := htx.GetSwapAccountPositionInfo("")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	b, _ := json.Marshal(data)
	util.WriteTestJsonFile("Test_GetSwapAccountPositionInfo", b)
	t.Logf("data len: %v", len(data))
}

func Test_GetSwapAccountType(t *testing.T) {
	data, err := htx.GetSwapAccountType()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("acc type: %v", data)
}

func Test_GetSwapAccountBalance(t *testing.T) {
	data, err := htx.GetSwapAccountBalance()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("acc banlance: %v", data)
}
func Test_SwapSellOpen(t *testing.T) {
	// 卖出开空
	data, err := htx.SwapSellOpen("DOT", 1)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order res: %v", data)
}

func Test_SwapBuyClose(t *testing.T) {
	// 买入平空
	data, err := htx.SwapBuyClose("DOT", 1)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order res: %v", data)
}

func Test_SwapBuyOpen(t *testing.T) {
	// 买入开多
	data, err := htx.SwapBuyOpen("DOT", 1)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order res: %v", data)
}
func Test_SwapSellClose(t *testing.T) {
	// 卖出平多
	data, err := htx.SwapSellClose("DOT", 1)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order res: %v", data)
}
