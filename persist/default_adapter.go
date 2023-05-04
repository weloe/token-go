package persist

import (
	"errors"
	"fmt"
	"github.com/weloe/token-go/constant"
	"strings"
	"sync"
	"time"
)

type DefaultAdapter struct {
	dataMap   *sync.Map
	expireMap *sync.Map
}

var _ Adapter = (*DefaultAdapter)(nil)

func NewDefaultAdapter() *DefaultAdapter {
	return &DefaultAdapter{
		dataMap:   &sync.Map{},
		expireMap: &sync.Map{},
	}
}

// GetStr if key is expired delete it before get data
func (d *DefaultAdapter) GetStr(key string) string {
	_ = d.getExpireAndDelete(key)
	value, _ := d.dataMap.Load(key)
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func (d *DefaultAdapter) SetStr(key string, value string, timeout int64) error {
	if timeout == 0 || timeout <= constant.NotValueExpire {
		return errors.New("args timeout error")
	}
	d.dataMap.Store(key, value)

	if timeout == constant.NeverExpire {
		d.expireMap.Store(key, constant.NeverExpire)
	} else {
		d.expireMap.Store(key, time.Now().UnixMilli()+timeout*1000)
	}
	return nil
}

func (d *DefaultAdapter) UpdateStr(key string, value string) error {
	timeout := d.GetStrTimeout(key)
	if timeout == constant.NotValueExpire {
		return errors.New("does not exist")
	}
	d.dataMap.Store(key, value)
	return nil
}

func (d *DefaultAdapter) DeleteStr(key string) error {
	d.dataMap.Delete(key)
	d.expireMap.Delete(key)
	return nil
}

func (d *DefaultAdapter) GetStrTimeout(key string) int64 {
	return d.getTimeout(key)
}

func (d *DefaultAdapter) UpdateStrTimeout(key string, timeout int64) error {
	if timeout == constant.NeverExpire {
		d.expireMap.Store(key, constant.NeverExpire)
	} else {
		d.expireMap.Store(key, time.Now().UnixMilli()+timeout*1000)
	}
	return nil
}

// interface{} operation
//
//

func (d *DefaultAdapter) Get(key string) interface{} {
	d.getExpireAndDelete(key)
	value, _ := d.dataMap.Load(key)
	return value
}

func (d *DefaultAdapter) Set(key string, value interface{}, timeout int64) error {
	if timeout == 0 || timeout <= constant.NotValueExpire {
		return errors.New("args timeout error")
	}
	d.dataMap.Store(key, value)

	if timeout == constant.NeverExpire {
		d.expireMap.Store(key, constant.NeverExpire)
	} else {
		d.expireMap.Store(key, time.Now().UnixMilli()+timeout*1000)
	}
	return nil
}

func (d *DefaultAdapter) Update(key string, value interface{}) error {
	timeout := d.GetStrTimeout(key)
	if timeout == constant.NotValueExpire {
		return errors.New("key does not exist")
	}
	d.dataMap.Store(key, value)
	return nil
}

func (d *DefaultAdapter) GetTimeout(key string) int64 {
	return d.getTimeout(key)
}

func (d *DefaultAdapter) UpdateTimeout(key string, timeout int64) error {
	if timeout == constant.NeverExpire {
		d.expireMap.Store(key, constant.NeverExpire)
	} else {
		d.expireMap.Store(key, time.Now().UnixMilli()+timeout*1000)
	}
	return nil
}

func (d *DefaultAdapter) Delete(key string) error {
	d.dataMap.Delete(key)
	d.expireMap.Delete(key)
	return nil
}

func (d *DefaultAdapter) DeleteBatchFilteredKey(keyPrefix string) error {
	d.dataMap.Range(func(key, value any) bool {
		if strings.HasPrefix(key.(string), keyPrefix) {
			d.dataMap.Delete(key)
		}
		return true
	})
	return nil
}

// delete key when getValue is expired
func (d *DefaultAdapter) getExpireAndDelete(key string) int64 {
	expirationTime, _ := d.expireMap.Load(key)

	if expirationTime == nil {
		return 0
	}

	if expirationTime.(int64) != constant.NeverExpire && expirationTime.(int64) <= time.Now().UnixMilli() {
		d.dataMap.Delete(key)
		d.expireMap.Delete(key)
	}
	return expirationTime.(int64)
}

func (d *DefaultAdapter) getTimeout(key string) int64 {
	expirationTime := d.getExpireAndDelete(key)
	if expirationTime == 0 {
		return constant.NotValueExpire
	}
	if expirationTime == constant.NeverExpire {
		return constant.NeverExpire
	}
	timeout := (expirationTime - time.Now().UnixMilli()) / 1000
	if timeout <= 0 {
		d.dataMap.Delete(key)
		d.expireMap.Delete(key)
		return constant.NotValueExpire
	}
	return timeout
}
