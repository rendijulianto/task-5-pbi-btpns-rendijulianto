package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword menghash kata sandi dengan bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash memeriksa apakah kata sandi yang diberikan cocok dengan hash yang disimpan.
func CheckPasswordHash(password, hash string)(bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}