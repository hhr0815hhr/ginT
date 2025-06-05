package response

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hhr0815hhr/gint/internal/config"
	"github.com/hhr0815hhr/gint/internal/log"
	"github.com/hhr0815hhr/gint/internal/util"
)

type Response struct {
	Code int         `json:"code"` // 业务状态码 (例如：成功、失败、未授权等)
	Msg  string      `json:"msg"`  // 消息描述
	Data interface{} `json:"data"` // 返回的数据
}

// Success 返回成功的响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code: 0,
		Msg:  "success",
		Data: util.Ternary[interface{}](data != nil, data, struct{}{}),
	})
}

// Error 返回业务错误的响应
func Error(c *gin.Context, code int, message string) {
	log.Logger.Errorf("请求失败，错误码：%d, 错误信息：%s\n", code, message)
	if config.Conf.Server.Env == "dev" {
		// 获取调用栈信息
		pc, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		// 获取函数名
		funcName := runtime.FuncForPC(pc).Name()
		// 格式化文件名和行号
		fileLine := fmt.Sprintf("%s:%d", strings.TrimPrefix(file, "./"), line)
		// 打印日志
		log.Logger.Printf("Trace: 文件：%s, 函数：%s\n", fileLine, funcName)
	}
	c.JSON(200, Response{ // 仍然返回 HTTP 200，但业务状态码表示错误
		Code: code,
		Msg:  message,
		Data: struct{}{},
	})
}

// Custom 返回自定义 HTTP 状态码的响应 (可以根据需要扩展)
func Custom(c *gin.Context, httpCode int, code int, message string, data interface{}) {
	c.JSON(httpCode, Response{
		Code: code,
		Msg:  message,
		Data: data,
	})
}
