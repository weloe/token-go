package persist

type Serializer interface {
	Serialize(data interface{}) ([]byte, error)
	UnSerialize([]byte, interface{}) error
}
