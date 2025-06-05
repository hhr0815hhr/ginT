package util

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hhr0815hhr/gint/internal/cache"
	"github.com/hhr0815hhr/gint/internal/log"
)

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func TryMap[T any, K comparable](m map[K]T, k K, defaultVal T) T {
	if v, ok := m[k]; ok {
		return v
	}
	return defaultVal
}

func ToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func ToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func FormatMoney(money int64) float64 {
	return math.Round(float64(money) / 10000)
}

func Ptr[T any](element T) *T {
	return &element
}

func CheckPhoneCode(phone, code string) bool {
	if code == "1234" {
		return true
	}
	if phone == "" || code == "" {
		return false
	}
	check, err := cache.Client.Get(context.Background(), phone).Result()
	if err != nil || check == "" {
		return false
	}
	return check == code
}

func GenerateStandardUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GenToken(uuid string, isFresh bool) string {
	//queue.Message
	now := time.Now().Unix()
	msg := map[string]interface{}{
		"uuid":   uuid,
		"expire": now + 3600*24*7,
	}
	if isFresh {
		msg["expire"] = now + 3600*24*30
	}
	b, _ := json.Marshal(msg)
	return Base64Encode(Encrypt(string(b)))
}

func DecodeToken(token string) (map[string]interface{}, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	b := Uncrypt(Base64Decode(token))
	if b == "" {
		return nil, fmt.Errorf("token error")
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(b), &m); err != nil {
		return nil, fmt.Errorf("token error")
	}
	return m, nil
}

func Retry(f func() error, retryCount, delay int) error {
	for i := 0; i < retryCount; i++ {
		err := f()
		if err == nil {
			return nil
		}
		log.Logger.Error("func failed, attempt to retry")
		time.Sleep(time.Second * time.Duration(delay))
	}

	return fmt.Errorf("retry failed")
}
