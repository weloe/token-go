package persist

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

	// Get get interface{}
	Get(key string) interface{}
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

	// DeleteBatchFilteredKey delete data by keyPrefix
	DeleteBatchFilteredKey(filterKeyPrefix string) error
}
