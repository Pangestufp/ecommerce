package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	TransactionID string `gorm:"primaryKey"`
	BatchID       string
	Type          string
	Quantity      int
	ReferenceType string
	ReferenceID   string
	Note          string
	Cost          *decimal.Decimal
	TotalCost     *decimal.Decimal
	CreatedAt     time.Time
}
