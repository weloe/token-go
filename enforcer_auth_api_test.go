package token_go

import (
	"github.com/weloe/token-go/model"
	"testing"
)

type MockRbacAuth struct {
}

func (m *MockRbacAuth) GetRole(id string) []string {
	var arr = make([]string, 2)
	arr[1] = "user"
	return arr
}

type MockAclAuth struct {
}

func (m *MockAclAuth) GetPermission(id string) []string {
	var arr = make([]string, 2)
	arr[1] = "user::get"
	return arr
}

func TestEnforcer_GetRole(t *testing.T) {
	err, enforcer, ctx := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("NewTestEnforcer() failed: %v", err)
	}
	m := &MockRbacAuth{}
	enforcer.SetAuth(m)
	loginModel := model.DefaultLoginModel()
	loginModel.Token = "233"
	_, err = enforcer.LoginByModel("id", loginModel, ctx)
	if err != nil {
		t.Errorf("Login() failed: %v", err)
	}

	err = enforcer.CheckRole(ctx, "user")
	if err != nil {
		t.Errorf("CheckRole() failed: %v", err)
	}

	err = enforcer.CheckPermission(ctx, "user::get")
	if err == nil {
		t.Errorf("CheckRole() failed")
	}
	t.Logf("CheckPermission() return %v", err)
}

func TestEnforcer_CheckPermission(t *testing.T) {
	err, enforcer, ctx := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("NewTestEnforcer() failed: %v", err)
	}
	m := &MockAclAuth{}
	enforcer.SetAuth(m)
	loginModel := model.DefaultLoginModel()
	loginModel.Token = "233"
	_, err = enforcer.LoginByModel("id", loginModel, ctx)
	if err != nil {
		t.Errorf("Login() failed: %v", err)
	}

	err = enforcer.CheckRole(ctx, "user")
	if err == nil {
		t.Errorf("CheckRole() failed")
	}
	t.Logf("CheckRole() return %v", err)

	err = enforcer.CheckPermission(ctx, "user::get")
	if err != nil {
		t.Errorf("CheckRole() failed: %v", err)
	}
}