package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Discount struct {
	DiscountID    string `gorm:"primaryKey"`
	ProductID     string
	DiscountName  string
	DiscountType  string
	DiscountValue decimal.Decimal
	StartAt       time.Time
	ExpiredAt     time.Time
	CreatedBy     string
	CreatedName   string
	Status        int
	CreatedAt     time.Time
}
