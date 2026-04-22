package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type productPriceHandler struct {
	service service.ProductPriceService
}

func NewProductPriceHandler(service service.ProductPriceService) *productPriceHandler {
	return &productPriceHandler{service: service}
}

func (h *productPriceHandler) Create(c *gin.Context) {
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

	var req dto.CreateProductPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.Create(&req, userID)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusCreated,
		Message:    "Success Create Product Price",
		Data:       data,
	}))
}

func (h *productPriceHandler) GetAll(c *gin.Context) {
	productID := c.Param("id")

	direction := c.Query("direction")
	cursorID := c.Query("id")
	createdAt := c.Query("created_at")

	limit := 5
	var cursor *dto.Paginate

	if cursorID != "" && createdAt != "" {
		t, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			c.JSON(400, dto.ResponseParam{StatusCode: 400, Message: "invalid created_at format"})
			return
		}

		cursor = &dto.Paginate{}
		if direction == "prev" {
			cursor.Direction = &direction
			cursor.FirstID = &cursorID
			cursor.FirstCreatedAt = &t
		} else {
			dirNext := "next"
			cursor.Direction = &dirNext
			cursor.LastID = &cursorID
			cursor.LastCreatedAt = &t
		}
	}

	data, paginate, err := h.service.GetAllByProductID(productID, cursor, limit)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Paginate:   paginate,
		Data:       data,
	}))
}
