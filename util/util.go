package util

func HasNil(arr []interface{}) bool {
	for _, elem := range arr {
		if elem == nil {
			return true
		}
	}
	return false
}
