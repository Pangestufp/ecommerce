package repository

import (
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"errors"

	"gorm.io/gorm"
)

type InventoryRepository interface {
	Create(inventory *entity.Inventory) error
	Update(inventory *entity.Inventory) error
	GetByID(batchID string) (*entity.Inventory, error)
	GetAllByProductID(productID string) ([]entity.Inventory, error)
	GetNextSeq(productID string) (int, string, error)
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

func (r *inventoryRepository) GetNextSeq(productID string) (int, string, error) {
	yearMonth := helper.TimeNowWIB().Format("200601")

	query := `
		INSERT INTO product_batch_sequences (product_id, year_month, last_seq)
		VALUES (?, ?, 1)
		ON CONFLICT (product_id, year_month) DO UPDATE
		SET last_seq = product_batch_sequences.last_seq + 1
		RETURNING last_seq
	`

	var seq int
	err := r.db.Raw(query, productID, yearMonth).Scan(&seq).Error
	if err != nil {
		return 0, "", err
	}

	return seq, yearMonth, nil
}
