package persist

var _ Adapter = (*EmptyAdapter)(nil)

// EmptyAdapter empty adapter for extension to init enforcer
type EmptyAdapter struct {
}

func (e *EmptyAdapter) GetStr(key string) string {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) SetStr(key string, value string, timeout int64) error {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) UpdateStr(key string, value string) error {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) DeleteStr(key string) error {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) GetStrTimeout(key string) int64 {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) UpdateStrTimeout(key string, timeout int64) error {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) Get(key string) interface{} {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) Set(key string, value interface{}, timeout int64) error {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) Update(key string, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) Delete(key string) error {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) GetTimeout(key string) int64 {
	//TODO implement me
	panic("implement me")
}

func (e *EmptyAdapter) UpdateTimeout(key string, timeout int64) error {
	//TODO implement me
	panic("implement me")
}

