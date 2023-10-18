package persist

// UpdatableWatcher called when data updated
type UpdatableWatcher interface {
	UpdateForSetStr(key string, value interface{}, timeout int64) error
	UpdateForUpdateStr(key string, value interface{}) error
	UpdateForSet(key string, value interface{}, timeout int64) error
	UpdateForUpdate(key string, value interface{}) error
	UpdateForDelete(key string) error
	UpdateForUpdateTimeout(key string, timeout int64) error
}
