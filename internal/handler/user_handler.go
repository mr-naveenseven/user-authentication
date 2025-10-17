package handler

import (
	"net/http"
	"strconv"
	"user-authentication/internal/core/user"

	"github.com/gin-gonic/gin"
)

type UserHandlerPort interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
}

type UserHandler struct {
	userService user.UserServicePort
}

func NewUserHandler(userService user.UserServicePort) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (handler *UserHandler) Create(c *gin.Context) {
	var user user.User

	// Bind JSON body to struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid request body",
			"detail": err.Error(),
		})
		return
	}

	// Call service layer to create user
	user, err := handler.userService.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create user",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})

}

func (handler *UserHandler) Get(c *gin.Context) {
	// Call service layer to fetch all users
	users, err := handler.userService.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create user",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User fetch successful", "users": users})
}

func (handler *UserHandler) GetByID(c *gin.Context) {

	userID, _ := strconv.Atoi(c.Param("id"))

	// Call service layer to fetch user details by userID
	user, err := handler.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create user",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User fetch successful", "users": user})
}

func (handler *UserHandler) Update(c *gin.Context) {
	var user user.User

	userID, _ := strconv.Atoi(c.Query("id"))

	// Bind JSON body to struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid request body",
			"detail": err.Error(),
		})
		return
	}

	user.ID = uint(userID)

	// Call service layer to create user
	user, err := handler.userService.Update(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create user",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})

}
