package dto

type StoreConfigRequest struct {
	ShopName          string `json:"shop_name"`
	Phone             string `json:"phone"`
	ProvinceID        string `json:"province_id"`
	CityID            string `json:"city_id"`
	DistrictID        string `json:"district_id"`
	SubDistrictID     string `json:"sub_district_id"`
	AdditionalAddress string `json:"additional_address"`
}

type StoreConfigResponse struct {
	ConfigID          string `json:"config_id"`
	ShopName          string `json:"shop_name"`
	Phone             string `json:"phone"`
	ProvinceID        string `json:"province_id"`
	ProvinceName      string `json:"province_name"`
	CityID            string `json:"city_id"`
	CityName          string `json:"city_name"`
	DistrictID        string `json:"district_id"`
	DistrictName      string `json:"district_name"`
	SubDistrictID     string `json:"sub_district_id"`
	SubDistrictName   string `json:"sub_district_name"`
	ZipCode           string `json:"zip_code"`
	AdditionalAddress string `json:"additional_address"`
}
