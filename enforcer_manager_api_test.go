package token_go

import (
	"testing"
)

func TestEnforcer_TempToken(t *testing.T) {
	enforcer, _ := NewTestEnforcer(t)
	service := "code"
	tempToken, err := enforcer.CreateTempToken("uuid", service, "1234", -1)
	if err != nil {
		t.Fatalf("CreateTempToken() failed: %v", err)
	}
	timeout := enforcer.GetTempTokenTimeout(service, tempToken)
	if timeout != -1 {
		t.Errorf("GetTempTokenTimeout() failed, unexpected timeout: %v", timeout)
	}
	codeValue := enforcer.ParseTempToken("code", tempToken)
	if codeValue != "1234" {
		t.Errorf("ParseTempToken() failed, unexpected codeValue: %v", codeValue)
	}

	// delete
	if enforcer.DeleteTempToken(service, tempToken) != nil {
		t.Fatalf("DeleteTempToken() failed: %v", err)
	}
	tokenTimeout := enforcer.GetTempTokenTimeout(service, tempToken)
	if tokenTimeout != -2 {
		t.Errorf("GetTempTokenTimeout() failed, unexpected tokenTimeout: %v", tokenTimeout)
	}
	codeValue = enforcer.ParseTempToken(service, tempToken)
	if codeValue != "" {
		t.Errorf("ParseTempToken() failed, unexpected codeValue: %v", codeValue)
	}
}
