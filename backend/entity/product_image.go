package entity

import "time"

type ProductImage struct {
	ImageID     string
	ProductID   string
	PicturePath string
	IsPrimary   int
	CreatedAt   time.Time
	VerifiedAt  *time.Time
}
