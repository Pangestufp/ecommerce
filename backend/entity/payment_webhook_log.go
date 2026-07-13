package entity

import "time"

type PaymentWebhookLog struct {
	LogID             string `gorm:"primaryKey"`
	SalesOrderID      string
	PaymentID         *string
	MidtransOrderID   string
	RawPayload        string     // raw JSON dari Midtrans, disimpan as-is
	TransactionStatus string     // di-extract dari payload untuk query mudah
	SignatureValid    bool       // hasil verifikasi signature Midtrans
	Processed         bool       // apakah webhook ini berhasil diproses
	ProcessedAt       *time.Time // kapan diproses
	ProcessNote       string     // catatan hasil proses atau error message
	IPAddress         string
	CreatedAt         time.Time
}
