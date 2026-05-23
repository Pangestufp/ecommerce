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

func CourierRouter(api *gin.RouterGroup) {
	CourierRepository := repository.NewCourierRepository(config.DB)
	CourierService := service.NewCourierService(CourierRepository)
	CourierHandler := handler.NewCourierHandler(CourierService)

	rlWrite := middleware.NewRateLimiter(20, time.Minute)
	rlRead := middleware.NewRateLimiter(60, time.Minute)

	courier := api.Group("/courier")
	courier.Use(middleware.JWTMiddleware())

	courier.POST("", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), CourierHandler.Create)
	courier.PUT("/:id", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), CourierHandler.Update)
	courier.PATCH("/:id/toggle", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), CourierHandler.Toggle)

	courier.GET("", rlRead.Middleware(), CourierHandler.GetAll)
}
