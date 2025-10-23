package user

import (
	"fmt"
	"time"
	"user-authentication/internal/postgres"
)

const (
	uaTable       = "user_accounts"
	uaColUserName = "username"
	uaColEmail    = "email"
	uaColPwdHash  = "password_hash"
	uaColIsActive = "is_active"
	uaColIsLocked = "is_locked"
)

type repoUser struct {
	ID                 int       `gorm:"column:id;primaryKey;autoIncrement"`
	Username           string    `gorm:"column:username"`
	Email              string    `gorm:"column:email"`
	PasswordHash       string    `gorm:"column:password_hash"`
	IsActive           bool      `gorm:"column:is_active"`
	IsLocked           bool      `gorm:"column:is_locked"`
	PasswordModifiedAt time.Time `gorm:"column:password_modified_at"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	ModifiedAt         time.Time `gorm:"column:modified_at"`
}

func (repoUser) TableName() string {
	return uaTable
}

// UserRepo represents the User Repository
type UserRepo struct {
	pgClient *postgres.PGClient
}

func NewUserRepo(pgClient *postgres.PGClient) *UserRepo {
	return &UserRepo{
		pgClient: pgClient,
	}
}

func toRepoUser(u User) repoUser {
	now := time.Now()

	return repoUser{
		ID:                 u.ID,
		Username:           u.Username,
		Email:              u.Email,
		PasswordHash:       u.PasswordHash,
		IsActive:           u.IsActive,
		IsLocked:           u.IsLocked,
		PasswordModifiedAt: time.Time{},
		CreatedAt:          now,
		ModifiedAt:         now,
	}
}

func toEntityUser(u repoUser) User {
	return User{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		IsActive:     u.IsActive,
		IsLocked:     u.IsLocked,
	}
}

func toEntityUsers(repoUsers []repoUser) []User {
	users := make([]User, 0, len(repoUsers))
	for _, u := range repoUsers {
		users = append(users, toEntityUser(u))
	}
	return users
}

func (repo *UserRepo) Create(user User) (User, error) {

	rUser := toRepoUser(user)
	rUser.PasswordModifiedAt = time.Now()

	if err := repo.pgClient.DB.Create(&rUser).Error; err != nil {
		return User{}, err
	}

	return toEntityUser(rUser), nil
}

func (repo *UserRepo) GetByID(userID int) (User, error) {

	var rUser repoUser
	err := repo.pgClient.DB.Select(uaColUserName, uaColEmail, uaColIsActive, uaColIsLocked).
		First(&rUser, userID).Error
	if err != nil {
		return User{}, fmt.Errorf("user fetch failed: %v", err)
	}

	return toEntityUser(rUser), nil
}

func (repo *UserRepo) GetByUsername(username string) (User, error) {
	var rUser repoUser
	err := repo.pgClient.DB.Where(&repoUser{
		Username: username,
	}).First(&rUser).Error
	if err != nil {
		return User{}, fmt.Errorf("user fetch failed: %v", err)
	}

	return toEntityUser(rUser), nil
}

func (repo *UserRepo) Get() ([]User, error) {

	var users []repoUser
	res := repo.pgClient.DB.Find(&users)
	if res.Error != nil {
		return []User{}, fmt.Errorf("users fetch failed: %v", res.Error)
	}

	if res.RowsAffected == 0 {
		return []User{}, fmt.Errorf("no users found")
	}

	return toEntityUsers(users), nil
}

func (repo *UserRepo) Update(user User) (User, error) {

	rUser := toRepoUser(user)

	res := repo.pgClient.DB.Model(&repoUser{}).Where(&repoUser{
		ID: user.ID,
	}).Select("username", "email").Updates(rUser)
	if res.Error != nil {
		return User{}, res.Error
	}

	if res.RowsAffected == 0 {
		return User{}, fmt.Errorf("user not found: %v", res.Error)
	}

	return user, nil
}
