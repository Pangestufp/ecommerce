package entity

import "time"

type ProductPrice struct {
	PriceID      string `gorm:"primaryKey"`
	ProductID    string
	ProductPrice float64
	CreatedAt    time.Time
	CreatedBy    string
	CreatedName  string
}
