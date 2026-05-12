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

func LogRouter(api *gin.RouterGroup) {

	LogRepository := repository.NewLogRepository(config.DB)
	ProductRepository := repository.NewProductRepository(config.DB)

	LogService := service.NewLogService(LogRepository, ProductRepository)

	// Inisialisasi Handler
	LogHandler := handler.NewLogHandler(LogService)

	// Grouping Route
	Log := api.Group("/logs")

	Log.Use(middleware.JWTMiddleware())

	rlRead := middleware.NewRateLimiter(60, time.Minute)

	// Endpoint untuk mengambil log berdasarkan Product ID
	Log.GET("/product/:id", rlRead.Middleware(), LogHandler.GetByProductID)

	// Endpoint untuk mengambil log berdasarkan Tipe
	Log.GET("/type", rlRead.Middleware(), LogHandler.GetByReferenceType)
}
