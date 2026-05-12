package dto

import "time"

type CreateAddressRequest struct {
	Label             string `json:"label"`
	RecipientName     string `json:"recipient_name"`
	Phone             string `json:"phone"`
	ProvinceID        string `json:"province_id"`
	CityID            string `json:"city_id"`
	DistrictID        string `json:"district_id"`
	SubDistrictID     string `json:"sub_district_id"`
	AdditionalAddress string `json:"additional_address"`
	IsPrimary         int    `json:"is_primary"`
}

type UpdateAddressRequest struct {
	Label             string `json:"label"`
	RecipientName     string `json:"recipient_name"`
	Phone             string `json:"phone"`
	ProvinceID        string `json:"province_id"`
	CityID            string `json:"city_id"`
	DistrictID        string `json:"district_id"`
	SubDistrictID     string `json:"sub_district_id"`
	AdditionalAddress string `json:"additional_address"`
	IsPrimary         int    `json:"is_primary"`
}

type AddressResponse struct {
	AddressID         string    `json:"address_id"`
	UserID            string    `json:"user_id"`
	Label             string    `json:"label"`
	RecipientName     string    `json:"recipient_name"`
	Phone             string    `json:"phone"`
	ProvinceID        string    `json:"province_id"`
	ProvinceName      string    `json:"province_name"`
	CityID            string    `json:"city_id"`
	CityName          string    `json:"city_name"`
	DistrictID        string    `json:"district_id"`
	DistrictName      string    `json:"district_name"`
	SubDistrictID     string    `json:"sub_district_id"`
	SubDistrictName   string    `json:"sub_district_name"`
	ZipCode           string    `json:"zip_code"`
	AdditionalAddress string    `json:"additional_address"`
	IsPrimary         int       `json:"is_primary"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
