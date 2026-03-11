package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

func (app *application) RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(5, 10)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "rate limit exceeded",
			})
			c.Abort()
		}
	}
}

// incrWithExpireScript atomically increments a counter key and sets its TTL on
// first creation. Using a Lua script ensures INCR and EXPIRE execute as a
// single atomic operation, eliminating the race where the key could persist
// forever if the process crashed between the two separate calls.
var incrWithExpireScript = redis.NewScript(`
local count = redis.call('INCR', KEYS[1])
if count == 1 then
  redis.call('EXPIRE', KEYS[1], ARGV[1])
end
return count
`)

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

		windowSecs := int(window.Seconds())
		count, err := incrWithExpireScript.Run(ctx, app.redisClient, []string{key}, windowSecs).Int64()
		if err != nil {
			// Redis failure: degrade gracefully rather than blocking the request.
			app.logger.Warnw("user rate limiter redis error", "userId", user.Id, "error", err)
			c.Next()
			return
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
