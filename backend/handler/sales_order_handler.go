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

type SalesOrderHandler struct {
	service service.SalesOrderService
}

func NewSalesOrderHandler(service service.SalesOrderService) *SalesOrderHandler {
	return &SalesOrderHandler{service: service}
}

func (h *SalesOrderHandler) GetAll(c *gin.Context) {
	limit, statuses, cursor, ok := parseListParams(c)
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: "Parameter tidak valid"})
		return
	}

	orders, paginate, err := h.service.GetAllPaginate(cursor, statuses, limit)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     orders,
		"paginate": paginate,
	})
}

func (h *SalesOrderHandler) GetMyOrders(c *gin.Context) {
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

	limit, statuses, cursor, ok := parseListParams(c)
	if !ok {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: "Parameter tidak valid"})
		return
	}

	orders, paginate, err := h.service.GetAllByUserPaginate(cursor, userID, statuses, limit)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     orders,
		"paginate": paginate,
	})
}

func (h *SalesOrderHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: "code tidak valid"})
		return
	}

	result, err := h.service.GetByCode(code)
	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func parseListParams(c *gin.Context) (limit int, statuses []string, cursor *dto.Paginate, ok bool) {
	limit = 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	statuses = c.QueryArray("statuses")
	if len(statuses) > 0 {
		if err := helper.ValidateStatuses(statuses); err != nil {
			errorhandler.ErrorHandler(c, err)
			return 0, nil, nil, false
		}
	}

	cursor = parseCursor(c)
	return limit, statuses, cursor, true
}

// parseCursor membaca query param cursor dari request Gin.
func parseCursor(c *gin.Context) *dto.Paginate {
	firstID := c.Query("first_id")
	lastID := c.Query("last_id")
	direction := c.Query("direction")
	hasNext := c.Query("has_next")
	hasPrev := c.Query("has_prev")

	if firstID == "" && lastID == "" && direction == "" {
		return nil
	}

	cursor := &dto.Paginate{}

	if firstID != "" {
		cursor.FirstID = &firstID
	}
	if lastID != "" {
		cursor.LastID = &lastID
	}
	if direction != "" {
		cursor.Direction = &direction
	}
	if hasNext != "" {
		cursor.HasNext = &hasNext
	}
	if hasPrev != "" {
		cursor.HasPrev = &hasPrev
	}

	if raw := c.Query("first_created_at"); raw != "" {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			cursor.FirstCreatedAt = &t
		}
	}
	if raw := c.Query("last_created_at"); raw != "" {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			cursor.LastCreatedAt = &t
		}
	}

	return cursor
}
