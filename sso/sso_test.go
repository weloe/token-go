package sso

import (
	tokenGo "github.com/weloe/token-go"
	"github.com/weloe/token-go/config"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/util"
	"testing"
)

func SendGetRequest(url string) (string, error) {
	return util.SendGetRequest(url)
}

func TestNewSsoServerEnforcer(t *testing.T) {
	var err error
	// use default adapter
	adapter := tokenGo.NewDefaultAdapter()
	enforcer, err := tokenGo.NewEnforcer(adapter)
	if err != nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}
	// enable logger
	enforcer.EnableLog()
	ssoOptions := &config.SsoOptions{
		Mode:          "",
		TicketTimeout: 300,
		AllowUrl:      "*",
		IsSlo:         true,

		IsHttp:    true,
		ServerUrl: "http://token-go-sso-server.com:9000",
		NotLoginView: func() interface{} {
			msg := "not log in SSO-Server, please visit <a href='/sso/doLogin?name=tokengo&pwd=123456' target='_blank'> doLogin </a>"
			return msg
		},
		DoLoginHandle: func(name string, pwd string, ctx ctx.Context) (interface{}, error) {
			if name != "tokengo" {
				return "name error", nil
			}
			if pwd != "123456" {
				return "pwd error", nil
			}
			token, err := enforcer.Login("1001", ctx)
			if err != nil {
				return nil, err
			}
			return model.Ok().SetData(token), nil
		},
		SendHttp: func(url string) (string, error) {
			return SendGetRequest(url)
		},
	}
	signOptions := &config.SignOptions{
		SecretKey:    "kQwIOrYvnXmSDkwEiFngrKidMcdrgKor",
		IsCheckNonce: true,
	}
	ssoEnforcer, err := NewSsoEnforcer(&Options{
		SsoOptions:  ssoOptions,
		SignOptions: signOptions,
		Enforcer:    enforcer,
	})
	if err != nil {
		t.Errorf("NewSsoEnforcer() failed: %v", err)
	}
	t.Logf("enforcer: %v", ssoEnforcer)
}
func TestNewSsoClient3Enforcer(t *testing.T) {
	var err error
	// use default adapter
	adapter := tokenGo.NewDefaultAdapter()
	enforcer, err := tokenGo.NewEnforcer(adapter)
	if err != nil {
		t.Errorf("NewEnforcer() failed: %v", err)
	}
	// enable logger
	enforcer.EnableLog()
	ssoOptions := &config.SsoOptions{
		AuthUrl:        "/sso/auth",
		IsSlo:          true,
		IsHttp:         true,
		SloUrl:         "/sso/signout",
		CheckTicketUrl: "/sso/checkTicket",
		ServerUrl:      "http://token-go-sso-server.com:9000",
		SendHttp: func(url string) (string, error) {
			return SendGetRequest(url)
		},
	}
	signOptions := &config.SignOptions{
		SecretKey:    "kQwIOrYvnXmSDkwEiFngrKidMcdrgKor",
		IsCheckNonce: true,
	}
	ssoEnforcer, err := NewSsoEnforcer(&Options{
		SsoOptions:  ssoOptions,
		SignOptions: signOptions,
		Enforcer:    enforcer,
	})
	if err != nil {
		t.Errorf("NewSsoEnforcer() failed: %v", err)
	}
	t.Logf("enforcer: %v", ssoEnforcer)
}
