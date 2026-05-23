package entity

import "time"

type RefreshToken struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	ExpiredAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
