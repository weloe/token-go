package persist

import "github.com/weloe/token-go/model"

type SerializerAdapter interface {
	Adapter
	Serialize(*model.Session) ([]byte, error)
	UnSerialize([]byte) (*model.Session, error)
}
