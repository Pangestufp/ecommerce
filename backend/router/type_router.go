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

func TypeRouter(api *gin.RouterGroup) {
	TypeRepository := repository.NewTypeRepository(config.DB)
	TypeService := service.NewTypeService(TypeRepository, config.RedisClient)
	TypeHandler := handler.NewTypeHandler(TypeService)

	Type := api.Group("/type")

	Type.Use(middleware.JWTMiddleware())

	Type.POST("", middleware.RoleMiddleware(helper.Admin()), TypeHandler.Create)
	Type.GET("", TypeHandler.GetAll)
	Type.GET("/:id", TypeHandler.GetByID)
	Type.PUT("/:id", middleware.RoleMiddleware(helper.Admin()), TypeHandler.Update)
	Type.DELETE("/:id", middleware.RoleMiddleware(helper.Admin()), TypeHandler.Delete)
}
