package util

import "fmt"

func HasNil(arr []interface{}) bool {
	for _, elem := range arr {
		if elem == nil {
			return true
		}
	}
	return false
}

func HasStr(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func InterfaceToBytes(data interface{}) ([]byte, error) {
	if b, ok := data.([]byte); ok {
		return b, nil
	}
	return nil, fmt.Errorf("unable to convert %T to []byte", data)
}
