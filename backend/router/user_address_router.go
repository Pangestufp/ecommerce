package router

import (
	"backend/config"
	"backend/handler"
	"backend/middleware"
	"backend/repository"
	"backend/service"
	"time"

	"github.com/gin-gonic/gin"
)

func UserAddressRouter(api *gin.RouterGroup) {
	UserAddressRepository := repository.NewAddressRepository(config.DB)
	UserRepository := repository.NewUserRepository(config.DB)
	CourierRepository := repository.NewCourierRepository(config.DB)
	OngkirService := service.NewRajaOngkirService(config.ENV.RajaOngkirAPIKey, config.ENV.RajaOngkirURL, config.RedisClient, CourierRepository)
	UserAddressService := service.NewAddressService(UserAddressRepository, UserRepository, OngkirService)
	UserAddressHandler := handler.NewAddressHandler(UserAddressService)

	UserAddress := api.Group("/address")

	UserAddress.Use(middleware.JWTMiddleware())

	rlWrite := middleware.NewRateLimiter(20, time.Minute)
	rlRead := middleware.NewRateLimiter(60, time.Minute)

	UserAddress.POST("", rlWrite.Middleware(), UserAddressHandler.Create)
	UserAddress.PUT("/:id", rlWrite.Middleware(), UserAddressHandler.Update)
	UserAddress.GET("/:id", rlRead.Middleware(), UserAddressHandler.GetMyAddresses)
	UserAddress.DELETE("/:id", rlWrite.Middleware(), UserAddressHandler.Delete)
}
