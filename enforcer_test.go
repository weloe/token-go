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
	if ctx == nil {
		t.Errorf("NewHttpContext failed: %v", ctx)
	}
	tokenConfig := &config.TokenConfig{
		TokenName:         "testToken",
		Timeout:           60,
		IsReadCookie:      true,
		IsReadHeader:      true,
		IsReadBody:        false,
		IsConcurrent:      false,
		IsShare:           false,
		MaxLoginCount:     -1,
		DataRefreshPeriod: -1,
	}
	logger := &log.DefaultLogger{}

	enforcer, err := NewEnforcer(adapter, tokenConfig)
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

}

func NewTestEnforcer(t *testing.T) (error, *Enforcer, ctx.Context) {
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

	enforcer, err := NewEnforcer(adapter, tokenConfig)
	return err, enforcer, ctx
}

func NewTestConcurrentEnforcer(t *testing.T) (error, *Enforcer, ctx.Context) {
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

	enforcer, err := NewEnforcer(adapter, tokenConfig)
	return err, enforcer, ctx
}

func NewTestNotConcurrentEnforcer(t *testing.T) (error, *Enforcer, ctx.Context) {
	reqBody := bytes.NewBufferString("test request body")
	req, err := http.NewRequest("POST", "/test", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req.Header.Set(constant.TokenName, "233")

	ctx := httpCtx.NewHttpContext(req, rr)
	t.Log(ctx)

	adapter := persist.NewDefaultAdapter()

	tokenConfig := config.DefaultTokenConfig()
	tokenConfig.IsConcurrent = false
	tokenConfig.IsShare = false
	tokenConfig.DataRefreshPeriod = 30

	enforcer, err := NewEnforcer(adapter, tokenConfig)
	return err, enforcer, ctx
}

func TestNewEnforcerByFile(t *testing.T) {
	err, _ := NewTestHttpContext(t)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	adapter := persist.NewDefaultAdapter()
	conf := "./examples/token_conf.yaml"

	enforcer, err := NewEnforcer(adapter, conf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if enforcer.conf != conf {
		t.Error("enforcer.conf should be equal to the passed conf parameter")
	}

	if enforcer.adapter != adapter {
		t.Error("enforcer.adapter should be equal to the passed adapter parameter")
	}

}

func TestEnforcer_Login(t *testing.T) {
	err, enforcer, ctx := NewTestEnforcer(t)
	enforcer.EnableLog()
	if err != nil {
		t.Errorf("InitWithConfig() failed: %v", err)
	}
	loginId := "1"
	_, err = enforcer.Login(loginId, ctx)
	if err != nil {
		t.Errorf("LoginByModel() failed: %v", err)
	}

	_, err = enforcer.LoginByModel(loginId, model.DefaultLoginModel(), ctx)
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
	t.Logf("id = %v  session.tokenSign len = %v", "1", session.TokenSignSize())

	err = enforcer.CheckLogin(ctx)
	if err == nil {
		t.Errorf("CheckLogin() failed: CheckLogin() return nil")
	}
	login, err = enforcer.IsLoginById(loginId)
	if err != nil {
		t.Logf("%v error: %v", login, err)
	}
	if login {
		t.Errorf("IsLoginById() failed: IsLoginById() = %v", login)
	}

}

func TestEnforcer_GetLoginId(t *testing.T) {
	err, enforcer, ctx := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("InitWithConfig() failed: %v", err)
	}
	loginModel := model.DefaultLoginModel()
	loginModel.Token = "233"
	_, err = enforcer.LoginByModel("id", loginModel, ctx)
	if err != nil {
		t.Errorf("Login() failed: %v", err)
	}

	id, err := enforcer.GetLoginId(ctx)
	if err != nil {
		t.Errorf("GetLoginId() failed: %v", err)
	}
	t.Logf("LoginId = %v", id)
	if id != "id" {
		t.Errorf("GetLoginId() failed: %v", err)
	}

}

func TestEnforcer_Logout(t *testing.T) {
	err, enforcer, ctx := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("InitWithConfig() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()
	loginModel.Token = "233"
	token, err := enforcer.LoginByModel("id", loginModel, ctx)
	if token != "233" {
		t.Errorf("LoginByModel() failed: unexpected token value %s, want '233' ", token)
	}
	if err != nil {
		t.Errorf("Login() failed: %v", err)
	}

	err = enforcer.Logout(ctx)
	if err != nil {
		t.Errorf("Logout() failed: %v", err)
	}

	login, err := enforcer.IsLogin(ctx)
	if login {
		t.Errorf("IsLogin() failed: unexpected value %v", login)
	}
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func TestEnforcer_Kickout(t *testing.T) {
	err, enforcer, ctx := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("InitWithConfig() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()
	loginModel.Token = "233"
	_, err = enforcer.LoginByModel("id", loginModel, ctx)
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
	login, err := enforcer.IsLogin(ctx)
	if login {
		t.Errorf("IsLogin() failed: unexpected value %v", login)
	}
	n := fmt.Sprintf("%v", err)
	if n != "this account is kicked out" {
		t.Errorf("IsLogin() failed: unexpected error value %v", err)
	}

}

func TestEnforcerNotConcurrentNotShareLogin(t *testing.T) {
	err, enforcer, ctx := NewTestNotConcurrentEnforcer(t)
	if err != nil {
		t.Errorf("InitWithConfig() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()

	for i := 0; i < 4; i++ {
		_, err = enforcer.LoginByModel("id", loginModel, ctx)
		if err != nil {
			t.Errorf("Login() failed: %v", err)
		}
	}
	session := enforcer.GetSession("id")
	if session.TokenSignSize() != 1 {
		t.Errorf("Login() failed: unexpected session.TokenSignList length = %v", session.TokenSignSize())
	}

}

func TestEnforcer_ConcurrentShare(t *testing.T) {
	err, enforcer, ctx := NewTestEnforcer(t)
	if err != nil {
		t.Errorf("InitWithConfig() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()
	for i := 0; i < 5; i++ {
		_, err = enforcer.LoginByModel("id", loginModel, ctx)
		if err != nil {
			t.Errorf("Login() failed: %v", err)
		}
	}
	session := enforcer.GetSession("id")
	t.Logf("Login(): session.TokenSignList length = %v", session.TokenSignSize())

	if session.TokenSignSize() != 1 {
		t.Errorf("Login() failed: unexpected session.TokenSignList length = %v", session.TokenSignSize())
	}

}
func TestEnforcer_ConcurrentNotShareMultiLogin(t *testing.T) {
	err, enforcer, ctx := NewTestConcurrentEnforcer(t)
	if err != nil {
		t.Errorf("InitWithConfig() failed: %v", err)
	}

	loginModel := model.DefaultLoginModel()
	for i := 0; i < 14; i++ {
		_, err = enforcer.LoginByModel("id", loginModel, ctx)
		if err != nil {
			t.Errorf("Login() failed: %v", err)
		}
	}
	session := enforcer.GetSession("id")
	if session.TokenSignSize() != 12 {
		t.Errorf("Login() failed: unexpected session.TokenSignList length = %v", session.TokenSignSize())
	}

}

func TestNewDefaultEnforcer(t *testing.T) {
	err, ctx := NewTestHttpContext(t)
	if ctx == nil {
		t.Errorf("NewTestHttpContext() failed: %v", err)
	}
	if err != nil {
		t.Errorf("NewTestHttpContext() failed: %v", err)
	}

	enforcer, err := NewEnforcer(persist.NewDefaultAdapter())

	if err != nil || enforcer == nil {
		t.Errorf("InitWithConfig() failed: %v", err)
	}
}

func TestNewEnforcer1(t *testing.T) {
	enforcer, err := NewEnforcer(NewDefaultAdapter())
	t.Log(err)
	t.Log(enforcer)
	enforcer, err = NewEnforcer(NewDefaultAdapter(), config.DefaultTokenConfig())
	t.Log(err)
	t.Log(enforcer)
}

func TestEnforcer_JsonAdapter(t *testing.T) {
	enforcer, err := NewEnforcer(persist.NewJsonAdapter(), config.DefaultTokenConfig())
	if err != nil {
		t.Fatalf("NewEnforcer() failed: %v", err)
	}
	newSession := model.NewSession("1", "2", "3")
	newSession.AddTokenSign(&model.TokenSign{
		Value:  "2",
		Device: "device",
	})
	newSession.AddTokenSign(&model.TokenSign{
		Value:  "3",
		Device: "device",
	})

	t.Log(newSession.Json())
	println(newSession.TokenSignSize())

	err = enforcer.SetSession("1", newSession, 565)
	if err != nil {
		t.Errorf("SetSession() failed: %v", err)
	}
	session := enforcer.GetSession("1")
	if id := session.Id; id != "1" {
		t.Errorf("GetSession() failed")
	}
	if num := len(session.TokenSignList); num != 2 {
		t.Fatalf("unexpected session tokenSignList length = %v", num)
	}

	err = enforcer.UpdateSession("1", model.NewSession("4", "5", "6"))
	if err != nil {
		t.Errorf("UpdateSession() failed: %v", err)
	}
	session = enforcer.GetSession("1")
	if id := session.Id; id != "4" {
		t.Errorf("GetSession() failed")
	}

}

func TestEnforcer_Banned(t *testing.T) {
	err, enforcer, _ := NewTestEnforcer(t)
	if err != nil {
		t.Fatalf("NewTestEnforcer() failed: %v", err)
	}
	err = enforcer.Banned("1", "comment", 1, 100)
	if err != nil {
		t.Fatalf("Banned() failed: %v", err)
	}
	isBanned := enforcer.IsBanned("1", "comment")
	if !isBanned {
		t.Errorf("unexpected isBanned is false")
	}
	level, err := enforcer.GetBannedLevel("1", "comment")
	if err != nil {
		t.Errorf("GetBannedLevel() failed: %v", err)
	}
	if level != 1 {
		t.Errorf("unexpected banned level = %v", level)
	}

	err = enforcer.UnBanned("1", "comment")
	if err != nil {
		t.Fatalf("UnBanned() failed: %v", err)
	}
	isBanned = enforcer.IsBanned("1", "comment")
	if isBanned {
		t.Errorf("unexpected isBanned is false")
	}
}

func TestEnforcer_GetBannedTime(t *testing.T) {
	err, enforcer, _ := NewTestEnforcer(t)
	if err != nil {
		t.Fatalf("NewTestEnforcer() failed: %v", err)
	}
	err = enforcer.Banned("1", "comment", 1, 100)
	if err != nil {
		t.Fatalf("Banned() failed: %v", err)
	}

	t.Logf("banned time = %v", enforcer.GetBannedTime("1", "comment"))

	err = enforcer.Banned("1", "comment", 1, -1)
	if err != nil {
		t.Fatalf("Banned() failed: %v", err)
	}

	t.Logf("banned time = %v", enforcer.GetBannedTime("1", "comment"))

}
