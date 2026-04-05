package helper

import (
	"backend/config"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password+config.ENV.SecretKey), bcrypt.DefaultCost)
	return string(passwordHash), err
}

func VerifyPassword(hashPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password+config.ENV.SecretKey))

	return err
}
