package entity

import "time"

type Product struct {
	ProductID   string `gorm:"primaryKey"`
	ProductCode string
	ProductName string
	ProductSlug string
	WeightGram  int
	TypeID      string
	Description string
	Status      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
