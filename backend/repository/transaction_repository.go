package repository

import (
	"backend/entity"
	"backend/dto"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(trx *entity.Transaction) error
	GetAllByBatchID(batchID string) ([]entity.Transaction, error)
	GetAllByBatchIDPaginate(batchID string, cursor *dto.Paginate, limit int) ([]entity.Transaction, error)
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

func (r *transactionRepository) GetAllByBatchIDPaginate(batchID string, cursor *dto.Paginate, limit int) ([]entity.Transaction, error) {
	if limit <= 0 {
		limit = 5
	}

	var transactions []entity.Transaction
	query := r.db.Model(&entity.Transaction{}).Where("batch_id = ?", batchID)

	if cursor != nil {
		if cursor.Direction != nil && *cursor.Direction == "prev" {
			if cursor.FirstID != nil && cursor.FirstCreatedAt != nil {
				query = query.Where("(created_at, transaction_id) > (?, ?)", cursor.FirstCreatedAt, cursor.FirstID).
					Order("created_at ASC, transaction_id ASC")
			}
		} else {
			if cursor.Direction != nil && *cursor.Direction == "next" {
				query = query.Where("(created_at, transaction_id) < (?, ?)", cursor.LastCreatedAt, cursor.LastID)
			}
			query = query.Order("created_at DESC, transaction_id DESC")
		}
	} else {
		query = query.Order("created_at DESC, transaction_id DESC")
	}

	
	err := query.Limit(limit + 1).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	if cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev" {
		for i, j := 0, len(transactions)-1; i < j; i, j = i+1, j-1 {
			transactions[i], transactions[j] = transactions[j], transactions[i]
		}
	}

	return transactions, nil
}

//buat service dan handler nya gunanay untuk menarik transaksi sesuai batch id tapi harus paginate 

