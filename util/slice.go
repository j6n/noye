package util

// Contains returns whether the input string is in the list
func Contains(input string, list ...string) bool {
	for _, item := range list {
		if item == input {
			return true
		}
	}

	return false
}
