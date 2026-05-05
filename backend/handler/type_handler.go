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

type typeHandler struct {
	service service.TypeService
}

func NewTypeHandler(service service.TypeService) *typeHandler {
	return &typeHandler{service: service}
}

func (h *typeHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	var req dto.TypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.CreateType(&req, userID)

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
	userID := c.MustGet("userID").(string)

	var req dto.TypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.UpdateType(id, &req, userID)
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
	userID := c.MustGet("userID").(string)

	if err := h.service.DeleteType(id,userID); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success Delete Type",
	}))
}

func (h *typeHandler) GetAll(c *gin.Context) {

	direction := c.Query("direction")
	id := c.Query("id")
	createdAt := c.Query("created_at")
	search := c.Query("search")

	if id == "" && createdAt == "" && direction == "" {
		types, err := h.service.GetAllType()
		if err != nil {
			errorhandler.ErrorHandler(c, err)
			return
		}
		c.JSON(200, dto.ResponseParam{
			StatusCode: 200,
			Message:    "success",
			Data:       types,
		})
		return
	}

	limit := 5
	var cursor *dto.Paginate

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

	products, paginate, err := h.service.GetAllTypePaginate(cursor, search, limit)
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
