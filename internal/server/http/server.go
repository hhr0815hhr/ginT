package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hhr0815hhr/gint/internal/controller"
	"github.com/hhr0815hhr/gint/internal/middleware"
)

type HTTPRoutes struct {
	//Port    string
	Routers []controller.Router
}

func NewHTTPServer(opts *HTTPRoutes) *gin.Engine {
	r := gin.New()
	r.Use(
		gin.Recovery(),
		middleware.Cors(),
		middleware.Locale(),
		gin.Logger(),
	)

	for _, router := range opts.Routers {
		router.RegisterRoute(r)
	}
	return r
}
