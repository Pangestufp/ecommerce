package router

import (
	"backend/config"
	"backend/handler"
	"backend/helper"
	"backend/middleware"
	"backend/repository"
	"backend/service"
	"time"

	"github.com/gin-gonic/gin"
)

func SalesOrderRouter(api *gin.RouterGroup) {

	salesOrderRepository := repository.NewSalesOrderRepository(config.DB)
	historyRepository := repository.NewOrderStatusHistoryRepository(config.DB)
	salesOrderService := service.NewSalesOrderService(salesOrderRepository, historyRepository, config.MinioClient, config.ENV.MinioBucket)
	salesOrderHandler := handler.NewSalesOrderHandler(salesOrderService)

	adminSalesOrder := api.Group("/admin/sales-orders")
	userSalesOrder := api.Group("/user/sales-orders")

	adminSalesOrder.Use(middleware.JWTMiddleware())
	userSalesOrder.Use(middleware.JWTMiddleware())

	rlRead := middleware.NewRateLimiter(60, time.Minute)

	adminSalesOrder.GET("", rlRead.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), salesOrderHandler.GetAll)
	adminSalesOrder.GET("/:code", rlRead.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), salesOrderHandler.GetByCodeAdmin)
	userSalesOrder.GET("", rlRead.Middleware(), salesOrderHandler.GetMyOrders)
	userSalesOrder.GET("/:code", rlRead.Middleware(), salesOrderHandler.GetByCodeUser)
}
