package main

import (
	"backend/config"
	"backend/middleware"
	"backend/router"
	"backend/server"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	config.LoadDB()

	config.ConnectRedis()

	config.ConnectMinio()

	config.ConnectElasticsearch()

	r := gin.Default()

	rl := middleware.NewRateLimiter(60, time.Minute)
	r.Use(rl.Middleware())
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")

	router.AuthRouter(api)
	router.StoreConfigRouter(api)
	router.TypeRouter(api)
	router.ProductRouter(api)
	router.ProductPriceRouter(api)
	router.DiscountRouter(api)
	router.InventoryRouter(api)
	router.CartRouter(api)

	server.Initialize(config.DB)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", config.ENV.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
