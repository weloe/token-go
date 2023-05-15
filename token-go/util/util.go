package util

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
