package router

import (
	"backend/config"
	"backend/handler"
	"backend/middleware"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func OngkirRouter(api *gin.RouterGroup) {
	OngkirService := service.NewRajaOngkirService(config.ENV.RajaOngkirAPIKey, config.ENV.RajaOngkirURL, config.RedisClient)
	OngkirHandler := handler.NewRajaOngkirHandler(OngkirService)

	ongkir := api.Group("/ongkir")
	ongkir.Use(middleware.JWTMiddleware())
	{
		ongkir.GET("/province", OngkirHandler.GetProvince)
		ongkir.GET("/city/:province_id", OngkirHandler.GetCity)
		ongkir.GET("/district/:city_id", OngkirHandler.GetDistrict)
	}
}
