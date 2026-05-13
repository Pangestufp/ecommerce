package entity

type StoreConfig struct {
	ConfigID          string `gorm:"primaryKey"`
	ShopName          string
	Phone             string
	ProvinceID        string
	ProvinceName      string
	CityID            string
	CityName          string
	DistrictID        string
	DistrictName      string
	SubDistrictID     string
	SubDistrictName   string
	ZipCode           string
	AdditionalAddress string
}
