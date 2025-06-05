package util

import "math/rand"

func RandomString(length int, strType int) string {
	const charset2 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const charset1 = "0123456789"
	var charset string
	if strType == 1 {
		charset = charset1
	} else if strType == 2 {
		charset = charset2
	} else {
		charset = charset1 + charset2
	}
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
