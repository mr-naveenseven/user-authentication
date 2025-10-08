package router

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func InitAuthRoutes(r *gin.Engine) {
	// authGroup := r.Group("/auth")
}

func InitRoutes(r *gin.Engine) {
	InitAuthRoutes(r)
}
