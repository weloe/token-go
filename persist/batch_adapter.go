package persist

type BatchAdapter interface {
	Adapter
	// DeleteBatchFilteredKey delete data by keyPrefix
	DeleteBatchFilteredKey(filterKeyPrefix string) error
}
