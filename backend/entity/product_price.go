package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProductPrice struct {
	PriceID      string `gorm:"primaryKey"`
	ProductID    string
	ProductPrice decimal.Decimal
	CreatedAt    time.Time
	CreatedBy    string
	CreatedName  string
}
