package repository

import (
	"backend/entity"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(trx *entity.Transaction) error
	GetAllByBatchID(batchID string) ([]entity.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(trx *entity.Transaction) error {
	return r.db.Create(trx).Error
}

func (r *transactionRepository) GetAllByBatchID(batchID string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.Where("batch_id = ?", batchID).
		Order("created_at DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
