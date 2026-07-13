package repository

import (
	"backend/entity"
	"backend/errorhandler"
	"errors"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(tx *gorm.DB, payment *entity.Payment) error
	UpdateByMidtransOrderID(tx *gorm.DB, midtransOrderID string, updates map[string]interface{}) error
	GetBySalesOrderID(salesOrderID string) (*entity.Payment, error)
	GetByMidtransOrderID(midtransOrderID string) (*entity.Payment, error)
	Transaction(fn func(tx *gorm.DB) error) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *paymentRepository) Create(tx *gorm.DB, payment *entity.Payment) error {
	return tx.Create(payment).Error
}

// UpdateByMidtransOrderID dipakai webhook handler untuk update status payment
// setelah notifikasi masuk dari Midtrans.
func (r *paymentRepository) UpdateByMidtransOrderID(tx *gorm.DB, midtransOrderID string, updates map[string]interface{}) error {
	result := tx.Model(&entity.Payment{}).
		Where("midtrans_order_id = ?", midtransOrderID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &errorhandler.NotFoundError{Message: "Payment tidak ditemukan"}
	}
	return nil
}

func (r *paymentRepository) GetBySalesOrderID(salesOrderID string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.Where("sales_order_id = ?", salesOrderID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Payment tidak ditemukan"}
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByMidtransOrderID(midtransOrderID string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.Where("midtrans_order_id = ?", midtransOrderID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Payment tidak ditemukan"}
		}
		return nil, err
	}
	return &payment, nil
}
