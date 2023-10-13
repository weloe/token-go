package cache

type CacheEx interface {
	Cache
	Range(func(key, value interface{}) bool)
	EnableCleanTimer(period int64)
}
