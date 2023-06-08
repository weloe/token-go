package sso

// ApiName  sso api name, used to dispatcher request.
type ApiName struct {
	// sso-server auth url
	SsoAuth string
	// sso-server rest api login url
	SsoDoLogin string
	// sso-server check ticket url
	SsoCheckTicket string
	// sso-server get user info url
	SsoUserInfo string
	// sso-server single logout url
	SsoSignout string
	// sso-client login url
	SsoLogin string
	// sso-client single logout url
	SsoLogout string
	// sso-client logout callback url
	SsoLogoutCall string
}

func DefaultApiName() *ApiName {
	return &ApiName{
		SsoAuth:        "/sso/auth",
		SsoDoLogin:     "/sso/doLogin",
		SsoCheckTicket: "/sso/checkTicket",
		SsoUserInfo:    "/sso/userInfo",
		SsoSignout:     "/sso/signout",
		SsoLogin:       "/sso/login",
		SsoLogout:      "/sso/logout",
		SsoLogoutCall:  "/sso/logoutCall",
	}
}
