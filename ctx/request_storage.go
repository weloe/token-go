package ctx

type StorageKey string

// ReqStorage
// Use to set-get-delete data,and it'll be cleaned after request
type ReqStorage interface {
	Source() interface{}
	Get(key StorageKey) interface{}
	Set(key StorageKey, value string)
	Delete(key StorageKey)
}
