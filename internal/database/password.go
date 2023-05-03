package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GetHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("Failed to hash password")
	}
	return string(hashedPassword), nil
}
