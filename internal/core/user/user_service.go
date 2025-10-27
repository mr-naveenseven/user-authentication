package user

import (
	"errors"
	"log"
	"user-authentication/pkg/password"
)

type UserServicePort interface {
	Create(user User) (User, error)
	Update(user User) (User, error)
	GetByID(userID int) (User, error)
	Get() ([]User, error)
}

type UserService struct {
	UserRepo UserRepoPort
}

var (
	ErrInvalidUserDetails = errors.New("invalid user details")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrInvalidUserPwd     = errors.New("invalid user password")
)

func NewUserService(userRepo UserRepoPort) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (service *UserService) Create(user User) (User, error) {
	if user.Email == "" || user.Username == "" || user.Password == "" {
		return User{}, ErrInvalidUserDetails
	}

	hashedPwd, err := password.HashPassword(user.Password)
	if err != nil {
		log.Printf("error hashing password %v", err)

		return User{}, err
	}

	user.PasswordHash = hashedPwd

	user, err = service.UserRepo.Create(user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (service *UserService) GetByID(userID int) (User, error) {
	if userID <= 0 {
		return User{}, ErrInvalidUserID
	}

	user, err := service.UserRepo.GetByID(userID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (service *UserService) Get() ([]User, error) {
	users, err := service.UserRepo.Get()
	if err != nil {
		return users, err
	}

	return users, nil
}

func (service *UserService) Update(user User) (User, error) {
	if user.Email == "" || user.Username == "" {
		return User{}, ErrInvalidUserDetails
	}

	user, err := service.UserRepo.Update(user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
