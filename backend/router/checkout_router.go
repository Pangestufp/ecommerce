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

func CheckoutRouter(api *gin.RouterGroup) {
	ProductRepository := repository.NewProductRepository(config.DB)
	DiscountRepository := repository.NewDiscountRepository(config.DB)
	ProductPriceRepository := repository.NewProductPriceRepository(config.DB)
	UserAddressRepository := repository.NewAddressRepository(config.DB)
	CheckoutService := service.NewCheckoutService(ProductRepository, DiscountRepository, ProductPriceRepository, UserAddressRepository, config.MinioClient, config.RedisClient, config.ENV.MinioBucket)
	CheckoutHandler := handler.NewCheckoutHandler(CheckoutService)

	Checkout := api.Group("/verify-checkout")

	Checkout.Use(middleware.JWTMiddleware())

	rlVerify := middleware.NewRateLimiter(20, time.Minute)
	Checkout.POST("", rlVerify.Middleware(), CheckoutHandler.Verify)
}
