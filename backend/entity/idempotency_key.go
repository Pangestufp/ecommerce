package entity

import "time"

type CheckoutIdempotencyKey struct {
	IdempotencyID  string `gorm:"primaryKey"`
	IdempotencyKey string `gorm:"uniqueIndex"`
	UserID         string
	SalesOrderID   *string
	Status         string
	ErrorMessage   string
	ResponseBody   string
	StatusCode     int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
