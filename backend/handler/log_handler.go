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

type logHandler struct {
	service service.LogService
}

func NewLogHandler(service service.LogService) *logHandler {
	return &logHandler{service: service}
}


func (h *logHandler) GetByProductID(c *gin.Context) {
	productID := c.Param("id") 

	direction := c.Query("direction")
	cursorID := c.Query("id")
	createdAt := c.Query("created_at")

	limit := 5
	var cursor *dto.Paginate

	if cursorID != "" && createdAt != "" {
		t, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: "invalid created_at format"})
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

	data, paginate, err := h.service.GetByProductID(productID, cursor, limit)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Get Logs by Product ID",
		Paginate:   paginate,
		Data:       data,
	}))
}


func (h *logHandler) GetByReferenceType(c *gin.Context) { 
	direction := c.Query("direction")
	cursorID := c.Query("id")
	createdAt := c.Query("created_at")

	limit := 5
	var cursor *dto.Paginate

	if cursorID != "" && createdAt != "" {
		t, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: "invalid created_at format"})
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

	data, paginate, err := h.service.GetByReferenceType(cursor, limit)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Get Logs by Reference Type",
		Paginate:   paginate,
		Data:       data,
	}))
}