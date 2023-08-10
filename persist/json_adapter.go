package persist

type JsonAdapter struct {
}

func NewJsonAdapter() *DefaultAdapter {
	d := NewDefaultAdapter()
	d.SetSerializer(NewJsonSerializer())
	return d
}
