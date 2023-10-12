package token_go

import "github.com/weloe/token-go/constant"

func (e *Enforcer) CreateTempToken(style string, service string, value string, timeout int64) (string, error) {
	token, err := e.generateFunc.Exec(style)
	if err != nil {
		return "", err
	}
	err = e.adapter.SetStr(e.spliceTempTokenKey(service, token), value, timeout)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (e *Enforcer) GetTempTokenTimeout(service string, tempToken string) int64 {
	if tempToken == "" {
		return constant.NotValueExpire
	}
	return e.adapter.GetTimeout(e.spliceTempTokenKey(service, tempToken))
}

func (e *Enforcer) ParseTempToken(service string, tempToken string) string {
	if tempToken == "" {
		return ""
	}
	return e.adapter.GetStr(e.spliceTempTokenKey(service, tempToken))
}

func (e *Enforcer) DeleteTempToken(service string, tempToken string) error {
	return e.adapter.DeleteStr(e.spliceTempTokenKey(service, tempToken))
}
