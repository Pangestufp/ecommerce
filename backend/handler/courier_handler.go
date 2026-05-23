package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type courierHandler struct {
	service service.CourierService
}

func NewCourierHandler(service service.CourierService) *courierHandler {
	return &courierHandler{service: service}
}

func (h *courierHandler) Create(c *gin.Context) {
	var req dto.CreateCourierRequest
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
		Message:    "Kurir berhasil ditambahkan",
		Data:       data,
	}))
}

func (h *courierHandler) GetAll(c *gin.Context) {
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

func (h *courierHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateCourierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.Update(id, &req)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Kurir berhasil diubah",
		Data:       data,
	}))
}

func (h *courierHandler) Toggle(c *gin.Context) {
	id := c.Param("id")

	data, err := h.service.Toggle(id)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	status := "dinonaktifkan"
	if data.Status == 1 {
		status = "diaktifkan"
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Kurir berhasil " + status,
		Data:       data,
	}))
}
