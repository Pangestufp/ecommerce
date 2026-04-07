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
	DiscountService := service.NewDiscountService(DiscountRepository, ProductRepository)
	DiscountHandler := handler.NewDiscountHandler(DiscountService)

	Discount := api.Group("/discount")

	Discount.Use(middleware.JWTMiddleware())

	Discount.POST("", middleware.RoleMiddleware(helper.Admin()), DiscountHandler.Create)
	Discount.DELETE("/:id", middleware.RoleMiddleware(helper.Admin()), DiscountHandler.Delete)
}
