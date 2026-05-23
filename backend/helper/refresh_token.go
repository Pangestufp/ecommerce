package helper

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func HashRefreshToken(token string) (string, error) {
	tokenHash, err := bcrypt.GenerateFromPassword(
		[]byte(token),
		bcrypt.DefaultCost,
	)

	return string(tokenHash), err
}

func VerifyRefreshToken(hashToken, token string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashToken),
		[]byte(token),
	)
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
