package token_go

import (
	"bytes"
	"fmt"
	"github.com/weloe/token-go/config"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/ctx"
	httpCtx "github.com/weloe/token-go/ctx/go-http-context"
	"github.com/weloe/token-go/log"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/persist"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewTestHttpContext(t *testing.T) (error, ctx.Context) {
	reqBody := bytes.NewBufferString("test request body")
	req, err := http.NewRequest("POST", "/test", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req.Header.Set(constant.TokenName, "233")

	ctx := NewHttpContext(req, rr)
	return err, ctx
}

func TestNewEnforcer(t *testing.T) {
	adapter := NewDefaultAdapter()
	ctx := httpCtx.NewHttpContext(nil, nil)

	tokenConfig := &config.TokenConfig{
		TokenName:     "testToken",
		Timeout:       60,
		IsReadCookie:  true,
		IsReadHeader:  true,
		IsReadBody:    false,
		IsConcurrent:  false,
		IsShare:       false,
		MaxLoginCount: -1,
	}
	logger := &log.DefaultLogger{}

	enforcer, err := NewEnforcer(tokenConfig, adapter, ctx)
	enforcer.SetType("u")
	if enforcer.GetType() != "u" {
		t.Error("enforcer.loginType should be user")
	}
	enforcer.SetAdapter(adapter)
	enforcer.SetLogger(logger)
	enforcer.SetWatcher(nil)
	enforcer.EnableLog()
	if !enforcer.IsLogEnable() {
		t.Errorf("enforcer.IsLogEnable() should be %v", enforcer.IsLogEnable())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if enforcer.config != *tokenConfig {
		t.Error("enforcer.config should be equal to the passed tokenConfig parameter")
	}
	if enforcer.GetAdapter() != adapter {
		t.Error("enforcer.adapter should be equal to the passed adapter parameter")
	}

	if enforcer.webCtx != ctx {
		t.Error("enforcer.webCtx should be equal to the passed ctx parameter")
	}

}

func NewTestEnforcer(t *testing.T) (error, *Enforcer) {
	reqBody := bytes.NewBufferString("test request body")
	req, err := http.NewRequest("POST", "/test", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req.Header.Set(constant.TokenName, "233")

	ctx := httpCtx.NewHttpContext(req, rr)

	adapter := persist.NewDefaultAdapter()

	tokenConfig := config.DefaultTokenConfig()

	enforcer, err := NewEnforcer(tokenConfig, adapter, ctx)
	return err, enforcer
}

func NewTestConcurrentEnforcer(t *testing.T) (error, *Enforcer) {
	reqBody := bytes.NewBufferString("test request body")
	req, err := http.NewRequest("POST", "/test", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req.Header.Set(constant.TokenName, "233")

	ctx := httpCtx.NewHttpContext(req, rr)

	adapter := persist.NewDefaultAdapter()

	tokenConfig := config.DefaultTokenConfig()
	tokenConfig.IsConcurrent = true
	tokenConfig.IsShare = false

	enforcer, err := NewEnforcer(tokenConfig, adapter, ctx)
	return err, enforcer
}

func NewTestNotConcurrentEnforcer(t *testing.T) (error, *Enforcer) {
	reqBody := bytes.NewBufferString("test request body")
	req, err := http.NewRequest("POST", "/test", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req.Header.Set(constant.TokenName, "233")

	ctx := httpCtx.NewHttpContext(req, rr)

	adapter := persist.NewDefaultAdapter()

	tokenConfig := config.DefaultTokenConfig()
	tokenConfig.IsConcurrent = false
	tokenConfig.IsShare = false

	enforcer, err := NewEnforcer(tokenConfig, adapter, ctx)
	return err, enforcer
}

func TestNewEnforcerByFile(t *testing.T) {
	err, ctx := NewTestHttpContext(t)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	adapter := persist.NewDefaultAdapter()
	conf := "testConf"

	enforcer, err := NewEnforcerByFile(conf, adapter, ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if enforcer.conf != conf {
		t.Error("enforcer.conf should be equal to the passed conf parameter")
	}

	if enforcer.adapter != adapter {
		t.Error("enforcer.adapter should be equal to the passed adapter parameter")
	}
	if enforcer.webCtx != ctx {
		t.Error("enforcer.webCtx should be equal to the passed ctx parameter")
	}
}

func TestEnforcer_Login(t *testing.T) {
	err, enforcer := NewTestEnforcer(t)
	enforcer.EnableLog()
	if err != nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}
	loginId := "1"
	_, err = enforcer.Login(loginId)
	if err != nil {
		t.Errorf("LoginByModel() failed: %v", err)
	}

	_, err = enforcer.LoginByModel(loginId, model.DefaultLoginModel())
	if err != nil {
		t.Errorf("LoginByModel() failed: %v", err)
	}
	login, err := enforcer.IsLoginById(loginId)
	if err != nil {
		t.Errorf("IsLoginById() failed: err should be nil now: %v", err)
	}
	if !login {
		t.Errorf("IsLoginById() failed: IsLoginById() = %v", login)
	}
	err = enforcer.Replaced("1", "")
	if err != nil {
		t.Errorf("Replaced() failed: %v", err)
	}
	session := enforcer.GetSession("1")
	t.Logf("id = %v  session.tokenSign len = %v", "1", session.TokenSignList.Len())

	login, err = enforcer.IsLoginById(loginId)
	if err != nil {
		t.Logf("%v error: %v", login, err)
	}
	if login {
		t.Errorf("IsLoginById() failed: IsLoginById() = %v", login)
	}

}

func TestEnforcer_GetLoginId(t *testing.T) {
	err, enforcer := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}
	loginModel := model.DefaultLoginModel()
	loginModel.Token = "233"
	_, err = enforcer.LoginByModel("id", loginModel)
	if err != nil {
		t.Errorf("Login() failed: %v", err)
	}

	id, err := enforcer.GetLoginId()
	if err != nil {
		t.Errorf("GetLoginId() failed: %v", err)
	}
	t.Logf("LoginId = %v", id)
	if id != "id" {
		t.Errorf("GetLoginId() failed: %v", err)
	}

}

func TestEnforcer_Logout(t *testing.T) {
	err, enforcer := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()
	loginModel.Token = "233"
	token, err := enforcer.LoginByModel("id", loginModel)
	if token != "233" {
		t.Errorf("LoginByModel() failed: unexpected token value %s, want '233' ", token)
	}
	if err != nil {
		t.Errorf("Login() failed: %v", err)
	}

	err = enforcer.Logout()
	if err != nil {
		t.Errorf("Logout() failed: %v", err)
	}

	login, err := enforcer.IsLogin()
	if login {
		t.Errorf("IsLogin() failed: unexpected value %v", login)
	}
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func TestEnforcer_Kickout(t *testing.T) {
	err, enforcer := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()
	loginModel.Token = "233"
	_, err = enforcer.LoginByModel("id", loginModel)
	if err != nil {
		t.Errorf("Login() failed: %v", err)
	}

	err = enforcer.Kickout("id", "")
	if err != nil {
		t.Errorf("Kickout() failed %v", err)
	}

	session := enforcer.GetSession("id")
	if session != nil {
		t.Errorf("unexpected session value %v", session)
	}
	login, err := enforcer.IsLogin()
	if login {
		t.Errorf("IsLogin() failed: unexpected value %v", login)
	}
	n := fmt.Sprintf("%v", err)
	if n != "this account is kicked out" {
		t.Errorf("IsLogin() failed: unexpected error value %v", err)
	}

}

func TestEnforcerNotConcurrentNotShareLogin(t *testing.T) {
	err, enforcer := NewTestNotConcurrentEnforcer(t)
	if err != nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()

	for i := 0; i < 4; i++ {
		_, err = enforcer.LoginByModel("id", loginModel)
		if err != nil {
			t.Errorf("Login() failed: %v", err)
		}
	}
	session := enforcer.GetSession("id")
	if session.TokenSignList.Len() != 1 {
		t.Errorf("Login() failed: unexpected session.TokenSignList length = %v", session.TokenSignList.Len())
	}

}

func TestEnforcer_ConcurrentShare(t *testing.T) {
	err, enforcer := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()
	for i := 0; i < 5; i++ {
		_, err = enforcer.LoginByModel("id", loginModel)
		if err != nil {
			t.Errorf("Login() failed: %v", err)
		}
	}
	session := enforcer.GetSession("id")
	if session.TokenSignList.Len() != 1 {
		t.Errorf("Login() failed: unexpected session.TokenSignList length = %v", session.TokenSignList.Len())
	}

}
func TestEnforcer_ConcurrentNotShareMultiLogin(t *testing.T) {
	err, enforcer := NewTestConcurrentEnforcer(t)
	if err != nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()
	for i := 0; i < 14; i++ {
		_, err = enforcer.LoginByModel("id", loginModel)
		if err != nil {
			t.Errorf("Login() failed: %v", err)
		}
	}
	session := enforcer.GetSession("id")
	if session.TokenSignList.Len() != 12 {
		t.Errorf("Login() failed: unexpected session.TokenSignList length = %v", session.TokenSignList.Len())
	}

}

func TestNewDefaultEnforcer(t *testing.T) {
	err, ctx := NewTestHttpContext(t)
	if err != nil {
		t.Errorf("NewTestHttpContext() failed: %v", err)
	}

	enforcer, err := NewDefaultEnforcer(persist.NewDefaultAdapter(), ctx)
	if err != nil || enforcer == nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}
}
