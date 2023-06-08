package sso

// ParamName http request param name.
type ParamName struct {
	Redirect      string
	Ticket        string
	Back          string
	Mode          string
	LoginId       string
	Client        string
	SsoLogoutCall string
	Name          string
	Pwd           string

	//=== sign param

	TimeStamp string
	Nonce     string
	Sign      string
	SecretKet string
}

func DefaultParamName() *ParamName {
	return &ParamName{
		Redirect:      "redirect",
		Ticket:        "ticket",
		Back:          "back",
		Mode:          "mode",
		LoginId:       "loginId",
		Client:        "client",
		SsoLogoutCall: "ssoLogoutCall",
		Name:          "name",
		Pwd:           "pwd",
		TimeStamp:     "timestamp",
		Nonce:         "nonce",
		Sign:          "sign",
		SecretKet:     "key",
	}
}
