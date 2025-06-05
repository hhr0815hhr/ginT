package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hhr0815hhr/gint/internal/util"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		authInfo, err := util.DecodeToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
		if authInfo["expire"].(float64) < float64(time.Now().Unix()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"data": nil,
				"msg":  "token expired",
			})
		}
		c.Set("uuid", authInfo["uuid"].(string))
	}
}
