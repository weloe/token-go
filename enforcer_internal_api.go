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
			var tokenSignList []*model.TokenSign
			// if device is empty, get all tokenSign
			if loginModel.Device == "" {
				tokenSignList = v.TokenSignList
			} else {
				tokenSignList = v.GetFilterTokenSignSlice(loginModel.Device)
			}
			// get the last token value
			if len(tokenSignList) != 0 && tokenSignList[len(tokenSignList)-1] != nil {
				tokenValue = tokenSignList[len(tokenSignList)-1].Value
				if tokenValue != "" {
					return tokenValue, nil
				}
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
	if e.config.IsWriteHeader {
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

func (e *Enforcer) createRefreshToken(id string, tokenValue string, loginModel *model.Login) (string, error) {
	// create refreshToken
	var err error
	if loginModel.RefreshToken == "" {
		loginModel.RefreshToken, err = e.generateFunc.Exec(e.config.TokenStyle)
		if err != nil {
			return "", err
		}
	}
	if loginModel.RefreshTokenTimeout == 0 {
		loginModel.RefreshTokenTimeout = e.config.RefreshTokenTimeout
	}
	err = e.setRefreshToken(&model.RefreshTokenSign{
		Id:           id,
		Token:        tokenValue,
		RefreshValue: loginModel.RefreshToken,
		Device:       loginModel.Device,
	}, loginModel.RefreshTokenTimeout)

	if err != nil {
		return "", err
	}
	return loginModel.RefreshToken, nil
}

// responseRefreshToken set token to cookie or header
func (e *Enforcer) responseRefreshToken(refreshTokenValue string, loginModel *model.Login, ctx ctx.Context) error {
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
			if loginModel.RefreshTokenTimeout != 0 {
				cookieTimeout = loginModel.RefreshTokenTimeout
			} else {
				cookieTimeout = tokenConfig.Timeout * 2
			}
			if cookieTimeout == constant.NeverExpire {
				cookieTimeout = math.MaxInt64
			}
		}

		if tokenConfig.CookieConfig.Path == "" {
			tokenConfig.CookieConfig.Path = "/"
		}

		// add cookie use tokenConfig.CookieConfig
		ctx.Response().AddCookie(tokenConfig.RefreshTokenName,
			refreshTokenValue,
			tokenConfig.CookieConfig.Path,
			tokenConfig.CookieConfig.Domain,
			cookieTimeout)
	}
	// set token to header
	if tokenConfig.IsWriteHeader {
		ctx.Response().SetHeader(tokenConfig.RefreshTokenName, refreshTokenValue)
		ctx.Response().AddHeader(constant.AccessControlExposeHeaders, tokenConfig.RefreshTokenName)
	}

	return nil
}

func (e *Enforcer) SetIdByToken(id string, tokenValue string, timeout int64) error {
	err := e.notifySetStr(e.spliceTokenKey(tokenValue), id, timeout)
	return err
}

func (e *Enforcer) getIdByToken(token string) string {
	return e.adapter.GetStr(e.spliceTokenKey(token))
}

func (e *Enforcer) deleteIdByToken(tokenValue string) error {
	err := e.notifyDelete(e.spliceTokenKey(tokenValue))
	return err
}

func (e *Enforcer) updateIdByToken(tokenValue string, id string) error {
	err := e.notifyUpdateStr(e.spliceTokenKey(tokenValue), id)
	return err
}

func (e *Enforcer) updateTokenTimeout(token string, timeout int64) error {
	err := e.notifyUpdateTimeout(e.spliceTokenKey(token), timeout)
	return err
}

func (e *Enforcer) deleteRefreshToken(tokenValue string) error {
	refreshToken := e.getRefreshTokenValue(tokenValue)
	err := e.notifyDelete(e.spliceRefreshTokenKey(tokenValue))
	if err != nil {
		return err
	}
	err = e.notifyDelete(e.spliceRefreshTokenSignKey(refreshToken))
	return err
}

func (e *Enforcer) setRefreshToken(refreshTokenSign *model.RefreshTokenSign, timeout int64) error {
	err := e.notifySet(e.spliceRefreshTokenSignKey(refreshTokenSign.RefreshValue), refreshTokenSign, timeout)
	if err != nil {
		return err
	}
	err = e.notifySetStr(e.spliceRefreshTokenKey(refreshTokenSign.Token), refreshTokenSign.RefreshValue, timeout)
	return err
}

func (e *Enforcer) getRefreshTokenValue(tokenValue string) string {
	return e.adapter.GetStr(e.spliceRefreshTokenKey(tokenValue))
}

func (e *Enforcer) getRefreshTokenSign(refreshToken string) *model.RefreshTokenSign {
	get := e.adapter.Get(e.spliceRefreshTokenSignKey(refreshToken), util.GetType(&model.RefreshTokenSign{}))
	if get != nil {
		return get.(*model.RefreshTokenSign)
	}
	return nil
}

func (e *Enforcer) setBanned(id string, service string, level int, time int64) error {
	err := e.notifySetStr(e.spliceBannedKey(id, service), strconv.Itoa(level), time)
	return err
}

func (e *Enforcer) deleteBanned(id string, service string) error {
	err := e.notifyDelete(e.spliceBannedKey(id, service))
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
	err := e.notifySetStr(e.spliceSecSafeKey(token, service), constant.DefaultSecondAuthValue, time)
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
	err := e.notifyDelete(e.spliceSecSafeKey(token, service))
	return err
}

func (e *Enforcer) setTempToken(service string, token string, value string, timeout int64) error {
	err := e.notifySetStr(e.spliceTempTokenKey(service, token), value, timeout)
	return err
}

func (e *Enforcer) getTimeoutByTempToken(service string, token string) int64 {
	return e.adapter.GetTimeout(e.spliceTempTokenKey(service, token))
}

func (e *Enforcer) deleteByTempToken(service string, tempToken string) error {
	return e.notifyDelete(e.spliceTempTokenKey(service, tempToken))
}

func (e *Enforcer) createQRCode(id string, timeout int64) error {
	return e.notifySet(e.spliceQRCodeKey(id), model.NewQRCode(id), timeout)
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
	return e.notifyUpdate(e.spliceQRCodeKey(id), qrCode)
}

func (e *Enforcer) deleteQRCode(id string) error {
	return e.notifyDelete(e.spliceQRCodeKey(id))
}

func (e *Enforcer) getByTempToken(service string, tempToken string) string {
	return e.adapter.GetStr(e.spliceTempTokenKey(service, tempToken))
}

// spliceSessionKey splice id-session key
func (e *Enforcer) spliceSessionKey(id string) string {
	return e.config.TokenName + ":" + e.loginType + ":session:" + id
}

// spliceRefreshTokenSignKey splice refreshToken-refreshToken key
func (e *Enforcer) spliceRefreshTokenSignKey(refreshToken string) string {
	return e.config.TokenName + ":" + e.loginType + ":refreshSign:" + refreshToken
}

// spliceRefreshTokenKey splice token-refreshToken key
func (e *Enforcer) spliceRefreshTokenKey(token string) string {
	return e.config.TokenName + ":" + e.loginType + ":refresh:" + token
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

func (e *Enforcer) notifySetStr(key string, value string, timeout int64) error {
	if e.shouldNotifyDispatcher() {
		return e.dispatcher.SetAllStr(key, value, timeout)
	}
	err := e.adapter.SetStr(key, value, timeout)
	if err != nil {
		return err
	}
	if e.shouldNotifyUpdatableWatcher() {
		return e.updatableWatcher.UpdateForSetStr(key, value, timeout)
	}
	return nil
}

func (e *Enforcer) notifyUpdateStr(key string, value string) error {
	if e.shouldNotifyDispatcher() {
		return e.dispatcher.UpdateAllStr(key, value)
	}
	err := e.adapter.UpdateStr(key, value)
	if err != nil {
		return err
	}
	if e.shouldNotifyUpdatableWatcher() {
		return e.updatableWatcher.UpdateForUpdateStr(key, value)
	}
	return nil
}

func (e *Enforcer) notifySet(key string, value interface{}, timeout int64) error {
	if e.shouldNotifyDispatcher() {
		return e.dispatcher.SetAll(key, value, timeout)
	}
	err := e.adapter.Set(key, value, timeout)
	if err != nil {
		return err
	}
	if e.shouldNotifyUpdatableWatcher() {
		return e.updatableWatcher.UpdateForSet(key, value, timeout)
	}
	return nil
}

func (e *Enforcer) notifyUpdate(key string, value interface{}) error {
	if e.shouldNotifyDispatcher() {
		return e.dispatcher.UpdateAll(key, value)
	}
	err := e.adapter.Update(key, value)
	if err != nil {
		return err
	}
	if e.shouldNotifyUpdatableWatcher() {
		return e.updatableWatcher.UpdateForUpdate(key, value)
	}
	return nil
}

func (e *Enforcer) notifyDelete(key string) error {
	if e.shouldNotifyDispatcher() {
		return e.dispatcher.DeleteAll(key)
	}
	err := e.adapter.Delete(key)
	if err != nil {
		return err
	}
	if e.shouldNotifyUpdatableWatcher() {
		return e.updatableWatcher.UpdateForDelete(key)
	}
	return nil
}

// nolint:golint,unused
func (e *Enforcer) notifyUpdateTimeout(key string, timeout int64) error {
	if e.shouldNotifyDispatcher() {
		return e.dispatcher.UpdateAllTimeout(key, timeout)
	}
	err := e.adapter.UpdateTimeout(key, timeout)
	if err != nil {
		return err
	}
	if e.shouldNotifyUpdatableWatcher() {
		return e.updatableWatcher.UpdateForUpdateTimeout(key, timeout)
	}
	return nil
}

func (e *Enforcer) shouldNotifyDispatcher() bool {
	if e.dispatcher != nil && e.notifyDispatcher {
		return true
	}
	return false
}

func (e *Enforcer) shouldNotifyUpdatableWatcher() bool {
	if e.updatableWatcher != nil && e.notifyUpdatableWatcher {
		return true
	}
	return false
}

// nolint:golint,unused
func (e *Enforcer) shouldPersist() bool {
	return e.adapter != nil
}
