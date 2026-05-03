package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type cartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) *cartHandler {
	return &cartHandler{service: service}
}

func (h *cartHandler) Verify(c *gin.Context) {
	var req dto.CartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.VerifyCart(&req)
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
