package password

import (
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// Hash generates a bcrypt hash from a plain text password
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify checks if a plain text password matches a hash
func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
