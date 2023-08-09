package persist

import "encoding/json"

type JsonSerializer struct {
}

func NewJsonSerializer() *JsonSerializer {
	return &JsonSerializer{}
}

func (j *JsonSerializer) Serialize(data interface{}) ([]byte, error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return serializedData, nil
}

func (j *JsonSerializer) UnSerialize(data []byte, result interface{}) error {
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	return nil
}
