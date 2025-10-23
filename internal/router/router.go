package router

import (
	"user-authentication/internal/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	e           *gin.Engine
	port        string
	userHandler handler.UserHandlerPort
	authHandler handler.AuthHandlerPort
}

func NewRouter(
	port string,
	userHandler handler.UserHandlerPort,
	authHandler handler.AuthHandlerPort,
) *Router {
	return &Router{
		e:           nil,
		port:        port,
		userHandler: userHandler,
		authHandler: authHandler,
	}
}

func (r *Router) InitRouter() {
	r.e = gin.Default()
	r.initRoutes()
}

func (r *Router) initRoutes() {
	r.registerUserRoutes()
	r.registerAuthRoutes()
}

func (r *Router) Run() {
	r.e.Run(":" + r.port)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	return r
}

// registerUserRoutes registers the user routes
func (r *Router) registerUserRoutes() {
	userGroup := r.e.Group("/user")
	userGroup.POST("", r.userHandler.Create)
	userGroup.GET("", r.userHandler.Get)
	userGroup.GET("/:id", r.userHandler.GetByID)
	userGroup.PUT("/:id", r.userHandler.Update)
}

// registerAuthRoutes registers the authentication routes
func (r *Router) registerAuthRoutes() {
	authGroup := r.e.Group("/auth")
	authGroup.GET("/login", r.authHandler.Login)
}
