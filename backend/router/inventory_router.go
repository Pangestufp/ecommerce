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

func InventoryRouter(api *gin.RouterGroup) {
	InventoryRepository := repository.NewInventoryRepository(config.DB)
	ProductRepository := repository.NewProductRepository(config.DB)
	UserRepository := repository.NewUserRepository(config.DB)
	LogRepository := repository.NewLogRepository(config.DB)
	InventoryService := service.NewInventoryService(InventoryRepository, ProductRepository, UserRepository, LogRepository, config.RedisClient)
	InventoryHandler := handler.NewInventoryHandler(InventoryService)

	Inventory := api.Group("/inventory")

	Inventory.Use(middleware.JWTMiddleware())
	rlWrite := middleware.NewRateLimiter(20, time.Minute)
	rlRead := middleware.NewRateLimiter(60, time.Minute)

	Inventory.POST("", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), InventoryHandler.Create)
	Inventory.GET("/:id", rlRead.Middleware(), InventoryHandler.GetAll)
	Inventory.PUT("/:id", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), InventoryHandler.Update)
}
