package repository

import (
	"backend/entity"
	"errors"

	"gorm.io/gorm"
)

type StoreConfigRepository interface {
	GetConfig() (*entity.StoreConfig, error)
	CreateConfig(config *entity.StoreConfig) error
	UpdateConfig(config *entity.StoreConfig) error
}

type storeConfigRepository struct {
	db *gorm.DB
}

func NewStoreConfigRepository(db *gorm.DB) *storeConfigRepository {
	return &storeConfigRepository{db: db}
}

func (r *storeConfigRepository) GetConfig() (*entity.StoreConfig, error) {
	var config entity.StoreConfig
	err := r.db.First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // belum ada config, return nil
		}
		return nil, err
	}
	return &config, nil
}

func (r *storeConfigRepository) CreateConfig(config *entity.StoreConfig) error {
	return r.db.Create(config).Error
}

func (r *storeConfigRepository) UpdateConfig(config *entity.StoreConfig) error {
	return r.db.Model(&entity.StoreConfig{}).
		Where("config_id = ?", config.ConfigID).
		Updates(map[string]interface{}{
			"origin":    config.Origin,
			"address":   config.Address,
			"shop_name": config.ShopName,
			"city_id":   config.CityID,
		}).Error
}
