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
	CourierRepository := repository.NewCourierRepository(config.DB)
	StoreConfigRepository := repository.NewStoreConfigRepository(config.DB)
	OngkirService := service.NewRajaOngkirService(config.ENV.RajaOngkirAPIKey, config.ENV.RajaOngkirURL, config.RedisClient, CourierRepository)
	StoreConfigService := service.NewStoreConfigService(StoreConfigRepository, OngkirService, config.RedisClient)
	CheckoutService := service.NewCheckoutService(ProductRepository, DiscountRepository, ProductPriceRepository, UserAddressRepository, StoreConfigService, OngkirService, config.MinioClient, config.RedisClient, config.ENV.MinioBucket)
	CheckoutHandler := handler.NewCheckoutHandler(CheckoutService)

	Checkout := api.Group("/checkout")

	Checkout.Use(middleware.JWTMiddleware())

	rlWrite := middleware.NewRateLimiter(20, time.Minute)
	rlRead := middleware.NewRateLimiter(60, time.Minute)
	Checkout.POST("", rlWrite.Middleware(), CheckoutHandler.CreateCheckout)
	Checkout.GET("/:id", rlRead.Middleware(), CheckoutHandler.GetCheckOut)

	Courier := api.Group("/courier-fee")

	Courier.Use(middleware.JWTMiddleware())
	Courier.POST("", rlWrite.Middleware(), CheckoutHandler.GetCourier)
}
