package log

import (
	"github.com/weloe/token-go/model"
	"log"
)

var _ Logger = (*DefaultLogger)(nil)

type DefaultLogger struct {
	enable bool
}

func (d *DefaultLogger) StartCleanTimer(period int64) {
	log.Printf("timer period = %v, timer start", period)
}

func (d *DefaultLogger) Enable(bool bool) {
	d.enable = bool
}

func (d *DefaultLogger) IsEnabled() bool {
	return d.enable
}

func (d *DefaultLogger) Login(loginType string, id interface{}, tokenValue string, loginModel *model.Login) {
	if !d.enable {
		return
	}
	log.Printf("Login: loginId = %v, loginType = %v, tokenValue = %v, "+
		"loginMode = %v", id, loginType, tokenValue, loginModel)

}

func (d *DefaultLogger) Logout(loginType string, id interface{}, tokenValue string) {
	if !d.enable {
		return
	}
	log.Printf("Logout: loginId = %v, loginType = %v, tokenValue = %v", id, loginType, tokenValue)
}

func (d *DefaultLogger) Kickout(loginType string, id interface{}, tokenValue string) {
	if !d.enable {
		return
	}
	log.Printf("Kickout: loginId = %v, loginType = %v, tokenValue = %v", id, loginType, tokenValue)
}

func (d *DefaultLogger) Replace(loginType string, id interface{}, tokenValue string) {
	if !d.enable {
		return
	}
	log.Printf("Replaced: loginId = %v, loginType = %v, tokenValue = %v", id, loginType, tokenValue)
}

func (d *DefaultLogger) Ban(loginType string, id interface{}, service string, level int, time int64) {
	if !d.enable {
		return
	}
	log.Printf("Banned: loginId = %v, loginType = %v, service = %v, level = %v, time = %v", id, loginType, service, level, time)
}

func (d *DefaultLogger) UnBan(loginType string, id interface{}, service string) {
	if !d.enable {
		return
	}
	log.Printf("UnBanned: loginId = %v, loginType = %v, service = %v", id, loginType, service)
}

func (d *DefaultLogger) RefreshToken(tokenValue string, id interface{}, timeout int64) {
	if !d.enable {
		return
	}
	log.Printf("RefreshToken: loginId = %v, tokenValue = %v, timeout = %v", id, tokenValue, timeout)
}

func (d *DefaultLogger) OpenSafe(loginType string, token string, service string, time int64) {
	if !d.enable {
		return
	}
	log.Printf("OpenSafe: loginType = %v, tokenValue = %v, service = %v, timeout = %v ", loginType, token, service, time)
}

func (d *DefaultLogger) CloseSafe(loginType string, token string, service string) {
	if !d.enable {
		return
	}
	log.Printf("CloseSafe: loginType = %v, tokenValue = %v, service = %v ", loginType, token, service)
}
