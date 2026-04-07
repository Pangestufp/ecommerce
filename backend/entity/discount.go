package entity

import "time"

type Discount struct {
	DiscountID    string
	ProductID     string
	DiscountName  string
	DiscountType  string
	DiscountValue float64
	StartAt       time.Time
	ExpiredAt     time.Time
	Status        int
	CreatedAt     time.Time
}
