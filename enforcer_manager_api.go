package token_go

import (
	"fmt"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/model"
	"strconv"
)

// Replaced replace other user
func (e *Enforcer) Replaced(id string, device string) error {
	var err error
	if session := e.GetSession(id); session != nil {
		// get by login device
		tokenSignList := session.GetFilterTokenSign(device)
		// sign account replaced
		for element := tokenSignList.Front(); element != nil; element = element.Next() {
			if tokenSign, ok := element.Value.(*model.TokenSign); ok {
				elementV := tokenSign.Value
				session.RemoveTokenSign(elementV)
				err = e.UpdateSession(id, session)
				if err != nil {
					return err
				}
				// sign token replaced
				err = e.updateIdByToken(elementV, strconv.Itoa(constant.BeReplaced))
				if err != nil {
					return err
				}

				// called logger
				e.logger.Replace(e.loginType, id, elementV)

				// called watcher
				if e.watcher != nil {
					e.watcher.Replace(e.loginType, id, elementV)
				}
			}
		}
	}
	return nil
}

// Kickout kickout user
func (e *Enforcer) Kickout(id string, device string) error {
	session := e.GetSession(id)
	if session != nil {
		// get by login device
		tokenSignList := session.GetFilterTokenSign(device)
		// sign account kicked
		for element := tokenSignList.Front(); element != nil; element = element.Next() {
			if tokenSign, ok := element.Value.(*model.TokenSign); ok {
				elementV := tokenSign.Value
				session.RemoveTokenSign(elementV)
				err := e.UpdateSession(id, session)
				if err != nil {
					return err
				}
				// sign token kicked
				err = e.updateIdByToken(elementV, strconv.Itoa(constant.BeKicked))
				if err != nil {
					return err
				}

				// called logger
				e.logger.Kickout(e.loginType, id, elementV)

				// called watcher
				if e.watcher != nil {
					e.watcher.Kickout(e.loginType, id, elementV)
				}
			}
		}

	}
	// check TokenSignList length, if length == 0, delete this session
	if session != nil && session.TokenSignSize() == 0 {
		err := e.DeleteSession(id)
		if err != nil {
			return err
		}
	}
	return nil
}

// Banned ban user, if time == 0,the timeout is not set
func (e *Enforcer) Banned(id string, service string, level int, time int64) error {
	if id == "" || service == "" {
		return fmt.Errorf("parameter cannot be nil")
	}
	if level < 1 {
		return fmt.Errorf("unexpected level = %v, level must large or equal 1", level)
	}
	err := e.setBanned(id, service, level, time)
	if err != nil {
		return err
	}

	// callback
	e.logger.Ban(e.loginType, id, service, level, time)
	if e.watcher != nil {
		e.watcher.Ban(e.loginType, id, service, level, time)
	}

	return nil
}

// UnBanned Unblock user account
func (e *Enforcer) UnBanned(id string, services ...string) error {
	if id == "" {
		return fmt.Errorf("parmeter id can not be nil")
	}
	if len(services) == 0 {
		return fmt.Errorf("parmeter services length can not be 0")
	}

	for _, service := range services {
		err := e.deleteBanned(id, service)
		if err != nil {
			return err
		}
		e.logger.UnBan(e.loginType, id, service)
		if e.watcher != nil {
			e.watcher.UnBan(e.loginType, id, service)
		}
	}
	return nil
}

// IsBanned if banned return true, else return false
func (e *Enforcer) IsBanned(id string, service string) bool {
	level := e.getBanned(id, service)
	return level != ""
}

// GetBannedLevel get banned level
func (e *Enforcer) GetBannedLevel(id string, service string) (int64, error) {
	str := e.getBanned(id, service)
	if str == "" {
		return 0, fmt.Errorf("loginId = %v, service = %v is not banned", id, service)
	}
	time, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return time, nil
}

func (e *Enforcer) OpenSafe(token string, service string, time int64) error {
	if time == 0 {
		return nil
	}
	err := e.CheckLoginByToken(token)
	if err != nil {
		return err
	}
	err = e.setSecSafe(token, service, time)
	if err != nil {
		return err
	}
	e.logger.OpenSafe(e.loginType, token, service, time)
	if e.watcher != nil {
		e.watcher.OpenSafe(e.loginType, token, service, time)
	}
	return nil
}

func (e *Enforcer) IsSafe(token string, service string) bool {
	if token == "" {
		return false
	}
	str := e.getSecSafe(token, service)
	return str != ""
}

func (e *Enforcer) GetSafeTime(token string, service string) int64 {
	if token == "" {
		return 0
	}
	timeout := e.getSafeTime(token, service)
	return timeout
}

func (e *Enforcer) CloseSafe(token string, service string) error {
	if token == "" {
		return nil
	}
	err := e.deleteSecSafe(token, service)
	if err != nil {
		return err
	}
	e.logger.CloseSafe(e.loginType, token, service)
	if e.watcher != nil {
		e.watcher.CloseSafe(e.loginType, token, service)
	}
	return nil
}

func (e *Enforcer) CreateTempToken(style string, service string, value string, timeout int64) (string, error) {
	token, err := e.generateFunc.Exec(style)
	if err != nil {
		return "", err
	}
	err = e.setTempToken(service, token, value, timeout)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (e *Enforcer) GetTempTokenTimeout(service string, tempToken string) int64 {
	if tempToken == "" {
		return constant.NotValueExpire
	}
	return e.getTimeoutByTempToken(service, tempToken)
}

func (e *Enforcer) ParseTempToken(service string, tempToken string) string {
	if tempToken == "" {
		return ""
	}
	return e.getByTempToken(service, tempToken)
}

func (e *Enforcer) DeleteTempToken(service string, tempToken string) error {
	return e.deleteByTempToken(service, tempToken)
}

func (e *Enforcer) CreateQRCodeState(QRCodeId string, timeout int64) error {
	return e.createQRCode(QRCodeId, timeout)
}

// Scanned update state to constant.WaitAuth, return tempToken
func (e *Enforcer) Scanned(QRCodeId string, loginId string) (string, error) {
	qrCode := e.getQRCode(QRCodeId)
	if qrCode == nil {
		return "", fmt.Errorf("QRCode doesn't exist")
	}
	if qrCode.State != model.WaitScan {
		return "", fmt.Errorf("QRCode state error: unexpected state value %v, want is %v", qrCode.State, model.WaitScan)
	}
	qrCode.State = model.WaitAuth
	qrCode.LoginId = loginId

	err := e.updateQRCode(QRCodeId, qrCode)
	if err != nil {
		return "", err
	}
	tempToken, err := e.CreateTempToken(e.config.TokenStyle, "qrCode", QRCodeId, e.config.Timeout)
	if err != nil {
		return "", err
	}
	return tempToken, nil
}

// ConfirmAuth update state to constant.ConfirmAuth
func (e *Enforcer) ConfirmAuth(tempToken string) error {
	qrCodeId := e.ParseTempToken("qrCode", tempToken)
	if qrCodeId == "" {
		return fmt.Errorf("confirm failed, tempToken error: %v", tempToken)
	}
	qrCode, err := e.getAndCheckQRCodeState(qrCodeId, model.WaitAuth)
	if err != nil {
		return err
	}

	qrCode.State = model.ConfirmAuth
	err = e.updateQRCode(qrCodeId, qrCode)
	if err != nil {
		return err
	}
	err = e.DeleteTempToken("qrCode", tempToken)
	if err != nil {
		return err
	}
	return err
}

// CancelAuth update state to constant.CancelAuth
func (e *Enforcer) CancelAuth(tempToken string) error {
	qrCodeId := e.ParseTempToken("qrCode", tempToken)
	if qrCodeId == "" {
		return fmt.Errorf("confirm failed, tempToken error: %v", tempToken)
	}
	qrCode, err := e.getAndCheckQRCodeState(qrCodeId, model.WaitAuth)
	if err != nil {
		return err
	}
	qrCode.State = model.CancelAuth
	err = e.updateQRCode(qrCodeId, qrCode)
	if err != nil {
		return err
	}
	err = e.DeleteTempToken("qrCode", tempToken)
	if err != nil {
		return err
	}
	return err
}

func (e *Enforcer) GetQRCode(QRCodeId string) *model.QRCode {
	return e.getQRCode(QRCodeId)
}

// GetQRCodeState
//	WaitScan   = 1
//	WaitAuth   = 2
//	ConfirmAuth  = 3
//	CancelAuth = 4
//	Expired    = 5
func (e *Enforcer) GetQRCodeState(QRCodeId string) model.QRCodeState {
	qrCode := e.getQRCode(QRCodeId)
	if qrCode == nil {
		return model.Expired
	}
	return qrCode.State
}

func (e *Enforcer) GetQRCodeTimeout(QRCodeId string) int64 {
	return e.getQRCodeTimeout(QRCodeId)
}

func (e *Enforcer) DeleteQRCode(QRCodeId string) error {
	return e.deleteQRCode(QRCodeId)
}
