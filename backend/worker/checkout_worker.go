package worker

import (
	"backend/config"
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CheckoutJob struct {
	Req            *dto.ConfirmCheckoutRequest
	UserID         string
	IdempotencyKey string
}

type CheckoutWorker struct {
	jobCh           chan CheckoutJob
	salesOrderRepo  repository.SalesOrderRepository
	idempotencyRepo repository.IdempotencyRepository
	productRepo     repository.ProductRepository
	discountRepo    repository.DiscountRepository
	priceRepo       repository.ProductPriceRepository
	addressRepo     repository.AddressRepository
	storeConfigRepo repository.StoreConfigRepository
	userRepo        repository.UserRepository
	historyRepo     repository.OrderStatusHistoryRepository
	redis           *redis.Client
}

var (
	InstanceCheckOut *CheckoutWorker
	onceCheckOut     sync.Once
)

// Initialize membuat singleton CheckoutWorker dan langsung menjalankan
// consumer goroutine-nya. Cukup dipanggil sekali di main/bootstrap, contoh:
// Setelah itu di service manapun tinggal panggil worker.Instance.Enqueue(job).
func InitializeCheckout(db *gorm.DB) {
	onceCheckOut.Do(func() {
		ProductRepository := repository.NewProductRepository(db)
		IdempotencyRepository := repository.NewIdempotencyRepository(db)
		DiscountRepository := repository.NewDiscountRepository(db)
		ProductPriceRepository := repository.NewProductPriceRepository(db)
		UserAddressRepository := repository.NewAddressRepository(db)
		StoreConfigRepository := repository.NewStoreConfigRepository(db)
		SalesOrderRepository := repository.NewSalesOrderRepository(db)
		UserRepository := repository.NewUserRepository(db)
		OrderStatusHistoryRepository := repository.NewOrderStatusHistoryRepository(db)

		InstanceCheckOut = &CheckoutWorker{
			jobCh:           make(chan CheckoutJob, 100),
			salesOrderRepo:  SalesOrderRepository,
			idempotencyRepo: IdempotencyRepository,
			productRepo:     ProductRepository,
			discountRepo:    DiscountRepository,
			priceRepo:       ProductPriceRepository,
			addressRepo:     UserAddressRepository,
			storeConfigRepo: StoreConfigRepository,
			userRepo:        UserRepository,
			historyRepo:     OrderStatusHistoryRepository,
			redis:           config.RedisClient,
		}
		InstanceCheckOut.start()
	})
}

// start menjalankan consumer goroutine. Dipanggil otomatis oleh Initialize,
// tidak perlu dipanggil manual dari luar package.
func (w *CheckoutWorker) start() {
	go func() {
		log.Println("[CheckoutWorker] started")
		for job := range w.jobCh {
			w.process(job)
		}
		log.Println("[CheckoutWorker] channel closed, shutdown")
	}()
}

// Shutdown menutup channel job. Panggil ini saat graceful shutdown aplikasi
// (opsional kalau tidak dipanggil, goroutine akan ikut mati saat proses exit).
func (w *CheckoutWorker) Shutdown(ctx context.Context) {
	close(w.jobCh)
}

// ─────────────────────────────────────────────
// ENQUEUE
// ─────────────────────────────────────────────

func (w *CheckoutWorker) Enqueue(job CheckoutJob) error {
	select {
	case w.jobCh <- job:
		return nil
	default:
		return &errorhandler.InternalServerError{
			Message: "Server sedang sibuk, silakan coba beberapa saat lagi",
		}
	}
}

// ─────────────────────────────────────────────
// PROCESS
// ─────────────────────────────────────────────

func (w *CheckoutWorker) process(job CheckoutJob) {
	ctx := context.Background()
	idemKey := job.IdempotencyKey
	statusKey := fmt.Sprintf("idempotency:status:%s", idemKey)

	log.Printf("[CheckoutWorker] processing idempotency_key=%s user=%s", idemKey, job.UserID)

	w.setRedisStatus(ctx, statusKey, dto.CheckoutStatusResponse{
		IdempotencyKey: idemKey,
		Status:         "PROCESSING",
	})

	orderSummary, err := w.execute(ctx, job)
	if err != nil {
		log.Printf("[CheckoutWorker] failed idempotency_key=%s err=%v", idemKey, err)

		// Simpan FAILED ke Redis
		w.setRedisStatus(ctx, statusKey, dto.CheckoutStatusResponse{
			IdempotencyKey: idemKey,
			Status:         "FAILED",
			ErrorMessage:   err.Error(),
		})

		// Best-effort update ke DB supaya state FAILED tidak hilang setelah restart
		now := helper.TimeNowWIB()
		failedRecord := &entity.CheckoutIdempotencyKey{
			IdempotencyID:  uuid.New().String(),
			IdempotencyKey: idemKey,
			UserID:         job.UserID,
			Status:         "FAILED",
			ErrorMessage:   err.Error(),
			StatusCode:     422,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if upsertErr := w.idempotencyRepo.CheckoutUpsertStatus(failedRecord); upsertErr != nil {
			log.Printf("[CheckoutWorker] gagal upsert FAILED ke DB idempotency_key=%s err=%v", idemKey, upsertErr)
		}
		return
	}

	successPayload := dto.CheckoutStatusResponse{
		IdempotencyKey: idemKey,
		Status:         "SUCCESS",
		Order:          orderSummary,
	}
	w.setRedisStatus(ctx, statusKey, successPayload)

	log.Printf("[CheckoutWorker] done idempotency_key=%s order=%s", idemKey, orderSummary.SalesOrderID)
}

func (w *CheckoutWorker) execute(ctx context.Context, job CheckoutJob) (*dto.OrderSummaryResponse, error) {
	req := job.Req
	userID := job.UserID
	idemKey := job.IdempotencyKey
	now := helper.TimeNowWIB()

	statusKey := fmt.Sprintf("idempotency:status:%s", idemKey)
	if cached, err := w.redis.Get(ctx, statusKey).Result(); err == nil {
		var existing dto.CheckoutStatusResponse
		if json.Unmarshal([]byte(cached), &existing) == nil && existing.Status == "SUCCESS" {
			return existing.Order, nil
		}
	}

	cacheKey := fmt.Sprintf("checkout:%s", req.CheckoutID)
	sessionRaw, err := w.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, &errorhandler.NotFoundError{Message: "Sesi checkout tidak ditemukan atau sudah kedaluwarsa"}
	}
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal membaca sesi checkout"}
	}

	var session dto.CheckoutRedisData
	if err := json.Unmarshal([]byte(sessionRaw), &session); err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Data sesi checkout tidak valid"}
	}
	if session.UserID != userID {
		return nil, &errorhandler.BadRequestError{Message: "Akses tidak diizinkan"}
	}

	user, err := w.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "User Invalid"}
	}

	productIDs := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	enriched, err := w.productRepo.GetProductsEnrichedBatch(productIDs)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal mengambil data produk"}
	}
	productMap := make(map[string]*dto.ProductEnrichedForES)
	for _, p := range enriched {
		productMap[p.ProductID] = p
	}

	allPrices, err := w.priceRepo.GetLatestByProductIDs(productIDs)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal mengambil harga produk"}
	}
	priceMap := make(map[string]entity.ProductPrice)
	for _, p := range allPrices {
		if _, exists := priceMap[p.ProductID]; !exists {
			priceMap[p.ProductID] = p
		}
	}

	allDiscounts, _ := w.discountRepo.GetActiveDiscountsByProductIDs(productIDs, now)
	discountMap := make(map[string][]entity.Discount)
	for _, d := range allDiscounts {
		discountMap[d.ProductID] = append(discountMap[d.ProductID], d)
	}

	weightMap := make(map[string]int)
	for _, p := range enriched {
		weightMap[p.ProductID] = p.WeightGram
	}

	// ── Ambil alamat dan konfigurasi toko ────────────────────────────────
	address, err := w.addressRepo.GetAddressByIDAndUserID(req.AddressID, userID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal mengambil alamat"}
	}
	if address == nil {
		return nil, &errorhandler.NotFoundError{Message: "Alamat tidak ditemukan"}
	}

	storeConfig, err := w.storeConfigRepo.GetConfig()
	if err != nil {
		return nil, err
	}

	salesOrderID := uuid.New().String()
	seq, yearMonth, err := w.salesOrderRepo.GetNextSeq(now)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Fail to make code"}
	}

	salesOrderCode := fmt.Sprintf("0000-%s-%06d", yearMonth, seq)
	confirmedAt := now

	var dbDetails []entity.SalesOrderDetail
	var respItems []dto.OrderItemResponse
	totalBeforeDiscount := decimal.Zero
	totalDiscount := decimal.Zero
	totalWeight := 0

	for _, reqItem := range req.Items {
		product, ok := productMap[reqItem.ProductID]
		if !ok || product.Available == 0 {
			return nil, &errorhandler.BadRequestError{
				Message: fmt.Sprintf("Produk %s sudah tidak tersedia", reqItem.ProductID),
			}
		}
		if reqItem.Qty > int(product.AvailableStock) {
			return nil, &errorhandler.BadRequestError{
				Message: fmt.Sprintf("Stok produk %s tidak mencukupi", product.ProductName),
			}
		}

		price, hasPrice := priceMap[reqItem.ProductID]
		if !hasPrice {
			return nil, &errorhandler.BadRequestError{
				Message: fmt.Sprintf("Harga produk %s tidak ditemukan", product.ProductName),
			}
		}

		unitPrice := price.ProductPrice
		discountPerUnit := decimal.Zero
		var discountID *string

		if reqItem.DiscountID != nil {
			for _, d := range discountMap[reqItem.ProductID] {
				if d.DiscountID == *reqItem.DiscountID {
					var calc decimal.Decimal
					if d.DiscountType == helper.Percentage() {
						calc = unitPrice.Mul(d.DiscountValue)
					} else {
						calc = d.DiscountValue
					}
					if !calc.IsNegative() && calc.LessThanOrEqual(unitPrice) {
						discountPerUnit = calc
						id := d.DiscountID
						discountID = &id
					}
					break
				}
			}
		}

		finalUnitPrice := unitPrice.Sub(discountPerUnit)
		qtyDec := decimal.NewFromInt(int64(reqItem.Qty))
		lineDiscount := discountPerUnit.Mul(qtyDec)
		lineBeforeDiscount := unitPrice.Mul(qtyDec)
		lineFinalTotal := finalUnitPrice.Mul(qtyDec)

		productWeight := weightMap[reqItem.ProductID]
		lineWeight := productWeight * reqItem.Qty

		totalWeight += lineWeight
		totalBeforeDiscount = totalBeforeDiscount.Add(lineBeforeDiscount)
		totalDiscount = totalDiscount.Add(lineDiscount)

		dbDetails = append(dbDetails, entity.SalesOrderDetail{
			DetailID:            uuid.New().String(),
			SalesOrderID:        salesOrderID,
			ProductID:           product.ProductID,
			ImageURL:            product.PrimaryImage,
			ProductName:         product.ProductName,
			Quantity:            reqItem.Qty,
			Weight:              productWeight,
			TotalWeight:         lineWeight,
			PriceBeforeDiscount: unitPrice,
			DiscountID:          discountID,
			DiscountAmount:      lineDiscount,
			FinalUnitPrice:      finalUnitPrice,
			FinalTotal:          lineFinalTotal,
			CreatedAt:           now,
			UpdatedAt:           now,
		})

		respItems = append(respItems, dto.OrderItemResponse{
			ProductID:      product.ProductID,
			ProductName:    product.ProductName,
			ImageURL:       product.PrimaryImage,
			Quantity:       reqItem.Qty,
			UnitPrice:      unitPrice,
			DiscountAmount: lineDiscount,
			FinalUnitPrice: finalUnitPrice,
			FinalTotal:     lineFinalTotal,
		})
	}

	finalTotal := totalBeforeDiscount.Sub(totalDiscount).Add(req.ShippingCost)

	subtotal := totalBeforeDiscount.Sub(totalDiscount)

	if !req.Subtotal.Equal(subtotal) {
		return nil, &errorhandler.BadRequestError{
			Message: "Total harga berubah ditemukan",
		}
	}

	customerAddress := fmt.Sprintf("%s, %s, %s, %s, %s %s",
		address.AdditionalAddress, address.SubDistrictName,
		address.DistrictName, address.CityName,
		address.ProvinceName, address.ZipCode,
	)
	originAddress := fmt.Sprintf("%s, %s, %s, %s, %s %s",
		storeConfig.AdditionalAddress, storeConfig.SubDistrictName,
		storeConfig.DistrictName, storeConfig.CityName,
		storeConfig.ProvinceName, storeConfig.ZipCode,
	)

	dbOrder := entity.SalesOrder{
		SalesOrderID:            salesOrderID,
		SalesOrderCode:          salesOrderCode,
		UserID:                  userID,
		CustomerName:            address.RecipientName,
		CustomerEmail:           user.Email,
		CustomerPhone:           address.Phone,
		CustomerAddressSnapshot: customerAddress,
		OriginAddressSnapshot:   originAddress,
		ShippingCourier:         req.CourierCode,
		ShippingDisplayName:     req.CourierName,
		ShippingService:         req.CourierService,
		ShippingFee:             req.ShippingCost,
		ShippingETD:             req.ShippingETD,
		ShippingDescription:     req.ShippingDescription,
		OriginDistrictID:        storeConfig.DistrictID,
		DestinationDistrictID:   address.DistrictID,
		TotalWeight:             totalWeight,
		TotalBeforeDiscount:     totalBeforeDiscount,
		DiscountAmount:          totalDiscount,
		FinalTotal:              finalTotal,
		Status:                  "PENDING_PAYMENT",
		Notes:                   &req.Note,
		ConfirmedAt:             &confirmedAt,
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	orderSummary := &dto.OrderSummaryResponse{
		SalesOrderID:        salesOrderID,
		SalesOrderCode:      salesOrderCode,
		Status:              "PENDING_PAYMENT",
		CustomerName:        address.RecipientName,
		CustomerPhone:       address.Phone,
		ShippingAddress:     customerAddress,
		ShippingCourier:     req.CourierCode,
		ShippingDisplayName: req.CourierName,
		ShippingService:     req.CourierService,
		ShippingFee:         req.ShippingCost,
		TotalBeforeDiscount: totalBeforeDiscount,
		DiscountAmount:      totalDiscount,
		FinalTotal:          finalTotal,
		Notes:               req.Note,
		Items:               respItems,
		ConfirmedAt:         &confirmedAt,
		CreatedAt:           now,
	}

	orderHistory := &entity.OrderStatusHistory{
		HistoryID:      uuid.NewString(),
		SalesOrderID:   salesOrderID,
		PreviousStatus: nil,
		Status:         helper.GetPendingPaymentStatus(),
		Note:           fmt.Sprintf("Membuat order %s", salesOrderCode),
		CreatedAt:      now,
		CreatedBy:      userID,
		CreatedName:    user.Name,
	}

	txErr := w.salesOrderRepo.Transaction(func(tx *gorm.DB) error {
		existing, err := w.idempotencyRepo.CheckoutFindByKey(tx, idemKey)
		if err != nil {
			return &errorhandler.InternalServerError{Message: "Gagal memeriksa idempotency key"}
		}
		if existing != nil && existing.Status == "SUCCESS" {
			if existing.ResponseBody != "" {
				var prev dto.CheckoutStatusResponse
				if json.Unmarshal([]byte(existing.ResponseBody), &prev) == nil && prev.Order != nil {
					orderSummary = prev.Order
				}
			}
			return nil // sudah pernah sukses, tidak perlu proses lagi
		}

		// Reserve stok tiap item - atomik di level DB (lihat komentar di repository)

		if err := w.salesOrderRepo.CreateOrder(tx, &dbOrder); err != nil {
			return &errorhandler.InternalServerError{Message: "Gagal menyimpan order"}
		}
		if err := w.salesOrderRepo.CreateOrderDetails(tx, dbDetails); err != nil {
			return &errorhandler.InternalServerError{Message: "Gagal menyimpan detail order"}
		}

		if err := w.salesOrderRepo.ReserveStock(tx, dbDetails); err != nil {
			return err
		}

		if err := w.historyRepo.Create(tx, orderHistory); err != nil {
			return &errorhandler.InternalServerError{Message: "Gagal menyimpan log"}
		}

		fullPayload := dto.CheckoutStatusResponse{
			IdempotencyKey: idemKey,
			Status:         "SUCCESS",
			Order:          orderSummary,
		}
		respJSON, _ := json.Marshal(fullPayload)

		idemRecord := &entity.CheckoutIdempotencyKey{
			IdempotencyID:  uuid.New().String(),
			IdempotencyKey: idemKey,
			UserID:         userID,
			SalesOrderID:   &salesOrderID,
			Status:         "SUCCESS",
			ResponseBody:   string(respJSON),
			StatusCode:     201,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		return w.idempotencyRepo.CheckoutCreate(tx, idemRecord)
	})

	if txErr != nil {
		return nil, txErr
	}

	// Hapus session checkout dari Redis
	w.redis.Del(ctx, cacheKey)

	w.notifyProductStockChanged(dbDetails)

	return orderSummary, nil
}

// ─────────────────────────────────────────────
// HELPER
// ─────────────────────────────────────────────

func (w *CheckoutWorker) setRedisStatus(ctx context.Context, key string, payload dto.CheckoutStatusResponse) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[CheckoutWorker] gagal marshal status key=%s: %v", key, err)
		return
	}
	if err := w.redis.Set(ctx, key, string(encoded), 24*time.Hour).Err(); err != nil {
		log.Printf("[CheckoutWorker] gagal set redis key=%s: %v", key, err)
	}
}

func (w *CheckoutWorker) notifyProductStockChanged(details []entity.SalesOrderDetail) {
	if Instance == nil || Instance.ProductEventChan == nil {
		log.Printf("[CheckoutWorker] ES writer belum di-initialize, skip %d event produk", len(details))
		return
	}
	for _, d := range details {
		select {
		case Instance.ProductEventChan <- &dto.ProductEvent{
			ProductID: d.ProductID,
			Type:      "create product price",
		}:
		default:
			log.Printf("[CheckoutWorker] ProductEventChan penuh, skip event product_id=%s", d.ProductID)
		}
	}
}
