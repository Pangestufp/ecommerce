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

func AuthRouter(api *gin.RouterGroup) {
	AuthRepository := repository.NewAuthRepository(config.DB)
	AuthService := service.NewAuthService(AuthRepository, config.RedisClient)
	AuthHandler := handler.NewAuthHandler(AuthService)

	rlLogin := middleware.NewRateLimiter(10, 15*time.Minute)
	rlRegister := middleware.NewRateLimiter(5, time.Hour)
	rlRefresh := middleware.NewRateLimiter(10, 15*time.Minute)

	auth := api.Group("/auth")
	auth.POST("/register", rlRegister.Middleware(), AuthHandler.RegisterCustomer)
	auth.POST("/login", rlLogin.Middleware(), AuthHandler.Login)
	auth.POST("/refresh", rlRefresh.Middleware(), AuthHandler.Refresh)
}
