package entity

import "time"

type Discount struct {
	DiscountID    string `gorm:"primaryKey"`
	ProductID     string
	DiscountName  string
	DiscountType  string
	DiscountValue float64
	StartAt       time.Time
	ExpiredAt     time.Time
	CreatedBy     string
	CreatedName   string
	Status        int
	CreatedAt     time.Time
}
