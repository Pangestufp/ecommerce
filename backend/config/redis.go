package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", ENV.RedisUrl, ENV.RedisPort),
		DB:   0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Cannot connect to Redis:", err)
	}

	log.Println("Connected to Redis at", ENV.RedisUrl)
}
