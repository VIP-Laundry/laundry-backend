package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a Bcrypt hash from a plain-text password.
func HashPassword(password string) (string, error) {
	// Using DefaultCost (10) as a balanced security/performance trade-off
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword compares a plain-text password with its hashed version.
// Returns nil on success, or an error if they don't match.
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
