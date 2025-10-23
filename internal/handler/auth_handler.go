package handler

import (
	"net/http"
	"user-authentication/internal/core/auth"
	"user-authentication/internal/core/user"

	"github.com/gin-gonic/gin"
)

// AuthHandlerPort represents the authentication handler port
type AuthHandlerPort interface {
	Login(c *gin.Context)
}

// AuthHandler represents the authenctication handler
type AuthHandler struct {
	authService auth.AuthServicePort
}

// NewAuthHandler create a new authentication handler to be used by the router
func NewAuthHandler(authService auth.AuthServicePort) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login authenticates the user with username and password
func (handler *AuthHandler) Login(c *gin.Context) {
	var user user.User

	// Bind JSON body to struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid request body",
			"detail": err.Error(),
		})
		return
	}

	accessToken, err := handler.authService.Login(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to login",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Login successfully", "token": accessToken})

}
