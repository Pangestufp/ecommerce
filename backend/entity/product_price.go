package entity

import "time"

type ProductPrice struct {
	PriceID      string
	ProductID    string
	ProductPrice float64
	CreatedAt    time.Time
}
