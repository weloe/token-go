package ctx

// ReqStorage
// Use to set-get-delete data,and it'll be cleaned after request
type ReqStorage interface {
	Source() interface{}
	Get(key string) interface{}
	Set(key string, value string)
	Delete(key string)
}
