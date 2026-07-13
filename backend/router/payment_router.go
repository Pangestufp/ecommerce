package router

import (
	"backend/config"
	"backend/handler"
	"backend/middleware"
	"backend/repository"
	"backend/service"
	"time"

	"github.com/gin-gonic/gin"
)

func PaymentRouter(api *gin.RouterGroup) {
	paymentRepository := repository.NewPaymentRepository(config.DB)
	webhookLogRepository := repository.NewPaymentWebhookLogRepository(config.DB)
	salesOrderRepository := repository.NewSalesOrderRepository(config.DB)
	inventoryRepository := repository.NewInventoryRepository(config.DB)
	historyRepository := repository.NewOrderStatusHistoryRepository(config.DB)

	paymentService := service.NewPaymentService(
		paymentRepository,
		webhookLogRepository,
		salesOrderRepository,
		historyRepository,
		inventoryRepository,
		config.RedisClient,
		config.ENV.MidtransProduction,
		config.ENV.SandBoxSnapURL,
		config.ENV.ProductionSnapURL,
		config.ENV.ServerKey,
		config.ENV.ENVPro,
	)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	rlWrite := middleware.NewRateLimiter(20, time.Minute)
	rlRead := middleware.NewRateLimiter(60, time.Minute)

	// User routes — butuh JWT
	payment := api.Group("/payment")
	payment.Use(middleware.JWTMiddleware())
	{
		payment.POST("", rlWrite.Middleware(), paymentHandler.CreateSnapToken)
		payment.GET("/status/:code", rlRead.Middleware(), paymentHandler.GetPaymentStatus)
	}

	webhook := api.Group("/midtrans")
	{
		webhook.POST("/notification", paymentHandler.MidtransWebhook)
	}
}
