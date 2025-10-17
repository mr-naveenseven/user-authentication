package router

import (
	"user-authentication/internal/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	e           *gin.Engine
	port        string
	userHandler handler.UserHandlerPort
}

func NewRouter(port string, userHandler handler.UserHandlerPort) *Router {
	return &Router{
		e:           nil,
		port:        port,
		userHandler: userHandler,
	}
}

func (r *Router) InitRouter() {
	r.e = gin.Default()
	r.initRoutes()
}

func (r *Router) initRoutes() {
	r.initUserRoutes()
}

func (r *Router) Run() {
	r.e.Run(":" + r.port)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	return r
}

func (r *Router) initUserRoutes() {
	userGroup := r.e.Group("/user")
	userGroup.POST("", r.userHandler.Create)
	userGroup.GET("", r.userHandler.Get)
	userGroup.GET("/:id", r.userHandler.GetByID)
	userGroup.PUT("/:id", r.userHandler.Update)
}
