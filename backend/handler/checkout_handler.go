package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type checkoutHandler struct {
	service service.CheckoutService
}

func NewCheckoutHandler(service service.CheckoutService) *checkoutHandler {
	return &checkoutHandler{service: service}
}

func (h *checkoutHandler) CreateCheckout(c *gin.Context) {
	var req dto.CartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	userIDVal, ok := c.Get("userID")
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{
			Message: "unauthorized",
		})
		return
	}

	userID, ok := userIDVal.(string)

	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.InternalServerError{
			Message: "invalid user id",
		})
		return
	}

	data, err := h.service.CreateCheckout(&req, userID)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Verify Cart",
		Data:       data,
	}))
}

func (h *checkoutHandler) GetCheckOut(c *gin.Context) {
	checkoutID := c.Param("id")

	userIDVal, ok := c.Get("userID")
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{
			Message: "unauthorized",
		})
		return
	}

	userID, ok := userIDVal.(string)

	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.InternalServerError{
			Message: "invalid user id",
		})
		return
	}

	data, err := h.service.GetCheckout(checkoutID, userID)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Verify Cart",
		Data:       data,
	}))
}

func (h *checkoutHandler) GetCourier(c *gin.Context) {
	var req dto.ShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	userIDVal, ok := c.Get("userID")
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{
			Message: "unauthorized",
		})
		return
	}

	userID, ok := userIDVal.(string)

	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.InternalServerError{
			Message: "invalid user id",
		})
		return
	}

	data, err := h.service.CalculateShippingFromAddress(&req, userID)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Verify Cart",
		Data:       data,
	}))

}

func (h *checkoutHandler) ConfirmCheckout(c *gin.Context) {
	var req dto.ConfirmCheckoutRequest
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

	data, err := h.service.ConfirmCheckout(&req, userID, idempotencyKey)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	// 202 Accepted: request diterima tapi belum selesai diproses
	c.JSON(http.StatusAccepted, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusAccepted,
		Message:    "Checkout diterima, silakan polling status",
		Data:       data,
	}))
}

func (h *checkoutHandler) GetCheckoutStatus(c *gin.Context) {
	idempotencyKey := c.Param("key")
	if idempotencyKey == "" {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{
			Message: "Idempotency key tidak valid",
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

	data, err := h.service.GetCheckoutStatus(idempotencyKey, userID)
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
