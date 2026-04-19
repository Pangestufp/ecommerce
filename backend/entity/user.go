package entity

import "time"

type User struct {
	UserID     string `gorm:"primaryKey"`
	Name       string
	Email      string
	Status     int
	Role       string
	Password   string
	Address    string
	Phone      string
	PostalCode string
	VerifiedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
