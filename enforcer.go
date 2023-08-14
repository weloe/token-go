package token_go

import (
	"errors"
	"fmt"
	"github.com/weloe/token-go/config"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/ctx"
	httpCtx "github.com/weloe/token-go/ctx/go-http-context"
	"github.com/weloe/token-go/log"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/persist"
	"github.com/weloe/token-go/util"
	"net/http"
	"strconv"
)

type Enforcer struct {
	conf         string
	loginType    string
	config       config.TokenConfig
	generateFunc model.GenerateTokenFunc
	adapter      persist.Adapter
	watcher      persist.Watcher
	logger       log.Logger
	authManager  interface{}
}

func NewDefaultAdapter() persist.Adapter {
	return persist.NewDefaultAdapter()
}

func NewHttpContext(req *http.Request, writer http.ResponseWriter) ctx.Context {
	return httpCtx.NewHttpContext(req, writer)
}

func NewEnforcer(adapter persist.Adapter, args ...interface{}) (*Enforcer, error) {
	var err error
	var enforcer *Enforcer
	if len(args) > 2 {
		return nil, fmt.Errorf("NewEnforcer() failed: unexpected args length = %v, it should be less than or equal to 2", len(args))
	}
	if util.HasNil(args) {
		return nil, errors.New("NewEnforcer() failed: parameters cannot be nil")
	}

	if len(args) == 0 {
		enforcer, err = InitWithDefaultConfig(adapter)
	} else if len(args) == 1 {
		switch args[0].(type) {
		case *config.TokenConfig:
			enforcer, err = InitWithConfig(args[0].(*config.TokenConfig), adapter)
		case string:
			enforcer, err = InitWithFile(args[0].(string), adapter)
		default:
			return nil, errors.New("NewEnforcer() failed: the second parameter should be *TokenConfig or string")
		}
	}

	return enforcer, err
}

func InitWithDefaultConfig(adapter persist.Adapter) (*Enforcer, error) {
	if adapter == nil {
		return nil, errors.New("InitWithDefaultConfig() failed: parameters cannot be nil")
	}
	return InitWithConfig(config.DefaultTokenConfig(), adapter)
}

func InitWithFile(conf string, adapter persist.Adapter) (*Enforcer, error) {
	if conf == "" || adapter == nil {
		return nil, errors.New("InitWithFile() failed: parameters cannot be nil")
	}
	newConfig, err := config.ReadConfig(conf)
	if err != nil {
		return nil, err
	}
	enforcer, err := InitWithConfig(newConfig.TokenConfig, adapter)
	enforcer.conf = conf
	return enforcer, err
}

func InitWithConfig(tokenConfig *config.TokenConfig, adapter persist.Adapter) (*Enforcer, error) {
	fm := model.LoadFunctionMap()
	if tokenConfig == nil || adapter == nil {
		return nil, errors.New("InitWithConfig() failed: parameters cannot be nil")
	}
	e := &Enforcer{
		loginType:    "user",
		config:       *tokenConfig,
		generateFunc: fm,
		adapter:      adapter,
		logger:       &log.DefaultLogger{},
	}

	e.startCleanTimer()

	return e, nil
}

// if e.adapter.(type) == *persist.DefaultAdapter, can start cleanTimer
func (e *Enforcer) startCleanTimer() {
	defaultAdapter, ok := e.adapter.(*persist.DefaultAdapter)
	if ok {
		dataRefreshPeriod := e.config.DataRefreshPeriod
		if period := dataRefreshPeriod; period >= 0 {
			defaultAdapter.StartCleanTimer(dataRefreshPeriod)
			e.logger.StartCleanTimer(dataRefreshPeriod)
		}
	}
}

func (e *Enforcer) SetType(t string) {
	e.loginType = t
}

func (e *Enforcer) GetType() string {
	return e.loginType
}

func (e *Enforcer) GetAdapter() persist.Adapter {
	return e.adapter
}

func (e *Enforcer) SetAdapter(adapter persist.Adapter) {
	e.adapter = adapter
}

func (e *Enforcer) GetWatcher() persist.Watcher {
	return e.watcher
}

func (e *Enforcer) SetWatcher(watcher persist.Watcher) {
	e.watcher = watcher
}

func (e *Enforcer) GetLogger() log.Logger {
	return e.logger
}

func (e *Enforcer) SetLogger(logger log.Logger) {
	e.logger = logger
}

func (e *Enforcer) EnableLog() {
	e.logger.Enable(true)
}

func (e *Enforcer) IsLogEnable() bool {
	return e.logger.IsEnabled()
}

// Login login by id and default loginModel, return tokenValue and error
func (e *Enforcer) Login(id string, ctx ctx.Context) (string, error) {
	return e.LoginByModel(id, model.DefaultLoginModel(), ctx)
}

func (e *Enforcer) LoginById(id string) (string, error) {
	return e.Login(id, nil)
}

// LoginByModel login by id and loginModel, return tokenValue and error
func (e *Enforcer) LoginByModel(id string, loginModel *model.Login, ctx ctx.Context) (string, error) {
	if loginModel == nil {
		return "", errors.New("arg loginModel can not be nil")
	}
	var err error
	var session *model.Session
	var tokenValue string
	tokenConfig := e.config

	// allocate token
	tokenValue, err = e.createLoginToken(id, loginModel)

	if err != nil {
		return "", err
	}

	// add tokenSign
	if session = e.GetSession(id); session == nil {
		session = model.NewSession(e.spliceSessionKey(id), "account-session", id)
	}
	session.AddTokenSign(&model.TokenSign{
		Value:  tokenValue,
		Device: loginModel.Device,
	})

	// reset session
	err = e.SetSession(id, session, loginModel.Timeout)
	if err != nil {
		return "", err
	}

	// set token-id
	err = e.adapter.SetStr(e.spliceTokenKey(tokenValue), id, loginModel.Timeout)
	if err != nil {
		return "", err
	}

	// response token
	err = e.responseToken(tokenValue, loginModel, ctx)
	if err != nil {
		return "", err
	}

	// called watcher
	m := &model.Login{
		Device:          loginModel.Device,
		IsLastingCookie: loginModel.IsLastingCookie,
		Timeout:         loginModel.Timeout,
		JwtData:         loginModel.JwtData,
		Token:           tokenValue,
		IsWriteHeader:   loginModel.IsWriteHeader,
	}

	// called logger
	e.logger.Login(e.loginType, id, tokenValue, m)

	if e.watcher != nil {
		e.watcher.Login(e.loginType, id, tokenValue, m)
	}

	// if login success check it
	if tokenConfig.IsConcurrent && !tokenConfig.IsShare {
		// check if the number of sessions for this account exceeds the maximum limit.
		if tokenConfig.MaxLoginCount != -1 {
			if session = e.GetSession(id); session != nil {
				// logout account until loginCount == maxLoginCount if loginCount > maxLoginCount
				for _, tokenSign := range session.TokenSignList {
					if session.TokenSignSize() > int(tokenConfig.MaxLoginCount) {
						// delete tokenSign
						session.RemoveTokenSign(tokenSign.Value)
						err = e.UpdateSession(id, session)
						if err != nil {
							return "", err
						}
						// delete token-id
						err = e.adapter.Delete(e.spliceTokenKey(tokenSign.Value))
						if err != nil {
							return "", err
						}
					}
				}

				// check TokenSignList length, if length == 0, delete this session
				if session != nil && session.TokenSignSize() == 0 {
					err = e.DeleteSession(id)
					if err != nil {
						return "", err
					}
				}
			}
		}

	}

	return tokenValue, nil
}

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
				err = e.adapter.UpdateStr(e.spliceTokenKey(elementV), strconv.Itoa(constant.BeReplaced))
				if err != nil {
					return err
				}

				// called logger
				e.logger.Replace(e.loginType, id, tokenSign.Value)

				// called watcher
				if e.watcher != nil {
					e.watcher.Replace(e.loginType, id, tokenSign.Value)
				}
			}
		}
	}
	return nil
}

// Logout user logout
func (e *Enforcer) Logout(ctx ctx.Context) error {
	tokenConfig := e.config

	token := e.GetRequestToken(ctx)
	if token == "" {
		return errors.New("logout() failed: token doesn't exist")
	}
	if e.config.IsReadCookie {
		ctx.Response().DeleteCookie(tokenConfig.TokenName,
			tokenConfig.CookieConfig.Path,
			tokenConfig.CookieConfig.Domain)
	}

	err := e.LogoutByToken(token)

	if err != nil {
		return err
	}
	return nil
}

// LogoutById force user to logout
func (e *Enforcer) LogoutById(id string) error {
	session := e.GetSession(id)
	if session != nil {
		for _, tokenSign := range session.TokenSignList {
			err := e.LogoutByToken(tokenSign.Value)
			if err != nil {
				return err
			}
		}
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
		err = e.UpdateSession(id, session)
		if err != nil {
			return err
		}
	}
	// check TokenSignList length, if length == 0, delete this session
	if session != nil && session.TokenSignSize() == 0 {
		err = e.DeleteSession(id)
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

// IsLoginById check if user logged in by loginId.
// check all tokenValue and if one is validated return true
func (e *Enforcer) IsLoginById(id string) (bool, error) {
	var error error
	session := e.GetSession(id)
	if session != nil {
		l := session.TokenSignList
		for _, tokenSign := range l {
			str := e.GetIdByToken(tokenSign.Value)
			if str == "" {
				continue
			}
			value, err := e.validateValue(str)
			if err != nil {
				error = err
				continue
			}
			if value {
				return true, nil
			}
		}
	}

	return false, error
}

func (e *Enforcer) GetIdByToken(token string) string {
	return e.adapter.GetStr(e.spliceTokenKey(token))
}

// IsLogin check if user logged in by token.
func (e *Enforcer) IsLogin(ctx ctx.Context) (bool, error) {
	tokenValue := e.GetRequestToken(ctx)
	return e.IsLoginByToken(tokenValue)
}

func (e *Enforcer) IsLoginByToken(tokenValue string) (bool, error) {
	if tokenValue == "" {
		return false, nil
	}
	str := e.GetIdByToken(tokenValue)
	if str == "" {
		return false, nil
	}

	return e.validateValue(str)
}

func (e *Enforcer) CheckLogin(ctx ctx.Context) error {
	_, err := e.GetLoginId(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e *Enforcer) CheckLoginByToken(token string) error {
	_, err := e.GetLoginIdByToken(token)
	if err != nil {
		return err
	}
	return nil
}

func (e *Enforcer) GetLoginId(ctx ctx.Context) (string, error) {
	tokenValue := e.GetRequestToken(ctx)
	return e.GetLoginIdByToken(tokenValue)
}

func (e *Enforcer) GetLoginIdByToken(token string) (string, error) {
	str := e.GetIdByToken(token)
	if str == "" {
		return "", errors.New("GetLoginId() failed: not logged in")
	}
	validate, err := e.validateValue(str)
	if !validate {
		return "", err
	}

	return str, nil
}

func (e *Enforcer) GetLoginCount(id string) int {
	if session := e.GetSession(id); session != nil {
		return session.TokenSignSize()
	}
	return 0
}

// Banned ban user, if time == 0,the timeout is not set
func (e *Enforcer) Banned(id string, service string, level int, time int64) error {
	if id == "" || service == "" {
		return fmt.Errorf("parameter cannot be nil")
	}
	if level < 1 {
		return fmt.Errorf("unexpected level = %v, level must large or equal 1", level)
	}
	err := e.adapter.SetStr(e.spliceBannedKey(id, service), strconv.Itoa(level), time)
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
		err := e.adapter.DeleteStr(e.spliceBannedKey(id, service))
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
func (e *Enforcer) IsBanned(id string, services string) bool {
	s := e.adapter.GetStr(e.spliceBannedKey(id, services))

	return s != ""
}

// GetBannedLevel get banned level
func (e *Enforcer) GetBannedLevel(id string, service string) (int64, error) {
	str := e.adapter.GetStr(e.spliceBannedKey(id, service))
	if str == "" {
		return 0, fmt.Errorf("loginId = %v, service = %v is not banned", id, service)
	}
	time, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return time, nil
}

// GetBannedTime get banned time
func (e *Enforcer) GetBannedTime(id string, service string) int64 {
	timeout := e.adapter.GetStrTimeout(e.spliceBannedKey(id, service))
	return timeout
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
				err = e.adapter.UpdateStr(e.spliceTokenKey(elementV), strconv.Itoa(constant.BeKicked))
				if err != nil {
					return err
				}

				// called logger
				e.logger.Kickout(e.loginType, id, tokenSign.Value)

				// called watcher
				if e.watcher != nil {
					e.watcher.Kickout(e.loginType, id, tokenSign.Value)
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

// GetRequestToken read token from requestHeader | cookie | requestBody
func (e *Enforcer) GetRequestToken(ctx ctx.Context) string {
	var tokenValue string
	if ctx == nil {
		return ""
	}
	if e.config.IsReadHeader {
		if tokenValue = ctx.Request().Header(e.config.TokenName); tokenValue != "" {
			return tokenValue
		}
	}
	if e.config.IsReadCookie {
		if tokenValue = ctx.Request().Cookie(e.config.TokenName); tokenValue != "" {
			return tokenValue
		}

	}
	if e.config.IsReadBody {
		if tokenValue = ctx.Request().PostForm(e.config.TokenName); tokenValue != "" {
			return tokenValue
		}
	}
	return tokenValue
}

// AddTokenGenerateFun add token generate strategy
func (e *Enforcer) AddTokenGenerateFun(tokenStyle string, f model.GenerateFunc) error {
	e.generateFunc.AddFunc(tokenStyle, f)
	return nil
}

func (e *Enforcer) GetSession(id string) *model.Session {
	if v := e.adapter.Get(e.spliceSessionKey(id), util.GetType(&model.Session{})); v != nil {
		return v.(*model.Session)
	}
	return nil
}

func (e *Enforcer) SetSession(id string, session *model.Session, timeout int64) error {
	err := e.adapter.Set(e.spliceSessionKey(id), session, timeout)
	if err != nil {
		return err
	}
	return nil
}

func (e *Enforcer) DeleteSession(id string) error {
	err := e.adapter.Delete(e.spliceSessionKey(id))
	if err != nil {
		return err
	}
	return nil
}

func (e *Enforcer) UpdateSession(id string, session *model.Session) error {
	err := e.adapter.Update(e.spliceSessionKey(id), session)
	if err != nil {
		return err
	}
	return nil
}

func (e *Enforcer) GetTokenConfig() config.TokenConfig {
	return e.config
}

func (e *Enforcer) OpenSafe(token string, service string, time int64) error {
	if time == 0 {
		return nil
	}
	err := e.CheckLoginByToken(token)
	if err != nil {
		return err
	}
	err = e.adapter.SetStr(e.spliceSecSafeKey(token, service), constant.DefaultSecondAuthValue, time)
	if err != nil {
		return err
	}
	if e.watcher != nil {
		e.watcher.OpenSafe(e.loginType, token, service, time)
	}
	return nil
}

func (e *Enforcer) IsSafe(token string, service string) bool {
	if token == "" {
		return false
	}
	str := e.adapter.GetStr(e.spliceSecSafeKey(token, service))
	return str != ""
}

func (e *Enforcer) GetSafeTime(token string, service string) int64 {
	if token == "" {
		return 0
	}
	timeout := e.adapter.GetTimeout(e.spliceSecSafeKey(token, service))
	return timeout
}

func (e *Enforcer) CloseSafe(token string, service string) error {
	if token == "" {
		return nil
	}
	err := e.adapter.DeleteStr(e.spliceSecSafeKey(token, service))
	if err != nil {
		return err
	}
	if e.watcher != nil {
		e.watcher.CloseSafe(e.loginType, token, service)
	}
	return nil
}
