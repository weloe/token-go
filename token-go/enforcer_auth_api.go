package token_go

import (
	"errors"
	"fmt"
	"github.com/weloe/token-go/auth"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/util"
)

func (e *Enforcer) SetAuth(manager interface{}) {
	e.authManager = manager
}

func (e *Enforcer) CheckRole(ctx ctx.Context, role string) error {
	if e.authManager == nil {
		return errors.New("authManager is nil")
	}
	rbac, ok := e.authManager.(auth.RBAC)
	if !ok {
		return errors.New("authManager doesn't implement RBAC interface")
	}
	id, err := e.GetLoginId(ctx)
	if err != nil {
		return err
	}
	roles := rbac.GetRole(id)
	if util.HasStr(roles, role) {
		return nil
	}
	return fmt.Errorf("id %v doesn't has role %v", id, role)
}

func (e *Enforcer) CheckPermission(ctx ctx.Context, permission string) error {
	if e.authManager == nil {
		return errors.New("authManager is nil")
	}
	acl, ok := e.authManager.(auth.ACL)
	if !ok {
		return errors.New("authManager doesn't implement ACL interface")
	}
	id, err := e.GetLoginId(ctx)
	if err != nil {
		return err
	}
	permissions := acl.GetPermission(id)
	if util.HasStr(permissions, permission) {
		return nil
	}
	return fmt.Errorf("id %v doesn't has permission %v", id, permission)
}
