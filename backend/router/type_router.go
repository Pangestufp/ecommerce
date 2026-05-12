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

func TypeRouter(api *gin.RouterGroup) {
	TypeRepository := repository.NewTypeRepository(config.DB)
	UserRepository := repository.NewUserRepository(config.DB)
	LogRepository := repository.NewLogRepository(config.DB)
	TypeService := service.NewTypeService(TypeRepository, config.RedisClient, UserRepository, LogRepository)
	TypeHandler := handler.NewTypeHandler(TypeService)

	Type := api.Group("/type")

	Type.Use(middleware.JWTMiddleware())

	rlWrite := middleware.NewRateLimiter(20, time.Minute)
	rlRead := middleware.NewRateLimiter(60, time.Minute)

	Type.POST("", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), TypeHandler.Create)
	Type.GET("", rlRead.Middleware(), TypeHandler.GetAll)
	Type.GET("/:id", rlRead.Middleware(), TypeHandler.GetByID)
	Type.PUT("/:id", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), TypeHandler.Update)
	Type.DELETE("/:id", rlWrite.Middleware(), middleware.RoleMiddleware([]string{helper.Admin()}), TypeHandler.Delete)
}
