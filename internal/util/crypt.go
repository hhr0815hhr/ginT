package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/hhr0815hhr/gint/internal/pkg/aes"
)

func MD5(data string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, data)
	if err != nil {
		// 在实际应用中可能需要更完善的错误处理
		fmt.Println("写入字符串错误:", err)
		return ""
	}
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(decodedBytes)
}

// ----------------AES CBC START----------------
func Encrypt(str string) string {
	return string(aes.AesEncryptCBC([]byte(str)))
}

func Uncrypt(str string) string {
	return string(aes.AesDecryptCBC([]byte(str)))
}

// ----------------AES CBC END----------------
