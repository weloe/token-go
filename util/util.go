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
	switch v := data.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return []byte(fmt.Sprintf("%d", v)), nil
	default:
		return nil, fmt.Errorf("unable to convert %T to []byte", data)
	}
}

// AppendStr do not add repeated str.
// If old slice has newStr, return directly, else append
func AppendStr(old []string, newStr string) []string {
	if HasStr(old, newStr) {
		return old
	}
	return append(old, newStr)
}
