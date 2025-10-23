package auth

import (
	"errors"
	"log"
	"user-authentication/internal/core/user"
	"user-authentication/pkg/jwt"
	"user-authentication/pkg/password"
)

// return errors for the authentication service layer
var (
	ErrIncorrectPwd  = errors.New("incorrect user password")
	ErrTokenCreation = errors.New("creating access token failed")
)

// empty token to return on the user authentication failure
const emptyToken = ""

// AuthServicePort represents the authentication service port
type AuthServicePort interface {
	Login(u user.User) (string, error)
}

// AuthService represents the authentication service
type AuthService struct {
	config   jwt.AuthTokenConfig
	userRepo user.UserRepoPort
}

// NewAuthService create a new authentication service to be used in the other layers
func NewAuthService(config jwt.AuthTokenConfig, userRepo user.UserRepoPort) *AuthService {
	return &AuthService{
		config:   config,
		userRepo: userRepo,
	}
}

// Login authenticates the user with username and password
func (service *AuthService) Login(u user.User) (string, error) {
	// username and password validation
	if u.Username == "" {
		log.Println("user details validation failed: invalid username, cannot be empty")
		return emptyToken, user.ErrInvalidUserDetails
	}
	if u.Password == "" {
		log.Println("user details validation failed: invalid user password")
		return emptyToken, user.ErrInvalidUserPwd
	}

	// fetching user details from the user repository layer
	dbUser, err := service.userRepo.GetByUsername(u.Username)
	if err != nil {
		log.Printf("user fetch failed: %v", err)
		return emptyToken, err
	}

	// user password verification from the database and handler
	isValid := password.VerifyPassword(u.Password, dbUser.PasswordHash)
	if !isValid {
		log.Println("password verficiation: wrong password")
		return emptyToken, ErrIncorrectPwd
	}

	// creating a new accesstoken
	authToken := jwt.NewAuthToken(service.config)
	err = authToken.Create(dbUser.ID, u.Username)
	if err != nil {
		log.Printf("access token creation failed: %v", err)
		return emptyToken, ErrTokenCreation
	}

	return authToken.EncodedAccessToken, nil
}
