package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hhr0815hhr/gint/internal/pkg/i18n"
)

func Locale() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("Accept-Language")
		if locale == "" {
			locale = i18n.DefaultLocale
		}
		if strings.Contains(locale, "zh") {
			locale = "zh"
		} else {
			locale = "en"
		}
		c.Set("locale", locale)
	}
}
