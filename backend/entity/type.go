package entity

import "time"

type Type struct {
	TypeID    string `gorm:"primaryKey"`
	TypeCode  string
	TypeName  string
	TypeDesc  string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
