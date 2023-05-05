package token_go

import (
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/log"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/persist"
)

var _ IEnforcer = &Enforcer{}

type IEnforcer interface {
	Login(id string) (string, error)
	LoginByModel(id string, loginModel *model.Login) (string, error)
	Logout() error
	IsLogin() (bool, error)
	IsLoginById(id string) (bool, error)
	GetLoginId() (string, error)

	Replaced(id string, device string) error
	// Banned TODO
	Banned(id string, service string) error
	Kickout(id string, device string) error

	GetRequestToken() string

	SetType(t string)
	GetType() string
	SetContext(ctx ctx.Context)
	GetAdapter() persist.Adapter
	SetAdapter(adapter persist.Adapter)
	SetWatcher(watcher persist.Watcher)
	SetLogger(logger log.Logger)
	EnableLog()
	IsLogEnable() bool
	GetSession(id string) *model.Session
	SetSession(id string, session *model.Session, timeout int64) error
}
