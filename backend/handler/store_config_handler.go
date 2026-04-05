package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type storeConfigHandler struct {
	service service.StoreConfigService
}

func NewStoreConfigHandler(service service.StoreConfigService) *storeConfigHandler {
	return &storeConfigHandler{service: service}
}

func (h *storeConfigHandler) Upsert(c *gin.Context) {
	var req dto.StoreConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := h.service.Upsert(&req); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success",
	}))
}

func (h *storeConfigHandler) GetConfig(c *gin.Context) {
	data, err := h.service.GetConfig()
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
