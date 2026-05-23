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

func CartRouter(api *gin.RouterGroup) {
	ProductRepository := repository.NewProductRepository(config.DB)
	CartService := service.NewCartService(ProductRepository, config.MinioClient, config.RedisClient, config.ENV.MinioBucket)
	CartHandler := handler.NewCartHandler(CartService)

	Cart := api.Group("/verify-cart")

	Cart.Use(middleware.JWTMiddleware())

	rlVerify := middleware.NewRateLimiter(20, time.Minute)
	Cart.POST("", rlVerify.Middleware(), CartHandler.Verify)
}
