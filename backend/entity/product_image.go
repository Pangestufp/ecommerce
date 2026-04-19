package entity

import "time"

type ProductImage struct {
	ImageID     string `gorm:"primaryKey"`
	ProductID   string
	PicturePath string
	IsPrimary   int
	CreatedAt   time.Time
	VerifiedAt  *time.Time
}
