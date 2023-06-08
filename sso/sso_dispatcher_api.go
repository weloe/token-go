package sso

import (
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/model"
)

/**
=========dispatcher api
*/

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
