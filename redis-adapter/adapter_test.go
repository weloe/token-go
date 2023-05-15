package redis_adapter

import (
	"github.com/weloe/token-go/token-go/persist"
	"testing"
	"time"
)

func NewTestRedisAdapter(t *testing.T) persist.Adapter {
	addr := "127.0.0.1:6379"
	pwd := ""
	db := 1
	adapter, err := NewRedisAdapter(addr, pwd, db)
	if err != nil {
		t.Fatalf("NewRedisAdapter() failed: %v", err)
	}
	return adapter
}

func TestDefaultAdapter_StrOperation(t *testing.T) {
	defaultAdapter := NewTestRedisAdapter(t)

	if v := defaultAdapter.GetStrTimeout("unExist"); v != 0 {
		t.Fatalf("GetStrTimeout() failed: timeout is %v,want 0 ", v)
	}

	if err := defaultAdapter.SetStr("k1", "v1", 0); err != nil {
		t.Fatalf("SetStr() failed: set timeout = 0")
	}

	if err := defaultAdapter.SetStr("k2", "v2", -1); err != nil {
		t.Fatalf("SetStr() failed: can't set data")
	}

	if v := defaultAdapter.GetStr("k2"); v != "v2" {
		t.Fatalf("GetStr() failed: value is %s, want 'v2' ", v)
	}

	if v := defaultAdapter.GetStrTimeout("k2"); v != 0 {
		t.Fatalf("GetStrTimeout() failed: timeout is %v,want 0 ", v)
	}

	if err := defaultAdapter.SetStr("k1", "v1", 1); err != nil {
		t.Fatalf("SetStr() failed: can't set data")
	}
	time.Sleep(2 * time.Second)
	if v := defaultAdapter.Get("k1"); v != nil {
		t.Fatalf("getExpireAndDelete() faliled: get expired value")
	}

	err1 := defaultAdapter.SetStr("k", "v", -1)
	if err1 != nil {
		t.Fatalf("SetStr() failed: %v", err1)
	}

	if err := defaultAdapter.UpdateStrTimeout("k", 9); err != nil {
		t.Fatalf("UpdateStrTimeout() failed: %v", err)
	}

	timeout := defaultAdapter.GetStrTimeout("k")
	t.Logf("get timeout = %v", timeout)

	getRes := defaultAdapter.GetStr("k")
	if getRes != "v" {
		t.Fatalf("GetStr() failed: %v", getRes)
	}

	err3 := defaultAdapter.UpdateStr("k", "L")
	if err3 != nil {
		t.Fatalf("UpdateStr() failed: %v", err3)
	}

	getRes = defaultAdapter.GetStr("k")
	if getRes != "L" {
		t.Fatalf("GetStr() failed: GetStr() =  %v want 'L' ", getRes)
	}

	err4 := defaultAdapter.DeleteStr("k")
	if err4 != nil {
		t.Fatalf("DeleteStr() failed: %v", err4)
	}

	getRes = defaultAdapter.GetStr("k")
	if getRes != "" {
		t.Fatalf("GetStr() failed: %v", getRes)
	}
}

func TestDefaultAdapter_InterfaceOperation(t *testing.T) {
	defaultAdapter := NewTestRedisAdapter(t)

	if v := defaultAdapter.GetTimeout("unExist"); v != 0 {
		t.Fatalf("GetTimeout() failed: timeout is %v,want 0 ", v)
	}

	if err := defaultAdapter.Set("k1", "v1", 0); err != nil {
		t.Fatalf("Set() failed: set timeout = 0")
	}

	if err := defaultAdapter.Set("k2", "v2", -1); err != nil {
		t.Fatalf("Set() failed: can't set data")
	}

	if v := defaultAdapter.Get("k2"); v != nil {
		t.Fatalf("Get() failed: value is %s, want 'v2' ", v)
	}

	if v := defaultAdapter.GetTimeout("k2"); v != 0 {
		t.Fatalf("GetTimeout() failed: timeout is %v,want 0 ", v)
	}

	if err := defaultAdapter.Set("k1", "v1", 1); err != nil {
		t.Fatalf("Set() failed: can't set data")
	}
	time.Sleep(2 * time.Second)
	if v := defaultAdapter.Get("k1"); v != nil {
		t.Fatalf("Get() faliled: get expired value")
	}

	err1 := defaultAdapter.Set("k", "v", -1)
	if err1 != nil {
		t.Fatalf("Set() failed: %v", err1)
	}

	if err := defaultAdapter.UpdateTimeout("k", 100); err != nil {
		t.Fatalf("UpdateTimeout() failed: %v", err)
	}

	timeout := defaultAdapter.GetTimeout("k")
	t.Logf("get timeout = %v", timeout)

	getRes := defaultAdapter.Get("k")
	if getRes != nil {
		t.Fatalf("GetGetStr() failed: %v", getRes)
	}

	err3 := defaultAdapter.Update("k", "L")
	if err3 != nil {
		t.Fatalf("Update() failed: %v", err3)
	}

	getRes = defaultAdapter.Get("k")
	if getRes != nil {
		t.Fatalf("Get() failed: GetStr() =  %v want 'L' ", getRes)
	}

	err4 := defaultAdapter.Delete("k")
	if err4 != nil {
		t.Fatalf("Delete() failed: %v", err4)
	}

	getRes = defaultAdapter.Get("k")
	if getRes != nil {
		t.Fatalf("Get() failed: %v", getRes)
	}
}

func TestDefaultAdapter_DeleteBatchFilteredValue(t *testing.T) {
	adapter := NewTestRedisAdapter(t)
	if err := adapter.SetStr("k_1", "v", -1); err != nil {
		t.Errorf("SetStr() failed: %v", err)
	}
	if err := adapter.SetStr("k_2", "v", -1); err != nil {
		t.Errorf("SetStr() failed: %v", err)
	}
	if err := adapter.SetStr("k_3", "v", -1); err != nil {
		t.Errorf("SetStr() failed: %v", err)
	}
	err := adapter.(persist.BatchAdapter).DeleteBatchFilteredKey("k_")
	if err != nil {
		t.Errorf("DeleteBatchFilteredKey() failed: %v", err)
	}
	str := adapter.GetStr("k_2")
	if str != "" {
		t.Errorf("DeleteBatchFilteredKey() failed")
	}
}
