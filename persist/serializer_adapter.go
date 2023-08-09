package persist

type SerializerAdapter interface {
	Adapter
	Serialize(data interface{}) ([]byte, error)
	UnSerialize([]byte, interface{}) error
}
