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

func InventoryRouter(api *gin.RouterGroup) {
	InventoryRepository := repository.NewInventoryRepository(config.DB)
	ProductRepository := repository.NewProductRepository(config.DB)
	UserRepository := repository.NewUserRepository(config.DB)
	LogRepository := repository.NewLogRepository(config.DB)
	InventoryService := service.NewInventoryService(InventoryRepository, ProductRepository, UserRepository, LogRepository, config.RedisClient)
	InventoryHandler := handler.NewInventoryHandler(InventoryService)

	Inventory := api.Group("/inventory")

	Inventory.Use(middleware.JWTMiddleware())

	Inventory.POST("", middleware.RoleMiddleware([]string{helper.Admin()}), InventoryHandler.Create)
	Inventory.GET("/:id", InventoryHandler.GetAll)
	Inventory.PUT("/:id", middleware.RoleMiddleware([]string{helper.Admin()}), InventoryHandler.Update)
}
