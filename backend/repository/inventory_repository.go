package repository

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	Create(inventory *entity.Inventory, userID string, userName string) error
	Update(inventory *entity.Inventory, userID string, userName string) error
	GetByID(batchID string) (*entity.Inventory, error)
	GetAllByProductID(productID string, cursor *dto.Paginate, search string, limit int) ([]entity.Inventory, error)
	GetNextSeq(productID string) (int, string, error)
	GetHighestCostByProductID(productID string) (*entity.Inventory, error)
	CommitReservations(tx *gorm.DB, salesOrderID string, salesOrderCode string, now time.Time) error
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

// fungsi GetHighestCostByProductID
func (r *inventoryRepository) GetHighestCostByProductID(productID string) (*entity.Inventory, error) {
	var inv entity.Inventory
	err := r.db.Where("product_id = ? and stock > 0", productID).Order("cost_price Desc").First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepository) CommitReservations(tx *gorm.DB, salesOrderID string, salesOrderCode string, now time.Time) error {
	// Ambil semua reservation milik order ini
	var reservations []entity.StockReservation
	if err := tx.
		Where("sales_order_id = ? AND status = ?", salesOrderID, "RESERVED").
		Find(&reservations).Error; err != nil {
		return err
	}

	if len(reservations) == 0 {
		return &errorhandler.NotFoundError{Message: "Transaksi tidak valid"}
	}

	// Kumpulkan batch IDs untuk ambil cost price sekaligus
	batchIDs := make([]string, 0, len(reservations))
	for _, r := range reservations {
		batchIDs = append(batchIDs, r.BatchID)
	}

	var inventories []entity.Inventory
	if err := tx.
		Where("batch_id IN ?", batchIDs).
		Find(&inventories).Error; err != nil {
		return err
	}

	// Map batch_id → inventory untuk lookup O(1)
	inventoryMap := make(map[string]entity.Inventory, len(inventories))
	for _, inv := range inventories {
		inventoryMap[inv.BatchID] = inv
	}

	transactions := make([]entity.Transaction, 0, len(reservations))

	for _, res := range reservations {
		inv, ok := inventoryMap[res.BatchID]
		if !ok {
			return fmt.Errorf("inventory batch %s tidak ditemukan", res.BatchID)
		}

		// Kurangi stock dan reserved_stock secara atomik
		if err := tx.Model(&entity.Inventory{}).
			Where("batch_id = ?", res.BatchID).
			Updates(map[string]interface{}{
				"stock":          gorm.Expr("stock - ?", res.Quantity),
				"reserved_stock": gorm.Expr("reserved_stock - ?", res.Quantity),
				"updated_at":     now,
			}).Error; err != nil {
			return err
		}

		// Hitung cost
		qty := decimal.NewFromInt(int64(res.Quantity))
		totalCost := inv.CostPrice.Mul(qty)

		transactions = append(transactions, entity.Transaction{
			TransactionID: uuid.NewString(),
			BatchID:       res.BatchID,
			Type:          "Out",
			Quantity:      res.Quantity,
			ReferenceType: "Sales-Order",
			ReferenceID:   salesOrderID,
			Note:          fmt.Sprintf("Penjualan order %s", salesOrderCode),
			Cost:          &inv.CostPrice,
			TotalCost:     &totalCost,
			CreatedAt:     now,
		})
	}

	// Bulk insert transactions
	if err := tx.Create(&transactions).Error; err != nil {
		return err
	}

	// Update semua reservation jadi COMMITTED
	if err := tx.Model(&entity.StockReservation{}).
		Where("sales_order_id = ? AND status = ?", salesOrderID, "RESERVED").
		Updates(map[string]interface{}{
			"status":     "COMMITTED",
			"updated_at": now,
		}).Error; err != nil {
		return err
	}

	return nil
}
