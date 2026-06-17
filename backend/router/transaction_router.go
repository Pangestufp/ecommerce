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

func TransactionRouter(api *gin.RouterGroup) {

	transactionRepo := repository.NewTransactionRepository(config.DB)
	inventoryRepo := repository.NewInventoryRepository(config.DB)

	transactionService := service.NewTransactionService(transactionRepo, inventoryRepo)

	transactionHandler := handler.NewTransactionHandler(transactionService)

	transaction := api.Group("/transaction")
	rlRead := middleware.NewRateLimiter(60, time.Minute)

	transaction.Use(middleware.JWTMiddleware())
	transaction.GET("/batch/:batchId", rlRead.Middleware(), transactionHandler.GetAllByBatchID)
}
