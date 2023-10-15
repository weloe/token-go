package token_go

import (
	"github.com/weloe/token-go/model"
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

func TestEnforcer_ConfirmQRCode(t *testing.T) {
	enforcer, _ := NewTestEnforcer(t)
	// in APP
	token, err := enforcer.LoginById("1")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	t.Logf("login token: %v", token)

	qrCodeId := "q1"

	err = enforcer.CreateQRCodeState(qrCodeId, -1)
	if err != nil {
		t.Fatalf("CreateQRCodeState() failed: %v", err)
	}
	t.Logf("After CreateQRCodeState(), current QRCode state: %v", enforcer.GetQRCodeState(qrCodeId))
	tempToken, err := enforcer.Scanned(qrCodeId, token)
	if err != nil {
		t.Fatalf("Scanned() failed: %v", err)
	}
	t.Logf("After Scanned(), current QRCode state: %v", enforcer.GetQRCodeState(qrCodeId))
	t.Logf("tempToken: %v", tempToken)
	err = enforcer.ConfirmAuth(tempToken)
	if err != nil {
		t.Fatalf("ConfirmAuth() failed: %v", err)
	}
	t.Logf("After ConfirmAuth(), current QRCode state: %v", enforcer.GetQRCodeState(qrCodeId))
	if enforcer.GetQRCodeState(qrCodeId) == model.ConfirmAuth {
		t.Logf(" QRCode login successfully.")
	}
}

func TestEnforcer_CancelAuthQRCode(t *testing.T) {
	enforcer, _ := NewTestEnforcer(t)
	// in APP
	token, err := enforcer.LoginById("1")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	t.Logf("login token: %v", token)

	qrCodeId := "q1"

	err = enforcer.CreateQRCodeState(qrCodeId, -1)
	if err != nil {
		t.Fatalf("CreateQRCodeState() failed: %v", err)
	}
	t.Logf("After CreateQRCodeState(), current QRCode state: %v", enforcer.GetQRCodeState(qrCodeId))
	tempToken, err := enforcer.Scanned(qrCodeId, token)
	if err != nil {
		t.Fatalf("Scanned() failed: %v", err)
	}
	t.Logf("After Scanned(), current QRCode state: %v", enforcer.GetQRCodeState(qrCodeId))
	t.Logf("tempToken: %v", tempToken)
	err = enforcer.CancelAuth(tempToken)
	if err != nil {
		t.Fatalf("CancelAuth() failed: %v", err)
	}
	t.Logf("After CancelAuth(), current QRCode state: %v", enforcer.GetQRCodeState(qrCodeId))
	if enforcer.GetQRCodeState(qrCodeId) == model.CancelAuth {
		t.Logf(" QRCode login is cancelled.")
	}
}
