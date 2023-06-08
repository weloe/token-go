package sso

import (
	"errors"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/model"
)

/**
=========processor SSO-Server api
*/

// SsoAuth SSO-Server: auth.
func (s *SsoEnforcer) SsoAuth(ctx ctx.Context) (interface{}, error) {
	request := ctx.Request()
	response := ctx.Response()

	isLogin, err := s.enforcer.IsLogin(ctx)
	if err != nil {
		return nil, err
	}
	// if you have not logged in to the SSO-Server, need to log in first
	if !isLogin {
		return s.config.NotLoginView(), nil
	}
	// if you have logged, check the mode
	mode := request.Query(s.paramName.Mode)

	redirect := request.Query(s.paramName.Redirect)

	// if mode == simple, redirect to client directly
	if mode == constant.MODE_SIMPLE {
		err = s.CheckRedirectUrl(redirect)
		if err != nil {
			return nil, err
		}
		response.Redirect(redirect)
		return nil, nil
	} else {
		// mode = ticket, redirect to client login with new ticket
		id, err := s.enforcer.GetLoginId(ctx)
		if err != nil {
			return nil, err
		}
		redirectUrl, err := s.buildRedirectUrl(id, request.Query(s.paramName.Client), redirect)
		if err != nil {
			return nil, err
		}
		response.Redirect(redirectUrl)
		return nil, nil
	}
}

// SsoDoLogin SSO-Server: rest login api.
func (s *SsoEnforcer) SsoDoLogin(ctx ctx.Context) (interface{}, error) {
	request := ctx.Request()
	paramName := s.paramName
	if s.config.DoLoginHandle == nil {
		return nil, errors.New("SsoConfig.DoLoginHandle is nil")
	}
	resp, err := s.config.DoLoginHandle(request.Query(paramName.Name), request.Query(paramName.Pwd), ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SsoCheckTicket  SSO-Server: check ticket to get loginId, returns loginId
func (s *SsoEnforcer) SsoCheckTicket(ctx ctx.Context) (interface{}, error) {
	paramName := s.paramName
	request := ctx.Request()
	client := request.Query(paramName.Client)
	ticket := request.Query(paramName.Ticket)
	if ticket == "" {
		return nil, errors.New("ticket can not be nil")
	}
	sloCallback := request.Query(paramName.SsoLogoutCall)

	// check ticket
	loginId, err := s.CheckTicketByClient(ticket, client)
	if err != nil {
		return nil, err
	}

	// register single sign out callback url
	err = s.RegisterSloCallbackUrl(loginId, sloCallback)
	if err != nil {
		return nil, err
	}

	if loginId == "" {
		return nil, errors.New("invalid ticket: " + ticket)
	}
	return model.Ok().SetData(loginId), nil
}

// SsoSignOut SSO-Server: single sign-out.
func (s *SsoEnforcer) SsoSignOut(ctx ctx.Context) (interface{}, error) {
	request := ctx.Request()
	paramName := s.paramName

	// if enable single sign-out and request param has loginId
	reqLoginId := request.Query(paramName.LoginId)
	if s.config.IsSlo && reqLoginId == "" {
		loginId, err := s.enforcer.GetLoginId(ctx)
		if err != nil {
			return nil, err
		}
		if loginId != "" {
			err = s.ssoSignOutById(loginId)
			if err != nil {
				return nil, err
			}
			// callback
			return s.ssoLogoutBack(ctx)
		}
	}

	// if enable http,single sign-out and request param has loginId
	if s.config.IsHttp && s.config.IsSlo && reqLoginId != "" {
		err := s.checkRequest(request)
		if err != nil {
			return nil, err
		}
		// Use loginId to get single sign-out urls from session, traverse the urls to send the request to notify client
		err = s.ssoSignOutById(reqLoginId)
		if err != nil {
			return nil, err
		}
		return model.Ok().SetMsg("sso sign-out success"), nil
	}

	return nil, errors.New("not handle")
}
