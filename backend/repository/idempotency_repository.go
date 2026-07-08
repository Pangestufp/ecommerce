package repository

import (
	"backend/entity"
	"errors"

	"gorm.io/gorm"
)

type IdempotencyRepository interface {
	CheckoutFindByKey(tx *gorm.DB, key string) (*entity.CheckoutIdempotencyKey, error)
	CheckoutFindByKeyForPolling(key string, userID string) (*entity.CheckoutIdempotencyKey, error)
	CheckoutCreate(tx *gorm.DB, record *entity.CheckoutIdempotencyKey) error
	CheckoutUpsertStatus(record *entity.CheckoutIdempotencyKey) error
}

type idempotencyRepository struct {
	db *gorm.DB
}

func NewIdempotencyRepository(db *gorm.DB) *idempotencyRepository {
	return &idempotencyRepository{db: db}
}

func (r *idempotencyRepository) CheckoutFindByKey(tx *gorm.DB, key string) (*entity.CheckoutIdempotencyKey, error) {
	var record entity.CheckoutIdempotencyKey
	err := tx.Where("idempotency_key = ?", key).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *idempotencyRepository) CheckoutFindByKeyForPolling(key string, userID string) (*entity.CheckoutIdempotencyKey, error) {
	var record entity.CheckoutIdempotencyKey
	err := r.db.
		Where("idempotency_key = ? AND user_id = ?", key, userID).
		First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *idempotencyRepository) CheckoutCreate(tx *gorm.DB, record *entity.CheckoutIdempotencyKey) error {
	return tx.Create(record).Error
}

// UpsertStatus pakai OnConflict supaya aman dipanggil berulang kali
// (misal worker retry update status FAILED).
func (r *idempotencyRepository) CheckoutUpsertStatus(record *entity.CheckoutIdempotencyKey) error {
	return r.db.
		Where("idempotency_key = ?", record.IdempotencyKey).
		Assign(map[string]any{
			"status":        record.Status,
			"error_message": record.ErrorMessage,
			"response_body": record.ResponseBody,
			"status_code":   record.StatusCode,
			"updated_at":    record.UpdatedAt,
		}).
		FirstOrCreate(record).Error
}
