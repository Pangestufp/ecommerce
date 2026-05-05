package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Inventory struct {
	BatchID       string `gorm:"primaryKey"`
	BatchCode     string
	ProductID     string
	CostPrice     decimal.Decimal
	Stock         int
	ReservedStock int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
