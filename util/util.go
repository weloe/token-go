package util

import (
	"fmt"
	"reflect"
)

func GetType(i any) reflect.Type {
	return reflect.TypeOf(i)
}

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

// AppendStr do not add repeated str.
// If old slice has newStr, return directly, else append
func AppendStr(old []string, newStr string) []string {
	if HasStr(old, newStr) {
		return old
	}
	return append(old, newStr)
}
