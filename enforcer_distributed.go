package token_go

import "github.com/weloe/token-go/persist"

type DistributedEnforcer struct {
	*Enforcer
}

func NewDistributedEnforcer(enforcer *Enforcer) *DistributedEnforcer {
	return &DistributedEnforcer{enforcer}
}

func (e *DistributedEnforcer) SetStrSelf(key string, value string, timeout int64) error {
	return e.adapter.SetStr(key, value, timeout)
}

func (e *DistributedEnforcer) UpdateStrSelf(key string, value string) error {
	return e.adapter.UpdateStr(key, value)
}

func (e *DistributedEnforcer) SetSelf(key string, value interface{}, timeout int64) error {
	return e.adapter.Set(key, value, timeout)
}

func (e *DistributedEnforcer) UpdateSelf(key string, value interface{}) error {
	return e.adapter.Update(key, value)
}

func (e *DistributedEnforcer) DeleteSelf(key string) error {
	return e.adapter.DeleteStr(key)
}

func (e *DistributedEnforcer) UpdateTimeoutSelf(key string, timeout int64) error {
	return e.adapter.UpdateTimeout(key, timeout)
}

func (e *DistributedEnforcer) EnableDispatcher(b bool) {
	if e.dispatcher == nil {
		return
	}
	e.notifyDispatcher = b
}

func (e *Enforcer) SetDispatcher(dispatcher persist.Dispatcher) {
	e.dispatcher = dispatcher
}
