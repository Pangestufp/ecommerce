package router

import (
	"backend/config"
	"backend/handler"
	"backend/helper"
	"backend/middleware"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func ProductPriceRouter(api *gin.RouterGroup) {
	ProductPriceRepository := repository.NewProductPriceRepository(config.DB)
	ProductRepository := repository.NewProductRepository(config.DB)
	UserRepository := repository.NewUserRepository(config.DB)
	InventoryRepository := repository.NewInventoryRepository(config.DB)
	LogRepository := repository.NewLogRepository(config.DB)
	ProductPriceService := service.NewProductPriceService(ProductPriceRepository, ProductRepository, UserRepository,InventoryRepository, LogRepository, config.RedisClient)
	ProductPriceHandler := handler.NewProductPriceHandler(ProductPriceService)

	ProductPrice := api.Group("/product-price")

	ProductPrice.Use(middleware.JWTMiddleware())

	ProductPrice.POST("", middleware.RoleMiddleware([]string{helper.Admin()}), ProductPriceHandler.Create)
	ProductPrice.GET("/:id", ProductPriceHandler.GetAll)
}
