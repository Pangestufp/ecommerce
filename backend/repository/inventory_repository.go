package repository

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	Create(inventory *entity.Inventory, userID string, userName string) error
	Update(inventory *entity.Inventory, userID string, userName string) error
	GetByID(batchID string) (*entity.Inventory, error)
	GetAllByProductID(productID string, cursor *dto.Paginate, search string, limit int) ([]entity.Inventory, error)
	GetNextSeq(productID string) (int, string, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *inventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(inventory *entity.Inventory, userID string, userName string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(inventory).Error; err != nil {
			return err
		}

		transaction := &entity.Transaction{
			TransactionID: uuid.New().String(),
			BatchID:       inventory.BatchID,
			Type:          helper.In(),
			Quantity:      inventory.Stock,
			ReferenceType: helper.BatchCreate(),
			ReferenceID:   userID,
			Note:          "Pembuatan batch - " + inventory.BatchCode + " Oleh " + userName,
			CreatedAt:     time.Now(),
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *inventoryRepository) Update(inventory *entity.Inventory, userID string, userName string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var current entity.Inventory
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			First(&current, "batch_id = ?", inventory.BatchID).Error; err != nil {
			return err
		}

		if current.ReservedStock > inventory.Stock {
			return &errorhandler.InternalServerError{Message: "Jumlah reservasi melebihi dari stock"}
		}

		diff := inventory.Stock - current.Stock

		result := tx.Save(inventory)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return &errorhandler.InternalServerError{Message: "No Row Effect"}
		}

		if diff != 0 {
			txType := helper.In()
			if diff < 0 {
				txType = helper.Out()
				diff = -diff
			}

			correction := &entity.Transaction{
				TransactionID: uuid.New().String(),
				BatchID:       inventory.BatchID,
				Type:          txType,
				Quantity:      diff,
				ReferenceType: helper.StockAdjust(),
				ReferenceID:   userID,
				Note:          "Koreksi perubahan - " + inventory.BatchCode + " Oleh " + userName,
				CreatedAt:     time.Now(),
			}

			if err := tx.Create(correction).Error; err != nil {
				return err
			}
		}

		return nil
	})
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

func (r *inventoryRepository) GetAllByProductID(productID string, cursor *dto.Paginate, search string, limit int) ([]entity.Inventory, error) {
	if limit <= 0 {
		limit = 5
	}

	var inventories []entity.Inventory

	query := r.db.Model(&entity.Inventory{}).
		Where("product_id = ?", productID)

	if search != "" {
		query = query.Where("batch_code ILIKE ?", "%"+search+"%")
	}

	if cursor != nil {
		if cursor.Direction != nil && *cursor.Direction == "prev" {
			if cursor.FirstID != nil && cursor.FirstCreatedAt != nil {
				query = query.Where("(created_at, batch_id) > (?, ?)", cursor.FirstCreatedAt, cursor.FirstID)
				query = query.Order("created_at ASC, batch_id ASC")
			}
		} else {
			if cursor.Direction != nil && *cursor.Direction == "next" {
				query = query.Where("(created_at, batch_id) < (?, ?)", cursor.LastCreatedAt, cursor.LastID)
			}
			query = query.Order("created_at DESC, batch_id DESC")
		}
	} else {
		query = query.Order("created_at DESC, batch_id DESC")
	}

	err := query.Limit(limit + 1).Find(&inventories).Error
	if err != nil {
		return nil, err
	}

	if cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev" {
		for i, j := 0, len(inventories)-1; i < j; i, j = i+1, j-1 {
			inventories[i], inventories[j] = inventories[j], inventories[i]
		}
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
