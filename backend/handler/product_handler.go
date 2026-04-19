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

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GeneratePresignedURLs(c *gin.Context) {
	var req dto.PresignedURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	res, err := h.service.GeneratePresignedURLs(req)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	data, err := h.service.Create(req)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusCreated,
		Message:    "Success",
		Data:       data,
	}))
}

func (h *ProductHandler) Update(c *gin.Context) {
	productID := c.Param("id")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	data, err := h.service.Update(productID, req)
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

func (h *ProductHandler) GetByID(c *gin.Context) {
	productID := c.Param("id")

	data, err := h.service.GetByID(productID)
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

func (h *ProductHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll()
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

func (h *ProductHandler) Delete(c *gin.Context) {
	productID := c.Param("id")

	if err := h.service.Delete(productID); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success",
	}))
}

func (h *ProductHandler) GetAllPaginated(c *gin.Context) {
	limit := 5

	var cursor *dto.Paginate

	direction := c.Query("direction")
	id := c.Query("id")
	createdAt := c.Query("created_at")

	if id != "" && createdAt != "" {
		t, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			c.JSON(400, dto.ResponseParam{
				StatusCode: 400,
				Message:    "invalid created_at format",
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

	products, paginate, err := h.service.GetAllPaginated(cursor, limit)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(200, dto.ResponseParam{
		StatusCode: 200,
		Message:    "success",
		Paginate:   paginate,
		Data:       products,
	})
}
