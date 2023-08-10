package persist

import "reflect"

type Adapter interface {

	// GetStr string operate string value
	GetStr(key string) string
	// SetStr set store value and timeout
	SetStr(key string, value string, timeout int64) error
	// UpdateStr only update value
	UpdateStr(key string, value string) error
	// DeleteStr delete string value
	DeleteStr(key string) error
	// GetStrTimeout get expire
	GetStrTimeout(key string) int64
	// UpdateStrTimeout update expire time
	UpdateStrTimeout(key string, timeout int64) error

	// Get returns interface{}
	// If serializer != nil, need to input reflect.Type, used to serializer to deserialize,
	// if ( serializer == nil || t == nil || len(t) == 0 ), returns value directly.
	Get(key string, t ...reflect.Type) interface{}
	// Set store interface{}
	Set(key string, value interface{}, timeout int64) error
	// Update only update interface{} value
	Update(key string, value interface{}) error
	// Delete delete interface{} value
	Delete(key string) error
	// GetTimeout get expire
	GetTimeout(key string) int64
	// UpdateTimeout update timeout
	UpdateTimeout(key string, timeout int64) error

	// SetSerializer used to serialize and deserialize
	// Serialize when call Set() or Update(), deserialize when call Get(key,t)
	SetSerializer(serializer Serializer)
}
