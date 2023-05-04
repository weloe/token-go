package persist

import (
	"testing"
	"time"
)

func NewTestDefaultAdapter() Adapter {
	return NewDefaultAdapter()
}

func TestDefaultAdapter_StrOperation(t *testing.T) {
	defaultAdapter := NewTestDefaultAdapter()

	if err := defaultAdapter.SetStr("k1", "v1", 0); err == nil {
		t.Errorf("SetStr() failed: set timeout = 0")
	}

	if err := defaultAdapter.SetStr("k1", "v1", -2); err == nil {
		t.Errorf("SetStr() failed: set timeout = -2")
	}

	if err := defaultAdapter.SetStr("k2", "v2", -1); err != nil {
		t.Errorf("SetStr() failed: can't set data")
	}
	if v := defaultAdapter.GetStr("k2"); v != "v2" {
		t.Errorf("GetStr() failed: value is %s, want 'v2' ", v)
	}
	if v := defaultAdapter.GetStrTimeout("k2"); v != -1 {
		t.Errorf("GetStrTimeout() failed: timeout is %v,want -1 ", v)
	}
	if v := defaultAdapter.GetStrTimeout("k3"); v != -2 {
		t.Errorf("GetStrTimeout() failed: timeout is %v,want -2 ", v)
	}

	if err := defaultAdapter.SetStr("k1", "v1", 1); err != nil {
		t.Errorf("SetStr() failed: can't set data")
	}
	time.Sleep(1 * time.Second)
	if v := defaultAdapter.Get("k1"); v != nil {
		t.Errorf("getExpireAndDelete() faliled: get expired value")
	}

	err1 := defaultAdapter.SetStr("k", "v", 9)
	if err1 != nil {
		t.Errorf("SetStr() failed: %v", err1)
	}
	time.Sleep(1 * time.Millisecond)
	timeout := defaultAdapter.GetStrTimeout("k")
	t.Logf("get timeout = %v", timeout)
	if timeout > 8 {
		t.Errorf("GetStrTimeout(） failed: %v", timeout)
	}

	if err := defaultAdapter.UpdateStrTimeout("k", -1); err != nil {
		t.Errorf("UpdateStrTimeout() failed: %v", err)
	}

	err2 := defaultAdapter.UpdateStrTimeout("k", 9)
	if err2 != nil {
		t.Errorf("UpdateStrTimeout() failed: %v", err2)
	}

	timeout = defaultAdapter.GetStrTimeout("k")
	t.Logf("get timeout = %v", timeout)

	getRes := defaultAdapter.GetStr("k")
	if getRes != "v" {
		t.Errorf("GetStr() failed: %v", getRes)
	}

	err3 := defaultAdapter.UpdateStr("k", "L")
	if err3 != nil {
		t.Errorf("UpdateStr() failed: %v", err3)
	}

	getRes = defaultAdapter.GetStr("k")
	if getRes != "L" {
		t.Errorf("GetStr() failed: GetStr() =  %v want 'L' ", getRes)
	}

	err4 := defaultAdapter.DeleteStr("k")
	if err4 != nil {
		t.Errorf("DeleteStr() failed: %v", err4)
	}
	err5 := defaultAdapter.UpdateStr("k", "L")
	if err5 == nil {
		t.Errorf("UpdateStr() failed: update not exist data")
	}

	getRes = defaultAdapter.GetStr("k")
	if getRes != "" {
		t.Errorf("GetStr() failed: %v", getRes)
	}
}

func TestDefaultAdapter_InterfaceOperation(t *testing.T) {
	defaultAdapter := NewTestDefaultAdapter()
	if err := defaultAdapter.Set("k1", "v1", 0); err == nil {
		t.Errorf("Set() failed: set timeout = 0")
	}

	if err := defaultAdapter.Set("k1", "v1", -2); err == nil {
		t.Errorf("Set() failed: set timeout = -2")
	}

	if err := defaultAdapter.Set("k2", "v2", -1); err != nil {
		t.Errorf("Set() failed: can't set data")
	}
	if v := defaultAdapter.Get("k2"); v.(string) != "v2" {
		t.Errorf("Get() failed: value is %s, want 'v2' ", v.(string))
	}
	if v := defaultAdapter.GetTimeout("k2"); v != -1 {
		t.Errorf("GetTimeout() failed: timeout is %v,want -1 ", v)
	}
	if v := defaultAdapter.GetTimeout("k3"); v != -2 {
		t.Errorf("GetTimeout() failed: timeout is %v,want -2 ", v)
	}

	if err := defaultAdapter.Set("k1", "v1", 1); err != nil {
		t.Errorf("Set() failed: can't set data")
	}
	time.Sleep(1 * time.Second)
	if v := defaultAdapter.Get("k1"); v != nil {
		t.Errorf("getExpireAndDelete() faliled: get expired value")
	}

	err1 := defaultAdapter.Set("k", "v", 9)
	if err1 != nil {
		t.Errorf("Set() failed: %v", err1)
	}
	time.Sleep(1 * time.Millisecond)

	timeout := defaultAdapter.GetTimeout("k")
	t.Logf("get timeout = %v", timeout)
	if timeout > 8 {
		t.Errorf("GetTimeout() failed: %v", timeout)
	}

	if err := defaultAdapter.UpdateTimeout("k", -1); err != nil {
		t.Errorf("UpdateTimeout() failed: %v", err)
	}
	err2 := defaultAdapter.UpdateTimeout("k", 9)
	if err2 != nil {
		t.Errorf("UpdateTimeout() failed: %v", err2)
	}

	timeout = defaultAdapter.GetTimeout("k")
	t.Logf("get timeout =  %v", timeout)

	getRes := defaultAdapter.Get("k")
	if getRes.(string) != "v" {
		t.Errorf("Get() failed: Get() = %v, want 'v' ", getRes.(string))
	}

	err3 := defaultAdapter.Update("k", "L")
	if err3 != nil {
		t.Errorf("Update() failed: %v", err3)
	}

	getRes = defaultAdapter.Get("k")
	if getRes.(string) != "L" {
		t.Errorf("Get() failed: Get() = %v want 'L' ", getRes.(string))
	}

	err4 := defaultAdapter.Delete("k")
	if err4 != nil {
		t.Errorf("Delete() failed: %v", err4)
	}
	err5 := defaultAdapter.Update("k", "L")
	if err5 == nil {
		t.Errorf("Update() failed: update not exist data")
	}

	getRes = defaultAdapter.Get("k")
	if getRes != nil {
		t.Errorf("Get() failed: %v s", getRes)
	}

}

func TestDefaultAdapter_DeleteBatchFilteredValue(t *testing.T) {
	adapter := NewTestDefaultAdapter()
	if err := adapter.SetStr("k_1", "v", -1); err != nil {
		t.Errorf("SetStr() failed: %v", err)
	}
	if err := adapter.SetStr("k_2", "v", -1); err != nil {
		t.Errorf("SetStr() failed: %v", err)
	}
	if err := adapter.SetStr("k_3", "v", -1); err != nil {
		t.Errorf("SetStr() failed: %v", err)
	}
	err := adapter.DeleteBatchFilteredKey("k_")
	if err != nil {
		t.Errorf("DeleteBatchFilteredKey() failed: %v", err)
	}
	str := adapter.GetStr("k_2")
	if str != "" {
		t.Errorf("DeleteBatchFilteredKey() failed")
	}
}
