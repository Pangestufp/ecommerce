package entity

import "time"

type Inventory struct {
	BatchID       string `gorm:"primaryKey"`
	BatchCode     string
	ProductID     string
	CostPrice     float64
	Stock         int
	ReservedStock int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
