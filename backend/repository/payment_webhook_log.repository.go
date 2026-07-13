package repository

import (
	"backend/entity"
	"backend/helper"

	"gorm.io/gorm"
)

type PaymentWebhookLogRepository interface {
	Create(log *entity.PaymentWebhookLog) error
	MarkProcessed(logID string, note string) error
}

type paymentWebhookLogRepository struct {
	db *gorm.DB
}

func NewPaymentWebhookLogRepository(db *gorm.DB) *paymentWebhookLogRepository {
	return &paymentWebhookLogRepository{db: db}
}

func (r *paymentWebhookLogRepository) Create(log *entity.PaymentWebhookLog) error {
	return r.db.Create(log).Error
}

func (r *paymentWebhookLogRepository) MarkProcessed(logID string, note string) error {
	return r.db.Model(&entity.PaymentWebhookLog{}).
		Where("log_id = ?", logID).
		Updates(map[string]interface{}{
			"processed":    true,
			"processed_at": helper.TimeNowWIB(),
			"process_note": note,
		}).Error
}
