package token_go

import (
	"github.com/weloe/token-go/config"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/log"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/persist"
)

var _ IEnforcer = &Enforcer{}

type IEnforcer interface {
	// Enforcer field api
	SetType(t string)
	GetType() string
	GetAdapter() persist.Adapter
	SetAdapter(adapter persist.Adapter)
	SetWatcher(watcher persist.Watcher)
	GetWatcher() persist.Watcher
	SetLogger(logger log.Logger)
	GetLogger() log.Logger
	EnableLog()
	IsLogEnable() bool
	GetTokenConfig() config.TokenConfig

	// Login login api
	Login(id string, ctx ...ctx.Context) (string, error)
	LoginById(id string, device ...string) (string, error)
	LoginByModel(id string, loginModel *model.Login, ctx ...ctx.Context) (string, error)

	Logout(ctx ctx.Context) error
	LogoutById(id string, device ...string) error
	LogoutByToken(token string) error

	IsLogin(ctx ctx.Context) (bool, error)
	IsLoginByToken(token string) (bool, error)
	IsLoginById(id string, device ...string) (bool, error)
	CheckLogin(ctx ctx.Context) error
	CheckLoginByToken(token string) error

	GetLoginId(ctx ctx.Context) (string, error)
	GetLoginIdByToken(token string) (string, error)
	GetId(ctx ctx.Context) string
	GetIdByToken(token string) string
	GetLoginCount(id string, device ...string) int

	// device manager api
	GetLoginDevices(id string) []string
	GetDeviceByToken(token string) string

	// refresh api
	GetRefreshToken(tokenValue string) string
	RefreshToken(refreshToken string, refreshModel ...*model.Refresh) (*model.RefreshRes, error)
	RefreshTokenByModel(refreshToken string, refreshModel *model.Refresh, ctx ...ctx.Context) (*model.RefreshRes, error)

	GetLoginCounts() (int, error)
	GetLoginTokenCounts() (int, error)

	Kickout(id string, device ...string) error
	Replaced(id string, device ...string) error

	// Banned banned api
	Banned(id string, service string, level int, time int64) error
	UnBanned(id string, services ...string) error
	IsBanned(id string, service string) bool
	GetBannedLevel(id string, service string) (int64, error)
	GetBannedTime(id string, service string) int64

	// Second auth api
	OpenSafe(token string, service string, time int64) error
	IsSafe(token string, service string) bool
	GetSafeTime(token string, service string) int64
	CloseSafe(token string, service string) error

	// Temp token api
	CreateTempToken(token string, service string, value string, timeout int64) (string, error)
	CreateTempTokenByStyle(style string, service string, value string, timeout int64) (string, error)
	GetTempTokenTimeout(service string, tempToken string) int64
	ParseTempToken(service string, tempToken string) string
	DeleteTempToken(service string, tempToken string) error

	GetRequestToken(ctx ctx.Context) string
	AddTokenGenerateFun(tokenStyle string, f model.HandlerFunc) error

	// QRCode api
	CreateQRCodeState(QRCodeId string, timeout int64) error
	GetQRCode(QRCodeId string) *model.QRCode
	GetQRCodeState(QRCodeId string) model.QRCodeState
	GetQRCodeTimeout(QRCodeId string) int64
	DeleteQRCode(QRCodeId string) error
	Scanned(QRCodeId string, loginId string) (string, error)
	ConfirmAuth(QRCodeTempToken string) error
	CancelAuth(QRCodeTempToken string) error

	// Access control api
	SetAuth(manager interface{})
	CheckRole(ctx ctx.Context, role string) error
	CheckPermission(ctx ctx.Context, permission string) error

	// Session api
	GetSession(id string) *model.Session
	DeleteSession(id string) error
	UpdateSession(id string, session *model.Session) error
	SetSession(id string, session *model.Session, timeout int64) error
}

var _ IDistributedEnforcer = &DistributedEnforcer{}

type IDistributedEnforcer interface {
	IEnforcer
	// SetStrSelf store string in all instances
	SetStrSelf(key string, value string, timeout int64) error
	// UpdateStrSelf only update string value in all instances
	UpdateStrSelf(key string, value string) error
	// SetSelf store interface{} in all instances
	SetSelf(key string, value interface{}, timeout int64) error
	// UpdateSelf only update interface{} value in all instances
	UpdateSelf(key string, value interface{}) error
	// DeleteSelf delete interface{} value in all instances
	DeleteSelf(key string) error
	// UpdateTimeoutSelf update timeout in all instances
	UpdateTimeoutSelf(key string, timeout int64) error
}
