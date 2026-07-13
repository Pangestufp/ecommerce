package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// ─── Request ──────────────────────────────────────────────────────────────────

// CreatePaymentRequest dikirim frontend waktu user mau bayar.
// Cukup SalesOrderCode — backend yang generate snap token ke Midtrans.
type CreatePaymentRequest struct {
	SalesOrderCode string `json:"sales_order_code" binding:"required"`
}

// MidtransNotification adalah shape payload webhook dari Midtrans.
// Field ini sesuai dengan dokumentasi Midtrans notification payload.
type MidtransNotification struct {
	OrderID           string `json:"order_id"`
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
	GrossAmount       string `json:"gross_amount"`
	Currency          string `json:"currency"`
	SignatureKey      string `json:"signature_key"`
	StatusCode        string `json:"status_code"`

	// Bank transfer
	VANumbers []struct {
		Bank     string `json:"bank"`
		VANumber string `json:"va_number"`
	} `json:"va_numbers"`

	// Convenience store
	PaymentCode string `json:"payment_code"`
	Store       string `json:"store"`

	// Mandiri Bill
	BillerKey string `json:"biller_key"`
	BillKey   string `json:"bill_key"`

	// Timestamps
	SettlementTime string `json:"settlement_time"`
	ExpiryTime     string `json:"expiry_time"`
}

// ─── Response ─────────────────────────────────────────────────────────────────

// CreatePaymentResponse dikembalikan ke frontend setelah snap token berhasil dibuat.
type CreatePaymentResponse struct {
	PaymentID   string     `json:"payment_id"`
	SnapToken   string     `json:"snap_token"`
	RedirectURL string     `json:"redirect_url"`
	ExpiredAt   *time.Time `json:"expired_at"`
}

// PaymentStatusResponse untuk polling status dari frontend.
type PaymentStatusResponse struct {
	SalesOrderCode    string          `json:"sales_order_code"`
	TransactionStatus string          `json:"transaction_status"`
	PaymentType       string          `json:"payment_type"`
	PaymentChannel    string          `json:"payment_channel"`
	GrossAmount       decimal.Decimal `json:"gross_amount"`
	GrossAmountFormat string          `json:"gross_amount_format"`

	// Info bayar — hanya terisi sesuai metode pembayaran
	VANumber    *string `json:"va_number"`
	PaymentCode *string `json:"payment_code"`
	BillerKey   *string `json:"biller_key"`
	BillKey     *string `json:"bill_key"`

	SnapToken   string     `json:"snap_token"`
	RedirectURL string     `json:"redirect_url"`
	ExpiredAt   *time.Time `json:"expired_at"`
}
