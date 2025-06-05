package airwallex

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hhr0815hhr/gint/internal/config"
	"github.com/hhr0815hhr/gint/internal/pkg/i18n"
)

// Airwallex 客户端配置
type Airwallex struct {
	APIKey   string
	ClientID string
	BaseURL  string
}

var once sync.Once

// 创建一个新的 Airwallex 客户端
func Instance() *Airwallex {
	var a Airwallex
	once.Do(func() {
		a = Airwallex{
			APIKey:   config.Conf.Server.AirWallex.Key,
			ClientID: config.Conf.Server.AirWallex.Id,
			BaseURL:  config.Conf.Server.AirWallex.Url,
		}
	})
	return &a
}

// CreatePaymentIntent 创建支付意图
func (a *Airwallex) CreatePaymentIntent(amount float64, user, title string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/pa/payment_intents/create", a.BaseURL)
	data := map[string]interface{}{
		"amount":               amount,
		"currency":             "USD",
		"country_code":         "SG",
		"customer_id":          user,
		"description":          title,
		"statement_descriptor": title,
	}
	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("x-client-id", a.ClientID)
	req.Header.Set("x-api-key", a.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return result, nil
}

func (a *Airwallex) CheckCallback(ctx *gin.Context) (gin.H, error) {
	signature := ctx.GetHeader("x-signature")
	timestamp := ctx.GetHeader("x-timestamp")
	// 读取请求体
	body := ctx.Request.Body
	defer body.Close()
	var b []byte
	body.Read(b)
	b, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		err = fmt.Errorf(i18n.T(ctx, "error.invalidRequest"))
		return gin.H{}, err
	}

	// 验证签名
	if !verifySignature(signature, timestamp, string(b)) {
		err = fmt.Errorf(i18n.T(ctx, "error.invalidSignature"))
		return gin.H{}, err
	}

	// 解析回调数据
	var data = make(gin.H)
	if err = json.Unmarshal(b, &data); err != nil {
		err = fmt.Errorf(i18n.T(ctx, "error.invalidRequest"))
		return gin.H{}, err
	}
	return data, err

}

// GetToken 获取认证 Token
func (a *Airwallex) GetToken() (string, error) {
	url := fmt.Sprintf("%s/api/v1/authentication/login", a.BaseURL)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("x-client-id", a.ClientID)
	req.Header.Set("x-api-key", a.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResponse map[string]interface{}
	json.Unmarshal(body, &tokenResponse)
	return fmt.Sprintf("Bearer %s", tokenResponse["token"].(string)), nil
}

func verifySignature(signature, timestamp, body string) bool {
	// 构造签名验证字符串
	valueToDigest := timestamp + body

	// 使用 HMAC-SHA256 计算预期签名
	h := hmac.New(sha256.New, []byte(config.Conf.Server.AirWallex.WebhookSecret))
	h.Write([]byte(valueToDigest))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// 比较签名
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
