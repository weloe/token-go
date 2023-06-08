package config

import (
	"errors"
	"github.com/weloe/token-go/ctx"
	model2 "github.com/weloe/token-go/model"
	"github.com/weloe/token-go/util"
	"strings"
)

func DefaultSsoConfig(serverUrl string, notLoginView func() interface{},
	doLoginHandle func(name string, pwd string, ctx ctx.Context) (interface{}, error),
	ticketResultHandle func(o1 string, s string) (interface{}, error),
	sendHttp func(url string) (string, error)) *SsoConfig {
	if notLoginView == nil {
		notLoginView = func() interface{} {
			return "not logged in to the SSO-Server"
		}
	}
	return &SsoConfig{
		Mode:               "",
		TicketTimeout:      60 * 5,
		AllowUrl:           "*",
		IsSlo:              true,
		IsHttp:             false,
		Client:             "",
		AuthUrl:            "/sso/auth",
		CheckTicketUrl:     "/sso/checkTicket",
		GetDataUrl:         "/sso/getData",
		UserInfoUrl:        "/sso/userInfo",
		SloUrl:             "/sso/signOut",
		SsoLogoutCall:      "",
		ServerUrl:          serverUrl,
		NotLoginView:       notLoginView,
		DoLoginHandle:      doLoginHandle,
		TicketResultHandle: ticketResultHandle,
		SendHttp:           sendHttp,
	}
}

func NewSsoConfig(options *SsoOptions) (*SsoConfig, error) {
	if options == nil {
		options = &SsoOptions{}
	}
	if options.TicketTimeout == 0 {
		options.TicketTimeout = 60 * 5
	}
	if options.AllowUrl == "" {
		options.AllowUrl = "*"
	}
	if options.AuthUrl == "" {
		options.AuthUrl = "/sso/auth"
	}
	if options.CheckTicketUrl == "" {
		options.CheckTicketUrl = "/sso/checkTicket"
	}
	if options.GetDataUrl == "" {
		options.GetDataUrl = "/sso/getData"
	}
	if options.UserInfoUrl == "" {
		options.UserInfoUrl = "/sso/userInfo"
	}
	if options.SloUrl == "" {
		options.SloUrl = "/sso/signout"
	}

	if options.NotLoginView == nil {
		options.NotLoginView = func() interface{} {
			return "not logged in to the SSO-Server"
		}
	}
	if options.DoLoginHandle == nil {
		options.DoLoginHandle = func(name string, pwd string, ctx ctx.Context) (interface{}, error) {
			return model2.Error(), errors.New("SsoConfig.DoLoginHandle is nil")
		}
	}
	if options.IsHttp && options.SendHttp == nil {
		return nil, errors.New("please config SSO SentHttp")
	}

	return &SsoConfig{
		Mode:               options.Mode,
		TicketTimeout:      options.TicketTimeout,
		AllowUrl:           options.AllowUrl,
		IsSlo:              options.IsSlo,
		IsHttp:             options.IsHttp,
		Client:             options.Client,
		AuthUrl:            options.AuthUrl,
		CheckTicketUrl:     options.CheckTicketUrl,
		GetDataUrl:         options.GetDataUrl,
		UserInfoUrl:        options.UserInfoUrl,
		SloUrl:             options.SloUrl,
		SsoLogoutCall:      options.SsoLogoutCall,
		ServerUrl:          options.ServerUrl,
		NotLoginView:       options.NotLoginView,
		DoLoginHandle:      options.DoLoginHandle,
		TicketResultHandle: options.TicketResultHandle,
		SendHttp:           options.SendHttp,
	}, nil
}

type SsoConfig struct {

	// Mode sso mode
	Mode string
	// TicketTimeout ticket timeout
	TicketTimeout int64
	// AllowUrl All allowed authorization callback addresses, separated by ','
	AllowUrl string
	IsSlo    bool
	IsHttp   bool

	// SSO-Client current client name
	Client string
	// SSO-Server auth url
	AuthUrl string
	// SSO-Server check ticket url
	CheckTicketUrl string
	GetDataUrl     string
	UserInfoUrl    string
	SloUrl         string
	SsoLogoutCall  string

	ServerUrl string

	/**
	sso callback func
	*/
	// NotLoginView
	NotLoginView func() interface{}

	// DoLoginHandle login func
	DoLoginHandle func(name string, pwd string, ctx ctx.Context) (interface{}, error)

	// TicketResultHandle called each time the result of the validation ticket is obtained from the SSO-Server
	TicketResultHandle func(loginId string, back string) (interface{}, error)

	// SendHttp sent http
	SendHttp func(url string) (string, error)
}

// SpliceAuthUrl return Server-side single sign-on authorization address
func (c *SsoConfig) SpliceAuthUrl() string {
	return util.SpliceUrl(c.ServerUrl, c.AuthUrl)
}

// SpliceCheckTicketUrl return ticket verification address on the server side
func (c *SsoConfig) SpliceCheckTicketUrl() string {
	return util.SpliceUrl(c.ServerUrl, c.CheckTicketUrl)
}

func (c *SsoConfig) SpliceGetDataUrl() string {
	return util.SpliceUrl(c.ServerUrl, c.GetDataUrl)
}

func (c *SsoConfig) SpliceUserInfoUrl() string {
	return util.SpliceUrl(c.ServerUrl, c.UserInfoUrl)
}

func (c *SsoConfig) SpliceSloUrl() string {
	return util.SpliceUrl(c.ServerUrl, c.SloUrl)
}

// SetAllow set allow callback url
func (c *SsoConfig) SetAllow(urls ...string) *SsoConfig {
	c.AllowUrl = strings.Join(urls, ",")
	return c
}
