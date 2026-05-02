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
	DiscountRepository := repository.NewDiscountRepository(config.DB)
	TypeRepository := repository.NewTypeRepository(config.DB)
	ProductService := service.NewProductService(ProductRepository, TypeRepository, DiscountRepository, config.MinioClient, config.RedisClient, config.ENV.MinioBucket)
	ProductHandler := handler.NewProductHandler(ProductService)

	Product := api.Group("/product")

	Product.Use(middleware.JWTMiddleware())

	Product.POST("/presigned-urls", middleware.RoleMiddleware([]string{helper.Admin()}), ProductHandler.GeneratePresignedURLs)
	Product.POST("", middleware.RoleMiddleware([]string{helper.Admin()}), ProductHandler.Create)
	Product.GET("", ProductHandler.GetAllPaginated)
	Product.GET("/:id", ProductHandler.GetByID)
	Product.PUT("/:id", middleware.RoleMiddleware([]string{helper.Admin()}), ProductHandler.Update)
	Product.DELETE("/:id", middleware.RoleMiddleware([]string{helper.Admin()}), ProductHandler.Delete)

	ProductCatalog := api.Group("/catalog")
	ProductCatalog.Use(middleware.JWTMiddleware())
	ProductCatalog.GET("", ProductHandler.GetProductBySearch)
	ProductCatalog.GET("/:slug", ProductHandler.GetProductBySlug)
}
