package token_go

import (
	"errors"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/model"
	"strconv"
)

// createLoginToken create by config.TokenConfig and model.Login
func (e *Enforcer) createLoginToken(id string, loginModel *model.Login) (string, error) {
	tokenConfig := e.config
	var tokenValue string
	var err error
	// if isConcurrent is false,
	if !tokenConfig.IsConcurrent {
		err = e.Replaced(id, loginModel.Device)
		if err != nil {
			return "", err
		}
	}

	// if loginModel set token, return directly
	if loginModel.Token != "" {
		return loginModel.Token, nil
	}

	// if share token
	if tokenConfig.IsConcurrent && tokenConfig.IsShare {
		// reuse the previous token.
		if v := e.GetSession(id); v != nil {
			tokenValue = v.GetLastTokenByDevice(loginModel.Device)
			if tokenValue != "" {
				return tokenValue, nil
			}

		}
	}

	// create new token
	tokenValue, err = e.generateFunc.Exec(tokenConfig.TokenStyle)
	if err != nil {
		return "", err
	}

	return tokenValue, nil
}

// ResponseToken set token to cookie or header
func (e *Enforcer) ResponseToken(tokenValue string, loginModel *model.Login, ctx ctx.Context) error {
	if ctx == nil {
		return nil
	}
	tokenConfig := e.config

	// set token to cookie
	if tokenConfig.IsReadCookie {
		cookieTimeout := tokenConfig.Timeout
		if loginModel.IsLastingCookie {
			cookieTimeout = -1
		}
		// add cookie use tokenConfig.CookieConfig
		ctx.Response().AddCookie(tokenConfig.TokenName,
			tokenValue,
			tokenConfig.CookieConfig.Path,
			tokenConfig.CookieConfig.Domain,
			cookieTimeout)
	}

	// set token to header
	if loginModel.IsWriteHeader {
		ctx.Response().SetHeader(tokenConfig.TokenName, tokenValue)
	}

	return nil
}

// LogoutByToken clear token info
func (e *Enforcer) LogoutByToken(token string) error {
	var err error
	// delete token-id
	id := e.GetIdByToken(token)
	if id == "" {
		return errors.New("not logged in")
	}
	// delete token-id
	err = e.adapter.Delete(e.spliceTokenKey(token))
	if err != nil {
		return err
	}
	session := e.GetSession(id)
	if session != nil {
		// delete tokenSign
		session.RemoveTokenSign(token)
		err = e.updateSession(id, session)
		if err != nil {
			return err
		}
	}
	// check TokenSignList length, if length == 0, delete this session
	if session != nil && session.TokenSignSize() == 0 {
		err = e.deleteSession(id)
		if err != nil {
			return err
		}
	}

	e.logger.Logout(e.loginType, id, token)

	if e.watcher != nil {
		e.watcher.Logout(e.loginType, id, token)
	}

	return nil
}

// validateValue validate if value is proper
func (e *Enforcer) validateValue(str string) (bool, error) {
	i, err := strconv.Atoi(str)
	// if convert err return true
	if err != nil {
		return true, nil
	}
	if i == constant.BeReplaced {
		return false, errors.New("this account is replaced")
	}
	if i == constant.BeKicked {
		return false, errors.New("this account is kicked out")
	}
	if i == constant.BeBanned {
		return false, errors.New("this account is banned")
	}
	return true, nil
}

// spliceSessionKey splice session-id key
func (e *Enforcer) spliceSessionKey(id string) string {
	return e.config.TokenName + ":" + e.loginType + ":session:" + id
}

// spliceTokenKey splice token-id key
func (e *Enforcer) spliceTokenKey(id string) string {
	return e.config.TokenName + ":" + e.loginType + ":token:" + id
}
