package controller

import "github.com/gin-gonic/gin"

type Router interface {
	RegisterRoute(r *gin.Engine)
}
