package sso

import (
	"errors"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/model"
	"strings"
)

/**
=========processor SSO-Client api
*/

// SsoClientLogin SSO-Client login.
func (s *SsoEnforcer) SsoClientLogin(ctx ctx.Context) (interface{}, error) {
	request := ctx.Request()
	response := ctx.Response()
	paramName := s.paramName
	apiName := s.apiName

	// get back value and redirect
	back := request.Query(paramName.Back)
	if back == "" {
		back = "/"
	}
	isLogin, err := s.enforcer.IsLogin(ctx)
	if err != nil {
		return nil, err
	}
	// if the current client is already logged in, there is no need to redirect to the SSO-Server
	if isLogin {
		response.Redirect(back)
		return nil, nil
	}

	// if isLogin == false, attempt to get ticket
	ticket := request.Query(paramName.Ticket)

	// if ticket == "", redirect to SSO-Server, the default path is /sso/auth
	if ticket == "" {
		serverAuthUrl, err := s.buildServerAuthUrl(request.UrlNoQuery(), back)
		if err != nil {
			return nil, err
		}
		response.Redirect(serverAuthUrl)
		return nil, nil
	}

	// else ticket != "", need to log in by ticket
	var loginId string

	// get current path
	ssoLoginUrl := apiName.SsoLogin
	if s.config.IsHttp {
		// use http request to SSO-Server to check ticket in SSO-Server
		var ssoLogoutCall string
		if s.config.IsSlo {
			// get logout callback url
			if s.config.SsoLogoutCall != "" {
				ssoLogoutCall = s.config.SsoLogoutCall
			} else if ssoLoginUrl != "" {
				ssoLogoutCall = strings.ReplaceAll(request.UrlNoQuery(), ssoLoginUrl, apiName.SsoLogoutCall)
			}
		}
		// send http to check
		checkTicketUrl, err := s.buildCheckTicketUrl(ticket, ssoLogoutCall)
		if err != nil {
			return nil, err
		}
		// send http request
		resp, err := s.request(checkTicketUrl)
		if err != nil {
			return nil, err
		}
		if resp.Code == model.ERROR {
			return nil, errors.New("request failed: " + resp.Msg)
		}
		loginId = resp.Data.(string)
	} else {
		// use adapter to check ticket
		loginId, err = s.CheckTicket(ticket)
		if err != nil {
			return nil, err
		}
	}

	// if set callback TicketResultHandle
	if s.config.TicketResultHandle != nil {
		return s.config.TicketResultHandle(loginId, back)
	}

	// if loginId == "", return error
	if loginId == "" {
		return nil, errors.New("invalid ticket: " + ticket)
	}
	// login in client actually
	token, err := s.enforcer.Login(loginId, ctx)
	if err != nil {
		return nil, err
	}

	// redirect to back
	response.Redirect(back)
	return model.Ok().SetData(token), nil
}

// SsoClientLogout SSO-Client single-logout.
func (s *SsoEnforcer) SsoClientLogout(ctx ctx.Context) (interface{}, error) {
	ssoConfig := s.config

	// enable single-logout and isHttp == false
	if ssoConfig.IsSlo && !ssoConfig.IsHttp {
		isLogin, err := s.enforcer.IsLogin(ctx)
		if err != nil {
			return nil, err
		}
		// check if you are logged in
		if isLogin {
			var id string
			id, err = s.enforcer.GetLoginId(ctx)
			if err != nil {
				return nil, err
			}
			// client logout
			err = s.enforcer.LogoutById(id)
			if err != nil {
				return nil, err
			}
			// callback
			return s.ssoLogoutBack(ctx)
		}

	}

	// enable single-logout and isHttp
	if ssoConfig.IsSlo && ssoConfig.IsHttp {
		isLogin, err := s.enforcer.IsLogin(ctx)
		if err != nil {
			return nil, err
		}
		if !isLogin {
			return s.ssoLogoutBack(ctx)
		}

		// get id
		var id string
		id, err = s.enforcer.GetLoginId(ctx)
		if err != nil {
			return nil, err
		}
		// use SSO-Server single-logout api to logout
		sloUrl, err := s.buildSloUrl(id)
		if err != nil {
			return nil, err
		}
		res, err := s.request(sloUrl)
		if err != nil {
			return nil, err
		}
		// check response data
		if res.Code == model.SUCCESS {
			login, _ := s.enforcer.IsLogin(ctx)
			if login {
				err := s.enforcer.Logout(ctx)
				if err != nil {
					return nil, err
				}
			}
			return s.ssoLogoutBack(ctx)
		} else {
			return nil, errors.New("request failed: " + res.Msg)
		}
	}

	return nil, errors.New("not handle")
}

// SsoClientLogoutCall client logout callback.
func (s *SsoEnforcer) SsoClientLogoutCall(ctx ctx.Context) (interface{}, error) {
	request := ctx.Request()
	loginId := request.Query(s.paramName.LoginId)
	if loginId == "" {
		return nil, errors.New("request param must include loginId")
	}
	// check request param
	err := s.checkRequest(request)
	if err != nil {
		return nil, err
	}
	// logout
	err = s.enforcer.LogoutById(loginId)
	if err != nil {
		return nil, err
	}
	return model.Ok().SetMsg("logout callback success"), nil
}

// GetData client build url to sent http to get data from SSO-Server.
func (s *SsoEnforcer) GetData(paramMap map[string]string) (interface{}, error) {
	finalUrl, err := s.buildGetDataUrl(paramMap)
	if err != nil {
		return nil, err
	}
	res, err := s.config.SendHttp(finalUrl)
	return res, err
}
