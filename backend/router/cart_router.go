package router

import (
	"backend/config"
	"backend/handler"
	"backend/middleware"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func CartRouter(api *gin.RouterGroup) {
	ProductRepository := repository.NewProductRepository(config.DB)
	CartService := service.NewCartService(ProductRepository, config.MinioClient, config.RedisClient, config.ENV.MinioBucket)
	CartHandler := handler.NewCartHandler(CartService)

	Cart := api.Group("/verify-cart")

	Cart.Use(middleware.JWTMiddleware())

	Cart.POST("", CartHandler.Verify)
}
