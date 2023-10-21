package persist

import (
	"errors"
	"fmt"
	"github.com/weloe/token-go/persist/cache"
	"github.com/weloe/token-go/util"
	"log"
	"reflect"
	"strings"
)

type DefaultAdapter struct {
	cache              cache.Cache
	serializer         Serializer
	enableRefreshTimer bool
}

var _ Adapter = (*DefaultAdapter)(nil)

func NewDefaultAdapter() *DefaultAdapter {
	return &DefaultAdapter{
		cache:              cache.NewDefaultLocalCache(),
		enableRefreshTimer: true,
	}
}

func (d *DefaultAdapter) SetSerializer(serializer Serializer) {
	d.serializer = serializer
}

// GetStr if key is expired delete it before get data
func (d *DefaultAdapter) GetStr(key string) string {
	value := d.cache.Get(key)
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func (d *DefaultAdapter) SetStr(key string, value string, timeout int64) error {
	return d.cache.Set(key, value, timeout)
}

func (d *DefaultAdapter) UpdateStr(key string, value string) error {
	return d.cache.Update(key, value)
}

func (d *DefaultAdapter) DeleteStr(key string) error {
	return d.cache.Delete(key)
}

func (d *DefaultAdapter) GetStrTimeout(key string) int64 {
	return d.cache.GetTimeout(key)
}

func (d *DefaultAdapter) UpdateStrTimeout(key string, timeout int64) error {
	return d.cache.UpdateTimeout(key, timeout)
}

// interface{} operation
//
//

func (d *DefaultAdapter) Get(key string, t ...reflect.Type) interface{} {
	value := d.cache.Get(key)

	if d.serializer == nil || t == nil || len(t) == 0 {
		return value
	}
	bytes, err := util.InterfaceToBytes(value)
	if err != nil {
		log.Printf("Adapter.Get() failed: %v", err)
		return nil
	}
	instance := reflect.New(t[0].Elem()).Interface()
	err = d.serializer.UnSerialize(bytes, instance)
	if err != nil {
		log.Printf("Adapter.Get() failed: %v", err)
		return nil
	}

	return instance
}

func (d *DefaultAdapter) Set(key string, value interface{}, timeout int64) error {
	if d.serializer == nil {
		return d.cache.Set(key, value, timeout)
	}
	bytes, err := d.serializer.Serialize(value)
	if err != nil {
		return err
	}
	return d.cache.Set(key, bytes, timeout)
}

func (d *DefaultAdapter) Update(key string, value interface{}) error {
	if d.serializer == nil {
		return d.cache.Update(key, value)
	}
	bytes, err := d.serializer.Serialize(value)
	if err != nil {
		return err
	}
	return d.cache.Update(key, bytes)
}

func (d *DefaultAdapter) GetTimeout(key string) int64 {
	return d.cache.GetTimeout(key)
}

func (d *DefaultAdapter) UpdateTimeout(key string, timeout int64) error {
	return d.cache.UpdateTimeout(key, timeout)
}

func (d *DefaultAdapter) Delete(key string) error {
	return d.cache.Delete(key)
}

func (d *DefaultAdapter) DeleteBatchFilteredKey(keyPrefix string) error {
	var err error
	cacheEx, ok := d.cache.(cache.CacheEx)
	if !ok {
		return errors.New("the cache does not implement the Range method")
	}
	cacheEx.Range(func(key, value any) bool {
		if strings.HasPrefix(key.(string), keyPrefix) {
			err = d.cache.Delete(key.(string))
			if err != nil {
				return false
			}
		}
		return true
	})
	return err
}

func (d *DefaultAdapter) GetCountsFilteredKey(keyPrefix string) (int, error) {
	cacheEx, ok := d.cache.(cache.CacheEx)
	if !ok {
		return 0, errors.New("the cache does not implement the Range method")
	}
	var counts int
	cacheEx.Range(func(key, value any) bool {
		if strings.HasPrefix(key.(string), keyPrefix) {
			counts++
		}
		return true
	})
	return counts, nil
}

func (d *DefaultAdapter) EnableCleanTimer(b bool) {
	d.enableRefreshTimer = b
}

func (d *DefaultAdapter) GetCleanTimer() bool {
	return d.enableRefreshTimer
}

func (d *DefaultAdapter) StartCleanTimer(period int64) error {
	cacheEx, ok := d.cache.(cache.CacheEx)
	if !ok {
		return errors.New("the Cache does not implement the StartCleanTimer method")
	}
	cacheEx.EnableCleanTimer(period)
	return nil
}
