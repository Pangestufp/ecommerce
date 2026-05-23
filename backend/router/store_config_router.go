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

func StoreConfigRouter(api *gin.RouterGroup) {
	StoreConfigRepository := repository.NewStoreConfigRepository(config.DB)
	CourierRepository := repository.NewCourierRepository(config.DB)
	OngkirService := service.NewRajaOngkirService(config.ENV.RajaOngkirAPIKey, config.ENV.RajaOngkirURL, config.RedisClient, CourierRepository)
	StoreConfigService := service.NewStoreConfigService(StoreConfigRepository, OngkirService, config.RedisClient)
	StoreConfigHandler := handler.NewStoreConfigHandler(StoreConfigService)

	StoreConfig := api.Group("/store-config")

	StoreConfig.Use(middleware.JWTMiddleware())

	rlWrite := middleware.NewRateLimiter(20, time.Minute)
	rlRead := middleware.NewRateLimiter(60, time.Minute)

	StoreConfig.PUT("", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), StoreConfigHandler.Upsert)
	StoreConfig.GET("", rlRead.Middleware(), StoreConfigHandler.GetConfig)

}
