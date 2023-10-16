package token_go

import (
	"fmt"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/errors"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/util"
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

func (e *Enforcer) SetIdByToken(id string, tokenValue string, timeout int64) error {
	err := e.adapter.SetStr(e.spliceTokenKey(tokenValue), id, timeout)
	return err
}

func (e *Enforcer) getIdByToken(token string) string {
	return e.adapter.GetStr(e.spliceTokenKey(token))
}

func (e *Enforcer) deleteIdByToken(tokenValue string) error {
	err := e.adapter.DeleteStr(e.spliceTokenKey(tokenValue))
	return err
}

func (e *Enforcer) updateIdByToken(tokenValue string, id string) error {
	err := e.adapter.UpdateStr(e.spliceTokenKey(tokenValue), id)
	return err
}

func (e *Enforcer) setBanned(id string, service string, level int, time int64) error {
	err := e.adapter.SetStr(e.spliceBannedKey(id, service), strconv.Itoa(level), time)
	return err
}

func (e *Enforcer) deleteBanned(id string, service string) error {
	err := e.adapter.DeleteStr(e.spliceBannedKey(id, service))
	return err
}

func (e *Enforcer) getBanned(id string, services string) string {
	s := e.adapter.GetStr(e.spliceBannedKey(id, services))
	return s
}

func (e *Enforcer) getBannedTime(id string, service string) int64 {
	timeout := e.adapter.GetStrTimeout(e.spliceBannedKey(id, service))
	return timeout
}

func (e *Enforcer) setSecSafe(token string, service string, time int64) error {
	err := e.adapter.SetStr(e.spliceSecSafeKey(token, service), constant.DefaultSecondAuthValue, time)
	return err
}

func (e *Enforcer) getSafeTime(token string, service string) int64 {
	timeout := e.adapter.GetTimeout(e.spliceSecSafeKey(token, service))
	return timeout
}

func (e *Enforcer) getSecSafe(token string, service string) string {
	str := e.adapter.GetStr(e.spliceSecSafeKey(token, service))
	return str
}

func (e *Enforcer) deleteSecSafe(token string, service string) error {
	err := e.adapter.DeleteStr(e.spliceSecSafeKey(token, service))
	return err
}

func (e *Enforcer) setTempToken(service string, token string, value string, timeout int64) error {
	err := e.adapter.SetStr(e.spliceTempTokenKey(service, token), value, timeout)
	return err
}

func (e *Enforcer) getTimeoutByTempToken(service string, token string) int64 {
	return e.adapter.GetTimeout(e.spliceTempTokenKey(service, token))
}

func (e *Enforcer) deleteByTempToken(service string, tempToken string) error {
	return e.adapter.DeleteStr(e.spliceTempTokenKey(service, tempToken))
}

func (e *Enforcer) createQRCode(id string, timeout int64) error {
	return e.adapter.Set(e.spliceQRCodeKey(id), model.NewQRCode(id), timeout)
}

func (e *Enforcer) getQRCode(id string) *model.QRCode {
	i := e.adapter.Get(e.spliceQRCodeKey(id), util.GetType(&model.QRCode{}))
	if i == nil {
		return nil
	}
	return i.(*model.QRCode)
}

func (e *Enforcer) getAndCheckQRCodeState(QRCodeId string, want model.QRCodeState) (*model.QRCode, error) {
	qrCode := e.getQRCode(QRCodeId)
	if qrCode == nil {
		return nil, fmt.Errorf("QRCode doesn't exist: %v", QRCodeId)
	}
	if s := e.GetQRCodeState(QRCodeId); s != want {
		return nil, fmt.Errorf("QRCode %v state error: unexpected state value %v, want is %v", QRCodeId, s, want)
	}
	return qrCode, nil
}

func (e *Enforcer) getQRCodeTimeout(id string) int64 {
	return e.adapter.GetTimeout(e.spliceQRCodeKey(id))
}

func (e *Enforcer) updateQRCode(id string, qrCode *model.QRCode) error {
	return e.adapter.Update(e.spliceQRCodeKey(id), qrCode)
}

func (e *Enforcer) deleteQRCode(id string) error {
	return e.adapter.Delete(e.spliceQRCodeKey(id))
}

func (e *Enforcer) getByTempToken(service string, tempToken string) string {
	return e.adapter.GetStr(e.spliceTempTokenKey(service, tempToken))
}

// spliceSessionKey splice session-id key
func (e *Enforcer) spliceSessionKey(id string) string {
	return e.config.TokenName + ":" + e.loginType + ":session:" + id
}

// spliceTokenKey splice token-id key
func (e *Enforcer) spliceTokenKey(token string) string {
	return e.config.TokenName + ":" + e.loginType + ":token:" + token
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

func (e *Enforcer) spliceQRCodeKey(QRCodeId string) string {
	return e.config.TokenName + ":" + "QRCode" + ":QRCode" + QRCodeId
}

func (e *Enforcer) SetJwtSecretKey(key string) {
	e.config.JwtSecretKey = key
}
