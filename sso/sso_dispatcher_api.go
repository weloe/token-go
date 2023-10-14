package sso

import (
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/model"
)

/**
=========dispatcher api
*/
// Login
//
// 1. SSO-Client : User click login button in client.
// 2. SSO-Client : If not login, redirect to SSO-Server's ApiName.SsoAuth( called SsoEnforcer.SsoAuth() method in sso_server_api ).
// 3. SSO-Server : If not login, called config.NotLoginView() method in sso_server_api.
// 4. SSO-Server : In config.NotLoginView(), user entered username and password to login.
// 5. SSO-Server : If login successfully, called ApiName.SsoAuth( called SsoEnforcer.SsoAuth() method in sso_server_api ) again,
//                 then redirect SSO-Client's ApiName.DoLogin with ticket(random string value)
// 6. SSO-Client : Get id through checking ticket(random string).
//                 6.1: If config.IsHttp == true, send http request to SSO-Server. SSO-Server check ticket, register logoutCallback url and returns loginId.
//                      If config.IsHttp == false, use adapter to check ticket.
//                 After check ticket successfully, if you set config.TicketResultHandle, called it.
//                 If config.TicketResultHandle is nil, check id, user login in client actually if loginId doesn't be nil, redirect back url.

// Logout

// If user logout in SSO-Server, request logoutCallback url to notify client to logout
// If user logout in SSO-Client, if ssoConfig.IsSlo && ssoConfig.IsHttp, send request to SSO-Server ApiName.SsoSignout. SSO-Server notify all clients to logout
//
//

// ServerDisPatcher dispatcher SSO-Server api, returns model.Result or string.
func (s *SsoEnforcer) ServerDisPatcher(ctx ctx.Context) interface{} {
	request := ctx.Request()
	apiName := s.apiName
	path := request.Path()
	var res interface{}
	var err error
	if path == apiName.SsoAuth {
		res, err = s.SsoAuth(ctx)
	} else if path == apiName.SsoDoLogin {
		res, err = s.SsoDoLogin(ctx)
	} else if path == apiName.SsoCheckTicket && s.config.IsHttp {
		res, err = s.SsoCheckTicket(ctx)
	} else if path == apiName.SsoSignout {
		res, err = s.SsoSignOut(ctx)
	} else {
		return model.Error().SetMsg("not handle")
	}
	if err != nil {
		return model.Error().SetMsg(err.Error())
	}
	if res == nil {
		return model.Ok()
	}

	return res
}

// ClientDispatcher dispatcher Client api, returns model.Result or string.
func (s *SsoEnforcer) ClientDispatcher(ctx ctx.Context) interface{} {
	request := ctx.Request()
	apiName := s.apiName
	path := request.Path()
	var res interface{}
	var err error

	if path == apiName.SsoLogin {
		res, err = s.SsoClientLogin(ctx)
	} else if path == apiName.SsoLogout {
		res, err = s.SsoClientLogout(ctx)
	} else if path == apiName.SsoLogoutCall && s.config.IsSlo && s.config.IsHttp {
		res, err = s.SsoClientLogoutCall(ctx)
	} else {
		return model.Error().SetMsg("not handle")
	}
	if err != nil {
		return model.Error().SetMsg(err.Error())
	}
	if res == nil {
		return model.Ok()
	}
	return res
}
