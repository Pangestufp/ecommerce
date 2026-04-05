package dto

type StoreConfigRequest struct {
	Origin   string `json:"origin"`
	Address  string `json:"address"`
	ShopName string `json:"shop_name"`
	CityID   string `json:"city_id"`
}

type StoreConfigResponse struct {
	ConfigID string `json:"config_id"`
	Origin   string `json:"origin"`
	Address  string `json:"address"`
	ShopName string `json:"shop_name"`
	CityID   string `json:"city_id"`
}
