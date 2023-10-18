package persist

type Dispatcher interface {

	// SetAllStr store string in all instances
	SetAllStr(key string, value string, timeout int64) error
	// UpdateAllStr only update string value in all instances
	UpdateAllStr(key string, value string) error

	// SetAll store interface{} in all instances
	SetAll(key string, value interface{}, timeout int64) error
	// UpdateAll only update interface{} value in all instances
	UpdateAll(key string, value interface{}) error

	// DeleteAll delete interface{} value in all instances
	DeleteAll(key string) error
	// UpdateAllTimeout update timeout in all instances
	UpdateAllTimeout(key string, timeout int64) error
}
