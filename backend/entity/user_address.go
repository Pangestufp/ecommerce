package entity

import "time"

type UserAddress struct {
	AddressID         string `gorm:"primaryKey"`
	UserID            string
	Label             string
	RecipientName     string
	Phone             string
	ProvinceID        string
	ProvinceName      string
	CityID            string
	CityName          string
	DistrictID        string
	DistrictName      string
	ZipCode           string
	AdditionalAddress string
	IsPrimary         int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
