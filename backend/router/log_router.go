package router

import (
	"backend/config"
	"backend/handler"
	"backend/repository"
	"backend/service"

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

	// Endpoint untuk mengambil log berdasarkan Product ID 
	Log.GET("/product/:id", LogHandler.GetByProductID)

	// Endpoint untuk mengambil log berdasarkan Tipe 
	Log.GET("/type/:type", LogHandler.GetByReferenceType)
}