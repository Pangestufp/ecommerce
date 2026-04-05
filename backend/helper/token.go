package helper

import (
	"backend/config"
	"backend/entity"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateToken(user *entity.User) (string, error) {

	var mySigningKey = []byte(config.ENV.SecretKey)

	claims := JWTClaims{
		user.UserID,
		user.Name,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(mySigningKey)

	return ss, err
}

func ValidateToken(tokenString string) (*string, *string, error) {
	var mySigningKey = []byte(config.ENV.SecretKey)

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, nil, errors.New("invalid token signature")
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, nil, errors.New("token expired")
		}
		return nil, nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, nil, errors.New("your token was expired")
	}

	return &claims.ID, &claims.Role, nil
}
