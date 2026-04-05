package router

import (
	"backend/config"
	"backend/handler"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func AuthRouter(api *gin.RouterGroup) {
	AuthRepository := repository.NewAuthRepository(config.DB)
	AuthService := service.NewAuthService(AuthRepository, config.RedisClient)
	AuthHandler := handler.NewAuthHandler(AuthService)

	api.POST("/register", AuthHandler.RegisterCustomer)
	api.POST("/login", AuthHandler.Login)
}
