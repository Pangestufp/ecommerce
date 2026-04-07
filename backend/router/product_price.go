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
	ProductPriceService := service.NewProductPriceService(ProductPriceRepository, ProductRepository, config.RedisClient)
	ProductPriceHandler := handler.NewProductPriceHandler(ProductPriceService)

	ProductPrice := api.Group("/productPrice")

	ProductPrice.Use(middleware.JWTMiddleware())

	ProductPrice.POST("", middleware.RoleMiddleware(helper.Admin()), ProductPriceHandler.Create)
	ProductPrice.GET("/:id", ProductPriceHandler.GetAll)
}
