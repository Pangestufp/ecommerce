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

func DiscountRouter(api *gin.RouterGroup) {
	DiscountRepository := repository.NewDiscountRepository(config.DB)
	ProductRepository := repository.NewProductRepository(config.DB)
	UserRepository := repository.NewUserRepository(config.DB)
	ProductPriceRepository := repository.NewProductPriceRepository(config.DB)
	DiscountService := service.NewDiscountService(DiscountRepository, ProductRepository, UserRepository, ProductPriceRepository, config.RedisClient)
	DiscountHandler := handler.NewDiscountHandler(DiscountService)

	Discount := api.Group("/discount")

	Discount.Use(middleware.JWTMiddleware())

	Discount.POST("", middleware.RoleMiddleware(helper.Admin()), DiscountHandler.Create)
	Discount.DELETE("/:id", middleware.RoleMiddleware(helper.Admin()), DiscountHandler.Delete)
	Discount.GET("/:id", DiscountHandler.GetAll)

	DiscountType := api.Group("/discountType")
	DiscountType.Use(middleware.JWTMiddleware())
	DiscountType.GET("", DiscountHandler.GetAllDiscountType)
}
