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
	newConfig, err := config.NewConfig(conf)
	if err != nil {
		return nil, err
	}
	enforcer, err := InitWithConfig(newConfig.(*config.FileConfig).TokenConfig, adapter)
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

func (e *Enforcer) SetWatcher(watcher persist.Watcher) {
	e.watcher = watcher
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
				for element, i := session.TokenSignList.Front(), 0; element != nil && i < session.TokenSignList.Len()-int(tokenConfig.MaxLoginCount); element, i = element.Next(), i+1 {
					tokenSign := element.Value.(*model.TokenSign)
					// delete tokenSign
					session.RemoveTokenSign(tokenSign.Value)
					err = e.updateSession(id, session)
					if err != nil {
						return "", err
					}
					// delete token-id
					err = e.adapter.Delete(e.spliceTokenKey(tokenSign.Value))
					if err != nil {
						return "", err
					}
				}
				// check TokenSignList length, if length == 0, delete this session
				if session != nil && session.TokenSignList.Len() == 0 {
					err = e.deleteSession(id)
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
				err = e.updateSession(id, session)
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

	err := e.logoutByToken(token)

	if err != nil {
		return err
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
		for element := l.Back(); element != nil; element = element.Prev() {
			tokenSign := element.Value.(*model.TokenSign)
			str := e.adapter.GetStr(e.spliceTokenKey(tokenSign.Value))
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

// IsLogin check if user logged in by token.
func (e *Enforcer) IsLogin(ctx ctx.Context) (bool, error) {
	tokenValue := e.GetRequestToken(ctx)
	if tokenValue == "" {
		return false, nil
	}
	str := e.adapter.GetStr(e.spliceTokenKey(tokenValue))
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

func (e *Enforcer) GetLoginId(ctx ctx.Context) (string, error) {
	tokenValue := e.GetRequestToken(ctx)
	str := e.adapter.GetStr(e.spliceTokenKey(tokenValue))
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
		return session.TokenSignList.Len()
	}
	return 0
}

func (e *Enforcer) Banned(id string, service string) error {
	panic("implement me ...")
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
				err := e.updateSession(id, session)
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
	if session != nil && session.TokenSignList.Len() == 0 {
		err := e.deleteSession(id)
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
	if v := e.adapter.Get(e.spliceSessionKey(id)); v != nil {
		session := v.(*model.Session)
		return session
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

func (e *Enforcer) deleteSession(id string) error {
	err := e.adapter.Delete(e.spliceSessionKey(id))
	if err != nil {
		return err
	}
	return nil
}
