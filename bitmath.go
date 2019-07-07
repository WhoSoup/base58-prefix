package main

// carryover addition
func co(arr []byte, depth int) []byte {
	arr[depth]++
	if arr[depth] == 0 && depth > 0 {
		return co(arr, depth-1)
	}
	return arr
}

// carryover subtraction
func cos(arr []byte, depth int) []byte {
	arr[depth]--
	if arr[depth] == 255 && depth > 0 {
		return cos(arr, depth-1)
	}
	return arr
}

// match a prefix up to position i
func matchall(b, match string, i int) bool {
	for j := 0; j <= i; j++ {
		if b[j] != match[j] {
			return false
		}
	}
	return true
}
