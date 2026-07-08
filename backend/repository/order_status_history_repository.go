package repository

import (
	"backend/entity"

	"gorm.io/gorm"
)

type OrderStatusHistoryRepository interface {
	Create(tx *gorm.DB, history *entity.OrderStatusHistory) error
	FindBySalesOrderID(salesOrderID string) ([]entity.OrderStatusHistory, error)
}

type orderStatusHistoryRepository struct {
	db *gorm.DB
}

func NewOrderStatusHistoryRepository(db *gorm.DB) OrderStatusHistoryRepository {
	return &orderStatusHistoryRepository{db: db}
}

func (r *orderStatusHistoryRepository) Create(tx *gorm.DB, history *entity.OrderStatusHistory) error {
	return tx.Create(history).Error
}

func (r *orderStatusHistoryRepository) FindBySalesOrderID(salesOrderID string) ([]entity.OrderStatusHistory, error) {
	var histories []entity.OrderStatusHistory
	err := r.db.
		Where("sales_order_id = ?", salesOrderID).
		Order("created_at ASC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}
