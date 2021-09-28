package middleware

import "github.com/gin-gonic/gin"

func SetupCommonMiddleware(r *gin.Engine) *gin.Engine {
	r.Use(gin.Logger())
	return r
}
