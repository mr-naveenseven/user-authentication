package password

import "golang.org/x/crypto/bcrypt"

const (
	passwordHashCost = 10
)

// HashPassword hashes the provided password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)

	return string(bytes), err
}

// VerifyPassword compares a plaintext password with a hashed password and returns true if they match.
func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
