package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type productPriceHandler struct {
	service service.ProductPriceService
}

func NewProductPriceHandler(service service.ProductPriceService) *productPriceHandler {
	return &productPriceHandler{service: service}
}

func (h *productPriceHandler) Create(c *gin.Context) {
	var req dto.CreateProductPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.Create(&req)
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

	data, err := h.service.GetAllByProductID(productID)
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
