package main

func IsInArray[V comparable](array []V, value V) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}

	return false
}
