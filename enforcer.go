package token_go

import (
	"errors"
	"fmt"
	"github.com/weloe/token-go/config"
	"github.com/weloe/token-go/ctx"
	httpCtx "github.com/weloe/token-go/ctx/go-http-context"
	"github.com/weloe/token-go/log"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/persist"
	"github.com/weloe/token-go/util"
	log2 "log"
	"net/http"
)

type Enforcer struct {
	conf         string
	loginType    string
	config       config.TokenConfig
	generateFunc model.GenerateTokenFunc
	adapter      persist.Adapter
	watcher      persist.Watcher
	logger       log.Logger

	dispatcher       persist.Dispatcher
	notifyDispatcher bool

	updatableWatcher       persist.UpdatableWatcher
	notifyUpdatableWatcher bool

	authManager interface{}
}

func (e *Enforcer) EnableUpdatableWatcher(b bool) {
	if e.updatableWatcher == nil {
		return
	}
	e.notifyUpdatableWatcher = b
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
	tokenConfig.InitConfig()
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
		if !defaultAdapter.GetCleanTimer() {
			return
		}
		dataRefreshPeriod := e.config.DataRefreshPeriod
		if period := dataRefreshPeriod; period > 0 {
			err := defaultAdapter.StartCleanTimer(period)
			if err != nil {
				log2.Printf("enble adapter cleanTimer failed: %v", err)
				return
			}
			e.logger.StartCleanTimer(period)
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

func (e *Enforcer) SetUpdatableWatcher(watcher persist.UpdatableWatcher) {
	if watcher != nil {
		e.updatableWatcher = watcher
		e.notifyUpdatableWatcher = true
	}
}

func (e *Enforcer) GetUpdatableWatcher() persist.UpdatableWatcher {
	return e.updatableWatcher
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
func (e *Enforcer) Login(id string, ctx ...ctx.Context) (string, error) {
	return e.LoginByModel(id, model.CreateLoginModelByDevice(""), ctx...)
}

func (e *Enforcer) LoginById(id string, device ...string) (string, error) {
	if len(device) > 0 && device[0] != "" {
		return e.LoginByModel(id, model.CreateLoginModelByDevice(device[0]), nil)
	}

	return e.Login(id, nil)
}

// LoginByModel login by id and loginModel, return tokenValue and error
func (e *Enforcer) LoginByModel(id string, loginModel *model.Login, c ...ctx.Context) (string, error) {
	if loginModel == nil {
		return "", errors.New("arg loginModel can not be nil")
	}
	if len(c) == 0 {
		c = []ctx.Context{nil}
	}
	var err error
	var session *model.Session
	var tokenValue string
	tokenConfig := e.config

	// allocate token
	tokenValue, err = e.createLoginToken(id, loginModel)
	device := loginModel.Device
	if err != nil {
		return "", err
	}

	// add tokenSign
	if session = e.GetSession(id); session == nil {
		session = model.NewSession("0", "account-session", id)
	}
	session.AddDistinctValueTokenSign(&model.TokenSign{
		Value:  tokenValue,
		Device: device,
	})

	if e.config.DoubleToken {
		refreshToken, err := e.createRefreshToken(id, tokenValue, loginModel)
		if err != nil {
			return "", err
		}
		err = e.responseRefreshToken(refreshToken, loginModel, c[0])
		if err != nil {
			return "", err
		}
	}

	timeout := loginModel.Timeout
	// reset session
	err = e.SetSession(id, session, timeout)
	if err != nil {
		return "", err
	}

	// set token-id
	err = e.SetIdByToken(id, tokenValue, timeout)
	if err != nil {
		return "", err
	}

	// response token
	err = e.ResponseToken(tokenValue, loginModel, c[0])
	if err != nil {
		return "", err
	}

	// called watcher
	m := &model.Login{
		Device:          device,
		IsLastingCookie: loginModel.IsLastingCookie,
		Timeout:         timeout,
		JwtData:         loginModel.JwtData,
		Token:           tokenValue,
	}

	// called logger
	e.logger.Login(e.loginType, id, tokenValue, m)

	if e.watcher != nil {
		e.watcher.Login(e.loginType, id, tokenValue, m)
	}

	if device != "" && tokenConfig.DeviceMaxLoginCount != -1 {
		if session = e.GetSession(id); session != nil {
			// get by login device
			tokenSignList := session.GetFilterTokenSignSlice(device)
			if len(tokenSignList) > int(tokenConfig.DeviceMaxLoginCount) {
				err = e.deleteRedundantTokenSign(session, tokenConfig.DeviceMaxLoginCount)
				if err != nil {
					return "", err
				}
			}
		}
	}

	// check if the number of sessions for this account exceeds the maximum limit.
	if tokenConfig.MaxLoginCount != -1 {
		if session = e.GetSession(id); session != nil {
			if session.TokenSignSize() <= int(tokenConfig.MaxLoginCount) {
				return tokenValue, nil
			}
			err = e.deleteRedundantTokenSign(session, tokenConfig.MaxLoginCount)
			if err != nil {
				return "", err
			}
		}
	}

	return tokenValue, nil
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
func (e *Enforcer) LogoutById(id string, device ...string) error {
	session := e.GetSession(id)
	if session != nil {
		for _, tokenSign := range session.TokenSignList {
			if len(device) > 0 && device[0] != "" && tokenSign.Device == device[0] {
				err := e.LogoutByToken(tokenSign.Value)
				if err != nil {
					return err
				}
			} else {
				err := e.LogoutByToken(tokenSign.Value)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// LogoutByToken clear token info
func (e *Enforcer) LogoutByToken(token string) error {
	var err error
	// delete token-id
	id := e.getIdByToken(token)
	if id == "" {
		return errors.New("user not logged in")
	}
	// delete token-id
	err = e.deleteIdByToken(token)
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
	err = e.deleteRefreshToken(token)
	if err != nil {
		return err
	}

	e.logger.Logout(e.loginType, id, token)

	if e.watcher != nil {
		e.watcher.Logout(e.loginType, id, token)
	}

	return nil
}

// IsLoginById check if user logged in by loginId.
// check all tokenValue and if one is validated return true
func (e *Enforcer) IsLoginById(id string, device ...string) (bool, error) {
	var err error
	session := e.GetSession(id)
	if session != nil {
		var l []*model.TokenSign
		if len(device) > 0 && device[0] != "" {
			l = session.GetFilterTokenSignSlice(device[0])
		} else {
			l = session.TokenSignList
		}
		for _, tokenSign := range l {
			err = e.CheckLoginByToken(tokenSign.Value)
			if err != nil {
				continue
			}
			return true, nil
		}
	}

	return false, err
}

// GetId get the id from the Adapter, do not check the value
// 	if GetId()= -4, it means that user be replaced
//	if GetId()= -5, it means that user be kicked
//	if GetId()= -6, it means that user be banned
func (e *Enforcer) GetId(ctx ctx.Context) string {
	token := e.GetRequestToken(ctx)
	return e.GetIdByToken(token)
}

// GetIdByToken get the id from the Adapter
func (e *Enforcer) GetIdByToken(token string) string {
	if token == "" {
		return ""
	}
	loginId := e.getIdByToken(token)
	return loginId
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

	err := e.CheckLoginByToken(tokenValue)
	if err != nil {
		return false, err
	}

	return true, nil
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

// GetLoginId get id and check it
func (e *Enforcer) GetLoginId(ctx ctx.Context) (string, error) {
	tokenValue := e.GetRequestToken(ctx)
	return e.GetLoginIdByToken(tokenValue)
}

func (e *Enforcer) GetLoginIdByToken(token string) (string, error) {
	str := e.GetIdByToken(token)
	if str == "" {
		return "", errors.New("GetLoginId() failed: not logged in")
	}
	validate, err := e.checkId(str)
	if !validate {
		return "", err
	}
	// auto refresh timeout, When the user accesses
	if e.config.AutoRenew {
		_ = e.updateTokenTimeout(token, e.config.Timeout)
		_ = e.UpdateSessionTimeout(str, e.config.Timeout)
	}
	return str, nil
}

func (e *Enforcer) GetLoginCount(id string, device ...string) int {
	if session := e.GetSession(id); session != nil {
		if len(device) > 0 && device[0] != "" {
			return session.GetFilterTokenSign(device[0]).Len()
		}

		return session.TokenSignSize()
	}
	return 0
}

// GetBannedTime get banned time
func (e *Enforcer) GetBannedTime(id string, service string) int64 {
	timeout := e.getBannedTime(id, service)
	return timeout
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
func (e *Enforcer) AddTokenGenerateFun(tokenStyle string, f model.HandlerFunc) error {
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
	err := e.notifySet(e.spliceSessionKey(id), session, timeout)
	if err != nil {
		return err
	}
	return nil
}

func (e *Enforcer) DeleteSession(id string) error {
	err := e.notifyDelete(e.spliceSessionKey(id))
	if err != nil {
		return err
	}
	return nil
}

func (e *Enforcer) UpdateSession(id string, session *model.Session) error {
	err := e.notifyUpdate(e.spliceSessionKey(id), session)
	if err != nil {
		return err
	}
	return nil
}

func (e *Enforcer) UpdateSessionTimeout(id string, timeout int64) error {
	err := e.notifyUpdateTimeout(id, timeout)
	return err
}

func (e *Enforcer) GetTokenConfig() config.TokenConfig {
	return e.config
}

func (e *Enforcer) GetLoginCounts() (int, error) {
	adapter, ok := e.adapter.(persist.BatchAdapter)
	if !ok {
		return 0, fmt.Errorf("the adapter does not implement persist.BatchAdapter")
	}
	c, err := adapter.GetCountsFilteredKey(e.spliceSessionKey(""))
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (e *Enforcer) GetLoginTokenCounts() (int, error) {
	adapter, ok := e.adapter.(persist.BatchAdapter)
	if !ok {
		return 0, fmt.Errorf("the adapter does not implement persist.BatchAdapter")
	}
	c, err := adapter.GetCountsFilteredKey(e.spliceTokenKey(""))
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (e *Enforcer) GetRefreshToken(tokenValue string) string {
	return e.getRefreshTokenValue(tokenValue)
}

func (e *Enforcer) RefreshToken(refreshToken string, refreshModel ...*model.Refresh) (*model.RefreshRes, error) {
	var m *model.Refresh
	if len(refreshModel) != 0 {
		m = refreshModel[0]
	} else {
		m = model.DefaultRefresh()
	}
	return e.RefreshTokenByModel(refreshToken, m, nil)
}

func (e *Enforcer) RefreshTokenByModel(refreshToken string, refreshModel *model.Refresh, ctx ...ctx.Context) (*model.RefreshRes, error) {
	if refreshModel == nil {
		return nil, errors.New("arg refreshModel can not be nil")
	}
	if !e.config.DoubleToken {
		return nil, fmt.Errorf("double tokens are not enabled")
	}
	refreshTokenSign := e.getRefreshTokenSign(refreshToken)
	if refreshTokenSign == nil {
		return nil, fmt.Errorf("the refresh token does not exist: %v", refreshToken)
	}
	err := e.deleteRefreshToken(refreshTokenSign.Token)
	if err != nil {
		return nil, err
	}

	login := &model.Login{
		Device:              refreshTokenSign.Device,
		IsLastingCookie:     refreshModel.IsLastingCookie,
		Timeout:             refreshModel.Timeout,
		JwtData:             refreshModel.JwtData,
		Token:               refreshModel.Token,
		RefreshToken:        refreshModel.RefreshToken,
		RefreshTokenTimeout: refreshModel.RefreshTokenTimeout,
	}

	token, err := e.LoginByModel(refreshTokenSign.Id, login, ctx...)
	if err != nil {
		return nil, err
	}
	return &model.RefreshRes{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (e *Enforcer) GetLoginDevices(id string) []string {
	session := e.GetSession(id)
	if session == nil {
		return nil
	}
	return session.GetAllDevice()
}

func (e *Enforcer) GetDeviceByToken(token string) string {
	id := e.getIdByToken(token)
	session := e.GetSession(id)
	if session == nil {
		return ""
	}
	tokenSign := session.GetTokenSign(token)
	if tokenSign == nil {
		return ""
	}
	return tokenSign.Device
}
