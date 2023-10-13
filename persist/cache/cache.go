package cache

type Cache interface {
	Get(key string) interface{}
	// Set store interface{}
	Set(key string, value interface{}, timeout int64) error
	// Update update value
	Update(key string, value interface{}) error
	// Delete delete value
	Delete(key string) error
	// GetTimeout get expire
	GetTimeout(key string) int64
	// UpdateTimeout update timeout
	UpdateTimeout(key string, timeout int64) error
}
