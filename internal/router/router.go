package router

import (
	"net/http"
	"strings"
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

// authMiddleware checks for the presence and validity of the authentication header token
func (r *Router) authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	isValid, err := r.authHandler.ValidateAccessToken(authToken[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()

}

// registerUserRoutes registers the user routes
func (r *Router) registerUserRoutes() {
	userGroup := r.e.Group("/user")
	userGroup.Use(r.authMiddleware)
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
