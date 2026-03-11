package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func (app *application) RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(5, 10)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Limite exceeded",
			})
			c.Abort()
		}
	}
}

// UserRateLimiter limits requests per authenticated user using Redis.
// Allows up to 100 requests per minute per user.
func (app *application) UserRateLimiter() gin.HandlerFunc {
	const limit = 100
	const window = time.Minute

	return func(c *gin.Context) {
		if !app.config.Redis.Enabled || app.redisClient == nil {
			c.Next()
			return
		}

		user := app.GetUserFromContext(c)
		key := fmt.Sprintf("rate:user:%d", user.Id)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		count, err := app.redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			app.redisClient.Expire(ctx, key, window)
		}

		if count > limit {
			c.Header("Retry-After", fmt.Sprintf("%.0f", window.Seconds()))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
