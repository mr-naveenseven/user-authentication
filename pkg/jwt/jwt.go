package jwt

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidUserID = errors.New("token generation: Invalid user ID")
	ErrEmptyUserName = errors.New("token generation: Username cannot be empty")
)

type AuthToken struct {
	secretKey          []byte
	accessToken        *jwt.Token
	EncodedAccessToken string
	accessTokenExpiry  int
}

type AuthTokenConfig struct {
	SecretKey         []byte
	AccessTokenExpiry int
}

type AccessTokenClaims struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func NewAuthToken(config AuthTokenConfig) *AuthToken {

	return &AuthToken{
		secretKey:          []byte(config.SecretKey),
		accessToken:        nil,
		EncodedAccessToken: "",
		accessTokenExpiry:  config.AccessTokenExpiry,
	}
}

func (authToken *AuthToken) createAccessToken(userId int, userName string) error {
	claims := &AccessTokenClaims{
		UserID:   userId,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "mr-naveenseven/user-authentication",
			Subject:   "AccessToken",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	authToken.accessToken = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := authToken.accessToken.SignedString(authToken.secretKey)
	if err != nil {
		log.Println("Failed to sign access token:", err)

		return err
	}
	authToken.EncodedAccessToken = signedToken

	return nil
}

func (authToken *AuthToken) Create(userId int, userName string) error {

	if len(authToken.secretKey) == 0 {

		return jwt.ErrHashUnavailable
	}

	if userId <= 0 {

		return ErrInvalidUserID
	}
	if userName == "" {

		return ErrEmptyUserName
	}

	err := authToken.createAccessToken(userId, userName)
	if err != nil {

		return err
	}

	return nil
}

func (authToken *AuthToken) Validate(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(authToken.secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Println("Failed to parse/validate access token:", err)

		return false, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		log.Printf("Access token valid. UserID: %d, UserName: %s, ExpiresAt: %v\n",
			claims.UserID, claims.UserName, claims.ExpiresAt)

		return true, nil
	} else {
		log.Println("Invalid access token claims")

		return false, nil
	}
}
