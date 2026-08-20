package htx_test

import (
	"encoding/json"
	"time"

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
func Test_GetSwapAccountInfo(t *testing.T) {
	data, err := htx.GetSwapAccountInfo()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	b, _ := json.Marshal(data)
	util.WriteTestJsonFile("Test_GetSwapAccountInfo", b)
	t.Logf("details len: %v", len(data.Details))
}
func Test_GetSwapPositionInfo(t *testing.T) {
	data, err := htx.GetSwapPositionInfo("")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	b, _ := json.Marshal(data)
	util.WriteTestJsonFile("Test_GetSwapPositionInfo", b)
	t.Logf("data len: %v", len(data))
}
func Test_GetSwapAccountPositionInfo(t *testing.T) {
	data, err := htx.GetSwapAccountPositionInfo("TRX")
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
	data, err := htx.SwapSellOpen("TRX", 1, 2)
	if err != nil {
		t.Errorf("Error:1%v", err)
		return
	}
	t.Logf("order res: %v ts: %v", data, time.Now().UnixNano())
}

func Test_SwapBuyClose(t *testing.T) {
	// 买入平空
	data, err := htx.SwapBuyClose("TRX", 2, 3)
	if err != nil {
		t.Errorf("Error1: %v", err)
		return
	}
	t.Logf("order res: %v", data)
}

func Test_SwapBuyOpen(t *testing.T) {
	// 买入开多
	data, err := htx.SwapBuyOpen("DOT", 1, 2)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order res: %v", data)
}
func Test_SwapSellClose(t *testing.T) {
	// 卖出平多
	data, err := htx.SwapSellClose("DOT", 1, 2)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order res: %v", data)
}

func Test_SwapSellOpenV5(t *testing.T) {
	// 卖出开空
	data, err := htx.SwapOrderV5("TRX", "cross", 2, "sell", "short", "market")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order res: %v", data)
}

func Test_SwapSellCloseV5(t *testing.T) {
	// 卖出平多
	data, err := htx.SwapOrderV5("TRX", "cross", 1, "buy", "both", "market")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("order res: %+v", data)
}

func Test_SetSwapLeverage(t *testing.T) {
	// 设置杠杆倍数
	data, err := htx.SetSwapLeverageRate("TRX", "cross", "both", 5)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("set leverage res: %+v", data)
}

func Test_SetSwapPositionMode(t *testing.T) {

	data, err := htx.SetSwapPositionMode("dual_side")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("set position mode res: %+v", data)
}
