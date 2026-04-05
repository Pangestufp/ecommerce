package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type typeHandler struct {
	service service.TypeService
}

func NewTypeHandler(service service.TypeService) *typeHandler {
	return &typeHandler{service: service}
}

func (h *typeHandler) Create(c *gin.Context) {
	var req dto.TypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.CreateType(&req)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusCreated,
		Message:    "Success Create Type",
		Data:       data,
	}))
}

func (h *typeHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.TypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.UpdateType(id, &req)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Update Type",
		Data:       data,
	}))
}

func (h *typeHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteType(id); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Delete Type",
	}))
}

func (h *typeHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAllType()
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

func (h *typeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	data, err := h.service.GetTypeByID(id)
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
