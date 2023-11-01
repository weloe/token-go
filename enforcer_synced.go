package token_go

import (
	"github.com/weloe/token-go/config"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/log"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/persist"
	"sync"
)

// SyncedEnforcer wraps Enforcer and provides synchronized access
type SyncedEnforcer struct {
	*Enforcer
	m sync.RWMutex
}

// NewSyncedEnforcer creates a synchronized enforcer
func NewSyncedEnforcer(adapter persist.Adapter, params ...interface{}) (*SyncedEnforcer, error) {
	e := &SyncedEnforcer{}
	var err error
	e.Enforcer, err = NewEnforcer(adapter, params...)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// GetLock return the private RWMutex lock
func (e *SyncedEnforcer) GetLock() *sync.RWMutex {
	return &e.m
}

func (e *SyncedEnforcer) SetType(t string) {
	e.m.Lock()
	defer e.m.Unlock()
	e.Enforcer.SetType(t)
}

func (e *SyncedEnforcer) GetType() string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetType()
}

func (e *SyncedEnforcer) GetAdapter() persist.Adapter {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetAdapter()
}

func (e *SyncedEnforcer) SetAdapter(adapter persist.Adapter) {
	e.m.Lock()
	defer e.m.Unlock()
	e.Enforcer.SetAdapter(adapter)
}

func (e *SyncedEnforcer) SetWatcher(watcher persist.Watcher) {
	e.m.Lock()
	defer e.m.Unlock()
	e.Enforcer.SetWatcher(watcher)
}

func (e *SyncedEnforcer) GetWatcher() persist.Watcher {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetWatcher()
}

func (e *SyncedEnforcer) SetLogger(logger log.Logger) {
	e.m.Lock()
	defer e.m.Unlock()
	e.Enforcer.SetLogger(logger)
}

func (e *SyncedEnforcer) GetLogger() log.Logger {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetLogger()
}

func (e *SyncedEnforcer) EnableLog() {
	e.m.Lock()
	defer e.m.Unlock()
	e.Enforcer.EnableLog()
}

func (e *SyncedEnforcer) IsLogEnable() bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.IsLogEnable()
}

func (e *SyncedEnforcer) GetTokenConfig() config.TokenConfig {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetTokenConfig()
}

func (e *SyncedEnforcer) Login(id string, ctx ctx.Context) (string, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.Login(id, ctx)
}

func (e *SyncedEnforcer) LoginById(id string, device ...string) (string, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.LoginById(id, device...)
}

func (e *SyncedEnforcer) LoginByModel(id string, loginModel *model.Login, ctx ctx.Context) (string, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.LoginByModel(id, loginModel, ctx)
}

func (e *SyncedEnforcer) Logout(ctx ctx.Context) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.Logout(ctx)
}

func (e *SyncedEnforcer) LogoutById(id string, device ...string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.LogoutById(id, device...)
}

func (e *SyncedEnforcer) LogoutByToken(token string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.LogoutByToken(token)
}

func (e *SyncedEnforcer) IsLogin(ctx ctx.Context) (bool, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.IsLogin(ctx)
}

func (e *SyncedEnforcer) IsLoginByToken(token string) (bool, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.IsLoginByToken(token)
}

func (e *SyncedEnforcer) IsLoginById(id string, device ...string) (bool, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.IsLoginById(id, device...)
}

func (e *SyncedEnforcer) CheckLogin(ctx ctx.Context) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.CheckLogin(ctx)
}

func (e *SyncedEnforcer) CheckLoginByToken(token string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.CheckLoginByToken(token)
}

func (e *SyncedEnforcer) GetLoginId(ctx ctx.Context) (string, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetLoginId(ctx)
}

func (e *SyncedEnforcer) GetLoginIdByToken(token string) (string, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetLoginIdByToken(token)
}

func (e *SyncedEnforcer) GetId(ctx ctx.Context) string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetId(ctx)
}

func (e *SyncedEnforcer) GetIdByToken(token string) string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetIdByToken(token)
}

func (e *SyncedEnforcer) GetLoginCount(id string, device ...string) int {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetLoginCount(id, device...)
}

func (e *SyncedEnforcer) GetRefreshToken(tokenValue string) string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetRefreshToken(tokenValue)
}

func (e *SyncedEnforcer) RefreshToken(refreshToken string, refreshModel ...*model.Refresh) (*model.RefreshRes, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.RefreshToken(refreshToken, refreshModel...)
}

func (e *SyncedEnforcer) RefreshTokenByModel(refreshToken string, refreshModel *model.Refresh, ctx ctx.Context) (*model.RefreshRes, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.RefreshTokenByModel(refreshToken, refreshModel, ctx)
}

func (e *SyncedEnforcer) GetLoginCounts() (int, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetLoginCounts()
}

func (e *SyncedEnforcer) GetLoginTokenCounts() (int, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetLoginTokenCounts()
}

func (e *SyncedEnforcer) Kickout(id string, device ...string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.Kickout(id, device...)
}

func (e *SyncedEnforcer) Replaced(id string, device ...string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.Replaced(id, device...)
}

func (e *SyncedEnforcer) Banned(id string, service string, level int, time int64) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.Banned(id, service, level, time)
}

func (e *SyncedEnforcer) UnBanned(id string, services ...string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.UnBanned(id, services...)
}

func (e *SyncedEnforcer) IsBanned(id string, service string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.IsBanned(id, service)
}

func (e *SyncedEnforcer) GetBannedLevel(id string, service string) (int64, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetBannedLevel(id, service)
}

func (e *SyncedEnforcer) GetBannedTime(id string, service string) int64 {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.getBannedTime(id, service)
}

func (e *SyncedEnforcer) OpenSafe(token string, service string, time int64) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.OpenSafe(token, service, time)
}

func (e *SyncedEnforcer) IsSafe(token string, service string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.IsSafe(token, service)
}

func (e *SyncedEnforcer) GetSafeTime(token string, service string) int64 {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetSafeTime(token, service)
}

func (e *SyncedEnforcer) CloseSafe(token string, service string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.CloseSafe(token, service)
}

func (e *SyncedEnforcer) CreateTempToken(token string, service string, value string, timeout int64) (string, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.CreateTempToken(token, service, value, timeout)
}

func (e *SyncedEnforcer) CreateTempTokenByStyle(style string, service string, value string, timeout int64) (string, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.CreateTempTokenByStyle(style, service, value, timeout)
}

func (e *SyncedEnforcer) GetTempTokenTimeout(service string, tempToken string) int64 {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetTempTokenTimeout(service, tempToken)
}

func (e *SyncedEnforcer) ParseTempToken(service string, tempToken string) string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.ParseTempToken(service, tempToken)
}

func (e *SyncedEnforcer) DeleteTempToken(service string, tempToken string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.DeleteTempToken(service, tempToken)
}

func (e *SyncedEnforcer) GetRequestToken(ctx ctx.Context) string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetRequestToken(ctx)
}

func (e *SyncedEnforcer) AddTokenGenerateFun(tokenStyle string, f model.HandlerFunc) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.AddTokenGenerateFun(tokenStyle, f)
}

func (e *SyncedEnforcer) CreateQRCodeState(QRCodeId string, timeout int64) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.CreateQRCodeState(QRCodeId, timeout)
}

func (e *SyncedEnforcer) GetQRCode(QRCodeId string) *model.QRCode {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetQRCode(QRCodeId)
}

func (e *SyncedEnforcer) GetQRCodeState(QRCodeId string) model.QRCodeState {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetQRCodeState(QRCodeId)
}

func (e *SyncedEnforcer) GetQRCodeTimeout(QRCodeId string) int64 {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetQRCodeTimeout(QRCodeId)
}

func (e *SyncedEnforcer) DeleteQRCode(QRCodeId string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.DeleteQRCode(QRCodeId)
}

func (e *SyncedEnforcer) Scanned(QRCodeId string, loginId string) (string, error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.Scanned(QRCodeId, loginId)
}

func (e *SyncedEnforcer) ConfirmAuth(QRCodeTempToken string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.ConfirmAuth(QRCodeTempToken)
}

func (e *SyncedEnforcer) CancelAuth(QRCodeTempToken string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.CancelAuth(QRCodeTempToken)
}

func (e *SyncedEnforcer) SetAuth(manager interface{}) {
	e.m.Lock()
	defer e.m.Unlock()
	e.Enforcer.SetAuth(manager)
}

func (e *SyncedEnforcer) CheckRole(ctx ctx.Context, role string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.CheckRole(ctx, role)
}

func (e *SyncedEnforcer) CheckPermission(ctx ctx.Context, permission string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.CheckPermission(ctx, permission)
}

func (e *SyncedEnforcer) GetSession(id string) *model.Session {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetSession(id)
}

func (e *SyncedEnforcer) DeleteSession(id string) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.DeleteSession(id)
}

func (e *SyncedEnforcer) UpdateSession(id string, session *model.Session) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.UpdateSession(id, session)
}

func (e *SyncedEnforcer) SetSession(id string, session *model.Session, timeout int64) error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.SetSession(id, session, timeout)
}
