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

func ProductRouter(api *gin.RouterGroup) {
	ProductRepository := repository.NewProductRepository(config.DB)
	TypeRepository := repository.NewTypeRepository(config.DB)
	ProductService := service.NewProductService(ProductRepository, TypeRepository, config.MinioClient, config.RedisClient, config.ENV.MinioBucket)
	ProductHandler := handler.NewProductHandler(ProductService)

	Product := api.Group("/product")

	Product.Use(middleware.JWTMiddleware())

	Product.POST("/GeneratePresignedURLs", middleware.RoleMiddleware(helper.Admin()), ProductHandler.GeneratePresignedURLs)
	Product.POST("", middleware.RoleMiddleware(helper.Admin()), ProductHandler.Create)
	Product.GET("", ProductHandler.GetAllPaginated)
	Product.GET("/:id", ProductHandler.GetByID)
	Product.PUT("/:id", middleware.RoleMiddleware(helper.Admin()), ProductHandler.Update)
	Product.DELETE("/:id", middleware.RoleMiddleware(helper.Admin()), ProductHandler.Delete)
}
