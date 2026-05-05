package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type rajaOngkirHandler struct {
	service service.RajaOngkirService
}

func NewRajaOngkirHandler(service service.RajaOngkirService) *rajaOngkirHandler {
	return &rajaOngkirHandler{service: service}
}

func (h *rajaOngkirHandler) GetProvince(c *gin.Context) {
	data, err := h.service.GetProvince()
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Get Province",
		Data:       data,
	}))
}

func (h *rajaOngkirHandler) GetCity(c *gin.Context) {
	provinceID := c.Param("province_id")

	data, err := h.service.GetCity(provinceID)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Get City",
		Data:       data,
	}))
}

func (h *rajaOngkirHandler) GetDistrict(c *gin.Context) {
	cityID := c.Param("city_id")

	data, err := h.service.GetDistrict(cityID)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Get District",
		Data:       data,
	}))
}
