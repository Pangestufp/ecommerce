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

func DiscountRouter(api *gin.RouterGroup) {
	DiscountRepository := repository.NewDiscountRepository(config.DB)
	ProductRepository := repository.NewProductRepository(config.DB)
	UserRepository := repository.NewUserRepository(config.DB)
	ProductPriceRepository := repository.NewProductPriceRepository(config.DB)
	LogRepository := repository.NewLogRepository(config.DB)
	DiscountService := service.NewDiscountService(DiscountRepository, ProductRepository, UserRepository, ProductPriceRepository, LogRepository, config.RedisClient)
	DiscountHandler := handler.NewDiscountHandler(DiscountService)

	rlWrite := middleware.NewRateLimiter(20, time.Minute)
	rlRead := middleware.NewRateLimiter(60, time.Minute)

	Discount := api.Group("/discount")
	Discount.Use(middleware.JWTMiddleware())
	Discount.POST("", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), DiscountHandler.Create)
	Discount.DELETE("/:id", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), DiscountHandler.Delete)
	Discount.GET("/:id", rlRead.Middleware(), DiscountHandler.GetAll)

	DiscountType := api.Group("/discount-type")
	DiscountType.Use(middleware.JWTMiddleware())
	DiscountType.GET("", rlRead.Middleware(), DiscountHandler.GetAllDiscountType)
}
