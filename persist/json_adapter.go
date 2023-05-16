package persist

import (
	"encoding/json"
	"github.com/weloe/token-go/model"
)

type JsonAdapter struct {
	*DefaultAdapter
}

func NewJsonAdapter() *JsonAdapter {
	return &JsonAdapter{NewDefaultAdapter()}
}

func (j *JsonAdapter) Serialize(session *model.Session) ([]byte, error) {
	return json.Marshal(session)
}

func (j *JsonAdapter) UnSerialize(bytes []byte) (*model.Session, error) {
	s := &model.Session{}
	err := json.Unmarshal(bytes, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
