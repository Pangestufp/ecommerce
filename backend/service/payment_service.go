package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreateSnapToken(req *dto.CreatePaymentRequest, userID string, idempotencyKey string) (*dto.CreatePaymentResponse, error)
	HandleWebhook(notification *dto.MidtransNotification, rawPayload string, ipAddress string) error
	GetPaymentStatus(salesOrderCode string, userID string) (*dto.PaymentStatusResponse, error)
}

type paymentService struct {
	paymentRepo        repository.PaymentRepository
	webhookLogRepo     repository.PaymentWebhookLogRepository
	salesOrderRepo     repository.SalesOrderRepository
	historyRepo        repository.OrderStatusHistoryRepository
	inventoryRepo      repository.InventoryRepository
	redis              *redis.Client
	midtransProduction bool
	sandBoxSnapURL     string
	productionSnapURL  string
	serverKey          string
	envPro             bool
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	webhookLogRepo repository.PaymentWebhookLogRepository,
	salesOrderRepo repository.SalesOrderRepository,
	historyRepo repository.OrderStatusHistoryRepository,
	inventoryRepo repository.InventoryRepository,
	redis *redis.Client,
	midtransProduction bool,
	sandBoxSnapURL string,
	productionSnapURL string,
	serverKey string,
	envPro bool,
) PaymentService {
	return &paymentService{
		paymentRepo:        paymentRepo,
		webhookLogRepo:     webhookLogRepo,
		salesOrderRepo:     salesOrderRepo,
		historyRepo:        historyRepo,
		inventoryRepo:      inventoryRepo,
		redis:              redis,
		midtransProduction: midtransProduction,
		sandBoxSnapURL:     sandBoxSnapURL,
		productionSnapURL:  productionSnapURL,
		serverKey:          serverKey,
		envPro:             envPro,
	}
}

func (s *paymentService) CreateSnapToken(req *dto.CreatePaymentRequest, userID string, idempotencyKey string) (*dto.CreatePaymentResponse, error) {
	if idempotencyKey == "" {
		return nil, &errorhandler.BadRequestError{Message: "Header Idempotency-Key wajib diisi"}
	}

	ctx := context.Background()
	statusKey := fmt.Sprintf("idempotency:payment:%s", idempotencyKey)

	if cached, err := s.redis.Get(ctx, statusKey).Result(); err == nil {
		if cached == "PROCESSING" {
			return nil, &errorhandler.BadRequestError{Message: "Request dengan Idempotency-Key ini sedang diproses, coba lagi sebentar"}
		}
		var resp dto.CreatePaymentResponse
		if json.Unmarshal([]byte(cached), &resp) == nil {
			return &resp, nil
		}
	}

	acquired, err := s.redis.SetNX(ctx, statusKey, "PROCESSING", 30*time.Second).Result()
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal memvalidasi idempotency key"}
	}
	if !acquired {
		return nil, &errorhandler.BadRequestError{Message: "Request dengan Idempotency-Key ini sedang diproses, coba lagi sebentar"}
	}

	succeeded := false
	defer func() {
		if !succeeded {
			s.redis.Del(ctx, statusKey)
		}
	}()

	order, details, err := s.salesOrderRepo.GetByCode(req.SalesOrderCode)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, &errorhandler.ForbiddenError{Message: "Akses tidak diizinkan"}
	}
	if order.Status != helper.GetPendingPaymentStatus() {
		return nil, &errorhandler.BadRequestError{Message: "Order tidak dalam status menunggu pembayaran"}
	}

	// Kalau snap token sudah ada dan belum expired, kembalikan yang lama
	existing, err := s.paymentRepo.GetBySalesOrderID(order.SalesOrderID)
	if err == nil && existing != nil {
		if existing.ExpiredAt != nil && existing.ExpiredAt.After(helper.TimeNowWIB()) {
			resp := &dto.CreatePaymentResponse{
				PaymentID:   existing.PaymentID,
				SnapToken:   existing.SnapToken,
				RedirectURL: existing.RedirectURL,
				ExpiredAt:   existing.ExpiredAt,
			}
			s.cacheIdempotencyResult(ctx, statusKey, resp, time.Until(*existing.ExpiredAt))
			return resp, nil
		}
	}

	// Build item details untuk Midtrans
	itemDetails := make([]map[string]interface{}, 0, len(details)+2)
	for _, d := range details {
		itemDetails = append(itemDetails, map[string]interface{}{
			"id":       d.ProductID,
			"price":    d.FinalUnitPrice.IntPart(),
			"quantity": d.Quantity,
			"name":     d.ProductName,
		})
	}
	itemDetails = append(itemDetails, map[string]interface{}{
		"id":       "SHIPPING",
		"price":    order.ShippingFee.IntPart(),
		"quantity": 1,
		"name":     fmt.Sprintf("Ongkir %s", order.ShippingDisplayName),
	})

	// Expiry sinkron dengan order 1 jam dari sekarang
	expiredAt := helper.TimeNowWIB().Add(1 * time.Hour)

	midtransPayload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     order.SalesOrderCode,
			"gross_amount": order.FinalTotal.IntPart(),
		},
		"item_details": itemDetails,
		"customer_details": map[string]interface{}{
			"first_name": order.CustomerName,
			"email":      order.CustomerEmail,
			"phone":      order.CustomerPhone,
		},
		"expiry": map[string]interface{}{
			"start_time": expiredAt.Format("2006-01-02 15:04:05 +0700"),
			"unit":       "hours",
			"duration":   24,
		},
	}

	snapToken, redirectURL, err := s.hitMidtransSnap(midtransPayload)
	if err != nil {
		return nil, err
	}

	now := helper.TimeNowWIB()
	payment := &entity.Payment{
		PaymentID:         uuid.NewString(),
		SalesOrderID:      order.SalesOrderID,
		MidtransOrderID:   order.SalesOrderCode,
		SnapToken:         snapToken,
		RedirectURL:       redirectURL,
		ExpiredAt:         &expiredAt,
		GrossAmount:       order.FinalTotal,
		Currency:          "IDR",
		TransactionStatus: "pending",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.paymentRepo.Transaction(func(tx *gorm.DB) error {
		return s.paymentRepo.Create(tx, payment)
	}); err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal menyimpan data payment"}
	}

	resp := &dto.CreatePaymentResponse{
		PaymentID:   payment.PaymentID,
		SnapToken:   snapToken,
		RedirectURL: redirectURL,
		ExpiredAt:   &expiredAt,
	}

	s.cacheIdempotencyResult(ctx, statusKey, resp, time.Until(expiredAt))
	s.SimulateWebhookIfDev(order.SalesOrderCode)
	return resp, nil
}

func (s *paymentService) cacheIdempotencyResult(ctx context.Context, key string, resp *dto.CreatePaymentResponse, ttl time.Duration) {
	encoded, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[CreateSnapToken] gagal marshal response key=%s: %v", key, err)
		return
	}
	if err := s.redis.Set(ctx, key, string(encoded), ttl).Err(); err != nil {
		log.Printf("[CreateSnapToken] gagal cache idempotency key=%s: %v", key, err)
	}
}

func (s *paymentService) HandleWebhook(notification *dto.MidtransNotification, rawPayload string, ipAddress string) error {
	now := helper.TimeNowWIB()

	signatureValid := verifyMidtransSignature(
		notification.OrderID,
		notification.StatusCode,
		notification.GrossAmount,
		notification.SignatureKey,
		s.serverKey,
		s.envPro,
	)

	payment, _ := s.paymentRepo.GetByMidtransOrderID(notification.OrderID)

	webhookLog := &entity.PaymentWebhookLog{
		LogID:             uuid.NewString(),
		MidtransOrderID:   notification.OrderID,
		RawPayload:        rawPayload,
		TransactionStatus: notification.TransactionStatus,
		SignatureValid:    signatureValid,
		Processed:         false,
		IPAddress:         ipAddress,
		CreatedAt:         now,
	}
	if payment != nil {
		webhookLog.SalesOrderID = payment.SalesOrderID
		webhookLog.PaymentID = &payment.PaymentID
	}
	if err := s.webhookLogRepo.Create(webhookLog); err != nil {
		log.Printf("[Webhook] gagal simpan log order_id=%s: %v", notification.OrderID, err)
	}

	if !signatureValid {
		log.Printf("[Webhook] signature tidak valid order_id=%s ip=%s", notification.OrderID, ipAddress)
		return nil
	}

	if payment == nil {
		log.Printf("[Webhook] payment tidak ditemukan order_id=%s", notification.OrderID)
		return nil
	}

	// Resolve aksi dari status Midtrans
	action, shouldProcess := resolveAction(notification.TransactionStatus, notification.FraudStatus)
	if !shouldProcess {
		_ = s.webhookLogRepo.MarkProcessed(webhookLog.LogID,
			fmt.Sprintf("status %s tidak perlu diproses", notification.TransactionStatus))
		return nil
	}

	// Semua update dalam satu transaksi DB
	processErr := s.paymentRepo.Transaction(func(tx *gorm.DB) error {
		// Update payment record
		if err := s.paymentRepo.UpdateByMidtransOrderID(tx, notification.OrderID,
			buildPaymentUpdates(notification, now),
		); err != nil {
			return err
		}

		order, _, err := s.salesOrderRepo.GetByCode(notification.OrderID)
		if err != nil {
			return err
		}

		nextStatus, err := helper.ApplyAction(order.Status, action)
		if err != nil {
			log.Printf("webhook aksi %s tidak valid untuk status %s order=%s",
				action, order.Status, notification.OrderID)
			return nil
		}

		// Update sales order status
		orderUpdates := map[string]interface{}{
			"status":     nextStatus,
			"updated_at": now,
		}

		if action == helper.GetPayAction() {
			orderUpdates["paid_at"] = now
		} else if action == helper.GetExpireAction() {
			orderUpdates["cancelled_at"] = now
		}

		if action == helper.GetPayAction() {
			if err := s.inventoryRepo.CommitReservations(tx, order.SalesOrderID, order.SalesOrderCode, now); err != nil {
				return err
			}
		}

		if err := s.salesOrderRepo.UpdateStatus(tx, order.SalesOrderID, orderUpdates); err != nil {
			return err
		}

		// Catat history
		prev := order.Status
		return s.historyRepo.Create(tx, &entity.OrderStatusHistory{
			HistoryID:      uuid.NewString(),
			SalesOrderID:   order.SalesOrderID,
			PreviousStatus: &prev,
			Status:         nextStatus,
			Note:           fmt.Sprintf("Notifikasi Midtrans: %s", notification.TransactionStatus),
			CreatedBy:      "SYSTEM",
			CreatedName:    "System",
			CreatedAt:      now,
		})
	})

	if processErr != nil {
		log.Printf("[Webhook] gagal proses order_id=%s: %v", notification.OrderID, processErr)
		_ = s.webhookLogRepo.MarkProcessed(webhookLog.LogID,
			fmt.Sprintf("error: %v", processErr))
		return processErr
	}

	_ = s.webhookLogRepo.MarkProcessed(webhookLog.LogID, "berhasil diproses")
	return nil
}

func (s *paymentService) GetPaymentStatus(salesOrderCode string, userID string) (*dto.PaymentStatusResponse, error) {
	order, _, err := s.salesOrderRepo.GetByCode(salesOrderCode)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, &errorhandler.ForbiddenError{Message: "Akses tidak diizinkan"}
	}

	payment, err := s.paymentRepo.GetBySalesOrderID(order.SalesOrderID)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentStatusResponse{
		SalesOrderCode:    salesOrderCode,
		TransactionStatus: payment.TransactionStatus,
		PaymentType:       payment.PaymentType,
		PaymentChannel:    payment.PaymentChannel,
		GrossAmount:       payment.GrossAmount,
		GrossAmountFormat: helper.FormatRupiah(payment.GrossAmount),
		VANumber:          payment.VANumber,
		PaymentCode:       payment.PaymentCode,
		BillerKey:         payment.BillerKey,
		BillKey:           payment.BillKey,
		SnapToken:         payment.SnapToken,
		RedirectURL:       payment.RedirectURL,
		ExpiredAt:         payment.ExpiredAt,
	}, nil
}

func (s *paymentService) hitMidtransSnap(payload map[string]interface{}) (string, string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Midtrans Snap] gagal marshal payload: %v", err)
		return "", "", &errorhandler.InternalServerError{Message: "Gagal marshal payload"}
	}

	// LOG 1 — pastikan payload yang dikirim benar
	log.Printf("[Midtrans Snap] payload dikirim: %s", string(body))

	snapURL := s.sandBoxSnapURL
	if s.midtransProduction {
		snapURL = s.productionSnapURL
	}

	// LOG 2 — pastikan URL yang dipakai benar
	log.Printf("[Midtrans Snap] hit URL: %s", snapURL)

	req, err := http.NewRequest(http.MethodPost, snapURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[Midtrans Snap] gagal buat request: %v", err)
		return "", "", &errorhandler.InternalServerError{Message: "Gagal membuat request ke Midtrans"}
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.serverKey + ":"))

	// LOG 3 — pastikan server key tidak kosong
	log.Printf("[Midtrans Snap] server key prefix: %.10s... (panjang: %d)", s.serverKey, len(s.serverKey))

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Midtrans Snap] gagal hit Midtrans: %v", err)
		return "", "", &errorhandler.InternalServerError{Message: "Gagal menghubungi Midtrans"}
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// LOG 4 — lihat mentah response dari Midtrans
	log.Printf("[Midtrans Snap] status=%d response body=%s", resp.StatusCode, string(respBody))

	if resp.StatusCode != http.StatusCreated {
		log.Printf("[Midtrans Snap] status bukan 201, ditolak")
		return "", "", &errorhandler.InternalServerError{Message: "Midtrans menolak request"}
	}

	var result struct {
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("[Midtrans Snap] gagal parse response: %v | body: %s", err, string(respBody))
		return "", "", &errorhandler.InternalServerError{Message: "Gagal parse response Midtrans"}
	}

	// LOG 5 — pastikan hasil parse benar
	log.Printf("[Midtrans Snap] token=%s redirect_url=%s", result.Token, result.RedirectURL)

	return result.Token, result.RedirectURL, nil
}

// verifyMidtransSignature memverifikasi signature dari notifikasi Midtrans.
// Formula: SHA512(order_id + status_code + gross_amount + server_key)
func verifyMidtransSignature(orderID, statusCode, grossAmount, receivedSignature, serverKey string, isPro bool) bool {
	if !isPro {
		log.Printf("[Signature] bypass aktif untuk order=%s (dev only)", orderID)
		return true
	}

	raw := orderID + statusCode + grossAmount + serverKey
	hash := sha512.Sum512([]byte(raw))
	computed := fmt.Sprintf("%x", hash)
	return computed == receivedSignature
}

// resolveAction memetakan status Midtrans ke action di state machine order.
func resolveAction(transactionStatus, fraudStatus string) (string, bool) {
	switch transactionStatus {
	case "capture":
		if fraudStatus == "accept" {
			return helper.GetPayAction(), true
		}
		return "", false
	case "settlement":
		return helper.GetPayAction(), true
	case "expire":
		return helper.GetExpireAction(), true
	case "pending", "deny", "cancel":
		return "", false
	default:
		return "", false
	}
}

// buildPaymentUpdates menyiapkan kolom yang perlu di-update di tabel payments.
func buildPaymentUpdates(n *dto.MidtransNotification, now time.Time) map[string]interface{} {
	updates := map[string]interface{}{
		"transaction_status":      n.TransactionStatus,
		"midtrans_transaction_id": n.TransactionID,
		"payment_type":            n.PaymentType,
		"updated_at":              now,
	}

	if n.FraudStatus != "" {
		updates["fraud_status"] = n.FraudStatus
	}

	if len(n.VANumbers) > 0 {
		updates["va_number"] = n.VANumbers[0].VANumber
		updates["payment_channel"] = n.VANumbers[0].Bank
	}
	if n.PaymentCode != "" {
		updates["payment_code"] = n.PaymentCode
		updates["payment_channel"] = n.Store
	}
	if n.BillerKey != "" {
		updates["biller_key"] = n.BillerKey
		updates["bill_key"] = n.BillKey
		updates["payment_channel"] = "mandiri"
	}

	if n.SettlementTime != "" {
		loc := time.FixedZone("WIB", 7*3600)
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", n.SettlementTime, loc); err == nil {
			updates["settlement_time"] = t
		}
	}

	return updates
}

func (s *paymentService) SimulateWebhookIfDev(salesOrderCode string) {
	if s.envPro {
		return
	}

	go func() {
		// Jeda 5 detik — simulasi user "bayar" di Snap UI
		time.Sleep(5 * time.Second)

		fakeTransactionID := uuid.NewString()
		fakeGrossAmount := "0" // tidak dipakai untuk verify, simulasi saja

		notification := &dto.MidtransNotification{
			OrderID:           salesOrderCode,
			TransactionID:     fakeTransactionID,
			TransactionStatus: "settlement",
			FraudStatus:       "",
			PaymentType:       "simulation",
			GrossAmount:       fakeGrossAmount,
			Currency:          "IDR",
			SignatureKey:      "bypass", // akan di-bypass di verifyMidtransSignature
			SettlementTime:    helper.TimeNowWIB().Format("2006-01-02 15:04:05"),
		}

		// Marshal jadi raw payload seperti aslinya dari Midtrans
		raw, _ := json.Marshal(notification)

		log.Printf("[SimulateWebhook] menjalankan simulasi settlement untuk order=%s", salesOrderCode)

		if err := s.HandleWebhook(notification, string(raw), "127.0.0.1"); err != nil {
			log.Printf("[SimulateWebhook] gagal: %v", err)
		}
	}()
}
