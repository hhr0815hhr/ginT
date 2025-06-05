package util

import "testing"

func TestEncrypt(t *testing.T) {
	a := Encrypt("123")
	t.Log(a)
	t.Log(Uncrypt(a))
}
