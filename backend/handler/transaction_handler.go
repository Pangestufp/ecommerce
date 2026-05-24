package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *transactionHandler {
	return &transactionHandler{service: service}
}

func (h *transactionHandler) GetAllByBatchID(c *gin.Context) {
	batchID := c.Param("batchId") 

	id := c.Query("id")
	createdAt := c.Query("created_at")
	direction := c.Query("direction")
	limitStr := c.Query("limit")

	limit := 5 
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	var cursor *dto.Paginate
	if id != "" && createdAt != "" {
		t, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ResponseParam{
				StatusCode: http.StatusBadRequest,
				Message:    "Format created_at salah",
			})
			return
		}

		cursor = &dto.Paginate{}
		if direction == "prev" {
			cursor.Direction = &direction
			cursor.FirstID = &id
			cursor.FirstCreatedAt = &t
		} else {
			dirNext := "next"
			cursor.Direction = &dirNext
			cursor.LastID = &id
			cursor.LastCreatedAt = &t
		}
	}

	// Panggil service
	data, paginate, err := h.service.GetAllByBatchID(batchID, cursor, limit)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}


	c.JSON(http.StatusOK, dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "success",
		Paginate:   paginate,
		Data:       data,
	})
}