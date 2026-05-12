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

func ProductRouter(api *gin.RouterGroup) {
	ProductRepository := repository.NewProductRepository(config.DB)
	DiscountRepository := repository.NewDiscountRepository(config.DB)
	TypeRepository := repository.NewTypeRepository(config.DB)
	UserRepository := repository.NewUserRepository(config.DB)
	LogRepository := repository.NewLogRepository(config.DB)
	ProductService := service.NewProductService(ProductRepository, TypeRepository, DiscountRepository, LogRepository, UserRepository, config.MinioClient, config.RedisClient, config.ENV.MinioBucket)
	ProductHandler := handler.NewProductHandler(ProductService)

	Product := api.Group("/product")

	Product.Use(middleware.JWTMiddleware())

	rlWrite := middleware.NewRateLimiter(20, time.Minute)
	rlRead := middleware.NewRateLimiter(60, time.Minute)

	Product.POST("/presigned-urls", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), ProductHandler.GeneratePresignedURLs)
	Product.POST("", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), ProductHandler.Create)
	Product.GET("", rlRead.Middleware(), ProductHandler.GetAllPaginated)
	Product.GET("/:id", rlRead.Middleware(), ProductHandler.GetByID)
	Product.PUT("/:id", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), ProductHandler.Update)
	Product.DELETE("/:id", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), ProductHandler.Delete)

	ProductCatalog := api.Group("/catalog")
	ProductCatalog.Use(middleware.JWTMiddleware())
	ProductCatalog.GET("", rlRead.Middleware(), ProductHandler.GetProductBySearch)
	ProductCatalog.GET("/:slug", rlRead.Middleware(), ProductHandler.GetProductBySlug)
}
