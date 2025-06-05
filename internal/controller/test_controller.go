package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hhr0815hhr/gint/internal/logic"
	"github.com/hhr0815hhr/gint/internal/pkg/response"
)

type TestController struct {
	testLogic *logic.TestLogic
}

func NewTestController() *TestController {
	return &TestController{}
}

var _ Router = (*TestController)(nil)

func (c *TestController) RegisterRoute(r *gin.Engine) {
	// 在这里定义你的路由
	t := r.Group("/test")
	{
		t.GET("", c.Test)
	}
}

// Test
// @Summary 测试
// @Description 测试
// @Tags 测试
// @Produce json
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "失败"
// @Router /test [get]
func (c *TestController) Test(ctx *gin.Context) {
	err := c.testLogic.Test()
	if err != nil {
		response.Error(ctx, 400, err.Error())
		return
	}
	response.Success(ctx, nil)
}
