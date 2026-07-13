package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	PaymentID             string     `gorm:"primaryKey"`
	SalesOrderID          string     // FK sales_orders
	MidtransOrderID       string     // order_id yang dikirim ke Midtrans
	SnapToken             string     // token untuk window.snap.pay()
	RedirectURL           string     // URL untuk flow redirect (mobile, dll)
	ExpiredAt             *time.Time // expiry yang di-set ke Midtrans waktu buat token
	PaymentType           string     // istilah Midtrans: bank_transfer, gopay,
	PaymentChannel        string     // lebih spesifik: bca, bni, bri, mandiri,
	VANumber              *string    // virtual account number untuk bank transfer
	PaymentCode           *string    // kode bayar untuk convenience store (Indomaret dll)
	BillerKey             *string    // khusus Mandiri Bill — butuh dua kode
	BillKey               *string
	GrossAmount           decimal.Decimal // jumlah yang harus dibayar, dari Midtrans
	Currency              string          // IDR — simpan eksplisit untuk future-proof
	TransactionStatus     string          // status dari Midtrans: pending, capture, settlement, deny, expire, cancel
	FraudStatus           *string         // hanya ada untuk credit card: accept, challenge, deny
	MidtransTransactionID *string         // transaction_id dari Midtrans
	SettlementTime        *time.Time      // waktu dana settle, dari Midtrans
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
