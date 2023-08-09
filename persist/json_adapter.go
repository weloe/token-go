package persist

type JsonAdapter struct {
	*DefaultAdapter
	*JsonSerializer
}

func NewJsonAdapter() *JsonAdapter {
	return &JsonAdapter{NewDefaultAdapter(), NewJsonSerializer()}
}
