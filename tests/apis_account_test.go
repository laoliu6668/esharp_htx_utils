package htx_test

import (
	"testing"

	htx "github.com/laoliu6668/esharp_htx_utils/apis"
)

func Test_GetUserId(t *testing.T) {
	uid, err := htx.GetUserId()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("data: %v", uid)
}

func Test_GetUserAccount(t *testing.T) {
	data, err := htx.GetUserAccount()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("data: %v", data)
}

func Test_SpotToSwapTransfer(t *testing.T) {
	data, err := htx.SpotToSwapTransfer(20, "dot")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("data: %v", data)
}

func Test_SwapToSpotTransfer(t *testing.T) {
	data, err := htx.SwapToSpotTransfer(10, "dot")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("data: %v", data)
}
