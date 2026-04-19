package entity

type StoreConfig struct {
	ConfigID string `gorm:"primaryKey"`
	Origin   string
	Address  string
	ShopName string
	CityID   string
}
