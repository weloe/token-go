package persist

import (
	"encoding/json"
)

type JsonAdapter struct {
	*DefaultAdapter
}

func NewJsonAdapter() *JsonAdapter {
	return &JsonAdapter{NewDefaultAdapter()}
}

func (j *JsonAdapter) Serialize(data interface{}) ([]byte, error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return serializedData, nil
}

func (j *JsonAdapter) UnSerialize(data []byte, result interface{}) error {
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	return nil
}
