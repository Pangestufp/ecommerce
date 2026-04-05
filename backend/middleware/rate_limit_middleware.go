package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	windowStart time.Time
	count       int
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for key, v := range rl.visitors {
			if time.Since(v.windowStart) > rl.window {
				delete(rl.visitors, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ip:" + c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			if id, ok := userID.(int); ok {
				key = "user:" + strconv.Itoa(id)
			}
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		v, exists := rl.visitors[key]
		if !exists || time.Since(v.windowStart) > rl.window {
			rl.visitors[key] = &visitor{
				windowStart: time.Now(),
				count:       1,
			}
			c.Next()
			return
		}

		if v.count >= rl.limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
			return
		}

		v.count++
		c.Next()
	}
}
