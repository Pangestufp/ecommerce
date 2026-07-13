package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type paymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *paymentHandler {
	return &paymentHandler{service: service}
}

func (h *paymentHandler) CreateSnapToken(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{
			Message: "Header Idempotency-Key wajib diisi",
		})
		return
	}

	val, ok := c.Get("userID")
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{Message: "unauthorized"})
		return
	}
	userID, ok := val.(string)
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.InternalServerError{Message: "invalid user id"})
		return
	}

	data, err := h.service.CreateSnapToken(&req, userID, idempotencyKey)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusCreated,
		Message:    "Snap token berhasil dibuat",
		Data:       data,
	}))
}

func (h *paymentHandler) GetPaymentStatus(c *gin.Context) {
	salesOrderCode := c.Param("code")
	if salesOrderCode == "" {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: "Sales order code diperlukan"})
		return
	}

	val, ok := c.Get("userID")
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{Message: "unauthorized"})
		return
	}
	userID, ok := val.(string)
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.InternalServerError{Message: "invalid user id"})
		return
	}

	data, err := h.service.GetPaymentStatus(salesOrderCode, userID)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Data:       data,
	}))
}

// MidtransWebhook tidak pakai JWT — keamanan dari verifikasi signature di service.
// Selalu return 200 supaya Midtrans tidak retry, meski ada error internal.
func (h *paymentHandler) MidtransWebhook(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		// Tetap 200 — kalau 4xx/5xx Midtrans akan retry terus
		c.JSON(http.StatusOK, gin.H{"message": "diterima"})
		return
	}

	var notification dto.MidtransNotification
	if err := json.Unmarshal(rawBody, &notification); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "diterima"})
		return
	}

	// Jalankan async — webhook harus response cepat (< 5 detik) ke Midtrans
	// Kalau proses berat, pindahkan ke goroutine atau worker
	_ = h.service.HandleWebhook(&notification, string(rawBody), c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": "diterima"})
}
