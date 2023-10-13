package token_go

import (
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/errors"
	"github.com/weloe/token-go/model"
	"math"
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
		var cookieTimeout int64
		if !loginModel.IsLastingCookie {
			cookieTimeout = -1
		} else {
			if loginModel.Timeout != 0 {
				cookieTimeout = loginModel.Timeout
			} else {
				cookieTimeout = tokenConfig.Timeout
			}
			if cookieTimeout == constant.NeverExpire {
				cookieTimeout = math.MaxInt64
			}
		}

		if tokenConfig.CookieConfig.Path == "" {
			tokenConfig.CookieConfig.Path = "/"
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
		ctx.Response().AddHeader(constant.AccessControlExposeHeaders, tokenConfig.TokenName)
	}

	return nil
}

// checkId check id
func (e *Enforcer) checkId(str string) (bool, error) {
	i, err := strconv.Atoi(str)
	// if convert err return true
	if err != nil {
		return true, nil
	}
	if i == constant.BeReplaced {
		return false, errors.BeReplaced
	}
	if i == constant.BeKicked {
		return false, errors.BeKicked
	}
	if i == constant.BeBanned {
		return false, errors.BeBanned
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

func (e *Enforcer) spliceBannedKey(id string, service string) string {
	return e.config.TokenName + ":" + e.loginType + ":ban:" + service + ":" + id
}

func (e *Enforcer) spliceSecSafeKey(token string, service string) string {
	return e.config.TokenName + ":" + e.loginType + ":safe:" + service + ":" + token
}

func (e *Enforcer) spliceTempTokenKey(service string, token string) string {
	return e.config.TokenName + ":" + "temp-token" + ":temp:" + service + ":" + token
}

func (e *Enforcer) SetJwtSecretKey(key string) {
	e.config.JwtSecretKey = key
}
