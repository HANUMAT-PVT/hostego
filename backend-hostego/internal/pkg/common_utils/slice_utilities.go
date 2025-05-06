package common_utils

func Contains(interfaceSlice []interface{}, value interface{}) bool {
	for _, value1 := range interfaceSlice {
		if value == value1 {
			return true
		}
	}
	return false
}

func ContainsString(sliceString []string, value string) bool {
	for _, value1 := range sliceString {
		if value == value1 {
			return true
		}
	}
	return false
}
