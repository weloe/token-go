package cache

import (
	"errors"
	"fmt"
	"github.com/weloe/token-go/constant"
	"sync"
	"time"
)

var _ Cache = (*DefaultLocalCache)(nil)

type DefaultLocalCache struct {
	dataMap   *sync.Map
	expireMap *sync.Map
	once      sync.Once
}

func NewDefaultLocalCache() *DefaultLocalCache {
	d := &DefaultLocalCache{dataMap: &sync.Map{}, expireMap: &sync.Map{}}
	return d
}

func (c *DefaultLocalCache) Get(key string) interface{} {
	_ = c.getExpireAndDelete(key)
	value, _ := c.dataMap.Load(key)
	return value
}

func (c *DefaultLocalCache) Set(key string, value interface{}, timeout int64) error {
	err := c.assertTimeout(timeout)
	if err != nil {
		return err
	}
	if timeout <= constant.NotValueExpire {
		timeout = constant.NotValueExpire
	}
	c.dataMap.Store(key, value)

	if timeout == constant.NeverExpire {
		c.expireMap.Store(key, constant.NeverExpire)
	} else {
		c.expireMap.Store(key, time.Now().UnixMilli()+timeout*1000)
	}
	return nil
}

func (c *DefaultLocalCache) assertTimeout(timeout int64) error {
	if timeout == 0 || timeout == -2 {
		return fmt.Errorf("timeout cannot be equal to %v", timeout)
	}
	return nil
}

func (c *DefaultLocalCache) Update(key string, value interface{}) error {
	timeout := c.GetTimeout(key)
	if timeout == constant.NotValueExpire {
		return errors.New("key does not exist")
	}
	c.dataMap.Store(key, value)
	return nil
}

func (c *DefaultLocalCache) Delete(key string) error {
	c.dataMap.Delete(key)
	c.expireMap.Delete(key)
	return nil
}

func (c *DefaultLocalCache) GetTimeout(key string) int64 {
	expirationTime := c.getExpireAndDelete(key)
	if expirationTime == 0 {
		return constant.NotValueExpire
	}
	if expirationTime == constant.NeverExpire {
		return constant.NeverExpire
	}
	timeout := (expirationTime - time.Now().UnixMilli()) / 1000
	if timeout <= 0 {
		c.dataMap.Delete(key)
		c.expireMap.Delete(key)
		return constant.NotValueExpire
	}
	return timeout
}

func (c *DefaultLocalCache) UpdateTimeout(key string, timeout int64) error {
	err := c.assertTimeout(timeout)
	if err != nil {
		return err
	}
	if timeout <= constant.NeverExpire {
		c.expireMap.Store(key, constant.NeverExpire)
	} else {
		c.expireMap.Store(key, time.Now().UnixMilli()+timeout*1000)
	}
	return nil
}

func (c *DefaultLocalCache) Range(f func(key, value any) bool) {
	c.dataMap.Range(f)
}

func (c *DefaultLocalCache) EnableCleanTimer(period int64) {
	c.once.Do(func() {
		go c.cleanTask(period)
	})
}

func (c *DefaultLocalCache) cleanTask(period int64) {
	if period < 0 {
		return
	}
	duration := period

	// create timer
	ticker := time.NewTicker(time.Duration(duration) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.expireMap.Range(func(key, value any) bool {
			_ = c.getExpireAndDelete(key.(string))
			return true
		})
	}
}

// getExpireAndDelete delete key when getValue is expired
func (c *DefaultLocalCache) getExpireAndDelete(key string) int64 {
	expirationTime, _ := c.expireMap.Load(key)

	if expirationTime == nil {
		return 0
	}

	if expirationTime.(int64) != constant.NeverExpire && expirationTime.(int64) <= time.Now().UnixMilli() {
		c.dataMap.Delete(key)
		c.expireMap.Delete(key)
	}
	return expirationTime.(int64)
}
