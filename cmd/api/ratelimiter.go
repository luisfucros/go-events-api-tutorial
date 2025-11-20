package main

import (
	"github.com/gin-gonic/gin"

	"golang.org/x/time/rate"
	"net/http"
)

func (app *application) RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(1, 4)
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
