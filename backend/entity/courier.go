package entity

import "time"

type Courier struct {
	ID        string `gorm:"primaryKey"`
	Code      string
	Name      string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
