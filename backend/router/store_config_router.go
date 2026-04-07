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

func StoreConfigRouter(api *gin.RouterGroup) {
	StoreConfigRepository := repository.NewStoreConfigRepository(config.DB)
	StoreConfigService := service.NewStoreConfigService(StoreConfigRepository, config.RedisClient)
	StoreConfigHandler := handler.NewStoreConfigHandler(StoreConfigService)

	StoreConfig := api.Group("/storeConfig")

	StoreConfig.Use(middleware.JWTMiddleware())

	StoreConfig.PUT("", middleware.RoleMiddleware(helper.Admin()), StoreConfigHandler.Upsert)
	StoreConfig.GET("", StoreConfigHandler.GetConfig)

}
