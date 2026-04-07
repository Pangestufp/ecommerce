package repository

import (
	"backend/entity"
	"backend/errorhandler"
	"errors"

	"gorm.io/gorm"
)

type InventoryRepository interface {
	Create(inventory *entity.Inventory) error
	Update(inventory *entity.Inventory) error
	GetByID(batchID string) (*entity.Inventory, error)
	GetAllByProductID(productID string) ([]entity.Inventory, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *inventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(inventory *entity.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *inventoryRepository) Update(inventory *entity.Inventory) error {
	result := r.db.Save(inventory)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &errorhandler.InternalServerError{Message: "No Row Effect"}
	}
	return nil
}

func (r *inventoryRepository) GetByID(batchID string) (*entity.Inventory, error) {
	var inventory entity.Inventory
	err := r.db.First(&inventory, "batch_id = ?", batchID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Inventory Not Found"}
		}
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) GetAllByProductID(productID string) ([]entity.Inventory, error) {
	var inventories []entity.Inventory
	err := r.db.Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&inventories).Error
	if err != nil {
		return nil, err
	}
	return inventories, nil
}
