package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var dailyRateLimitScript = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])

if count == 1 then
    redis.call("EXPIREAT", KEYS[1], ARGV[1])
end

return count
`)

type DailyRateLimiter struct {
	redis *redis.Client
	limit int
}

func NewDailyRateLimiter(redis *redis.Client, limit int) *DailyRateLimiter {
	return &DailyRateLimiter{
		redis: redis,
		limit: limit,
	}
}

func (rl *DailyRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Default key menggunakan IP
		key := fmt.Sprintf("daily_limit:%s:ip:%s",
			c.FullPath(),
			c.ClientIP(),
		)

		// Jika user sudah login gunakan userID
		if userID, exists := c.Get("userID"); exists {
			if id, ok := userID.(string); ok {
				key = fmt.Sprintf("daily_limit:%s:user:%s",
					c.FullPath(),
					id,
				)
			}
		}

		// Reset setiap pukul 00:00 WIB
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}

		now := time.Now().In(loc)

		nextMidnight := time.Date(
			now.Year(),
			now.Month(),
			now.Day()+1,
			0, 0, 0, 0,
			loc,
		)

		expireAt := nextMidnight.Unix()

		count, err := dailyRateLimitScript.Run(
			ctx,
			rl.redis,
			[]string{key},
			expireAt,
		).Int64()

		if err != nil {
			// Fail Closed
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}

		if count > int64(rl.limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Daily request limit exceeded",
			})
			return
		}

		c.Next()
	}
}
