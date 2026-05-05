package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"
	"strconv"
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
	userIDVal, ok := c.Get("userID")
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{
			Message: "unauthorized",
		})
		return
	}
	userID, ok := userIDVal.(string)

	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.InternalServerError{
			Message: "invalid user id",
		})
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	data, err := h.service.Create(req, userID)
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
	userIDVal, ok := c.Get("userID")
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{
			Message: "unauthorized",
		})
		return
	}
	userID, ok := userIDVal.(string)

	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.InternalServerError{
			Message: "invalid user id",
		})
		return
	}

	productID := c.Param("id")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	data, err := h.service.Update(productID, req, userID)
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
	userIDVal, ok := c.Get("userID")
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{
			Message: "unauthorized",
		})
		return
	}
	userID, ok := userIDVal.(string)

	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.InternalServerError{
			Message: "invalid user id",
		})
		return
	}

	productID := c.Param("id")

	if err := h.service.Delete(productID, userID); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "Success",
	}))
}

func (h *ProductHandler) GetAllPaginated(c *gin.Context) {
	direction := c.Query("direction")
	id := c.Query("id")
	createdAt := c.Query("created_at")
	search := c.Query("search")

	if id == "" && createdAt == "" && direction == "" {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: "Invalid direction format"})
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

	products, paginate, err := h.service.GetAllPaginated(cursor, search, limit)
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

func (h *ProductHandler) GetProductBySearch(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.service.GetProductBySearch(search, page, limit)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(200, dto.ResponseParam{
		StatusCode: 200,
		Message:    "success",
		Data:       products,
	})
}

func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	product, err := h.service.GetProductEnrichedBySlug(slug)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(200, dto.ResponseParam{
		StatusCode: 200,
		Message:    "success",
		Data:       product,
	})
}
