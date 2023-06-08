package config

import "github.com/weloe/token-go/ctx"

// SsoOptions new SsoConfig options.
type SsoOptions struct {
	CookieDomain string
	
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

// SignOptions SignConfig options
type SignOptions struct {
	SecretKey          string
	TimeStampDisparity int64
	IsCheckNonce       bool
}
