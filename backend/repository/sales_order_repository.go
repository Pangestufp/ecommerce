package repository

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SalesOrderRepository interface {
	Transaction(fn func(tx *gorm.DB) error) error
	ReserveStock(tx *gorm.DB, dbDetails []entity.SalesOrderDetail) error
	CreateOrder(tx *gorm.DB, order *entity.SalesOrder) error
	CreateOrderDetails(tx *gorm.DB, details []entity.SalesOrderDetail) error
	GetNextSeq(now time.Time) (int, string, error)

	GetAllPaginate(cursor *dto.Paginate, statuses []string, limit int) ([]entity.SalesOrder, error)
	GetAllByUserPaginate(cursor *dto.Paginate, userID string, statuses []string, limit int) ([]entity.SalesOrder, error)
	GetByCode(salesOrderCode string) (*entity.SalesOrder, []entity.SalesOrderDetail, error)
}

type salesOrderRepository struct {
	db *gorm.DB
}

func NewSalesOrderRepository(db *gorm.DB) *salesOrderRepository {
	return &salesOrderRepository{db: db}
}

func (r *salesOrderRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// ReserveStock mengurangi stok tersedia secara atomik.
func (r *salesOrderRepository) ReserveStock(tx *gorm.DB, dbDetails []entity.SalesOrderDetail) error {
	now := time.Now()

	for _, detail := range dbDetails {

		var inventories []entity.Inventory

		err := tx.
			Clauses(clause.Locking{
				Strength: "UPDATE",
			}).
			Where("product_id = ?", detail.ProductID).
			Where("stock > reserved_stock").
			Order("created_at ASC").
			Order("batch_id ASC").
			Find(&inventories).Error

		if err != nil {
			return err
		}

		need := detail.Quantity
		reservations := make([]entity.StockReservation, 0)

		for _, inventory := range inventories {

			available := inventory.Stock - inventory.ReservedStock
			if available <= 0 {
				continue
			}

			reserveQty := available
			if reserveQty > need {
				reserveQty = need
			}

			result := tx.Model(&entity.Inventory{}).
				Where("batch_id = ?", inventory.BatchID).
				Update(
					"reserved_stock",
					gorm.Expr("reserved_stock + ?", reserveQty),
				)

			if result.Error != nil {
				return result.Error
			}

			reservations = append(reservations, entity.StockReservation{
				ReservedID:   uuid.NewString(),
				BatchID:      inventory.BatchID,
				DetailID:     detail.DetailID,
				SalesOrderID: detail.SalesOrderID,
				Quantity:     reserveQty,
				Status:       "RESERVED",
				CreatedAt:    now,
				UpdatedAt:    now,
			})

			need -= reserveQty

			if need == 0 {
				break
			}
		}

		if need > 0 {
			return &errorhandler.BadRequestError{
				Message: fmt.Sprintf(
					"stok produk %s tidak mencukupi",
					detail.ProductName,
				),
			}
		}

		if len(reservations) > 0 {
			if err := tx.Create(&reservations).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *salesOrderRepository) CreateOrder(tx *gorm.DB, order *entity.SalesOrder) error {
	return tx.Create(order).Error
}

func (r *salesOrderRepository) CreateOrderDetails(tx *gorm.DB, details []entity.SalesOrderDetail) error {
	if len(details) == 0 {
		return nil
	}
	return tx.Create(&details).Error
}

func (r *salesOrderRepository) GetNextSeq(now time.Time) (int, string, error) {
	yearMonth := now.Format("200601")

	query := `
		INSERT INTO sales_order_batch_sequences (base_id, year_month, last_seq)
		VALUES ('0000', ?, 1)
		ON CONFLICT (base_id, year_month) DO UPDATE
		SET last_seq = sales_order_batch_sequences.last_seq + 1
		RETURNING last_seq
	`

	var seq int
	err := r.db.Raw(query, yearMonth).Scan(&seq).Error
	if err != nil {
		return 0, "", err
	}

	return seq, yearMonth, nil
}

func (r *salesOrderRepository) GetAllPaginate(cursor *dto.Paginate, statuses []string, limit int) ([]entity.SalesOrder, error) {
	if limit <= 0 {
		limit = 10
	}

	var orders []entity.SalesOrder

	query := r.db.Model(&entity.SalesOrder{})

	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	query = applyCursorSalesOrder(query, cursor)

	err := query.Limit(limit + 1).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return reverseSalesOrderIfPrev(orders, cursor), nil
}

// GetAllByUserPaginate untuk user normal — hanya order milik userID tersebut.
// statuses opsional, kalau kosong semua status ditampilkan.
func (r *salesOrderRepository) GetAllByUserPaginate(cursor *dto.Paginate, userID string, statuses []string, limit int) ([]entity.SalesOrder, error) {
	if limit <= 0 {
		limit = 10
	}

	var orders []entity.SalesOrder

	query := r.db.Model(&entity.SalesOrder{}).
		Where("user_id = ?", userID)

	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	query = applyCursorSalesOrder(query, cursor)

	err := query.Limit(limit + 1).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return reverseSalesOrderIfPrev(orders, cursor), nil
}

// GetByCode mengambil sales order beserta detail item berdasarkan kode order.
func (r *salesOrderRepository) GetByCode(salesOrderCode string) (*entity.SalesOrder, []entity.SalesOrderDetail, error) {
	var order entity.SalesOrder
	err := r.db.Where("sales_order_code = ?", salesOrderCode).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, &errorhandler.NotFoundError{Message: "Sales order tidak ditemukan"}
		}
		return nil, nil, err
	}

	var details []entity.SalesOrderDetail
	err = r.db.
		Where("sales_order_id = ?", order.SalesOrderID).
		Order("created_at ASC").
		Find(&details).Error
	if err != nil {
		return nil, nil, err
	}

	return &order, details, nil
}

func applyCursorSalesOrder(query *gorm.DB, cursor *dto.Paginate) *gorm.DB {
	if cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev" {
		if cursor.FirstID != nil && cursor.FirstCreatedAt != nil {
			query = query.Where("(created_at, sales_order_id) > (?, ?)", cursor.FirstCreatedAt, cursor.FirstID)
			query = query.Order("created_at ASC, sales_order_id ASC")
		}
	} else {
		if cursor != nil && cursor.Direction != nil && *cursor.Direction == "next" {
			if cursor.LastCreatedAt != nil && cursor.LastID != nil {
				query = query.Where("(created_at, sales_order_id) < (?, ?)", cursor.LastCreatedAt, cursor.LastID)
			}
		}
		query = query.Order("created_at DESC, sales_order_id DESC")
	}
	return query
}

func reverseSalesOrderIfPrev(orders []entity.SalesOrder, cursor *dto.Paginate) []entity.SalesOrder {
	if cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev" {
		for i, j := 0, len(orders)-1; i < j; i, j = i+1, j-1 {
			orders[i], orders[j] = orders[j], orders[i]
		}
	}
	return orders
}
