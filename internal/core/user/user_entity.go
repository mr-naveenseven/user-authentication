package user

// User represents the requried user details
type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash"`
	IsActive     bool   `json:"is_active"`
	IsLocked     bool   `json:"is_locked"`
}
