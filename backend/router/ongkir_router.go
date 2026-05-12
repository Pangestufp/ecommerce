package router

import (
	"backend/config"
	"backend/handler"
	"backend/middleware"
	"backend/service"
	"time"

	"github.com/gin-gonic/gin"
)

func OngkirRouter(api *gin.RouterGroup) {
	OngkirService := service.NewRajaOngkirService(config.ENV.RajaOngkirAPIKey, config.ENV.RajaOngkirURL, config.RedisClient)
	OngkirHandler := handler.NewRajaOngkirHandler(OngkirService)

	rlOngkir := middleware.NewRateLimiter(30, time.Minute)

	ongkir := api.Group("/ongkir")
	ongkir.Use(middleware.JWTMiddleware())
	{
		ongkir.GET("/province", rlOngkir.Middleware(), OngkirHandler.GetProvince)
		ongkir.GET("/city/:province_id", rlOngkir.Middleware(), OngkirHandler.GetCity)
		ongkir.GET("/district/:city_id", rlOngkir.Middleware(), OngkirHandler.GetDistrict)
	}
}
